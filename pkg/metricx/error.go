package metricx

import "sync/atomic"

type ErrorStat struct {
	val uint64

	svcErr map[string]uint64
}

func NewErrorStat() ErrorMetricx {
	return &ErrorStat{
		svcErr: map[string]uint64{},
	}
}

func (e *ErrorStat) Add(delta int) {

	if delta < 0 {
		panic("error stat cannot increase nagative num")
	}

	v := uint64(delta)
	if int(v) == delta {
		atomic.AddUint64(&e.val, v)
	}

}

func (e *ErrorStat) Inc() {
	atomic.AddUint64(&e.val, 1)
}

func (e *ErrorStat) IncWithService(svc string) {

	if svc == "" {
		// TODO: 记录error日志
		return
	}

	v := e.svcErr[svc]
	atomic.AddUint64(&v, 1)
}

func (e *ErrorStat) Count() uint64 {
	return atomic.LoadUint64(&e.val)
}

func (e *ErrorStat) CountWithService(svc string) uint64 {

	if svc == "" {
		return e.Count()
	}

	if v, ok := e.svcErr[svc]; ok {
		return v
	} else {
		return uint64(0)
	}
}

func (e *ErrorStat) Rate(base uint64) float64 {
	if base <= 0 {
		// TODO: 记录日志
		panic("error metricx rate base num must greater than 0") // return float64(0)
	}

	ret := float64(e.Count()) / float64(base)

	return ret
}

func (e *ErrorStat) RateWithService(svc string, base uint64) float64 {
	if svc == "" {
		return e.Rate(base)
	}

	if v, ok := e.svcErr[svc]; ok {
		ret := float64(v) / float64(base)
		return ret
	} else {
		return float64(0)
	}

}
