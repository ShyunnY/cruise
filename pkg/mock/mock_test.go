package mock

import (
	"context"
	"github.com/ShyunnY/cruise/pkg/clog"
	pipe "github.com/ShyunnY/cruise/pkg/pipeline"
	"github.com/ShyunnY/cruise/pkg/pool"
	"github.com/ShyunnY/cruise/pkg/reader"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"log"
	"testing"
)

// 功能测试文件

func TestWork(t *testing.T) {
	clog.SetLogger()
	r, err := reader.NewGrpcReader(reader.GrpcReaderConfig{
		Address: "192.168.136.134",
		Port:    16685,
	})
	if err != nil {
		log.Fatal(" jaeger server unstart error")
	}

	wp := pool.NewWorkPool(pool.WorkConfig{
		Interval: 10,
		Read:     r,
		Sink: pipe.SinkPipeFunc(func(span *v1.ResourceSpans) error {
			log.Println(span)
			return nil
		}),
		Stages:  nil,
		BufSize: 1000,
	})

	wp.Work(context.TODO())
}
