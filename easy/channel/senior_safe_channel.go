package channel

import "sync"

type SeniorSafeChannel[T any] struct {
	DataCh    chan T       // 用于传递监听事件的通道
	isClosed  bool         // 表示通道是否已关闭的标志
	needClear bool         // 关闭后是否需要清空
	lock      sync.RWMutex // 用于同步访问的读写锁
}

func NewSeniorSafeChannel[T any](sz int) *SeniorSafeChannel[T] {
	if sz < 0 {
		sz = 0
	}
	return &SeniorSafeChannel[T]{
		DataCh: make(chan T, sz),
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
	close(sch.DataCh)
	sch.isClosed = true
	if sch.needClear {
		for _, ok := <-sch.DataCh; ok; {
			_, ok = <-sch.DataCh
		}
	}
}

func (sch *SeniorSafeChannel[T]) IsClosed() bool {
	sch.lock.RLock()
	defer sch.lock.RUnlock()
	return sch.isClosed
}
