package pipe

import (
	"context"
	"github.com/ShyunnY/cruise/pkg/storage"
	"github.com/ShyunnY/cruise/pkg/storage/memory"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"log"
	"time"
)

const (
	defaultBatchSize = 2 << 7
	defaultInterval  = time.Minute * 5
)

type storageSink struct {
	store         storage.Storage
	batchSize     int
	flushInterval time.Duration
}

type StorageSinkConfig struct {
	Store     storage.Storage
	BatchSize int
	Interval  time.Duration
}

func NewStorageSink(conf StorageSinkConfig) SinkPipe {
	return &storageSink{
		store:         conf.Store,
		batchSize:     conf.BatchSize,
		flushInterval: conf.Interval,
	}
}

// Sink 使用批处理的方式进行传递
func (ss *storageSink) Sink(ctx context.Context, spanCh chan *v1.ResourceSpans) error {

	batch := make([]*v1.ResourceSpans, 0, ss.batchSize)

	for {

		finish := false
		send := false

		timer := time.After(ss.flushInterval)
		last := time.Now()

		select {
		case span := <-spanCh:
			batch = append(batch, span)
			send = len(batch) == ss.batchSize

			if send {
				// metrics
				// log.Println("span batch reach preset capacity,will send spans for store backend")

				log.Println("span满足数量,开始刷往存储后端memory...")
			}
		case <-timer:
			timer = time.After(ss.flushInterval)

			// once the refresh interval is up, it will all refresh to the storage backend
			send = time.Since(last) > (ss.flushInterval) && len(batch) > 0
			if send {
				// metrics
				// log.Println("satisfy refresh cycle, will send spans for store backend")

				log.Println("interval满足时间周期,开始刷往存储后端memory...")
			}

		case <-ctx.Done():
			finish = true
			send = len(batch) > 0
		}

		if send {
			// log.Println("sink has plush spans for store")
			if err := ss.store.PutSpan(batch); err != nil {
				panic(err)
			}

			// reset batch slice
			batch = make([]*v1.ResourceSpans, 0, ss.batchSize)
			send = false
		}

		if finish {
			break
		}

	}

	return nil
}

func (sc *StorageSinkConfig) setDefault() {

	if sc.BatchSize == 0 {
		// set default size = 256
		sc.BatchSize = defaultBatchSize
	}

	if sc.Interval == 0 {
		// set default interval = 5min
		sc.Interval = defaultInterval
	}

	if sc.Store == nil {
		// set default memory store
		sc.Store = memory.NewStoreMemory()
	}

}
