package mock

import (
	"context"
	"fmt"
	"github.com/ShyunnY/cruise/pkg/clog"
	pipe "github.com/ShyunnY/cruise/pkg/pipeline"
	"github.com/ShyunnY/cruise/pkg/pool"
	"github.com/ShyunnY/cruise/pkg/reader"
	"github.com/ShyunnY/cruise/pkg/storage"
	"github.com/ShyunnY/cruise/pkg/storage/memory"
	"log"
	"testing"
	"time"
)

// 功能测试文件

func TestWork(t *testing.T) {

	clog.SetLogger()

	// New Reader
	r, err := reader.NewGrpcReader(reader.GrpcReaderConfig{
		Host: "192.168.136.134",
		Port: 16685,
	})
	if err != nil {
		log.Fatal(" jaeger server unstart error")
	}

	// NewMemory
	sm := memory.NewStoreMemory()

	// New Sink
	sink := pipe.NewStorageSink(pipe.StorageSinkConfig{})

	// New WorkPool
	wp := pool.NewWorkPool(pool.WorkConfig{
		Interval: 10,
		Read:     r,
		Sink:     sink,
		Stages:   nil,
		BufSize:  1000,
	})

	go func() {

		time.Sleep(time.Second * 5)
		result := sm.ListTrace(storage.TraceParameters{
			SvcName: "orange",
		})

		log.Println(len(result))

	}()

	wp.Work(context.TODO())

}

func TestSlice(t *testing.T) {

	sli := []int{1, 2, 3, 4, 5}
	fmt.Printf("sli: %+v\n", sli)

	clear(sli)
	fmt.Printf("sli: %+v\n", sli)

}
