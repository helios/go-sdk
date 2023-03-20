package heliosgrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	pb "github.com/helios/go-sdk/proxy-libs/heliosgrpc/chat"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type GrpcServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *GrpcServer) SayHello(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &pb.Message{Body: "Hello From the Server!"}, nil
}

func validateAttributes(attrs []attribute.KeyValue, t *testing.T) {
	for _, value := range attrs {
		key := value.Key
		if key == semconv.RPCSystemKey {
			assert.Equal(t, "grpc", value.Value.AsString())
		} else if key == semconv.RPCMethodKey {
			assert.Equal(t, "SayHello", value.Value.AsString())
		} else if key == semconv.RPCServiceKey {
			assert.Equal(t, "chat.ChatService", value.Value.AsString())
		}
	}
}

var unaryInterceptorCalled bool = false

func noopUnaryInterceptor(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error) {
	unaryInterceptorCalled = true
	return handler(ctx, req)
}

func noopstreamInterceptor(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error {
	return nil
}

func initTracing(t *testing.T) *tracetest.SpanRecorder{
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return sr
}

func registerServerAndPerformCall(t *testing.T, port string) *pb.Message {
	go func() {
		lis, _ := net.Listen("tcp", ":" + port)
		grpcServer := NewServer(UnaryInterceptor(noopUnaryInterceptor), StreamInterceptor(noopstreamInterceptor))
		s := GrpcServer{}
		pb.RegisterChatServiceServer(grpcServer, &s)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	var conn *ClientConn
	conn, err := Dial(fmt.Sprintf("localhost:%s", port), WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewChatServiceClient(conn)

	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	response, _ := client.SayHello(ctx, &pb.Message{Body: "helios"})
	assert.Equal(t, "Hello From the Server!", response.Body)
	return response
}

func TestServerInstrumentation(t *testing.T) {
	sr := initTracing(t)

	port := "3030"
	registerServerAndPerformCall(t, port)

	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 2, len(spans))

	serverSpan := spans[0]
	assert.Equal(t, trace.SpanKind(2), serverSpan.SpanKind())
	validateAttributes(serverSpan.Attributes(), t)
	clientSpan := spans[1]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.False(t, clientSpan.Parent().HasTraceID())
	validateAttributes(clientSpan.Attributes(), t)
	assert.Equal(t, serverSpan.Parent().SpanID(), clientSpan.SpanContext().SpanID())
	assert.True(t, unaryInterceptorCalled)
}

func TestDisableInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	defer os.Setenv("HS_DISABLED", "")

	sr := initTracing(t)

	port := "3031"
	registerServerAndPerformCall(t, port)

	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 0, len(spans))
}
