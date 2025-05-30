package easychannel

import (
	"fmt"
	"testing"
	"time"
)

func TestSeniorSafeChannel(t *testing.T) {
	ch := NewSeniorSafeChannel[int](3).SetClear(false)
	// write channel
	go func() {
		for i := 0; i < 10; i++ {
			if !ch.IsClosed() {
				ch.DataChan <- i
				fmt.Println("write chan success", i)
			} else {
				fmt.Println("write chan failed", i)
			}
		}
	}()

	// read channel
	go func() {
		tm := time.NewTimer(2 * time.Second)
		run := true
		for run {
			select {
			case v, ok := <-ch.DataChan:
				if ok {
					fmt.Println("read chan success", v)
				} else {
					fmt.Println("read chan failed because closed")
					run = false
					break
				}
			case <-tm.C:
				fmt.Println("read chan timeout")
				run = false
			}
		}
	}()

	ch.Close()
	time.Sleep(5 * time.Second)
}
