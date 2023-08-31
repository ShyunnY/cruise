package metricx

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestElapsedMetricx(t *testing.T) {

	em := NewElapsedStat()
	t1 := time.Second
	t2 := time.Millisecond * 500
	t3 := time.Millisecond * 300
	em.Add(t1)
	em.Add(t2)
	em.Add(t3)

	assert.Equal(t, time.Millisecond*600, em.Avg())
}
