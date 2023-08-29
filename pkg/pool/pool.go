package pool

import (
	"context"
	"github.com/ShyunnY/cruise/pkg/clog"
	"github.com/ShyunnY/cruise/pkg/pipeline"
	"github.com/ShyunnY/cruise/pkg/reader"
	"github.com/gogo/protobuf/types"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"sync/atomic"
	"time"
)

type placeholder struct{}

type WorkPool struct {
	delay time.Duration

	pipe   pipe.Pipeline
	reader reader.Reader

	done   chan placeholder
	chunks chan []*v1.ResourceSpans
	send   chan *v1.ResourceSpans

	pass atomic.Bool
}

type WorkConfig struct {
	Interval time.Duration
	Read     reader.Reader
	Sink     pipe.SinkPipe
	Stages   []pipe.StagePipe
	BufSize  int
}

func NewWorkPool(conf WorkConfig) *WorkPool {

	w := &WorkPool{
		delay:  conf.Interval,
		reader: conf.Read,
		chunks: make(chan []*v1.ResourceSpans, 10),
		send:   make(chan *v1.ResourceSpans, conf.BufSize),
	}

	pip := w.newPipeline(conf.Sink, conf.Stages...)
	w.pipe = pip

	return w
}

func (w *WorkPool) Work(ctx context.Context) {
	go w.backgroupSend()
	timer := make(<-chan time.Time)

	for {
		// mock param
		begin, _ := types.TimestampProto(time.Now().Add(-time.Hour))
		end, _ := types.TimestampProto(time.Now())

		// query reader trace data
		res, err := w.reader.SearchTraces(context.TODO(), reader.SearchTracesRequest{
			SearchParam: &api_v3.TraceQueryParameters{
				ServiceName:  "orange",
				StartTimeMin: begin,
				StartTimeMax: end,
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

		timer = time.After(w.delay * time.Second)
		select {
		case <-ctx.Done():
			close(w.done)
		case <-w.done:
			return
		case <-timer:
			// reset timer
			timer = time.After(w.delay)
		}
	}
}

func (w *WorkPool) backgroupSend() {
	for {
		select {
		case spans := <-w.chunks:
			if spans == nil {
				return
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

		}
	}()

	return pip
}
