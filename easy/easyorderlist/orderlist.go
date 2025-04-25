package easyorderlist

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	defaultListTimeout = 300000
	defaultNodeTimeout = 15000
)

type AllDoneCallback func(data any, params ...any) (any, error)
type List struct {
	ctx     context.Context
	head    *Node
	tail    *Node
	count   int
	timeout int // ms

	allDoneCh chan bool
	callback  AllDoneCallback
}

func NewList(ctx context.Context, p ...any) *List {
	var ob *Node
	if len(p) > 0 {
		if v, ok := p[0].(*Node); ok {
			ob = v
		}
	}
	if ob == nil {
		ob = NewNode(nil)
	}
	return &List{
		ctx:       ctx,
		head:      ob,
		tail:      ob,
		allDoneCh: make(chan bool),
		timeout:   defaultListTimeout,
	}
}

func (l *List) SetTimeout(tm int) *List {
	l.timeout = tm
	return l
}

// SetAllDoneCallback
// e.g. NewList(n).SetAllDoneCallback(f)
func (l *List) SetAllDoneCallback(cb AllDoneCallback) *List {
	l.callback = cb
	return l
}

func (l *List) Append(n *Node, p ...any) *Node {
	var isLast bool
	if len(p) > 0 {
		if v, ok := p[0].(bool); ok {
			isLast = v
		}
	}

	l.count++
	n.SetIndex(l.count)
	n.SetAllDoneCh(l.allDoneCh)
	if isLast {
		fmt.Println("[list] receive the last node", n.Index())
		n.SetLast(true)
	}

	if l.tail == nil {
		l.tail = NewNode(nil)
	}
	l.tail = l.tail.Append(n)
	return l.Tail()
}

func (l *List) Start(data any, params ...any) {
	timeout := time.After(time.Duration(l.timeout) * time.Millisecond)
	go func() {
		if l.head != nil {
			l.head.Done(data, params)
		}
		select {
		case <-l.allDoneCh:
			l.AllDone(data, params)
		case <-timeout:
			l.AllDone(errors.New("the list timeout"))
		case <-l.ctx.Done():
			l.AllDone(errors.New("list ctx done before all done"))
		}
	}()
}

func (l *List) AllDone(data any, params ...any) {
	if l.callback != nil {
		_, _ = l.callback(data, params)
	}
}

func (l *List) Head() *Node {
	return l.head
}

func (l *List) Tail() *Node {
	return l.tail
}
