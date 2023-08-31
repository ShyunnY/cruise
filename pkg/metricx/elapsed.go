package metricx

import (
	"sync/atomic"
	"time"
)

type elapsedStat struct {
	valDuration uint64
	valCount    uint64

	svcDuration map[string]uint64
}

func NewElapsedStat() ElapsedMetricx {
	return &elapsedStat{
		svcDuration: map[string]uint64{},
	}
}

func (e *elapsedStat) AddWithService(duration time.Duration, svc string) {

	if duration < 0 {
		panic("elapsed metricx must need duration greater than 0")
	}

	if svc == "" {
		panic("elapsed metricx must need no-empty service name recode")
	}

	val := uint64(duration)
	if time.Duration(val) == duration {

		if v, ok := e.svcDuration[svc]; !ok {
			atomic.StoreUint64(&v, val)
		} else {
			atomic.AddUint64(&v, val)
		}

	}

}

func (e *elapsedStat) AvgWithService(svc string) time.Duration {

	// return overall avg duration
	if svc == "" {
		// TODO: 记录日志
		return e.Avg()
	}

	if v, ok := e.svcDuration[svc]; ok {
		return time.Duration(v)
	} else {
		return time.Duration(0)
	}
}

func (e *elapsedStat) Add(duration time.Duration) {

	if duration < 0 {
		panic("elapsed duration must more than 0")
	}

	// inc counter
	e.inc()

	val := uint64(duration)
	if time.Duration(val) == duration {
		atomic.AddUint64(&e.valDuration, val)
	}

}

func (e *elapsedStat) Avg() time.Duration {

	duration := atomic.LoadUint64(&e.valDuration)
	count := atomic.LoadUint64(&e.valCount)

	val := duration
	ret := val / count

	return time.Duration(ret)
}

func (e *elapsedStat) inc() {
	atomic.AddUint64(&e.valCount, 1)
}
