package reader

import (
	"context"
	"fmt"
	"github.com/jaegertracing/jaeger/proto-gen/api_v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 16685
)

type GrpcReader struct {
	cli api_v3.QueryServiceClient
}

type GrpcReaderConfig struct {
	Host string
	Port int
	opts []grpc.DialOption
}

// NewGrpcReader TODO: Config set default val
func NewGrpcReader(conf GrpcReaderConfig) (*GrpcReader, error) {

	conf.setDefault()

	address := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	conn, err := grpc.Dial(
		address,
		conf.opts...,
	)
	if err != nil {
		// TODO: 指数回退?
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

func (g *GrpcReaderConfig) setDefault() {

	if g.Host == "" {
		g.Host = defaultHost
	}

	if g.Port == 0 {
		g.Port = defaultPort
	}

	if len(g.opts) == 0 {
		g.opts = append(g.opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

}
