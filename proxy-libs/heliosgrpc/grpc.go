package heliosgrpc

import (
	"context"

	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	realGrpc "google.golang.org/grpc"
)

var InstrumentedSymbols = [...]string{"NewServer", "Dial", "DialContext"}

func Dial(target string, opts ...realGrpc.DialOption) (*realGrpc.ClientConn, error) {
	newOptions := append(opts, realGrpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()), realGrpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	return realGrpc.Dial(target, newOptions...)
}

func DialContext(ctx context.Context, target string, opts ...realGrpc.DialOption) (conn *realGrpc.ClientConn, err error) {
	newOptions := append(opts, realGrpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()), realGrpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	return realGrpc.DialContext(ctx, target, newOptions...)
}

func NewServer(opt ...realGrpc.ServerOption) *realGrpc.Server {
	newOptions := append(opt, realGrpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()), realGrpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
	return realGrpc.NewServer(newOptions...)
}
