package mock

import (
	"context"
	"github.com/ShyunnY/cruise/pkg/clog"
	"github.com/ShyunnY/cruise/pkg/metricx"
	pipe "github.com/ShyunnY/cruise/pkg/pipeline"
	"github.com/ShyunnY/cruise/pkg/pool"
	"github.com/ShyunnY/cruise/pkg/reader"
	"github.com/ShyunnY/cruise/pkg/storage/memory"
	"github.com/gogo/protobuf/types"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	cv1 "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"log"
	"testing"
	"time"
)

// 功能测试文件

func readTraceData() []*v1.ResourceSpans {
	clog.SetLogger()

	start := time.Now().Add(-time.Hour * 24)
	now := time.Now()

	st, _ := types.TimestampProto(start)
	ed, _ := types.TimestampProto(now)

	// New Reader
	r, err := reader.NewGrpcReader(reader.GrpcReaderConfig{
		Host: "192.168.136.134",
		Port: 16685,
	})
	if err != nil {
		log.Fatal("jaeger server unstart error", err)
	}

	traces, _ := r.SearchTraces(context.TODO(), reader.SearchTracesRequest{
		SearchParam: &api_v3.TraceQueryParameters{
			ServiceName:  "orange",
			StartTimeMin: st,
			StartTimeMax: ed,
		},
	})

	return traces.ResourceSpans
}

func TestMetricx(t *testing.T) {

	// mock tagMetricx
	tagMetricxFirst := metricx.NewTagsStat(&cv1.KeyValue{
		Key:   "password",
		Value: &cv1.AnyValue{Value: &cv1.AnyValue_StringValue{StringValue: "123456"}},
	})

	tagMetricxSecond := metricx.NewTagsStat(&cv1.KeyValue{
		Key:   "usr",
		Value: &cv1.AnyValue{Value: &cv1.AnyValue_StringValue{StringValue: "z3"}},
	})

	rs := readTraceData()
	manage := metricx.NewManage(metricx.ManageConfig{
		TagsM: []metricx.TagMetricx{tagMetricxFirst, tagMetricxSecond},
	})

	for _, r := range rs {
		manage.Handle(r)
	}

	log.Println("service span count: ", manage.GetSpanCount())
	log.Println("service avg duration: ", manage.GetElapsedAvg())
	log.Println("service error rate: ", manage.GetErrorRate(manage.GetSpanCount()))
	log.Println("service error count: ", manage.GetErrorCount())
	log.Println("customize tags count(usr): ", tagMetricxSecond.Count())
	log.Println("customize tags count(pwd): ", tagMetricxFirst.Count())
}

func TestPipeline(t *testing.T) {

	var (
		traces = readTraceData()
		ch     = make(chan *v1.ResourceSpans)
	)
	go func() {
		for _, trace := range traces {
			ch <- trace
		}
	}()

	pipeline := pipe.NewSpanPipeline()
	pipeline.AddStagePipes(pipe.NewMetricxStage(nil))
	pipeline.AddSinkPipe(pipe.NewConsoleSink())
	pipeline.AddSourcePipe(pipe.SourcePipeFunc(func(ctx context.Context) (<-chan *v1.ResourceSpans, error) {
		return ch, nil
	}))

	pipeline.Run(context.Background())

	time.Sleep(time.Minute)

}

func TestAllWork(t *testing.T) {

	// 1.set log
	clog.SetLogger()

	// 1.1 new reader √
	grpcReader, err := reader.NewGrpcReader(reader.GrpcReaderConfig{
		Host: "192.168.136.134",
		Port: 16685,
	})
	if err != nil {
		panic(err)
	}

	// 1.2 new store √
	storeMemory := memory.NewStoreMemory()

	// 1.3 new sink
	sink := pipe.NewStorageSink(pipe.StorageSinkConfig{
		Store:     storeMemory,
		BatchSize: 4,
		Interval:  time.Second * 30,
	})

	// 1.3.1 mock console
	//sink = pipe.NewConsoleSink()

	// 1.4 new stages √
	manage := metricx.NewManage(metricx.ManageConfig{})
	metricxStage := pipe.NewMetricxStage(manage)

	// 2.new pool √
	workPool := pool.NewWorkPool(pool.WorkConfig{
		Interval: time.Second * 20,
		Read:     grpcReader,
		Sink:     sink,
		Stages:   []pipe.StagePipe{metricxStage},
		BufSize:  2,
		Period:   time.Hour * 24,
	})

	// 4. metricx info
	go func() {

		time.Sleep(time.Second * 10)
		log.Println("service span count: ", manage.GetSpanCount())
		log.Println("service avg duration: ", manage.GetElapsedAvg())
		log.Println("service error rate: ", manage.GetErrorRate(manage.GetSpanCount()))
		log.Println("service error count: ", manage.GetErrorCount())
	}()

	// 3.work pool running...
	workPool.Work(context.TODO())
}
