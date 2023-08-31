package metricx

import (
	cv1 "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagStat(t *testing.T) {

	num := 10
	tag := &cv1.KeyValue{
		Key:   "usr",
		Value: &cv1.AnyValue{Value: &cv1.AnyValue_StringValue{StringValue: "z3"}},
	}
	tagMetricx := NewTagsStat(tag)

	for _, span := range generateSpans(num) {
		for _, attr := range span.Attributes {
			key := attr.GetKey()
			val := attr.GetValue().GetStringValue()

			if key == tag.GetKey() && val == tag.GetValue().GetStringValue() {
				tagMetricx.Inc()
			}
		}
	}

	assert.Equal(t, uint64(1), tagMetricx.Count())
	assert.Equal(t, float64(1), tagMetricx.Rate(tagMetricx.Count()))

}

func generateSpans(num int) []*v1.Span {
	var spans []*v1.Span

	for i := 0; i < num; i++ {
		span := &v1.Span{}
		if i == 0 {
			span.Attributes = append(span.Attributes, &cv1.KeyValue{
				Key:   "usr",
				Value: &cv1.AnyValue{Value: &cv1.AnyValue_StringValue{StringValue: "z3"}},
			})
		}
		spans = append(spans, span)
	}

	return spans
}
