package pool

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {

	after := time.After(time.Second * -1)
	fmt.Printf("after: %+v\n", <-after)
}

func TestA(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("ctx done")
			}
		}
	}()

	time.Sleep(time.Second * 3)
	cancel()
	time.Sleep(time.Second * 10)
}
