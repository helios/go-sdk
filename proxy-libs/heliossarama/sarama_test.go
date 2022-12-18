package heliossarama

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	exportsExtractor "github.com/helios/go-instrumentor/exports_extractor"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

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
	assert.Equal(handler.t, 2, len(spans))
	producerSpan := spans[0]
	consumerSpan := spans[1]
	assert.Equal(handler.t, fmt.Sprintf("%v send", Topic), producerSpan.Name())
	assert.Equal(handler.t, trace.SpanKindProducer, producerSpan.SpanKind())
	handler.assertAttributes(producerSpan, "send")
	assert.Equal(handler.t, fmt.Sprintf("%v receive", Topic), consumerSpan.Name())
	assert.Equal(handler.t, trace.SpanKindConsumer, consumerSpan.SpanKind())
	handler.assertAttributes(consumerSpan, "receive")
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

func TestNewAsyncProducerAndNewConsumerGroupInstrumentations(t *testing.T) {
	spanRecorder := getSpanRecorder()
	config := getConfig()
	deleteTopic(config)
	key := "0"
	value := "Hello, World!"

	asyncProducer, _ := NewAsyncProducer(Addresses, config)
	asyncProducer.Input() <- &ProducerMessage{
		Topic: Topic,
		Key:   StringEncoder(key),
		Value: StringEncoder(value),
	}
	<-asyncProducer.Successes()
	asyncProducer.Close()

	consumerGroup, _ := NewConsumerGroup(Addresses, "consumerGroup", config)
	consumerGroup.Consume(context.Background(), []string{Topic}, &TestConsumerGroupHandler{
		messageKey:   key,
		messageValue: value,
		spanRecorder: spanRecorder,
		t:            t,
	})
	// consumerGroup.Close()
}

func TestNewSyncProducerAndNewConsumerGroupFromClientInstrumentations(t *testing.T) {
	spanRecorder := getSpanRecorder()
	config := getConfig()
	deleteTopic(config)
	key := "1"
	value := "Welcome to Helios!"

	syncProducer, _ := NewSyncProducer(Addresses, config)
	syncProducer.SendMessage(&ProducerMessage{
		Topic: Topic,
		Key:   StringEncoder(key),
		Value: StringEncoder(value),
	})
	syncProducer.Close()

	client, _ := NewClient(Addresses, config)
	consumerGroup, _ := NewConsumerGroupFromClient("consumerGroupFromClient", client)
	consumerGroup.Consume(context.Background(), []string{Topic}, &TestConsumerGroupHandler{
		messageKey:   key,
		messageValue: value,
		spanRecorder: spanRecorder,
		t:            t,
	})
	// consumerGroup.Close()
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

	// Delete original exports that cannot be wrapped by proxy Helios exports.
	originalExports = delete(originalExports, "NewMockWrapper")
	originalExports = delete(originalExports, "Wrap")

	// Get Helios sarama exports.
	srcDir, _ := filepath.Abs(".")
	heliosExports := exportsExtractor.ExtractExports(srcDir, "heliossarama")
	sort.Slice(heliosExports, func(i int, j int) bool { return heliosExports[i].Name < heliosExports[j].Name })

	assert.Equal(t, len(originalExports), len(heliosExports))
	assert.EqualValues(t, originalExports, heliosExports)
}
