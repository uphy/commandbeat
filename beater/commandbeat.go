package beater

import (
	"errors"
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/uphy/commandbeat/command"
	"github.com/uphy/commandbeat/config"
)

// Commandbeat represents the client for commandbeat
type Commandbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
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

	bt := &Commandbeat{
		done:   make(chan struct{}),
		config: c,
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
		commandSpec, err := bt.createSpec(name, &task)
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
}

func (bt *Commandbeat) createSpec(name string, task *config.TaskConfig) (*command.Spec, error) {
	c, err := task.Command()
	if err != nil {
		return nil, err
	}
	if len(c) == 0 {
		return nil, errors.New("command empty")
	}
	commandName := c[0]
	commandArgs := []string{}
	if len(c) > 0 {
		commandArgs = c[1:]
	}
	parser, err := task.Parser()
	if err != nil {
		return nil, err
	}
	return command.NewSpec(name, commandName, parser, task.Debug, commandArgs...), nil
}
