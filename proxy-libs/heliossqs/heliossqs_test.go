package heliossqs

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func getSpanRecorder() *tracetest.SpanRecorder {
	spanRecorder := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return spanRecorder
}

func validateMessageBody(t *testing.T, attributes []attribute.KeyValue, expected string) {
	found := false
	for _, attr := range attributes {
		if attr.Key == "messaging.payload" {
			found = true
			assert.Equal(t, expected, attr.Value.AsString())
		}
	}

	assert.True(t, found)
}

func assertSpans(t *testing.T, spans []sdktrace.ReadOnlySpan, messageBody string) {
	assert.Equal(t, 4, len(spans))
	createQueueSpan := spans[0]
	assert.Equal(t, "SQS", createQueueSpan.Name())
	assert.Equal(t, trace.SpanKindClient, createQueueSpan.SpanKind())
	assertAttributes(t, createQueueSpan.Attributes(), "CreateQueue")

	purgeQueueSpan := spans[1]
	assert.Equal(t, "SQS", purgeQueueSpan.Name())
	assert.Equal(t, trace.SpanKindClient, purgeQueueSpan.SpanKind())
	assertAttributes(t, purgeQueueSpan.Attributes(), "PurgeQueue")

	sendMessageSpan := spans[2]
	assert.Equal(t, "SQS", sendMessageSpan.Name())
	assert.Equal(t, trace.SpanKindClient, sendMessageSpan.SpanKind())
	sendMessageAttrs := sendMessageSpan.Attributes()
	assertAttributes(t, sendMessageAttrs, "SendMessage")
	validateMessageBody(t, sendMessageAttrs, messageBody)

	receiveMessageSpan := spans[3]
	assert.Equal(t, "SQS", receiveMessageSpan.Name())
	assert.Equal(t, trace.SpanKindClient, receiveMessageSpan.SpanKind())
}

func assertAttributes(t *testing.T, attributes []attribute.KeyValue, operation string) {
	for _, attribute := range attributes {
		key := attribute.Key
		value := attribute.Value.AsString()

		switch key {
		case "aws.operation":
			assert.Equal(t, operation, value)
		case "aws.service":
			assert.Equal(t, "SQS", value)
		}
	}
}

func TestContextPropagation(t *testing.T) {
	spanRecorder := getSpanRecorder()
	// init aws config
	ctx := context.Background()
	newCreds := credentials.NewStaticCredentialsProvider("test", "test", "")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://localhost:4566",
			SigningRegion: "us-east-1",
		}, nil
	})

	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithCredentialsProvider(newCreds), awsConfig.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := NewFromConfig(cfg)
	queueName := "test_queue"
	createQueue := &sqs.CreateQueueInput{QueueName: &queueName}
	createQueueResult, err := client.CreateQueue(ctx, createQueue)
	if err != nil {
		panic("failed creating queue, " + err.Error())
	}

	queueUrl := createQueueResult.QueueUrl
	_, err = client.PurgeQueue(ctx, &sqs.PurgeQueueInput{QueueUrl: queueUrl})
	if err != nil {
		log.Fatalf("failed to purge queue, %v", err)
		return
	}

	messageBody := "message body"
	_, err = client.SendMessage(ctx, &sqs.SendMessageInput{MessageBody: &messageBody, QueueUrl: queueUrl})
	if err != nil {
		log.Fatalf("failed to send message, %v", err)
		return
	}

	receiveMessageResult, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: queueUrl})
	if err != nil {
		log.Fatalf("failed to receive message, %v", err)
		return
	}

	message := receiveMessageResult.Messages[0]
	assert.Equal(t, messageBody, *message.Body)
	spanRecorder.ForceFlush(context.Background())
	spans := spanRecorder.Ended()
	assertSpans(t, spans, messageBody)
	messageAttributes := message.MessageAttributes
	traceparent := messageAttributes["traceparent"].StringValue
	assert.NotNil(t, traceparent)
	sendSpan := spans[2]
	spanContext := sendSpan.SpanContext()
	assert.Equal(t, fmt.Sprintf("00-%s-%s-01", spanContext.TraceID().String(), spanContext.SpanID().String()), *traceparent)
}
