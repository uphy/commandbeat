package beater

import (
	"errors"

	"github.com/robfig/cron"
	"github.com/uphy/commandbeat/command"
	"github.com/uphy/commandbeat/config"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
)

type (
	taskScheduler struct {
		c      *cron.Cron
		runner *command.CommandRunner
	}
)

func newTaskSchedular(client beat.Client) *taskScheduler {
	c := cron.New()
	runner := command.NewCommandRunner(NewPublishHandler(newElasticsearchPublisher(client)))
	return &taskScheduler{c, runner}
}

func (t *taskScheduler) schedule(spec string, name string, task *config.TaskConfig) error {
	commandSpec, err := t.createCommandSpec(name, task)
	if err != nil {
		return err
	}
	return t.c.AddFunc(spec, func() {
		logp.Info("Running task... (name=%s)", name)
		t.runner.Run(commandSpec)
	})
}

func (t *taskScheduler) createCommandSpec(name string, task *config.TaskConfig) (*command.Spec, error) {
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
	return command.NewCommand(name, commandName, parser, task.Debug, commandArgs...), nil
}

func (t *taskScheduler) start() {
	t.c.Start()
}

func (t *taskScheduler) stop() {
	t.c.Stop()
}
