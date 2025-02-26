package context

import (
	"context"
	"fmt"
	"time"
)

func testCancel1() {
	ctx, cancel := context.WithCancel(context.Background())
	f := func(cancel *context.CancelFunc) {
		go func() {
			time.Sleep(2 * time.Second)
			(*cancel)()
			fmt.Println("do cancel")
		}()
	}
	for i := 0; i < 5; i++ {
		f(&cancel)
	}
	select {
	case <-ctx.Done():
		fmt.Println("ctx done")
	}

	time.Sleep(time.Second * 30)
	cancel()
}

func main() {
	testCancel1()
}
