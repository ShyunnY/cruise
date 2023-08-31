package metricx

import (
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	"time"
)

// Metricx
// TODO: 我们同时还需要一个可以外部进行拓展的Tag
// TODO: 根据需求可以对指定的tag进行统计次数
// TODO: 基础版本: error次数统计,duration平均耗时,预定义分位数,spanName统计...
type Metricx interface {
	Desc()
}

type Cardinality interface {
	Inc()
	Count()
	Rate(base uint64) float32
}

// ErrorMetricx
// use compute error count and rate
type ErrorMetricx interface {
	Cardinality
	IncWithService(svc string)
	RateWithService(svc string)

	// TODO: consider add operation inc() and rate() if need
}

// ElapsedMetricx
// use compute elapsed and avg elapsed
type ElapsedMetricx interface {
	Add(duration time.Duration)
	AddWithService(duration time.Duration, svc string)

	Avg() time.Duration
	AvgWithService(svc string) time.Duration
}

type TagMetricx interface {
	Cardinality
	GetTag() *v1.KeyValue
}

type SpanMetricx interface {
	Add(delta int)
	Count() uint64
}
