package pipe

import (
	"context"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"log"
)

type ConsoleSink struct {
}

func NewConsoleSink() SinkPipe {
	return &ConsoleSink{}
}

func (c *ConsoleSink) Sink(ctx context.Context, spanCh chan *v1.ResourceSpans) error {

	for span := range spanCh {
		for _, lbas := range span.InstrumentationLibrarySpans {

			for _, sp := range lbas.Spans {
				log.Printf("[ SpanName: %s, SpanKind: %s, SpanStatus: %s", sp.GetName(), sp.GetKind().String(), sp.Status.String())
			}

		}
	}
	return nil
}
