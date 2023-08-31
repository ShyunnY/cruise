package metricx

import (
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"time"
)

type Manage struct {
	errMetricx     ErrorMetricx
	elapsedMetricx ElapsedMetricx
	spanMetricx    SpanMetricx
	tagMetricx     []TagMetricx
	// TODO: other metricx
}

func (m *Manage) Handle(span *v1.ResourceSpans) {

	for _, lbras := range span.InstrumentationLibrarySpans {

		// recode span num
		m.spanMetricx.Add(len(lbras.Spans))

		for _, span := range lbras.Spans {
			for _, attr := range span.Attributes {

				// recoder err metricx
				if attr.GetKey() == "otel.status_code" && attr.GetValue().GetStringValue() == "ERROR" {
					m.errMetricx.Inc()
				}

				for _, tm := range m.tagMetricx {
					tmKey := tm.GetTag().GetKey()
					tmVal := tm.GetTag().GetValue().GetValue()
					if attr.GetKey() == tmKey && attr.GetValue().GetValue() == tmVal {
						tm.Inc()
					}
				}
			}

			// recoder duration
			if duration := span.EndTimeUnixNano - span.StartTimeUnixNano; duration > 0 {
				m.elapsedMetricx.Add(time.Duration(duration))
			}

		}
	}

}
