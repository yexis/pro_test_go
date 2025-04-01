package _select

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func sele() {
	errCh := make(chan error, 1)
	dataCh := make(chan interface{})

	go func() {
		dataCh <- true
		errCh <- errors.New("err occurred")
	}()

	select {
	case <-errCh:
		fmt.Println("get err")
	case <-dataCh:
		fmt.Println("get data")
	}
}

func TestSelect(t *testing.T) {
	for i := 0; i < 10000; i++ {
		sele()
		time.Sleep(500 * time.Millisecond)
	}
}
