package command

import (
	"errors"

	"github.com/robfig/cron"
	"github.com/uphy/commandbeat/config"
)

type (
	Scheduler struct {
		c      *cron.Cron
		runner *Runner
	}
)

func NewScheduler(runner *Runner) *Scheduler {
	return &Scheduler{cron.New(), runner}
}

func (t *Scheduler) Schedule(scheduleSpec string, name string, task *config.TaskConfig) error {
	spec, err := t.createSpec(name, task)
	if err != nil {
		return err
	}
	return t.c.AddFunc(scheduleSpec, func() {
		t.runner.Run(spec)
	})
}

func (t *Scheduler) createSpec(name string, task *config.TaskConfig) (*Spec, error) {
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
	return NewSpec(name, commandName, parser, task.Debug, commandArgs...), nil
}

func (t *Scheduler) Start() {
	t.c.Start()
}

func (t *Scheduler) Stop() {
	t.c.Stop()
}
