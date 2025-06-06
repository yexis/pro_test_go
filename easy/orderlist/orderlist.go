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
	defaultNodeCount   = 100
)

type Func func(data any, params ...any) (any, error)

type List struct {
	ctx     context.Context
	head    *Node
	tail    *Node
	count   int
	timeout int  // ms
	block   bool // whether block

	callback Func

	hook     *decorator.SeniorHook[HookKey]
	listener *easylistener.SeniorListeners[listEventType]

	ExeInfos []NodeStatistic
}

func NewList(ctx context.Context, p ...any) *List {
	var ob *Node
	if len(p) > 0 {
		if v, ok := p[0].(*Node); ok {
			ob = v
		}
	}
	if ob == nil {
		ob = NewNode(ctx, nil)
	}
	return &List{
		ctx:      ctx,
		head:     ob,
		tail:     ob,
		timeout:  defaultListTimeout,
		ExeInfos: make([]NodeStatistic, 0, defaultNodeCount),
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
	if isLast {
		fmt.Println("[list] receive the last node", n.GetIndex())
		n.SetLast(true)
	}

	if l.tail == nil {
		l.tail = NewNode(l.ctx, nil)
	}
	l.tail = l.tail.Append(n)
	return l.Tail()
}

func (l *List) Head() *Node {
	return l.head
}

func (l *List) Tail() *Node {
	return l.tail
}

func (l *List) Start(data any, params ...any) {
	eventCh := easylistener.NewSeniorEventChannel[listEventType](3)
	l.listener = easylistener.NewSeniorListeners[listEventType]()
	go func() {
		_, _ = l.listener.EasyListen(eventCh, []*decorator.Action{
			easylistener.WrapSeniorListener(l.dataHandler, DataType, false, false, false),
			easylistener.WrapSeniorListener(l.errorHandler, ErrorType, false, false, false),
			easylistener.WrapSeniorListener(l.endHandler, EndType, false, false, false),
		})
	}()

	l.hook = decorator.NewSeniorHook[HookKey]()
	l.hook.AddHook(costKey, &decorator.Action{C: l.doCostRecord, E: nil})
	l.hook.AddHook(errorKey, &decorator.Action{C: l.doErrorRecord, E: nil})

	timeout := time.After(time.Duration(l.timeout) * time.Millisecond)
	go func() {
		if l.head != nil {
			l.head.Done(data, params...)
		}
		select {
		case <-timeout:
			l.AllDone(errors.New("the list timeout"))
		case <-l.ctx.Done():
			l.AllDone(errors.New("list ctx done before all done"))
		}
	}()
}

func (l *List) AllDone(data any, params ...any) {
	if l.callback != nil {
		_, _ = l.callback(data, params...)
	}
}

// region: event handler start
func (l *List) dataHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	fmt.Printf("__orderlist__ data handler\n")
	event := input.(easylistener.SeniorListenersEvent[listEventType])
	r := event.Value.(Recorder)
	l.hook.DoHook(task, event.Value, stage, r.GetRecordKey())
	return nil, nil
}

// when error occurred
func (l *List) errorHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	event := input.(easylistener.SeniorListenersEvent[listEventType])
	err := event.Value.(error)
	fmt.Printf("__orderlist__ error handler:%s\n", err.Error())
	l.AllDone(nil, nil)
	return nil, nil
}

func (l *List) endHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	fmt.Printf("__orderlist__ end handler\n")
	l.AllDone(nil, nil)
	return nil, nil
}

// region: event handler end

// region: record start
func (l *List) doErrorRecord(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	er, ok := input.(errorRecord)
	if !ok {
		return nil, errors.New("input is not error_record")
	}
	// only node count le defaultNodeCount, record
	if er.index > defaultNodeCount {
		return nil, nil
	}
	for len(l.ExeInfos) <= er.index {
		l.ExeInfos = append(l.ExeInfos, NodeStatistic{})
	}
	l.ExeInfos[er.index].err = er.err
	fmt.Printf("do error record:%+v\n", er)

	return nil, nil
}

func (l *List) doCostRecord(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	cr, ok := input.(costRecord)
	if !ok {
		return nil, errors.New("input is not error_record")
	}
	// only node count le defaultNodeCount, record
	if cr.index > defaultNodeCount {
		return nil, nil
	}
	for len(l.ExeInfos) <= cr.index {
		l.ExeInfos = append(l.ExeInfos, NodeStatistic{})
	}
	l.ExeInfos[cr.index].ms = cr.ms
	fmt.Printf("do cost record:%+v\n", cr)
	return nil, nil
}

// region: record end
