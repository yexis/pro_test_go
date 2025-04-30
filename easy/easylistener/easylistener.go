package easylistener

import (
	"context"
	"errors"
	"fmt"
	"pro_test_go/decorator"
	"sync"
	"time"
)

// EventType ... event type
// data, error, close, end ...
type EventType string

const (
	listenerIEVClose   = 0
	listenerIEVRemove  = 1
	listenerIEVClear   = 2
	listenerIEVDestroy = 3
)

const (
	defaultTimeout = 300000 // ms
)

// ListenEvent ... listener event
type ListenEvent struct {
	Key   EventType
	Value interface{}
}

type ListenEventArgs struct {
	Once    bool // run once
	Async   bool // run sync or async
	EndLoop bool // run then end
}

type ListenerEventAction struct {
	decorator.Action
	Key EventType // event key
	ListenEventArgs
}

// listenerInternalEventAction ... listener internal events action
type listenerInternalEventAction struct {
	Type int
	Data interface{}
}

// Listeners ... Listener
type Listeners struct {
	running bool
	status  int
	ievChan chan *listenerInternalEventAction
	lock    sync.Mutex

	Timeout int
}

func (l *Listeners) Prepare() {
	l.running = true
}

// SetTimeout ... set timeout
func (l *Listeners) SetTimeout(t int) *Listeners {
	l.Timeout = t
	return l
}

func WrapListener(c decorator.Ctrl, event string, once bool, endLoop bool, async bool, params ...interface{}) *decorator.Action {
	return &decorator.Action{
		C: c,
		P: append([]interface{}{
			event,
			&ListenEventArgs{
				Once:    once,
				EndLoop: endLoop,
				Async:   async,
			},
		}, params...),
	}
}

func (l *Listeners) doEvent(task *decorator.Task, input interface{}, ps *decorator.Stage,
	a *ListenerEventAction) (interface{}, error) {
	i, e := a.C(task, input, ps, a.P...)
	if e != nil && a.E != nil {
		i, e = a.E(task, e, ps, a.P...)
	}
	if i != nil {
		//l.End(i)
	} else if e != nil {
		//l.End(e)
	}
	return i, e
}

// EasyListen ... easy listen
func (l *Listeners) EasyListen(eventChan chan *ListenEvent, actions []*decorator.Action) (interface{}, error) {
	return l.Listen(nil, nil, nil, eventChan, actions)
}

// Listen
// ch chan *ListenEvent: read-only?
func (l *Listeners) Listen(task *decorator.Task, input interface{}, ps *decorator.Stage,
	ch chan *ListenEvent, actions []*decorator.Action) (interface{}, error) {
	if len(actions) <= 0 {
		return nil, nil
	}
	l.Prepare()

	var ctx context.Context
	if task == nil || task.Context != nil {
		ctx = task.Context
	} else {
		ctx = context.Background()
	}

	// listener 超时逻辑
	if l.Timeout == 0 {
		l.Timeout = defaultTimeout
	}
	tm := time.After(time.Duration(l.Timeout) * time.Millisecond)

	mp := make(map[EventType]*ListenerEventAction)
	var err error
	var ret interface{}

	for _, action := range actions {
		size := len(action.P)
		if size < 2 {
			return nil, errors.New(decorator.EM1303MissingParams)
		}
		var key EventType
		switch v := action.P[0].(type) {
		case string:
			key = EventType(v)
		case EventType:
			key = v
		default:
			return nil, fmt.Errorf(decorator.EM1305WrongParams, "listen action params[0]")
		}
		args, ok := action.P[1].(*ListenEventArgs)
		if !ok {
			return nil, fmt.Errorf(decorator.EM1305WrongParams, "listen action params[1]")
		}
		if key == "" {
			continue
		}
		la := &ListenerEventAction{
			Action: decorator.Action{
				C: action.C,
				P: action.P[2:],
				E: action.E,
			},
			Key:             key,
			ListenEventArgs: *args,
		}
		mp[key] = la
	}

	for l.running {
		select {
		case iev := <-l.ievChan:
			// end
			if iev.Type == listenerIEVDestroy {
				mp = map[EventType]*ListenerEventAction{}
				l.Stop(ch)
			}
		case ev := <-ch:
			if ev == nil {
				break
			}
			a, ok := mp[ev.Key]
			if !ok {
				break
			}
			if a.Once {
				delete(mp, ev.Key)
			}
			if a.Async {
				go l.doEvent(task, input, ps, a)
				if a.EndLoop {

				}
			} else {
				go l.doEvent(task, input, ps, a)
				if a.EndLoop {

				}
			}
		case <-tm:
			l.Stop(ch)
		case <-ctx.Done():
			l.Stop(ch)
		}
	}
	return ret, err
}

func (l *Listeners) Destroy() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.ievChan != nil {
		l.ievChan <- &listenerInternalEventAction{
			Type: listenerIEVDestroy,
			Data: nil,
		}
	}
}

func (l *Listeners) Stop(ch chan *ListenEvent) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.running = false
	if l.ievChan != nil {
		close(l.ievChan)
	}
}
