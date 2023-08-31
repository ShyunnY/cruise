package metricx

import (
	"log"
	"testing"
)

func TestErrorMetricx(t *testing.T) {

	es := NewErrorStat()

	es.Inc()
	es.Inc()
	es.Inc()
	es.Inc()
	es.Add(10)

	log.Println(es.ErrRate(100))

}
