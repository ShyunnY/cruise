package metricx

import "sync/atomic"

type SpanStat struct {
	val uint64
}

func NewSpanStat() SpanMetricx {
	return &SpanStat{}
}

func (s *SpanStat) Add(delta uint64) {

	atomic.AddUint64(&s.val, delta)
}

func (s *SpanStat) Count() uint64 {

	return atomic.LoadUint64(&s.val)
}
