package easychannel

import "sync"

type SafeChannel struct {
	DataCh    chan interface{} // 用于传递监听事件的通道
	isClosed  bool             // 表示通道是否已关闭的标志
	needClear bool             // 关闭后是否需要清空
	lock      sync.RWMutex     // 用于同步访问的读写锁
}

func NewSafeChannel(sz int) *SafeChannel {
	if sz < 0 {
		sz = 0
	}
	return &SafeChannel{
		DataCh: make(chan interface{}, sz),
	}
}

func (sch *SafeChannel) SetClear(cl bool) *SafeChannel {
	sch.needClear = cl
	return sch
}

func (sch *SafeChannel) Close() {
	sch.lock.Lock()
	defer sch.lock.Unlock()
	if sch.isClosed {
		return
	}
	close(sch.DataCh)
	for _, ok := <-sch.DataCh; ok; {
		_, ok = <-sch.DataCh
	}
	sch.isClosed = true
}

func (sch *SafeChannel) IsClosed() bool {
	sch.lock.RLock()
	defer sch.lock.RUnlock()
	return sch.isClosed
}
