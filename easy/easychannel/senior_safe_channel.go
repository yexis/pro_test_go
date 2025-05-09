package easychannel

import (
	"fmt"
	"sync"
)

// SeniorSafeChannel ...
// safe channel that supports close check
// and do not close data-chan outside even thought it is exported
type SeniorSafeChannel[T any] struct {
	DataChan  chan T       // 用于传递监听事件的通道
	isClosed  bool         // 表示通道是否已关闭的标志
	needClear bool         // 关闭后是否需要清空
	lock      sync.RWMutex // 用于同步访问的读写锁
}

func NewSeniorSafeChannel[T any](sz int) *SeniorSafeChannel[T] {
	if sz < 0 {
		sz = 0
	}
	return &SeniorSafeChannel[T]{
		DataChan: make(chan T, sz),
	}
}

func (sch *SeniorSafeChannel[T]) SetClear(cl bool) *SeniorSafeChannel[T] {
	sch.needClear = cl
	return sch
}

func (sch *SeniorSafeChannel[T]) Close() {
	sch.lock.Lock()
	defer sch.lock.Unlock()
	if sch.isClosed {
		return
	}
	fmt.Println("__senior_safe_channel__ close")
	if sch.needClear {
		for {
			select {
			case <-sch.DataChan:
				// discard
			default:
				close(sch.DataChan)
				sch.isClosed = true
				return
			}
		}
	}
}

func (sch *SeniorSafeChannel[T]) IsClosed() bool {
	sch.lock.RLock()
	defer sch.lock.RUnlock()
	return sch.isClosed
}

func (sch *SeniorSafeChannel[T]) Send(data T) bool {
	sch.lock.RLock()
	defer sch.lock.RUnlock()
	if sch.isClosed {
		return false
	}
	sch.DataChan <- data
	return true
}

func (sch *SeniorSafeChannel[T]) Receive() <-chan T {
	return sch.DataChan
}
