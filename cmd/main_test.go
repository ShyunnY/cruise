package main

import (
	"context"
	"k8s.io/apimachinery/pkg/util/wait"
	"log"
	"testing"
	"time"
)

func TestWaitUtil(t *testing.T) {
	wait.PollUntilContextCancel(
		context.TODO(),
		time.Second,
		false,
		func(ctx context.Context) (done bool, err error) {
			log.Println("info...")
			return false, err
		},
	)

}
