package easylistener

import (
	"errors"
	"fmt"
	"pro_test_go/decorator"
	"testing"
	"time"
)

type MyEventType int

const (
	DataEventType MyEventType = iota
	ErrorEventType
	CloseEventType
)

func seniorDataHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	in := input.(*SeniorListenersEvent[MyEventType])
	fmt.Println("data-handler", in.Value.(string))
	return nil, nil
}

func seniorErrorHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	in := input.(*SeniorListenersEvent[MyEventType])
	fmt.Println("error-handler", in.Value.(error))
	return nil, nil
}

func seniorCloseHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	fmt.Println("close-handler")
	return nil, nil
}

func TestSeniorListenerEasyListen(t *testing.T) {
	eventChan := NewSeniorEventChannel[MyEventType](5)
	l := &SeniorListeners[MyEventType]{}
	l.SetTimeout(30000)
	go func() {
		_, err := l.EasyListen(
			eventChan,
			[]*decorator.Action{
				WrapSeniorListener(seniorDataHandler, DataEventType, false, false, false),
				WrapSeniorListener(seniorErrorHandler, ErrorEventType, true, false, false),
				WrapSeniorListener(seniorCloseHandler, CloseEventType, true, false, false),
			},
			TplEventTypeParser[MyEventType],
		)
		if err != nil {
			fmt.Println("listen err", err.Error())
		}
	}()

	emit := func(key MyEventType, value interface{}) {
		if eventChan.IsClosed() {
			fmt.Println("event chan was closed", key, value)
			return
		}
		eventChan.Send(
			&SeniorListenersEvent[MyEventType]{
				Key:   key,
				Value: value,
			},
		)
	}

	// mock request
	go func() {
		emit(DataEventType, "result_1")
		time.Sleep(10 * time.Millisecond)
		emit(DataEventType, "result_2")
		time.Sleep(10 * time.Millisecond)
		emit(DataEventType, "result_3")
		time.Sleep(10 * time.Millisecond)
		emit(ErrorEventType, errors.New("error_1"))
		time.Sleep(10 * time.Millisecond)
		emit(CloseEventType, nil)
		time.Sleep(10 * time.Millisecond)

		emit(DataEventType, "result_4")
		time.Sleep(10 * time.Millisecond)
		emit(DataEventType, "result_5")
		time.Sleep(10 * time.Millisecond)
		emit(DataEventType, "result_6")
		time.Sleep(10 * time.Millisecond)
		emit(ErrorEventType, errors.New("error_2"))
		time.Sleep(10 * time.Millisecond)
	}()

	time.Sleep(5 * time.Second)
	l.Destroy()
}
