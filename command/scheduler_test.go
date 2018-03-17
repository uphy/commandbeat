package command

import (
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	handler := &testHandler{}
	runner := NewRunner(handler)
	spec := &testCommand{"echo", []string{"aaa"}}
	scheduler := NewScheduler(runner)
	if err := scheduler.Schedule("@every 1s", spec); err != nil {
		t.Error(err)
	}
	scheduler.Start()
	defer scheduler.Stop()
	ticker := time.NewTicker(time.Second * 2)
	select {
	case <-ticker.C:
		if "aaa" != handler.stdout {
			t.Error("not scheduled.")
		}
	}
}

func TestSchedulerError(t *testing.T) {
	handler := &testHandler{}
	runner := NewRunner(handler)
	spec := &testCommand{"echoaaaaaa", []string{"aaa"}}
	scheduler := NewScheduler(runner)
	if err := scheduler.Schedule("@every 1s", spec); err != nil {
		t.Error(err)
	}
	scheduler.Start()
	defer scheduler.Stop()
	ticker := time.NewTicker(time.Second * 2)
	select {
	case <-ticker.C:
		if handler.handledErr == nil {
			t.Error("error not returned.")
		}
	}
}
