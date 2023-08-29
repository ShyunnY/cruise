package pipe

import (
	"context"
	"github.com/ShyunnY/cruise/pkg/clog"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"sync"
)

type spanPipeline struct {
	source SourcePipe
	sink   SinkPipe

	stages  []StagePipe
	workNum int

	mux sync.Mutex
}

func NewSpanPipeline() Pipeline {
	return &spanPipeline{
		mux:     sync.Mutex{},
		workNum: 1,
	}
}

func (p *spanPipeline) AddSourcePipe(source SourcePipe) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.source = source
}

func (p *spanPipeline) AddSinkPipe(sink SinkPipe) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.sink = sink
}

func (p *spanPipeline) AddStagePipes(stages ...StagePipe) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if stages == nil || len(stages) == 0 {
		return
	}
	p.stages = append(p.stages, stages...)
}

// Run
// TODO: 需要考虑一下 错误是应该立刻终止还是处理完剩下的再进行终止
// 发生错误 应该告诉上述 别发了 关闭通道
func (p *spanPipeline) Run(ctx context.Context) error {

	errCh := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(p.workNum)

	in, err := p.source.Input(ctx)
	if err != nil {
		// clog.error
		return err
	}

	handler := func(ch <-chan *v1.ResourceSpans) error {
		defer func() {
			wg.Done()
		}()

		var processErr error
		for span := range ch {

			select {
			case err := <-errCh:
				return err
			case <-ctx.Done():
				// 处理超时
				return nil
			default:
				for _, stage := range p.stages {
					span, processErr = stage.Process(span)
					if processErr != nil {
						clog.CL.Error("pipeline processor error")
						// 发生错误 关闭通道
						close(errCh)
						return processErr
					}
				}

				// 处理sink err
				err := p.sink.Sink(span)
				if err != nil {
					clog.CL.Error("pipeline sink receive error")
					close(errCh)
					return err
				}

			}

		}

		return nil
	}

	for i := 0; i < p.workNum; i++ {
		go func() {
			if err := handler(in); err != nil {
				errCh <- err
			}
		}()
	}

	return nil
}
