package beater

import (
	"errors"

	"github.com/robfig/cron"
	"github.com/uphy/commandbeat/config"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
)

type (
	taskScheduler struct {
		c      *cron.Cron
		runner *commandRunner
	}
)

func newTaskSchedular(client beat.Client) *taskScheduler {
	c := cron.New()
	runner := newCommandRunner(client)
	return &taskScheduler{c, runner}
}

func (t *taskScheduler) schedule(spec string, name string, task *config.TaskConfig) error {
	commandSpec, err := t.createCommandSpec(name, task)
	if err != nil {
		return err
	}
	parser, err := task.Parser()
	if err != nil {
		return err
	}
	return t.c.AddFunc(spec, func() {
		logp.Info("Running task... (name=%s)", name)
		t.runner.run(commandSpec, parser)
	})
}

func (t *taskScheduler) createCommandSpec(name string, task *config.TaskConfig) (*commandSpec, error) {
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
	return newCommand(name, commandName, task.Debug, commandArgs...), nil
}

func (t *taskScheduler) start() {
	t.c.Start()
}

func (t *taskScheduler) stop() {
	t.c.Stop()
}
