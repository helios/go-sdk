package heliosgrpc

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/helios/helios-go-instrumenter/exports_extractor"
	pb "github.com/helios/helios-proxy-libs/heliosgrpc/chat"

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

func TestServerInstrumentation(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	go func() {
		lis, _ := net.Listen("tcp", ":3030")
		grpcServer := NewServer()
		s := GrpcServer{}
		pb.RegisterChatServiceServer(grpcServer, &s)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	var conn *ClientConn
	conn, err := Dial("localhost:3030", WithTransportCredentials(insecure.NewCredentials()))

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
}

func TestInterfaceMatch(t *testing.T) {
	grpcRepoFolder := exports_extractor.CloneGitRepository("https://github.com/grpc/grpc-go", "v1.50.1")
	grpcPackageName := "grpc"
	grpcExports := exports_extractor.ExtractExports(grpcRepoFolder, grpcPackageName)
	os.RemoveAll(grpcRepoFolder)
	sort.Slice(grpcExports, func(i, j int) bool {
		return grpcExports[i].Name < grpcExports[j].Name
	})

	heliosGrpcRoot, _ := filepath.Abs(".")
	heliosGrpcPackageName := "heliosgrpc"
	heliosGrpcExports := exports_extractor.ExtractExports(heliosGrpcRoot, heliosGrpcPackageName)
	sort.Slice(heliosGrpcExports, func(i, j int) bool {
		return heliosGrpcExports[i].Name < heliosGrpcExports[j].Name
	})

	assert.EqualValues(t, grpcExports, heliosGrpcExports)
}
