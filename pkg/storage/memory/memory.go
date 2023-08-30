package memory

import (
	"errors"
	"github.com/ShyunnY/cruise/pkg/storage"
	"github.com/jaegertracing/jaeger/model"
	cv1 "github.com/jaegertracing/jaeger/proto-gen/otel/common/v1"
	rv1 "github.com/jaegertracing/jaeger/proto-gen/otel/resource/v1"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"sort"
)

var (
	serviceNameKey = "service.name"
)

type StoreMemory struct {
	traces   map[string]*v1.TracesData
	spans    map[string]*v1.Span
	services map[string][]string

	operations map[string][]storage.Operation
}

func NewStoreMemory() *StoreMemory {
	return &StoreMemory{
		traces:     map[string]*v1.TracesData{},
		spans:      map[string]*v1.Span{},
		services:   map[string][]string{},
		operations: map[string][]storage.Operation{},
	}
}

func (s *StoreMemory) GetTrace(traceID string) *v1.TracesData {

	if td, ok := s.traces[traceID]; !ok {
		return nil
	} else {
		return td
	}

}

func (s *StoreMemory) ListTrace(tp storage.TraceParameters) []*v1.TracesData {

	var retMe []*v1.TracesData

	// service and operation will first filter
	ids := s.matchSvcAndOp(tp)
	if ids == nil {
		return nil
	}

	// match resource(tags) and duraion,span time will second filter
	for _, id := range ids {
		if span, ok := s.spans[id]; ok {

			if matchOther(span, tp) {

				// get trace by span.traceID
				traceID, _ := model.TraceIDFromBytes(span.TraceId)

				if tc, ok := s.traces[traceID.String()]; ok {

					// match tags
					if matchTags(tc, tp) {
						retMe = append(retMe, tc)
					}
				}

			}
		}
	}

	// sort by start time
	if tp.TraceNum > 0 && len(retMe) > int(tp.TraceNum) {
		sort.Slice(retMe, func(i, j int) bool {
			return sortTraceByTime(retMe, i, j)
		})

		retMe = retMe[len(retMe)-int(tp.TraceNum):]
	}

	return retMe
}

func (s *StoreMemory) ListServices() []string {

	var serviceList []string

	clear(serviceList)
	for svc, _ := range s.services {
		serviceList = append(serviceList, svc)
	}

	return serviceList
}

func (s *StoreMemory) ListOperations(svc string) []string {

	ops := map[string]struct{}{}
	var opts []string

	for _, operation := range s.operations[svc] {
		// Deduplication
		if _, ok := ops[operation.Name]; !ok {
			ops[operation.Name] = struct{}{}
			opts = append(opts, operation.Name)
		}
	}

	return opts
}

// PutSpan
// add span
// TODO: i choice batch handler
func (s *StoreMemory) PutSpan(rspans []*v1.ResourceSpans) error {

	if rspans == nil {
		return errors.New("span cannot be nil")
	}

	for _, rspan := range rspans {
		var svcName string

		for _, attr := range rspan.Resource.Attributes {
			if attr.GetKey() == serviceNameKey {
				svcName = attr.GetValue().GetStringValue()
			}
		}

		for _, librarySpan := range rspan.InstrumentationLibrarySpans {
			for _, span := range librarySpan.Spans {

				spanID, _ := model.SpanIDFromBytes(span.GetSpanId())
				traceID, _ := model.TraceIDFromBytes(span.GetTraceId())

				// set traces
				if _, ok := s.traces[traceID.String()]; !ok {
					s.traces[traceID.String()] = &v1.TracesData{}
				}
				s.traces[traceID.String()].ResourceSpans = append(s.traces[traceID.String()].ResourceSpans, rspan)

				// set service
				s.services[svcName] = append(s.services[svcName], spanID.String())

				// set operation
				operation := storage.Operation{
					Name:     span.GetName(),
					SpanKind: span.Kind.String(),
					SpanID:   spanID.String(),
				}
				s.operations[svcName] = append(s.operations[svcName], operation)

				// set spans
				s.spans[spanID.String()] = span
			}
		}
	}

	return nil
}

