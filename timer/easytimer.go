package timer

import (
	"sync"
	"time"
)

const (
	TimerOnlyCheck        = -1
	TimerReady            = 0
	TimerStopped          = 1
	TimerStoppedOnTimeout = 2
	TimerStoppedByCannel  = 3
	TimerStoppedByClear   = 4
)

// EasyTimerCallback ... req timer method
type EasyTimerCallback func(rt *EasyTimer, args ...interface{}) error

// EmptyEasyTimerCallback ... 用于无需回调的场景
func EmptyEasyTimerCallback(rt *EasyTimer, args ...interface{}) error {
	return nil
}

var PureEmptyEasyTimerCallback EasyTimerCallback = nil

// EasyTimer ... request timer
type EasyTimer struct {
	errCh  chan error
	timer  *time.Timer
	lock   sync.Mutex
	status int
}

// Stop ... stop
// option ...
// -1 ... not stop and fast check;
// 0 ... not stop & normal check;
// 2 ... do stop;
// 3 ... do stop & timer
// 4 ... do stop & timer
func (et *EasyTimer) Stop(option int) (int, bool) {
	defer et.lock.Unlock()
	et.lock.Lock()

	if option == TimerOnlyCheck {
		return et.status, false
	}

	if et.status >= TimerStopped {
		return et.status, false
	}

	n := false
	if option == TimerStoppedOnTimeout || option == TimerStopped {
		et.status = option
		n = true
	} else if option == TimerStoppedByCannel || option == TimerStoppedByClear {
		et.status = option
		if et.timer != nil {
			et.timer.Stop()
			et.errCh <- nil
		}
		n = true
	}
	return et.status, n
}

// Reset ... reset if not stopped
func (et *EasyTimer) Reset(interval time.Duration) bool {
	b := false
	s, _ := et.Stop(TimerOnlyCheck)
	if s >= TimerStopped {
		return b
	}
	if et.timer != nil {
		b = et.timer.Reset(interval)
	}
	return b
}

// Start ... start
func (et *EasyTimer) Start(option int, interval time.Duration, method EasyTimerCallback, args ...interface{}) {
	if stopped, _ := et.Stop(option); stopped >= TimerStopped {
		return
	}

	et.errCh = make(chan error, 1)
	et.timer = time.AfterFunc(interval, func() {
		if _, curStopped := et.Stop(TimerStoppedOnTimeout); curStopped {
			var err error
			if method != nil {
				err = method(et, args...)
			}
			et.errCh <- err
		}
	})
}

func (et *EasyTimer) C() <-chan error {
	return et.errCh
}
