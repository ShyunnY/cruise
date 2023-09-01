package pipe

import (
	"context"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
)

type Pipeline interface {
	Run(ctx context.Context) error
	AddSourcePipe(source SourcePipe)
	AddSinkPipe(sink SinkPipe)
	AddStagePipes(stages ...StagePipe)
	ShutdownNotify() <-chan struct{}
}

// SourcePipe
// as a input source for pipeline
// TODO: 考虑加入泛型进行优化
type SourcePipe interface {
	Input(ctx context.Context) (<-chan *v1.ResourceSpans, error)
}

// SinkPipe
// as a sink source for pipeline
type SinkPipe interface {
	Sink(ctx context.Context, spanCh chan *v1.ResourceSpans) error
}

// StagePipe
// as a stage for pipeline
type StagePipe interface {
	Process(span *v1.ResourceSpans) (*v1.ResourceSpans, error)
}

type SourcePipeFunc func(ctx context.Context) (<-chan *v1.ResourceSpans, error)

func (f SourcePipeFunc) Input(ctx context.Context) (<-chan *v1.ResourceSpans, error) {
	return f(ctx)
}

type SinkPipeFunc func(ctx context.Context, spanCh chan *v1.ResourceSpans) error

func (f SinkPipeFunc) Sink(ctx context.Context, spanCh chan *v1.ResourceSpans) error {
	return f(ctx, spanCh)
}
