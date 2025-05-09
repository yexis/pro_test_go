package orderlist

import (
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
	index int
	next  *Node
	last  bool

	nextCh    chan bool
	startCh   chan bool
	allDoneCh chan bool

	// the timeout to wait next node
	timeout int // ms
	// the callback after node done
	doneCallback Func
	// report event channel
	eventCh *easylistener.SeniorEventChannel[listEventType]
}

func NewNode(o Object, p ...any) *Node {
	last := false
	if len(p) > 0 {
		if v, ok := p[0].(bool); ok {
			last = v
		}
	}
	n := &Node{
		Object:  o,
		next:    nil,
		nextCh:  make(chan bool),
		startCh: make(chan bool),
		last:    last,
		timeout: defaultNodeTimeout,
	}
	if last {
		n.notifyCanNext()
	}
	return n
}

func (n *Node) SetTimeout(tm int) *Node {
	n.timeout = tm
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

func (n *Node) SetAllDoneCh(ch chan bool) *Node {
	n.allDoneCh = ch
	return n
}

func (n *Node) SetIndex(idx int) *Node {
	n.index = idx
	return n
}

func (n *Node) Index() int {
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
		n.notifyCanNext()
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
		case <-n.Object.Done():
			fmt.Printf("node[%d] receive object done\n", n.Index())
			n.Done(data, params...)
		case e := <-n.Object.Error():
			fmt.Printf("node[%d] receive error:%s\n", n.Index(), e.Error())
			n.Done(e, params...)
		}
	}()
}

func (n *Node) Done(data any, params ...any) {
	if n.doneCallback != nil {
		_, _ = n.doneCallback(data, params...)
	}

	timeout := time.After(time.Duration(n.timeout) * time.Millisecond)
	ct := true
	for ct {
		select {
		case <-timeout:
			if n.allDoneCh != nil {
				fmt.Printf("node[%d] wait next timeout\n", n.Index())
				n.allDoneCh <- true
				close(n.allDoneCh)
			}
			ct = false
		case <-n.nextCh:
			if n.next != nil {
				n.ToNext()
				ct = false
				break
			}
			if n.last {
				if n.allDoneCh != nil {
					n.allDoneCh <- true
					close(n.allDoneCh)
				}
				ct = false
			}
		}
	}
}

// notifyCanNext ... next node arrive
func (n *Node) notifyCanNext() {
	close(n.nextCh)
}

// ToNext ... move to next node
func (n *Node) ToNext() {
	fmt.Printf("[Node] --------------------------------------- switch node[%d] to node[%d]\n", n.Index(), n.Index()+1)
	n.next.startCh <- true
}

func (n *Node) emitEvent(key listEventType, value interface{}) {
	n.eventCh.Send(&easylistener.SeniorListenersEvent[listEventType]{
		Key:   key,
		Value: value,
	})
}
