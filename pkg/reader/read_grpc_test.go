package reader

import (
	"context"
	"github.com/ShyunnY/cruise/pkg/storage/memory"
	"github.com/gogo/protobuf/types"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	"log"
	"testing"
	"time"
)

func TestGrpcReader(t *testing.T) {

	reader, err := NewGrpcReader(GrpcReaderConfig{
		Address: "192.168.136.134",
		Port:    16685,
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
	for _, span := range traces.ResourceSpans {
		if err := sm.PutSpan(span); err != nil {
			panic(err)
		}
	}

	for _, svc := range sm.ListServices() {
		log.Println(svc)
		for _, op := range sm.ListOperations(svc) {
			log.Println(op)
		}
	}

}
