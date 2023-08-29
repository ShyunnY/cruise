package reader

import (
	"context"
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

	start := time.Now().Add(-time.Hour)
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

	log.Println(traces)

}
