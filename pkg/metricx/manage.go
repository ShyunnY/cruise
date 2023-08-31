package metricx

import (
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"time"
)

var (
	serviceNameKey = "service.name"
)

type Manage struct {
	Errmetricx     ErrorMetricx
	elapsedMetricx ElapsedMetricx
	SpanMetricx    SpanMetricx
	TagsMetricx    []TagMetricx
	// TODO: other metricx

	useSvcGroup bool
}

type ManageConfig struct {
	Errm  ErrorMetricx
	ElapM ElapsedMetricx
	SpanM SpanMetricx
	TagsM []TagMetricx

	UseServiceGroup bool
}

func NewManage(conf ManageConfig) *Manage {

	// TODO: set default Manager Config val
	conf.setDefault()

	return &Manage{
		Errmetricx:     conf.Errm,
		elapsedMetricx: conf.ElapM,
		SpanMetricx:    conf.SpanM,
		TagsMetricx:    conf.TagsM,
		useSvcGroup:    conf.UseServiceGroup,
	}
}

// Handle
// use handle span recode metrics for all metricx
func (m *Manage) Handle(span *v1.ResourceSpans) {

	for _, lbras := range span.InstrumentationLibrarySpans {

		// recode span num
		m.SpanMetricx.Add(uint64(len(lbras.Spans)))

		var svcName string
		for _, attr := range span.Resource.Attributes {
			// get servcei name
			if attr.GetKey() == serviceNameKey {
				svcName = attr.GetValue().GetStringValue()
			}
		}

		for _, span := range lbras.Spans {
			for _, attr := range span.Attributes {

				// recoder err metricx
				if attr.GetKey() == "otel.status_code" && attr.GetValue().GetStringValue() == "ERROR" {
					m.Errmetricx.Inc()

					if m.useSvcGroup {
						m.Errmetricx.IncWithService(svcName)
					}

				}

				// recoder customize tag Metricx
				for _, tm := range m.TagsMetricx {
					tmKey := tm.GetTag().GetKey()
					tmVal := tm.GetTag().GetValue().GetStringValue()

					if attr.GetKey() == tmKey && attr.GetValue().GetStringValue() == tmVal {
						tm.Inc()
					}
				}

			}

			// recoder duration
			if duration := span.EndTimeUnixNano - span.StartTimeUnixNano; duration > 0 {
				m.elapsedMetricx.Add(time.Duration(duration))

				if m.useSvcGroup {
					m.elapsedMetricx.AddWithService(time.Duration(duration), svcName)
				}

			}

		}
	}

}

func (c *ManageConfig) setDefault() {

	// use default error metricx if not set
	if c.Errm == nil {
		es := NewErrorStat()
		c.Errm = es
	}

	// use default elapsed metricx if not set
	if c.ElapM == nil {
		els := NewElapsedStat()
		c.ElapM = els
	}

	// use default SpanMetricx if not set
	if c.SpanM == nil {
		sm := NewSpanStat()
		c.SpanM = sm
	}

}

func (m *Manage) GetElapsedAvg(svc ...string) time.Duration {

	svcName := overrideSvcName(svc)
	return m.elapsedMetricx.AvgWithService(svcName)
}

func (m *Manage) GetErrorRate(base uint64, svc ...string) float64 {

	svcName := overrideSvcName(svc)
	return m.Errmetricx.RateWithService(svcName, base)
}

func (m *Manage) GetErrorCount(svc ...string) uint64 {

	svcName := overrideSvcName(svc)
	return m.Errmetricx.CountWithService(svcName)
}

func (m *Manage) GetSpanCount() uint64 {
	// TODO: are statistics based on service considered ???
	return m.SpanMetricx.Count()
}

func (m *Manage) Process(span *v1.ResourceSpans) (*v1.ResourceSpans, error) {
	m.Handle(span)
	return span, nil
}

func overrideSvcName(svc []string) string {
	var svcName string
	// Override service name if provide
	if len(svc) > 0 {
		svcName = svc[0]
	}
	return svcName
}
