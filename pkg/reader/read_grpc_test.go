package reader

import (
	"context"
	"github.com/ShyunnY/cruise/pkg/storage"
	"github.com/ShyunnY/cruise/pkg/storage/memory"
	"github.com/gogo/protobuf/types"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	"log"
	"testing"
	"time"
)

func TestGrpcReader(t *testing.T) {

	reader, err := NewGrpcReader(GrpcReaderConfig{
		Host: "192.168.136.134",
		Port: 16685,
	})
	if err != nil {
		panic(err)
	}

	start := time.Now().Add(-time.Hour * 2)
	log.Println(start, "ts: ", start.UnixMilli())

	now := time.Now()
	log.Println(now, "ts: ", now.UnixMilli())

	st, err := types.TimestampProto(start)
	if err != nil {
		log.Println(err)
	}

	ed, err := types.TimestampProto(now)
	if err != nil {
		log.Println(err)
	}

	traces, err := reader.SearchTraces(context.TODO(), SearchTracesRequest{
		SearchParam: &api_v3.TraceQueryParameters{
			ServiceName:  "orange",
			StartTimeMin: st,
			StartTimeMax: ed,
		},
	})
	sm := memory.NewStoreMemory()

	// 测试PutSpan
	sm.PutSpan(traces.ResourceSpans)

	for _, svc := range sm.ListServices() {
		log.Println(svc)
		for _, op := range sm.ListOperations(svc) {
			log.Println(op)
		}
	}

	// 1. 测试 matchServiceOrOperation √

	// 2. 测试 matchOther
	// 2.1 ElapsedMin √
	// 2.2 ElapsedMax √
	// 2.3 ElapsedMin-ElapsedMax √
	// 2.4 StartTime √

	// 3. 测试matchTags(⭐)
	ts := sm.ListTrace(storage.TraceParameters{
		SvcName:       "orange",
		OperationName: "/usr",
		//BeginTime:     time.Now().Add(-time.Minute * 30),
		ElapsedMin: 0 * time.Millisecond,
		ElapsedMax: 200 * time.Millisecond,
		Resources: map[string]string{
			"password": "123456",
			"event":    "usr-event",
		},
	})

	log.Println("return trace num: ", len(ts))

}
