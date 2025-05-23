package orderlist

import (
	"context"
	"errors"
	"fmt"
	"pro_test_go/decorator"
	"pro_test_go/easy/easylistener"
	"time"
)

const (
	defaultListTimeout = 300000
	defaultNodeTimeout = 15000
)

type listEventType int

const (
	DataType listEventType = iota
	ErrorType
	EndType
)

type Func func(data any, params ...any) (any, error)
type List struct {
	ctx     context.Context
	head    *Node
	tail    *Node
	count   int
	timeout int  // ms
	block   bool // whether block

	allDoneCh chan bool
	callback  Func

	listener *easylistener.SeniorListeners[listEventType]
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
func (l *List) SetAllDoneCallback(cb Func) *List {
	l.callback = cb
	return l
}

// SetBlock ... if err, block or not
func (l *List) SetBlock(b bool) *List {
	l.block = b
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
		fmt.Println("[list] receive the last node", n.GetIndex())
		n.SetLast(true)
	}

	if l.tail == nil {
		l.tail = NewNode(nil)
	}
	l.tail = l.tail.Append(n)
	return l.Tail()
}

func (l *List) Start(data any, params ...any) {
	eventCh := easylistener.NewSeniorEventChannel[listEventType](3)
	l.listener = easylistener.NewSeniorListeners[listEventType]()
	go func() {
		_, _ = l.listener.EasyListen(eventCh, []*decorator.Action{
			easylistener.WrapSeniorListener(nil, DataType, false, false, false),
			easylistener.WrapSeniorListener(nil, ErrorType, false, false, false),
			easylistener.WrapSeniorListener(nil, EndType, false, false, false),
		})
	}()

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

func (l *List) DataHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) {

}

func (l *List) ErrorHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) {

}

func (l *List) EndHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) {

}
