package waitingqueue

import (
	"context"

	"github.com/yexis/pro_test_go/easy/easylogger"
	"github.com/yexis/pro_test_go/easy/easystate"

	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type QueueElement interface {
	GetDelay() time.Duration
	GetCallback() func()
	GetBeforeAfter() int8
}

type QueueElementImpl struct {
	Delay       time.Duration
	Callback    func()
	BeforeAfter int8 // 0: before 1: after
}

func (qe *QueueElementImpl) GetDelay() time.Duration { return qe.Delay }
func (qe *QueueElementImpl) GetCallback() func()     { return qe.Callback }
func (qe *QueueElementImpl) GetBeforeAfter() int8    { return qe.BeforeAfter }

type Status int32

const (
	NotRunning Status = iota
	Running
	PushDone
	DealDone
	Closed
)

type QueueState struct {
	State string
	Time  int64
	Cost  int64
}

type WaitingQueue struct {
	ctx    context.Context
	cancel context.CancelFunc
	q      chan QueueElement
	lock   sync.Mutex
	rst    []QueueElement
	ticker time.Duration

	dealDone chan struct{}
	state    atomic.Int32
	stateLog *easystate.EasyState[Status]

	delayTime time.Duration
	delayCost time.Duration

	easylogger.LoggerWriter
}

func NewWaitingQueue() *WaitingQueue {
	ctx, cancel := context.WithCancel(context.Background())
	wq := &WaitingQueue{
		ctx:      ctx,
		cancel:   cancel,
		q:        make(chan QueueElement, 16),
		rst:      make([]QueueElement, 0, 64),
		ticker:   24 * time.Millisecond,
		dealDone: make(chan struct{}),
	}
	wq.state.Store(int32(NotRunning))
	wq.stateLog = easystate.NewEasyState(
		func(state Status) string {
			return strconv.Itoa(int(state))
		},
	).AddState(NotRunning)
	return wq.run()
}

func (ap *WaitingQueue) SetTicker(ticker time.Duration) *WaitingQueue {
	ap.ticker = ticker
	return ap
}

func (ap *WaitingQueue) run() *WaitingQueue {
	ap.trySet(NotRunning, Running)
	go func() {
		startTime := time.Now()
		defer func() {
			ap.delayCost = time.Since(startTime)
			close(ap.dealDone)
		}()

		ticker := time.NewTicker(ap.ticker)
		defer ticker.Stop()

		for {
			select {
			case e := <-ap.q:
				if dur := e.GetDelay(); e.GetBeforeAfter() == 0 && dur > 0 {
					ap.Logger("__wq__ start before sleep:%v", e.GetDelay())
					time.Sleep(dur)
					ap.Logger("__wq__ end before sleep")
				}
				if cb := e.GetCallback(); cb != nil {
					cb()
				}
				if dur := e.GetDelay(); e.GetBeforeAfter() == 1 && dur > 0 {
					ap.Logger("__wq__ start after sleep:%v", e.GetDelay())
					time.Sleep(dur)
					ap.Logger("__wq__ end after sleep")
				}

			case <-ticker.C:
				if ap.drain() {
					return
				}
			case <-ap.ctx.Done():
				return
			}
		}
	}()
	return ap
}

// drain ...
// trans data from rst to q
func (ap *WaitingQueue) drain() bool {
	ap.lock.Lock()
	defer ap.lock.Unlock()

	running := true
	for len(ap.rst) > 0 && running {
		select {
		case ap.q <- ap.rst[0]:
			ap.rst = ap.rst[1:]
		default:
			running = false
			break
		}
	}
	if len(ap.q) == 0 && len(ap.rst) == 0 {
		if ap.trySet(PushDone, DealDone) {
			return true
		}
	}
	return false
}

func (ap *WaitingQueue) finish() {}

func (ap *WaitingQueue) trySet(old, new Status) bool {
	did := ap.state.CompareAndSwap(int32(old), int32(new))
	if did {
		ap.stateLog.AddState(new)
	}
	return did
}

func (ap *WaitingQueue) Push(result QueueElement) {
	ap.lock.Lock()
	ap.rst = append(ap.rst, result)
	ap.delayTime += result.GetDelay()
	ap.lock.Unlock()
}

func (ap *WaitingQueue) Cancel() {
	ap.cancel()
	ap.finish()
}

// PushDone ...
// tell queue that push done
func (ap *WaitingQueue) PushDone() {
	ap.trySet(Running, PushDone)
}

// AskDone ...
// ask queue whether deal done
func (ap *WaitingQueue) AskDone(dur ...time.Duration) *WaitingQueue {
	<-ap.dealDone
	ap.trySet(DealDone, Closed)
	return ap
}

func (ap *WaitingQueue) GetStates() QueueState {
	return QueueState{
		State: ap.stateLog.GetState(),
		Time:  ap.delayTime.Milliseconds(),
		Cost:  ap.delayCost.Milliseconds(),
	}
}
