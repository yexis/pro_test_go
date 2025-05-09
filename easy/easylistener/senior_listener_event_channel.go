package easylistener

import "pro_test_go/easy/easychannel"

type SeniorEventChannel[K comparable] struct {
	easychannel.SeniorSafeChannel[*SeniorListenersEvent[K]]
}

func NewSeniorEventChannel[K comparable](sz int) *SeniorEventChannel[K] {
	if sz < 0 {
		sz = 1
	}
	elc := &SeniorEventChannel[K]{
		SeniorSafeChannel: easychannel.SeniorSafeChannel[*SeniorListenersEvent[K]]{
			DataChan: make(chan *SeniorListenersEvent[K], sz),
		},
	}
	elc.SetClear(true)
	return elc
}
