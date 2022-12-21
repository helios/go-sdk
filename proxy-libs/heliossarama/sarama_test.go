package heliossarama

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/Shopify/sarama"
	exportsExtractor "github.com/helios/go-instrumentor/exports_extractor"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const testRootSpanName = "test_root"

var (
	Addresses []string = []string{"localhost:9093"}
	Topic     string   = "testTopic"
)

type TestConsumerGroupHandler struct {
	messageKey   string
	messageValue string
	spanRecorder *tracetest.SpanRecorder
	t            *testing.T
}

func (handler *TestConsumerGroupHandler) Setup(session ConsumerGroupSession) error {
	return nil
}

func (handler *TestConsumerGroupHandler) Cleanup(session ConsumerGroupSession) error {
	return nil
}

func (handler *TestConsumerGroupHandler) ConsumeClaim(session ConsumerGroupSession, claim ConsumerGroupClaim) error {
	message := <-claim.Messages()
	assert.Equal(handler.t, Topic, message.Topic)
	assert.Equal(handler.t, handler.messageKey, string(message.Key))
	assert.Equal(handler.t, handler.messageValue, string(message.Value))
	session.MarkMessage(message, "")
	handler.assertSpans()
	return nil
}

func (handler *TestConsumerGroupHandler) assertSpans() {
	handler.spanRecorder.ForceFlush(context.Background())
	spans := handler.spanRecorder.Ended()
	assert.Equal(handler.t, 3, len(spans))
	rootTestSpan := spans[0]
	producerSpan := spans[1]
	consumerSpan := spans[2]
	assert.Equal(handler.t, testRootSpanName, rootTestSpan.Name())
	assert.Equal(handler.t, fmt.Sprintf("%v send", Topic), producerSpan.Name())
	assert.Equal(handler.t, trace.SpanKindProducer, producerSpan.SpanKind())
	handler.assertAttributes(producerSpan, "send")
	assert.Equal(handler.t, fmt.Sprintf("%v receive", Topic), consumerSpan.Name())
	assert.Equal(handler.t, trace.SpanKindConsumer, consumerSpan.SpanKind())
	handler.assertAttributes(consumerSpan, "receive")
	assert.Equal(handler.t, rootTestSpan.SpanContext().SpanID(), producerSpan.Parent().SpanID())
	assert.Equal(handler.t, producerSpan.SpanContext().SpanID(), consumerSpan.Parent().SpanID())
}

func (handler *TestConsumerGroupHandler) assertAttributes(span sdktrace.ReadOnlySpan, messagingOperation string) {
	for _, attribute := range span.Attributes() {
		key := attribute.Key
		value := attribute.Value.AsString()

		switch key {
		case semconv.MessagingDestinationKey:
			assert.Equal(handler.t, Topic, value)
		case semconv.MessagingDestinationKindKey:
			assert.Equal(handler.t, "topic", value)
		case semconv.MessagingOperationKey:
			assert.Equal(handler.t, messagingOperation, value)
		case semconv.MessagingSystemKey:
			assert.Equal(handler.t, "kafka", value)
		}
	}
}

func getSpanRecorder() *tracetest.SpanRecorder {
	spanRecorder := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return spanRecorder
}

func getConfig() *Config {
	config := NewConfig()
	config.Version = V2_5_0_0
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = OffsetOldest
	return config
}

func deleteTopic(config *Config) {
	clusterAdmin, _ := sarama.NewClusterAdmin(Addresses, config)
	clusterAdmin.DeleteTopic(Topic)
	clusterAdmin.Close()
}

func createRootSpanAndInjectMessage(message *sarama.ProducerMessage) {
	ctx, span := otel.GetTracerProvider().Tracer("test").Start(context.Background(), testRootSpanName)
	span.End()

	InjectContextToMessage(ctx, message)
}

func TestNewAsyncProducerAndNewConsumerGroupInstrumentations(t *testing.T) {
	spanRecorder := getSpanRecorder()
	config := getConfig()
	deleteTopic(config)
	key := "0"
	value := "Hello, World!"

	asyncProducer, _ := NewAsyncProducer(Addresses, config)
	message := ProducerMessage{
		Topic: Topic,
		Key:   StringEncoder(key),
		Value: StringEncoder(value),
	}

	createRootSpanAndInjectMessage(&message)
	asyncProducer.Input() <- &message
	<-asyncProducer.Successes()
	asyncProducer.Close()

	consumerGroup, _ := NewConsumerGroup(Addresses, "consumerGroup", config)
	consumerGroup.Consume(context.Background(), []string{Topic}, &TestConsumerGroupHandler{
		messageKey:   key,
		messageValue: value,
		spanRecorder: spanRecorder,
		t:            t,
	})
}

func TestNewSyncProducerAndNewConsumerGroupFromClientInstrumentations(t *testing.T) {
	spanRecorder := getSpanRecorder()
	config := getConfig()
	deleteTopic(config)
	key := "1"
	value := "Welcome to Helios!"

	syncProducer, _ := NewSyncProducer(Addresses, config)
	message := ProducerMessage{
		Topic: Topic,
		Key:   StringEncoder(key),
		Value: StringEncoder(value),
	}
	createRootSpanAndInjectMessage(&message)
	syncProducer.SendMessage(&message)
	syncProducer.Close()

	client, _ := NewClient(Addresses, config)
	consumerGroup, _ := NewConsumerGroupFromClient("consumerGroupFromClient", client)
	consumerGroup.Consume(context.Background(), []string{Topic}, &TestConsumerGroupHandler{
		messageKey:   key,
		messageValue: value,
		spanRecorder: spanRecorder,
		t:            t,
	})
}

func TestInterfaceMatch(t *testing.T) {
	delete := func(exports []exportsExtractor.ExtractedObject, name string) []exportsExtractor.ExtractedObject {
		for i, export := range exports {
			if export.Name == name {
				return append(exports[:i], exports[i+1:]...)
			}
		}

		return exports
	}

	// Get original sarama exports.
	originalRepository := exportsExtractor.CloneGitRepository("https://github.com/Shopify/sarama", "v1.37.2")
	originalExports := exportsExtractor.ExtractExports(originalRepository, "sarama")
	os.RemoveAll(originalRepository)
	sort.Slice(originalExports, func(i int, j int) bool { return originalExports[i].Name < originalExports[j].Name })

	// Get Helios sarama exports.
	srcDir, _ := filepath.Abs(".")
	heliosExports := exportsExtractor.ExtractExports(srcDir, "heliossarama")
	sort.Slice(heliosExports, func(i int, j int) bool { return heliosExports[i].Name < heliosExports[j].Name })

	// "NewMockWrapper" cannot be wrapped because its parameter's type is private - Remove it from the expected list.
	originalExports = delete(originalExports, "NewMockWrapper")
	// The signature of "Wrap" was changed because the original return type is private - Remove it from both lists.
	originalExports = delete(originalExports, "Wrap")
	heliosExports = delete(heliosExports, "Wrap")
	// A helper method we've added to improve context propagation
	heliosExports = delete(heliosExports, "InjectContextToMessage")

	assert.Equal(t, len(originalExports), len(heliosExports))
	assert.EqualValues(t, originalExports, heliosExports)
}
