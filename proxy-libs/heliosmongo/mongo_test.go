package heliosmongo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

func getSpanRecorder() *tracetest.SpanRecorder {
	spanRecorder := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return spanRecorder
}

func getClientOptions() *options.ClientOptions {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://localhost:27017")
	return clientOptions
}

func insertUser(client *mongo.Client, id int, name string, role string) {
	users := client.Database("test").Collection("users")
	_, error := users.InsertOne(context.Background(), bson.D{
		{Key: "id", Value: id},
		{Key: "name", Value: name},
		{Key: "role", Value: role},
	})

	if error != nil {
		panic(error)
	}
}

func assertSpan(t *testing.T, spanRecorder *tracetest.SpanRecorder) []attribute.KeyValue {
	spanRecorder.ForceFlush(context.Background())
	spans := spanRecorder.Ended()
	assert.Equal(t, 1, len(spans))
	span := spans[0]
	assert.Equal(t, "users.insert", span.Name())
	assert.Equal(t, trace.SpanKindClient, span.SpanKind())
	return span.Attributes()
}

func assertAttributes(t *testing.T, attributes []attribute.KeyValue, id int, name string, role string) {
	for _, attribute := range attributes {
		key := attribute.Key
		value := attribute.Value.AsString()

		switch key {
		case semconv.DBMongoDBCollectionKey:
			assert.Equal(t, "users", value)
		case semconv.DBNameKey:
			assert.Equal(t, "test", value)
		case semconv.DBOperationKey:
			assert.Equal(t, "insert", value)
		case semconv.DBStatementKey:
			assert.Contains(t, value, fmt.Sprintf("\"id\":%v", id))
			assert.Contains(t, value, fmt.Sprintf("\"name\":\"%v\"", name))
			assert.Contains(t, value, fmt.Sprintf("\"role\":\"%v\"", role))
		case semconv.DBSystemKey:
			assert.Equal(t, "mongodb", value)
		}
	}
}

func registerMongoConnectAndInsertUser(t *testing.T, userId int, username string, company string) {
	clientOptions := getClientOptions()
	client, error := Connect(context.Background(), clientOptions)

	if error != nil {
		panic(error)
	}

	client.Connect(context.Background())
	insertUser(client, userId, username, company)
}

func TestConnectInstrumentation(t *testing.T) {
	spanRecorder := getSpanRecorder()

	registerMongoConnectAndInsertUser(t, 12345, "Lior Govrin", "Software Engineer")

	attributes := assertSpan(t, spanRecorder)
	assertAttributes(t, attributes, 12345, "Lior Govrin", "Software Engineer")
}

func TestDisableInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	defer os.Setenv("HS_DISABLED", "")

	spanRecorder := getSpanRecorder()

	registerMongoConnectAndInsertUser(t, 12345, "Lior Govrin", "Software Engineer")

	spanRecorder.ForceFlush(context.Background())
	spans := spanRecorder.Ended()
	assert.Equal(t, 0, len(spans))
}

func registerMongoClientAndInsertUser(t *testing.T, userId int, username string, company string) {
	clientOptions := getClientOptions()
	client, error := NewClient(clientOptions)

	if error != nil {
		panic(error)
	}

	client.Connect(context.Background())
	insertUser(client, userId, username, company)
}

func TestNewClientInstrumentation(t *testing.T) {
	spanRecorder := getSpanRecorder()

	registerMongoClientAndInsertUser(t, 67890, "Bob McClown", "Company Jester")

	attributes := assertSpan(t, spanRecorder)
	assertAttributes(t, attributes, 67890, "Bob McClown", "Company Jester")
}

func TestDisableClientInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	defer os.Setenv("HS_DISABLED", "")

	spanRecorder := getSpanRecorder()

	registerMongoClientAndInsertUser(t, 67890, "Bob McClown", "Company Jester")

	spanRecorder.ForceFlush(context.Background())
	spans := spanRecorder.Ended()
	assert.Equal(t, 0, len(spans))
}
