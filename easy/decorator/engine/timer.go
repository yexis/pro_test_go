package engine

import (
	"time"
)

type TimerJob struct {
	JobBase
	a ActionI
	t time.Duration
}

func newTimerJob(c ActionI, t time.Duration) *TimerJob {
	return &TimerJob{a: c, t: t}
}

func (cj *TimerJob) Do(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	done := make(chan struct{})
	timer := time.NewTimer(cj.t)
	defer timer.Stop()

	var e error
	var i interface{} = input
	go func() {
		defer close(done)
		i, e = cj.a.Do(task, input, stage, cj.Params()...)
	}()
	select {
	case <-done:
	case <-timer.C:
	}
	return i, e
}
