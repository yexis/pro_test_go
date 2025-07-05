package easylistener

import "github.com/yexis/pro_test_go/easy/easychannel"

type EventChannel struct {
	easychannel.SeniorSafeChannel[*ListenersEvent]
}

func NewEventChannel(sz int) *EventChannel {
	if sz < 0 {
		sz = 1
	}
	elc := &EventChannel{
		SeniorSafeChannel: easychannel.SeniorSafeChannel[*ListenersEvent]{
			DataChan: make(chan *ListenersEvent, sz),
		},
	}
	elc.SetClear(true)
	return elc
}
