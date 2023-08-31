package metricx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpanStat(t *testing.T) {

	spanStat := NewSpanStat()
	spanStat.Add(1)
	spanStat.Add(20)
	spanStat.Add(30)
	spanStat.Add(40)

	assert.Equal(t, uint64(91), spanStat.Count())

}
