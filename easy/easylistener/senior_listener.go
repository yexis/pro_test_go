package easylistener

import (
	"context"
	"errors"
	"fmt"
	"pro_test_go/decorator"
	"sync"
	"time"
)

type EventTypeParser[K comparable] func(interface{}) (K, bool)

func TplEventTypeParser[K comparable](val interface{}) (K, bool) {
	result, ok := val.(K)
	return result, ok
}

// SeniorListenersEvent ... senior listener event
type SeniorListenersEvent[K comparable] struct {
	Key   K
	Value interface{}
}

// SeniorListenersEventArgs ... senior listener event args
type SeniorListenersEventArgs struct {
	Once    bool // run once
	Async   bool // run sync or async, not used now
	EndLoop bool // run then end
}

type seniorListenersEventAction[K comparable] struct {
	decorator.Action
	Key K // event key
	SeniorListenersEventArgs
}

// listenersInternalEventAction ... listener internal events action
type seniorListenersInternalEventAction struct {
	Type int
	Data interface{}
}

// SeniorListeners ... senior listener
type SeniorListeners[K comparable] struct {
	running   bool
	status    int
	timeoutMs int
	ievChan   chan *seniorListenersInternalEventAction
	lock      sync.Mutex
}

// SetTimeout ... set timeoutMs
func (l *SeniorListeners[K]) SetTimeout(t int) *SeniorListeners[K] {
	l.timeoutMs = t
	return l
}

func WrapSeniorListener[K comparable](c decorator.Ctrl, tp K, once bool, endLoop bool, async bool, params ...interface{}) *decorator.Action {
	return &decorator.Action{
		C: c,
		P: append([]interface{}{
			tp,
			&SeniorListenersEventArgs{
				Once:    once,
				EndLoop: endLoop,
				Async:   async,
			},
		}, params...),
	}
}

func WrapDefaultSeniorListener[K comparable](c decorator.Ctrl, tp K, params ...interface{}) *decorator.Action {
	return &decorator.Action{
		C: c,
		P: append([]interface{}{
			tp,
			&SeniorListenersEventArgs{},
		}, params...),
	}
}

// EasyListen ... easy listen
func (l *SeniorListeners[K]) EasyListen(ch *SeniorEventChannel[K], actions []*decorator.Action, parser EventTypeParser[K]) (interface{}, error) {
	return l.Listen(nil, nil, nil, ch, actions, parser)
}

// Listen
// ch chan *ListenersEvent: read-only?
func (l *SeniorListeners[K]) Listen(task *decorator.Task, input interface{}, ps *decorator.Stage,
	ch *SeniorEventChannel[K], actions []*decorator.Action, parser EventTypeParser[K]) (interface{}, error) {
	if len(actions) <= 0 {
		return nil, errors.New(decorator.EM1301EmptyAction)
	}
	if !l.Prepare() {
		return nil, errors.New("prepare failed")
	}

	var ctx context.Context
	if task != nil && task.Context != nil {
		ctx = task.Context
	} else {
		ctx = context.Background()
	}

	// listener 超时逻辑
	if l.timeoutMs == 0 {
		l.timeoutMs = defaultTimeoutMs
	}

	tm := time.After(time.Duration(l.timeoutMs) * time.Millisecond)
	mp := make(map[K]*seniorListenersEventAction[K])
	var err error
	var ret interface{}
	for _, action := range actions {
		size := len(action.P)
		if size < 2 {
			return nil, errors.New(decorator.EM1303MissingParams)
		}
		key, ok := parser(action.P[0])
		if !ok {
			return nil, fmt.Errorf(decorator.EM1305WrongParams, "listen action params[0]")
		}
		args, ok := action.P[1].(*SeniorListenersEventArgs)
		if !ok {
			return nil, fmt.Errorf(decorator.EM1305WrongParams, "listen action params[1]")
		}
		la := &seniorListenersEventAction[K]{
			Action: decorator.Action{
				C: action.C,
				P: action.P[2:],
				E: action.E,
			},
			Key:                      key,
			SeniorListenersEventArgs: *args,
		}
		mp[key] = la
	}

	for l.running {
		select {
		case iev := <-l.ievChan:
			if iev.Type == listenerIEVClose {
				ret = iev.Data
				l.stopListen(ch)
			} else if iev.Type == listenerIEVRemove {
				if et, ok := parser(iev.Data); ok {
					delete(mp, et)
				}
			} else if iev.Type == listenerIEVClear {
				mp = map[K]*seniorListenersEventAction[K]{}
			} else if iev.Type == listenerIEVDestroy {
				mp = map[K]*seniorListenersEventAction[K]{}
				l.stopListen(ch)
			}
		case event := <-ch.Receive():
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
					_, _ = l.doEvent(task, event, ps, action)
					if action.EndLoop {
						l.end(nil)
					}
				}()
			} else {
				_, _ = l.doEvent(task, event, ps, action)
				if action.EndLoop {
					l.end(nil)
				}
			}
		case <-tm:
			l.end(errors.New("listeners timeout"))
		case <-ctx.Done():
			l.end(errors.New("listeners ctx done"))
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

func (l *SeniorListeners[K]) Prepare() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.status > listenerStReady {
		fmt.Println("__listeners__ prepared already")
		return false
	}
	fmt.Println("__listeners__ prepare")
	l.running = true
	l.ievChan = make(chan *seniorListenersInternalEventAction, 3)
	return true
}

func (l *SeniorListeners[K]) Destroy() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.status == listenerStReleased {
		fmt.Println("__listeners__ destroy skip")
		return
	}
	if l.status < listenerStToRelease {
		l.status = listenerStToRelease
	}
	fmt.Println("__listeners__ destroy")
	if l.ievChan != nil {
		l.ievChan <- &seniorListenersInternalEventAction{
			Type: listenerIEVDestroy,
			Data: nil,
		}
	}
}

// end ... end
func (l *SeniorListeners[K]) end(data interface{}) {
	defer l.lock.Unlock()
	l.lock.Lock()
	if l.status == listenerStReleased {
		fmt.Println("__listeners__ destroy end skip")
		return
	}
	fmt.Println("__listeners__ end")
	if l.status < listenerStToRelease {
		l.status = listenerStToRelease
	}
	if l.ievChan != nil {
		l.ievChan <- &seniorListenersInternalEventAction{
			Type: listenerIEVClose,
			Data: data,
		}
	}
}

func (l *SeniorListeners[K]) doEvent(task *decorator.Task, input interface{}, ps *decorator.Stage,
	a *seniorListenersEventAction[K]) (interface{}, error) {
	i, e := a.C(task, input, ps, a.P...)
	if e != nil && a.E != nil {
		i, e = a.E(task, e, ps, a.P...)
	}
	// i != nil as finished
	if i != nil {
		l.end(i)
	} else if e != nil {
		l.end(e)
	}
	return i, e
}

func (l *SeniorListeners[K]) stopListen(ch *SeniorEventChannel[K]) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.status == listenerStReleased {
		fmt.Println("__listeners__ stopListen already")
		return
	}
	fmt.Println("__listeners__ stopListen")
	l.running = false
	l.status = listenerStReleased
	if l.ievChan != nil {
		close(l.ievChan)
		l.ievChan = nil
	}
	ch.Close()
}
