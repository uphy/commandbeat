package command

import (
	"log"

	"github.com/robfig/cron"
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

func (t *Scheduler) Schedule(scheduleSpec string, commandSpec Spec) error {
	return t.c.AddFunc(scheduleSpec, func() {
		if err := t.runner.Run(commandSpec); err != nil {
			// TODO error handling
			log.Printf("failed to run command. (scheduleSpec=%s, commandSpec=%v, err=%v)", scheduleSpec, commandSpec, err)
		}
	})
}

func (t *Scheduler) Start() {
	t.c.Start()
}

func (t *Scheduler) Stop() {
	t.c.Stop()
}
