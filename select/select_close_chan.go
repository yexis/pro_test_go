package _select

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func SelectCloseChanBeforeSelect() {
	ctx, can := context.WithCancel(context.Background())
	ch := make(chan bool, 1)
	can()
	close(ch)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recover panic", r)
			}
		}()

		select {
		case <-ctx.Done():
			fmt.Println("ctx done")
		case ch <- true:
			fmt.Println("write success")
		default:
			fmt.Println("write full")
		}
	}()

	time.Sleep(100 * time.Millisecond)
}

func SelectCloseChanSyncSelect() {
	ctx, can := context.WithCancel(context.Background())
	ch := make(chan bool, 1)

	go func() {
		can()
		close(ch)
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recover panic", r)
			}
		}()

		select {
		case <-ctx.Done():
			fmt.Println("ctx done")
		case ch <- true:
			fmt.Println("write success")
		default:
			fmt.Println("write full")
		}
	}()

	time.Sleep(100 * time.Millisecond)
}

func SelectCloseChanSyncSelectByAtomic() {
	ctx, can := context.WithCancel(context.Background())
	ch := make(chan bool, 1)
	status := atomic.Int32{}

	go func() {
		can()
		close(ch)
		status.Add(1)
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recover panic", r)
			}
		}()

		cls := status.Load()
		if cls > 0 {
			fmt.Println("chan has been closed")
			return
		}
		select {
		case <-ctx.Done():
			fmt.Println("ctx done")
		case ch <- true:
			fmt.Println("write success")
		default:
			fmt.Println("write full")
		}
	}()

	time.Sleep(100 * time.Millisecond)
}

func SelectCloseChanSyncSelectByMutex() {
	ctx, can := context.WithCancel(context.Background())
	ch := make(chan bool, 1)
	status := 0
	lock := sync.Mutex{}

	go func() {
		can()
		lock.Lock()
		status++
		close(ch)
		lock.Unlock()
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recover panic", r)
			}
		}()

		var cls int
		lock.Lock()
		cls = status
		if cls > 0 {
			fmt.Println("chan has been closed")
			return
		}
		select {
		case <-ctx.Done():
			fmt.Println("ctx done")
		case ch <- true:
			fmt.Println("write success")
		default:
			fmt.Println("write full")
		}
		lock.Unlock()
	}()

	time.Sleep(100 * time.Millisecond)
}
