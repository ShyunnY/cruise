package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {

	after := time.After(time.Second * -1)
	fmt.Printf("after: %+v\n", <-after)
}
