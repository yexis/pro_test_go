package context

import (
	"context"
	"fmt"
	"testing"
	"time"
	"unsafe"
)

// ConcurrentCancelChan start ***********
type ConcurrentCancelChan chan bool

// ConcurrentCancelChan end *************

// ConcurrentContext start ***********
type concurrentContext struct {
	ctx context.Context
	can context.CancelFunc
}

func (cc *concurrentContext) Done() <-chan struct{} {
	return cc.ctx.Done()
}

func (cc *concurrentContext) Cancel() {
	cc.can()
}

func newConcurrentContext(fa context.Context) *concurrentContext {
	ctx, cancel := context.WithCancel(fa)
	return &concurrentContext{
		ctx,
		cancel,
	}
}

// ConcurrentContext end *************

// DegradeConcurrentContext start ***********
type DegradeConcurrentContext struct {
	concurrentContext
}

func NewDegradeConcurrentContext(fa context.Context) *DegradeConcurrentContext {
	cc := newConcurrentContext(fa)
	return (*DegradeConcurrentContext)(unsafe.Pointer(cc))
}

// DegradeConcurrentContext end *************

func TestDegradeContext(t *testing.T) {
	fa, faCan := context.WithCancel(context.Background())
	go func() {
		fmt.Println("sleep 5s")
		time.Sleep(2 * time.Second)
		faCan()
	}()

	dc := NewDegradeConcurrentContext(fa)
	go func() {
		fmt.Println("sleep 1s")
		time.Sleep(1 * time.Second)
		dc.Cancel()
	}()

	go func() {
		fmt.Println("sleep 3s")
		time.Sleep(3 * time.Second)
		dc.Cancel()
	}()

	select {
	case <-dc.Done():
		fmt.Println("ctx done")
	}
	dc.Cancel()
	faCan()
}
