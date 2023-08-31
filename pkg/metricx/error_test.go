package metricx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorMetricx(t *testing.T) {

	es := NewErrorStat()

	es.Inc()
	es.Inc()
	es.Inc()
	es.Inc()
	es.Inc()

	assert.Equal(t, 0.05, es.Rate(100))

}
