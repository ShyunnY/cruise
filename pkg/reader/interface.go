package reader

import (
	"context"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	v1 "github.com/jaegertracing/jaeger/proto-gen/otel/trace/v1"
)

type Reader interface {
	SearchTraces(ctx context.Context, in SearchTracesRequest) (*SpansResponse, error)
	QueryTrace(ctx context.Context, in QueryTraceRequest) (*SpansResponse, error)
	QueryServices(ctx context.Context, in QueryServicesRequest) (*ServicesResponse, error)
	QueryOperations(ctx context.Context, in QueryOperationsRequest) (*OperationsResponse, error)
}

type SearchTracesRequest struct {
	SearchParam *api_v3.TraceQueryParameters
}

type QueryTraceRequest struct {
	TraceID string
}

type QueryServicesRequest struct {
}

type QueryOperationsRequest struct {
	Service  string
	SpanKind string
}

type SpansResponse struct {
	ResourceSpans []*v1.ResourceSpans
}

type ServicesResponse struct {
	Services []string
}

type OperationsResponse struct {
	Operations []*api_v3.Operation
}
