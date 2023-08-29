package reader

import (
	"context"
	"fmt"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

type GrpcReader struct {
	cli api_v3.QueryServiceClient
}

type GrpcReaderConfig struct {
	Address string
	Port    int
}

func NewGrpcReader(conf GrpcReaderConfig) (*GrpcReader, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", conf.Address, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// clog.CL.Error("grpc reader unable to connect the target server")
		return nil, err
	}

	cli := api_v3.NewQueryServiceClient(conn)

	return &GrpcReader{
		cli: cli,
	}, nil
}

func (g *GrpcReader) SearchTraces(ctx context.Context, in SearchTracesRequest) (*SpansResponse, error) {

	trace, err := g.cli.FindTraces(ctx, &api_v3.FindTracesRequest{
		Query: in.SearchParam,
	})
	if err != nil {
		return nil, err
	}

	sr := &SpansResponse{}
	for {
		recv, err := trace.Recv()
		if err != nil {
			if err == io.EOF {
				return sr, nil
			} else {
				return nil, err
			}
		}

		sr.ResourceSpans = append(sr.ResourceSpans, recv.ResourceSpans...)
	}

}

func (g *GrpcReader) QueryTrace(ctx context.Context, in QueryTraceRequest) (*SpansResponse, error) {

	trace, err := g.cli.GetTrace(ctx, &api_v3.GetTraceRequest{TraceId: in.TraceID})
	if err != nil {
		return nil, err
	}

	sr := &SpansResponse{}
	for {
		recv, err := trace.Recv()
		if err != nil {
			if err == io.EOF {
				return sr, nil
			} else {
				return nil, err
			}
		}
		sr.ResourceSpans = append(sr.ResourceSpans, recv.ResourceSpans...)
	}

}

func (g *GrpcReader) QueryServices(ctx context.Context, in QueryServicesRequest) (*ServicesResponse, error) {
	return nil, nil
}

func (g *GrpcReader) QueryOperations(ctx context.Context, in QueryOperationsRequest) (*OperationsResponse, error) {
	return nil, nil
}
