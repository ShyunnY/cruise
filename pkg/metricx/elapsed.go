package metricx

import (
	"sync/atomic"
	"time"
)

type elapsedStat struct {
	valDuration int64
	valCount    uint64
}

func NewelapsedStat() ElapsedMetricx {
	return &elapsedStat{}
}

func (e *elapsedStat) Add(delta time.Duration) {

	if delta < 0 {
		panic("elapsed duration must more than 0")
	}

	// inc counter
	e.inc()

	val := int64(delta)
	if time.Duration(val) == delta {
		atomic.AddInt64(&e.valDuration, val)
	}

}

func (e *elapsedStat) Avg() time.Duration {

	duration := atomic.LoadInt64(&e.valDuration)
	count := atomic.LoadUint64(&e.valCount)

	val := uint64(duration)
	ret := val / count

	return time.Duration(ret)
}

func (e *elapsedStat) inc() {
	atomic.AddUint64(&e.valCount, 1)
}
