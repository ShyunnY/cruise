package metricx

import "sync/atomic"

type ErrorStat struct {
	val uint64
}

func NewErrorStat() ErrorMetricx {
	return &ErrorStat{}
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

func (e *ErrorStat) Count() uint64 {
	return atomic.LoadUint64(&e.val)
}

func (e *ErrorStat) ErrRate(base int64) float32 {

	if base <= 0 {
		panic("err rate base num must greater than 0")
	}

	b := uint64(base)

	ret := float32(e.Count()) / float32(b)

	return ret
}
