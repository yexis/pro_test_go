package orderlist

import (
	"fmt"
	"time"
)

type Object interface {
	Run(data any, params ...any) error
	AllDone(data any, params ...any) error
	Done() <-chan bool
	Error() <-chan error
}

type Node struct {
	Object
	index     int
	next      *Node
	last      bool
	startCh   chan bool
	allDoneCh chan bool

	// the timeout to wait next node
	timeout int // ms
}

func NewNode(o Object, p ...any) *Node {
	last := false
	if len(p) > 0 {
		if v, ok := p[0].(bool); ok {
			last = v
		}
	}
	return &Node{
		Object:  o,
		next:    nil,
		startCh: make(chan bool),
		last:    last,
		timeout: defaultNodeTimeout,
	}
}

func (n *Node) SetTimeout(tm int) *Node {
	n.timeout = tm
	return n
}

func (n *Node) SetLast(b bool) *Node {
	n.last = b
	return n
}

func (n *Node) Append(np *Node) *Node {
	n.next = np
	return n.next
}

func (n *Node) Ready() bool {
	return <-n.startCh
}

func (n *Node) Start(data any, params ...any) {
	n.AddListener(data, params)

	go func() {
		// wait for the previous paragraph to be completed
		if n.Ready() {
			n.Run(data, params)
		}
	}()
}

func (n *Node) Run(data any, params ...any) error {
	if n.Object != nil {
		return n.Object.Run(data, params)
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
			n.Done(data, params)
		case e := <-n.Object.Error():
			fmt.Printf("node[%d] receive error:%s\n", n.Index(), e.Error())
			n.Done(data, params)
		}
	}()
}

func (n *Node) Done(data any, params any) {
	timeout := time.After(time.Duration(n.timeout) * time.Millisecond)
	ct := true
	for ct {
		select {
		case <-timeout:
			if n.allDoneCh != nil {
				fmt.Printf("node[%d] wait next timeout\n", n.Index())
				n.AllDone()
			}
			ct = false
		default:
			if n.next != nil {
				n.ToNext()
				ct = false
			} else if n.last && n.allDoneCh != nil {
				n.AllDone()
				ct = false
			}
		}
	}
}

func (n *Node) ToNext() {
	fmt.Printf("[Node] --------------------------------------- switch node[%d] to node[%d]\n", n.Index(), n.Index()+1)
	n.next.startCh <- true
}

func (n *Node) AllDone() {
	fmt.Println("[Node] --------------------------------------- all done!!!")
	n.allDoneCh <- true
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
