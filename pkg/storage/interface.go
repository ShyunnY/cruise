package storage

import (
	"github.com/gogo/protobuf/types"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"time"
)

// Storage
// 存储接口
// TODO: 这部分是由logic进行调用获取的
// TODO: 考虑使用批处理 避免重复进行调用导致 栈创建销毁的开销
type Storage interface {
	GetTrace(traceID string) *v1.TracesData

	ListTrace(tp TraceParameters) []*v1.TracesData
	ListServices() []string
	ListOperations() []spanstore.Operation

	// PutSpan
	// add or update trace
	PutSpan(span []*v1.ResourceSpans) error

	// PutService add or update service
	PutService(service string, spanID string) error

	// PutOperation add or update operation
	PutOperation(service string, operation spanstore.Operation) error
}

type TraceParameters struct {
	SvcName       string
	OperationName string
	Resources     map[string]string
	BeginTime     time.Time
	EndTime       time.Time
	ElapsedMax    time.Duration
	ElapsedMin    time.Duration
	TraceNum      int32
}

func (t *TraceParameters) ConvertTraceQueryParameters() *api_v3.TraceQueryParameters {

	// convert types.Timestamp
	begin, _ := types.TimestampProto(t.BeginTime)
	end, _ := types.TimestampProto(t.EndTime)

	// convert types.Duration
	emax := types.DurationProto(t.ElapsedMax)
	emin := types.DurationProto(t.ElapsedMin)

	tqp := &api_v3.TraceQueryParameters{
		ServiceName:   t.SvcName,
		OperationName: t.OperationName,
		Attributes:    t.Resources,
		DurationMax:   emax,
		DurationMin:   emin,
		StartTimeMax:  begin,
		StartTimeMin:  end,
		NumTraces:     t.TraceNum,
	}
	return tqp
}

type Operation struct {
	Name     string
	SpanKind string
	SpanID   string
}
