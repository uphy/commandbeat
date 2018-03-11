package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

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

	bt := &Commandbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
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

	scheduler := newTaskSchedular(bt.client)
	defer scheduler.stop()
	for name, task := range bt.config.Tasks {
		spec := task.Schedule
		if spec == "" {
			spec = defaultSchedule
		}
		if err := scheduler.schedule(spec, name, &task); err != nil {
			return err
		}
	}
	scheduler.start()

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
