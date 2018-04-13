package command

import (
	"github.com/robfig/cron"
)

type (
	// Scheduler schedules command execution.
	Scheduler struct {
		c      *cron.Cron
		runner *Runner
	}
)

// NewScheduler creates a new Scheduler.
func NewScheduler(runner *Runner) *Scheduler {
	return &Scheduler{cron.New(), runner}
}

// Schedule schedules command execution.
func (t *Scheduler) Schedule(scheduleSpec string, commandSpec Spec) error {
	return t.c.AddFunc(scheduleSpec, func() {
		t.runner.Run(commandSpec)
	})
}

// Start starts the Scheduler.
func (t *Scheduler) Start() {
	t.c.Start()
}

// Stop stops the Scheduler.
func (t *Scheduler) Stop() {
	t.c.Stop()
}
