package pipe

import (
	"github.com/ShyunnY/cruise/pkg/metricx"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
)

type MetricxStage struct {

	// TODO: 后续可以考虑增加prometheus metrics进行监控
	manager *metricx.Manage
}

func NewMetricxStage(manage *metricx.Manage) StagePipe {

	if manage == nil {
		manage = metricx.NewManage(metricx.ManageConfig{})
	}

	return &MetricxStage{
		manager: manage,
	}
}

func (m *MetricxStage) Process(span *v1.ResourceSpans) (
	*v1.ResourceSpans, error) {

	return m.manager.Process(span)
}

func (m *MetricxStage) GetManager() *metricx.Manage {
	return m.manager
}
