package metricx

import (
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	"time"
)

// Metricx
// TODO: 我们同时还需要一个可以外部进行拓展的Tag
// TODO: 根据需求可以对指定的tag进行统计次数
// TODO: 基础版本: error次数统计,duration平均耗时,预定义分位数,spanName统计...
// TODO: ⭐: trace数据应该需要定期轮询外部
type Metricx interface {
	Desc()
}

type Cardinality interface {
	Inc()
	Count() uint64
	Rate(base uint64) float64
}

// ErrorMetricx
// use compute error count and rate
type ErrorMetricx interface {
	Cardinality

	IncWithService(svc string)
	RateWithService(svc string, base uint64) float64
	CountWithService(svc string) uint64
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
	Add(delta uint64)
	Count() uint64
}
