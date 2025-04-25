package easytimer

import (
	"fmt"
	"testing"
	"time"
)

func timeout(rt *EasyTimer, args ...interface{}) error {
	fmt.Println("time out")
	return nil
}

func TestEasyTimer(t *testing.T) {
	et := EasyTimer{}
	et.Start(TimerOnlyCheck, time.Duration(5)*time.Second, timeout)

	status, curStopped := et.Stop(TimerOnlyCheck)
	if curStopped {
		fmt.Println("curr stopped succeed")
	} else {
		fmt.Println("curr stopped failed")
	}
	fmt.Println("status:", status)

	status, curStopped = et.Stop(TimerStoppedByCannel)
	if curStopped {
		fmt.Println("curr stopped succeed")
	} else {
		fmt.Println("curr stopped failed")
	}
	fmt.Println("status:", status)
}

func TestEasyTimer2(t *testing.T) {
	et := EasyTimer{}
	et.Start(TimerOnlyCheck, time.Duration(5)*time.Second, EmptyEasyTimerCallback)

	select {
	case <-et.C():
		fmt.Println("time out")
	}
}