func (s *StoreMemory) PutService(service string, spanID string) error {
	return nil
}

func (s *StoreMemory) PutOperation(service string, operation spanstore.Operation) error {
	return nil
}

// match service and operation by traceParameter
func (s *StoreMemory) matchSvcAndOp(tp storage.TraceParameters) []string {

	spanIDMap := make(map[string]struct{})
	var ret []string
	// match svc
	if ids, ok := s.services[tp.SvcName]; !ok {
		return nil
	} else {
		for _, id := range ids {
			spanIDMap[id] = struct{}{}
		}
	}

	// match operation
	for _, op := range s.operations[tp.SvcName] {

		if tp.OperationName != "" && op.Name != tp.OperationName {
			delete(spanIDMap, op.SpanID)
		}

	}

	// at this point you can return
	// and these two operations can filter out most of the irrelevant spans

	for id, _ := range spanIDMap {
		ret = append(ret, id)
	}

	return ret
}

// match tags and duration and start-end time
func matchOther(span *v1.Span, tp storage.TraceParameters) bool {

	if !tp.BeginTime.IsZero() && span.StartTimeUnixNano < uint64(tp.BeginTime.UnixNano()) {
		return false
	}
	if !tp.EndTime.IsZero() && span.EndTimeUnixNano > uint64(tp.EndTime.UnixNano()) {
		return false
	}

	if ok := func() bool {
		interval := span.EndTimeUnixNano - span.StartTimeUnixNano

		if tp.ElapsedMin != 0 && uint64(tp.ElapsedMin) > interval {
			return false
		}
		if tp.ElapsedMax != 0 && uint64(tp.ElapsedMax) < interval {
			return false
		}
		return true

	}(); !ok {
		return false
	}

	// match resource tags

	return true
}

func matchTags(tc *v1.TracesData, tp storage.TraceParameters) bool {

	var res []*cv1.KeyValue

	for _, rspans := range tc.ResourceSpans {

		res = append(res, flattenResourceTags(rspans.Resource)...)

		for _, ls := range rspans.InstrumentationLibrarySpans {
			res = append(res, flattenlibrarySpanTags(ls)...)
		}

	}

	// match tags kv
	for queryK, queryV := range tp.Resources {

		// In order to meet the matching in multiple tags
		if !findKeyValMatch(res, queryK, queryV) {
			return false
		}
	}

	return true
}

func findKeyValMatch(res []*cv1.KeyValue, queryK string, queryV string) bool {

	for _, kv := range res {
		if kv.Key == queryK && kv.Value.GetStringValue() == queryV {
			return true
		}
	}

	return false
}

func flattenResourceTags(resource *rv1.Resource) []*cv1.KeyValue {

	var retMe []*cv1.KeyValue

	retMe = append(retMe, resource.Attributes...)
	return retMe
}

func flattenlibrarySpanTags(ls *v1.InstrumentationLibrarySpans) []*cv1.KeyValue {

	var retMe []*cv1.KeyValue

	// append InstrumentationLibrary tag
	retMe = append(retMe, &cv1.KeyValue{
		Key:   ls.InstrumentationLibrary.GetName(),
		Value: &cv1.AnyValue{Value: &cv1.AnyValue_StringValue{StringValue: ls.InstrumentationLibrary.String()}},
	})

	// append span attributes
	for _, span := range ls.Spans {
		retMe = append(retMe, span.Attributes...)

		// append span events
		for _, event := range span.Events {
			retMe = append(retMe, event.Attributes...)
		}
	}

	return retMe
}

func sortTraceByTime(retMe []*v1.TracesData, i int, j int) bool {
	head := retMe[i].ResourceSpans[0].InstrumentationLibrarySpans[0].Spans[0].StartTimeUnixNano
	tail := retMe[j].ResourceSpans[0].InstrumentationLibrarySpans[0].Spans[0].StartTimeUnixNano

	return head < tail
}
