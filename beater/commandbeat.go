package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/uphy/commandbeat/command"
	"github.com/uphy/commandbeat/config"
)

// Commandbeat represents the client for commandbeat
type Commandbeat struct {
	done          chan struct{}
	config        config.Config
	client        beat.Client
	scriptManager *scriptManager
}

const (
	defaultSchedule string = "@every 1m"
)

// New creates new beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	if err := parseConfig(&c); err != nil {
		return nil, fmt.Errorf("Error parsing config: %v", err)
	}
	scriptManager, err := newScriptManager(".scripts")
	if err != nil {
		return nil, fmt.Errorf("Script manager cannot be created: %v", err)
	}

	bt := &Commandbeat{
		done:          make(chan struct{}),
		config:        c,
		scriptManager: scriptManager,
	}
	return bt, nil
}

func parseConfig(c *config.Config) error {
	if c.Debug {
		for name, task := range c.Tasks {
			task.Debug = true
			c.Tasks[name] = task
		}
	}
	return nil
}

// Run starts the commandbeat application.
// This function blocks until Stop() call.
func (bt *Commandbeat) Run(b *beat.Beat) error {
	logp.Info("commandbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	runner := command.NewRunner(newPublishHandler(newElasticPublisher(bt.client)))
	scheduler := command.NewScheduler(runner)
	defer scheduler.Stop()
	for name, task := range bt.config.Tasks {
		scheduleSpec := task.Schedule
		if scheduleSpec == "" {
			scheduleSpec = defaultSchedule
		}
		commandSpec, err := bt.newSpec(name, &task)
		if err != nil {
			return err
		}
		if err := scheduler.Schedule(scheduleSpec, commandSpec); err != nil {
			return err
		}
	}
	scheduler.Start()

	for {
		select {
		case <-bt.done:
			return nil
		}
	}
}

// Stop stops the application.
func (bt *Commandbeat) Stop() {
	bt.client.Close()
	close(bt.done)

	logp.Info("Cleaning the work directory for script manager.")
	if err := bt.scriptManager.clean(); err != nil {
		logp.Err("script manager cannot be cleaned up. (dir=%s, err=%v)", bt.scriptManager.directory, err)
	}
}
