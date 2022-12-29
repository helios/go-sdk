package heliossarama

import (
	"context"
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
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
