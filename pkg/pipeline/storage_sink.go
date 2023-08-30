package pipe

import (
	"context"
	"github.com/ShyunnY/cruise/pkg/storage"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"log"
)

type storageSink struct {
	store     storage.Storage
	batchSize int
}

func NewStorageSink(store storage.Storage, batchSize int) SinkPipe {
	return &storageSink{
		store:     store,
		batchSize: batchSize,
	}
}

// Sink 使用批处理的方式进行传递
func (ss *storageSink) Sink(ctx context.Context, spanCh chan *v1.ResourceSpans) error {

	batch := make([]*v1.ResourceSpans, 0, ss.batchSize)

	for {

		finish := false
		send := false

		// TODO: consider add interval flush
		select {
		case span := <-spanCh:
			batch = append(batch, span)
			send = len(batch) == ss.batchSize

			if send {
				log.Println("Sink plush spans for store")
				// 指标数据
			}

		case <-ctx.Done():
			log.Println("Sink exit")
			finish = true
			send = len(batch) > 0
		}

		if send {
			if err := ss.store.PutSpan(batch); err != nil {
				panic(err)
			}
		}

		if finish {
			break
		}

	}

	return nil
}
