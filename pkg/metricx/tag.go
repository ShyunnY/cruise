package metricx

import (
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	"sync/atomic"
)

type TagStat struct {
	val uint64
	tag *v1.KeyValue
}

func NewTagsStat(t *v1.KeyValue) TagMetricx {
	return &TagStat{
		tag: t,
	}
}

func (t *TagStat) Inc() {
	atomic.AddUint64(&t.val, 1)
}

func (t *TagStat) Count() uint64 {
	return atomic.LoadUint64(&t.val)
}

func (t *TagStat) Rate(base uint64) float64 {

	ret := float64(base) / float64(t.Count())
	return ret
}

func (t *TagStat) GetTag() *v1.KeyValue {
	return t.tag
}
