package listener

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

func dissectEventType(d interface{}) (EventType, bool) {
	var et EventType
	var ok bool
	switch v := d.(type) {
	case EventType:
		et = v
		ok = true
	case string:
		et = EventType(v)
		ok = true
	}
	return et, ok
}

const (
	listenerIEVClose   = 0
	listenerIEVRemove  = 1
	listenerIEVClear   = 2
	listenerIEVDestroy = 3
)

const (
	listenerStReady     = 0
	listenerStToRelease = 1
	listenerStReleased  = 2
)

const (
	defaultTimeoutMs = 300000 // ms
)

// ListenersEvent ... listener event
type ListenersEvent struct {
	Key   EventType
	Value interface{}
}

type ListenersEventArgs struct {
	Once    bool // run once
	Async   bool // run sync or async
	EndLoop bool // run then end
}

type listenersEventAction struct {
	decorator.Action
	Key EventType // event key
	ListenersEventArgs
}

// listenersInternalEventAction ... listener internal events action
type listenersInternalEventAction struct {
	Type int
	Data interface{}
}

// Listeners ... Listener
type Listeners struct {
	running bool
	status  int
	timeout int
	ievChan chan *listenersInternalEventAction
	lock    sync.Mutex
}

// SetTimeout ... set timeout
func (l *Listeners) SetTimeout(t int) *Listeners {
	l.timeout = t
	return l
}

func WrapListener(c decorator.Ctrl, event EventType, once bool, endLoop bool, async bool, params ...interface{}) *decorator.Action {
	return &decorator.Action{
		C: c,
		P: append([]interface{}{
			event,
			&ListenersEventArgs{
				Once:    once,
				EndLoop: endLoop,
				Async:   async,
			},
		}, params...),
	}
}

func WrapDefaultListener(c decorator.Ctrl, event EventType, params ...interface{}) *decorator.Action {
	return &decorator.Action{
		C: c,
		P: append([]interface{}{
			event,
			&ListenersEventArgs{},
		}, params...),
	}
}

// EasyListen ... easy listen
func (l *Listeners) EasyListen(ch *EventChannel, actions []*decorator.Action) (interface{}, error) {
	return l.Listen(nil, nil, nil, ch, actions)
}

// Listen
// ch chan *ListenersEvent: read-only?
func (l *Listeners) Listen(task *decorator.Task, input interface{}, ps *decorator.Stage,
	ch *EventChannel, actions []*decorator.Action) (interface{}, error) {
	if len(actions) <= 0 {
		return nil, errors.New(decorator.EM1301EmptyAction)
	}
	if !l.Prepare() {
		return nil, errors.New("prepare failed")
	}

	var ctx context.Context
	if task == nil || task.Context != nil {
		ctx = task.Context
	} else {
		ctx = context.Background()
	}

	// listener 超时逻辑
	if l.timeout == 0 {
		l.timeout = defaultTimeoutMs
	}

	tm := time.After(time.Duration(l.timeout) * time.Millisecond)
	mp := make(map[EventType]*listenersEventAction)
	var err error
	var ret interface{}
	for _, action := range actions {
		size := len(action.P)
		if size < 2 {
			return nil, errors.New(decorator.EM1303MissingParams)
		}
		var key EventType
		var ok bool
		key, ok = dissectEventType(action.P[0])
		if !ok {
			return nil, fmt.Errorf(decorator.EM1305WrongParams, "listen action params[0]")
		}
		args, ok := action.P[1].(*ListenersEventArgs)
		if !ok {
			return nil, fmt.Errorf(decorator.EM1305WrongParams, "listen action params[1]")
		}
		if key == "" {
			continue
		}
		la := &listenersEventAction{
			Action: decorator.Action{
				C: action.C,
				P: action.P[2:],
				E: action.E,
			},
			Key:                key,
			ListenersEventArgs: *args,
		}
		mp[key] = la
	}

	for l.running {
		select {
		case iev := <-l.ievChan:
			if iev.Type == listenerIEVClose {
				ret = iev.Data
				l.Stop(ch)
			} else if iev.Type == listenerIEVRemove {
				et, ok := dissectEventType(iev.Data)
				if ok {
					delete(mp, et)
				}
			} else if iev.Type == listenerIEVClear {
				mp = map[EventType]*listenersEventAction{}
			} else if iev.Type == listenerIEVDestroy {
				mp = map[EventType]*listenersEventAction{}
				l.Stop(ch)
			}
		case event := <-ch.DataChan:
			if event == nil {
				break
			}
			action, exist := mp[event.Key]
			if !exist {
				break
			}
			if action.Once {
				delete(mp, event.Key)
			}
			if action.Async {
				go func() {
					l.doEvent(task, event, ps, action)
					if action.EndLoop {
						l.Stop(ch)
					}
				}()
			} else {
				l.doEvent(task, event, ps, action)
				if action.EndLoop {
					l.Stop(ch)
				}
			}
		case <-tm:
			l.Stop(ch)
		case <-ctx.Done():
			l.Stop(ch)
		}
	}

	switch ret.(type) {
	case error:
		err = ret.(error)
		ret = nil
	default:
		break
	}
	return ret, err
}

func (l *Listeners) doEvent(task *decorator.Task, input interface{}, ps *decorator.Stage,
	a *listenersEventAction) (interface{}, error) {
	i, e := a.C(task, input, ps, a.P...)
	if e != nil && a.E != nil {
		i, e = a.E(task, e, ps, a.P...)
	}
	// i != nil as finished
	if i != nil {
		l.End(i)
	} else if e != nil {
		l.End(e)
	}
	return i, e
}

func (l *Listeners) Prepare() bool {
	defer l.lock.Unlock()
	l.lock.Lock()
	if l.status > listenerStReady {
		return false
	}
	l.running = true
	l.ievChan = make(chan *listenersInternalEventAction, 3)
	return true
}

func (l *Listeners) Destroy() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.ievChan != nil {
		l.ievChan <- &listenersInternalEventAction{
			Type: listenerIEVDestroy,
			Data: nil,
		}
	}
}

func (l *Listeners) Stop(ch *EventChannel) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.status == listenerStReleased {
		return
	}
	l.running = false
	l.status = listenerStReleased
	if l.ievChan != nil {
		close(l.ievChan)
	}
}

// End ... end
func (l *Listeners) End(data interface{}) {
	defer l.lock.Unlock()
	l.lock.Lock()
	if l.status < listenerStToRelease {
		l.status = listenerStToRelease
	}
	if l.ievChan != nil {
		l.ievChan <- &listenersInternalEventAction{
			Type: listenerIEVClose,
			Data: data,
		}
	}
}
