package pool

import (
	"context"
	"github.com/ShyunnY/cruise/pkg/clog"
	"github.com/ShyunnY/cruise/pkg/pipeline"
	"github.com/ShyunnY/cruise/pkg/reader"
	"github.com/gogo/protobuf/types"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"log"
	"sync/atomic"
	"time"
)

const (
	defaultInterval = time.Minute * 10
	defaultPeriod   = time.Hour * 24
	defaultBufSize  = 2 << 9
)

type placeholder struct{}

type WorkPool struct {
	// pool get traces for reader interval
	delay  time.Duration
	period time.Duration

	pipe   pipe.Pipeline
	reader reader.Reader

	done   chan placeholder
	chunks chan []*v1.ResourceSpans
	send   chan *v1.ResourceSpans

	pass atomic.Bool
}

type WorkConfig struct {
	Interval time.Duration
	period   time.Duration

	Read    reader.Reader
	Sink    pipe.SinkPipe
	Stages  []pipe.StagePipe
	BufSize int
}

func NewWorkPool(conf WorkConfig) *WorkPool {

	conf.setDefault()

	w := &WorkPool{
		delay:  conf.Interval,
		period: conf.period,
		reader: conf.Read,
		chunks: make(chan []*v1.ResourceSpans),
		send:   make(chan *v1.ResourceSpans, conf.BufSize),
	}

	w.OnPass()
	pip := w.newPipeline(conf.Sink, conf.Stages...)
	w.pipe = pip

	return w
}

// Work
// TODO: 可以使用stopCh进行控制组件的关闭
func (w *WorkPool) Work(ctx context.Context) {
	go w.backgroupSend()

	timer := make(<-chan time.Time)
	last := time.Now()

	startTs := convertProtoTimestamp(last.Add(-w.period))
	endTs, _ := types.TimestampProto(last)
	for {
		// query reader trace data
		res, err := w.reader.SearchTraces(context.TODO(), reader.SearchTracesRequest{
			SearchParam: &api_v3.TraceQueryParameters{
				ServiceName:  "orange",
				StartTimeMin: startTs,
				StartTimeMax: endTs,
			},
		})
		if err != nil {
			// TODO: 考虑是否进行指数回退进行获取
			clog.CL.Error(err.Error())
		}

		// send pipeline channel
		if res != nil {
			w.chunks <- res.ResourceSpans
		}

		timer = time.After(w.delay)
		select {
		case <-ctx.Done():
			close(w.done)
		case <-w.done:
			return
		case <-timer:
			// reset timer
			timer = time.After(w.delay)

			// reset query timestamp
			last = time.Now()
			startTs = convertProtoTimestamp(last.Add(-w.delay))
			endTs = convertProtoTimestamp(last)
		}
	}
}

func (w *WorkPool) backgroupSend() {
	for {
		select {
		case spans := <-w.chunks:
			if spans == nil {
				continue
			}

			for _, span := range spans {
				if w.pass.Load() {
					w.send <- span
				}
			}
		case <-w.done:
			return
		}
	}
}

func (w *WorkPool) OnPass() {
	w.pass.Store(true)
}

func (w *WorkPool) OffPass() {
	w.pass.Store(false)
}

// create work pool internal pipeline
func (w *WorkPool) newPipeline(sink pipe.SinkPipe, stages ...pipe.StagePipe) pipe.Pipeline {
	pip := pipe.NewSpanPipeline()
	pip.AddSinkPipe(sink)
	pip.AddStagePipes(stages...)

	sourceFn := pipe.SourcePipeFunc(func(ctx context.Context) (<-chan *v1.ResourceSpans, error) {

		go func() {
			select {
			case <-ctx.Done():
				w.OffPass()
				close(w.send)
			}
		}()

		return w.send, nil
	})

	pip.AddSourcePipe(sourceFn)

	go func() {
		// TODO: 发生错误,需要通知 pool暂时停止发送信息
		err := pip.Run(context.Background())
		if err != nil {
			log.Println(err)
		}
	}()

	return pip
}

func (c *WorkConfig) setDefault() {

	if c.Stages == nil {
		// TODO: set default stages
	}

	if c.Interval <= 0 {
		// set default interval
		// If interval is not set, I consider 10 min as a time interval to get
		c.Interval = defaultInterval
	}

	if c.period <= 0 {
		// set default 1 day period
		c.period = defaultPeriod
	}

	if c.Read == nil {
		// set default grpc reader
		grpcReader, err := reader.NewGrpcReader(reader.GrpcReaderConfig{})
		if err != nil {
			// TODO: 尝试指数回退机制
			log.Println(err)
		}
		c.Read = grpcReader
	}

	if c.Sink == nil {
		// set default memorey store sink
		c.Sink = pipe.NewStorageSink(pipe.StorageSinkConfig{})
	}

	if c.BufSize <= 0 {
		// set default bufSize=1024
		c.BufSize = defaultBufSize
	}

}

func convertProtoTimestamp(t time.Time) *types.Timestamp {

	// we can ignore convert error
	ts, _ := types.TimestampProto(t)
	return ts
}
