package orderlist

import (
	"context"
	"errors"
	"fmt"
	"pro_test_go/easy/easylistener"
	"time"
)

type Object interface {
	Run(data any, params ...any) error
	Done() <-chan bool
	Error() <-chan error
}

type Node struct {
	Object
	ctx   context.Context
	index int
	last  bool
	next  *Node

	nextArrivedCh chan bool
	startCh       chan bool

	// the waitNextTimeout to wait next node
	waitNextTimeout int // ms
	// the callback after node done
	doneCallback Func
	// report event channel
	eventCh *easylistener.SeniorEventChannel[listEventType]
}

func NewNode(ctx context.Context, o Object, p ...any) *Node {
	last := false
	if len(p) > 0 {
		if v, ok := p[0].(bool); ok {
			last = v
		}
	}
	n := &Node{
		ctx:             ctx,
		Object:          o,
		next:            nil,
		nextArrivedCh:   make(chan bool),
		startCh:         make(chan bool),
		last:            last,
		waitNextTimeout: defaultNodeTimeout,
	}
	if last {
		n.notifyNextArrived()
	}
	return n
}

func (n *Node) SetTimeout(tm int) *Node {
	n.waitNextTimeout = tm
	return n
}

func (n *Node) SetLast(b bool) *Node {
	n.last = b
	return n
}

func (n *Node) SetDoneCallback(cb Func) *Node {
	n.doneCallback = cb
	return n
}

func (n *Node) SetIndex(idx int) *Node {
	n.index = idx
	return n
}

func (n *Node) GetIndex() int {
	return n.index
}

// Append ... append a node on tail
func (n *Node) Append(np *Node) *Node {
	// last node is not allowed to append
	if n.last {
		return n
	}

	n.next = np
	// tell to next node arrived
	if np != nil {
		n.notifyNextArrived()
	}
	return n.next
}

func (n *Node) JudgeReady() bool {
	return <-n.startCh
}

func (n *Node) Start(data any, params ...any) {
	n.AddListener(data, params)

	go func() {
		// wait for the previous node to be completed
		if n.JudgeReady() {
			_ = n.Run(data, params)
		}
	}()
}

func (n *Node) Run(data any, params ...any) error {
	if n.Object != nil {
		return n.Object.Run(data, params...)
	}
	return nil
}

func (n *Node) AddListener(data any, params ...any) {
	if n.Object == nil {
		return
	}

	go func() {
		select {
		case <-n.ctx.Done():
			fmt.Printf("__node[%d]__ ctx done and exit before object done\n", n.GetIndex())
			break
		case <-n.Object.Done():
			fmt.Printf("__node[%d]__ receive object done\n", n.GetIndex())
			n.Done(data, params...)
		case e := <-n.Object.Error():
			fmt.Printf("__node[%d]__ receive error:%s\n", n.GetIndex(), e.Error())
			n.Done(e, params...)
		}
	}()
}

func (n *Node) Done(data any, params ...any) {
	if n.doneCallback != nil {
		_, _ = n.doneCallback(data, params...)
	}

	// next arrived timeout
	dur := time.Duration(n.waitNextTimeout)
	C := time.After(dur * time.Millisecond)

	select {
	case <-n.ctx.Done():
		fmt.Printf("__node[%d]__ ctx done when wait next node\n", n.GetIndex())
		n.emitEvent(ErrorType, errors.New("ctx done when wait next node"))

	case <-C:
		fmt.Printf("__node[%d]__ timeout when wait next node\n", n.GetIndex())
		n.emitEvent(ErrorType, errors.New("timeout when wait next node"))

	case <-n.nextArrivedCh:
		if n.next != nil {
			n.toNext()
			break
		}
		if n.last {
			n.emitEvent(EndType, nil)
		}
	}

}

// notifyNextArrived ... next node arrived
func (n *Node) notifyNextArrived() {
	close(n.nextArrivedCh)
}

// toNext ... move to next node
func (n *Node) toNext() {
	fmt.Printf("__node[%d]__ switch to node[%d]\n", n.GetIndex(), n.GetIndex()+1)
	n.next.startCh <- true
}

func (n *Node) emitEvent(key listEventType, value interface{}) {
	if n.eventCh == nil {
		return
	}
	n.eventCh.Send(&easylistener.SeniorListenersEvent[listEventType]{
		Key:   key,
		Value: value,
	})
}
