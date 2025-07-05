package easylistener

import (
	"errors"
	"fmt"
	"pro_test_go/easy/decorator"
	"testing"
	"time"
)

func dataHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	in := input.(*ListenersEvent)
	fmt.Println("data-handler", in.Value.(string))
	return nil, nil
}

func errorHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	in := input.(*ListenersEvent)
	fmt.Println("error-handler", in.Value.(error))
	return nil, nil
}

func closeHandler(task *decorator.Task, input interface{}, stage *decorator.Stage, params ...interface{}) (interface{}, error) {
	fmt.Println("close-handler")
	return nil, nil
}

func TestListenerEasyListen(t *testing.T) {
	eventChan := NewEventChannel(5)
	l := Listeners{}
	l.SetTimeout(30000)
	go func() {
		l.EasyListen(eventChan, []*decorator.Action{
			WrapListener(dataHandler, "data", false, false, false),
			WrapListener(errorHandler, "error", true, false, false),
			WrapListener(closeHandler, "close", true, false, false),
		})
	}()

	emit := func(key string, value interface{}) {
		if eventChan.IsClosed() {
			fmt.Println("event chan was closed", key, value)
			return
		}
		eventChan.Send(
			&ListenersEvent{
				Key:   EventType(key),
				Value: value,
			},
		)
	}
	// mock request
	go func() {
		emit("data", "result_1")
		time.Sleep(10 * time.Millisecond)
		emit("data", "result_2")
		time.Sleep(10 * time.Millisecond)
		emit("data", "result_3")
		time.Sleep(10 * time.Millisecond)
		emit("error", errors.New("error_1"))
		time.Sleep(10 * time.Millisecond)
		emit("close", nil)
		time.Sleep(10 * time.Millisecond)

		emit("data", "result_4")
		time.Sleep(10 * time.Millisecond)
		emit("data", "result_5")
		time.Sleep(10 * time.Millisecond)
		emit("data", "result_6")
		time.Sleep(10 * time.Millisecond)
		emit("error", errors.New("error_2"))
		time.Sleep(10 * time.Millisecond)
	}()

	time.Sleep(5 * time.Second)
	l.Destroy()
}
