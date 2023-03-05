package heliossqs

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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

func validateTraceparentMessageAttribute(t *testing.T, message types.Message, sendSpan sdktrace.ReadOnlySpan) {
	messageAttributes := message.MessageAttributes
	traceparent := messageAttributes["traceparent"].StringValue
	assert.NotNil(t, traceparent)
	spanContext := sendSpan.SpanContext()
	assert.Equal(t, fmt.Sprintf("00-%s-%s-01", spanContext.TraceID().String(), spanContext.SpanID().String()), *traceparent)
}

func initAwsConfig(t *testing.T, ctx context.Context) aws.Config{

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

	return cfg
}

func initSqsClient(t *testing.T, ctx context.Context) *sqs.Client {
	cfg := initAwsConfig(t, ctx)
	client := NewFromConfig(cfg)
	return client
}

func createSqsQueue(t *testing.T, client *sqs.Client, ctx context.Context, queueName string) *sqs.CreateQueueOutput {
	createQueue := &sqs.CreateQueueInput{QueueName: &queueName}
	createQueueResult, err := client.CreateQueue(ctx, createQueue)
	if err != nil {
		panic("failed creating queue, " + err.Error())
	}

	queueUrl := createQueueResult.QueueUrl
	_, err = client.PurgeQueue(ctx, &sqs.PurgeQueueInput{QueueUrl: queueUrl})
	if err != nil {
		log.Fatalf("failed to purge queue, %v", err)
		return nil
	}

	return createQueueResult
}

func TestContextPropagation(t *testing.T) {
	spanRecorder := getSpanRecorder()
	ctx := context.Background()
	client := initSqsClient(t, ctx)

	createQueueResult := createSqsQueue(t, client, ctx, "test_queue")
	queueUrl := createQueueResult.QueueUrl
	
	messageBody := "message body"
	_, err := client.SendMessage(ctx, &sqs.SendMessageInput{MessageBody: &messageBody, QueueUrl: queueUrl})
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
	validateTraceparentMessageAttribute(t, message, spans[2])

	messageId := "abcd1234"
	messageBody2 := "messageBody2"
	entry := types.SendMessageBatchRequestEntry{Id: &messageId, MessageBody: &messageBody2}
	_, err = client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{QueueUrl: queueUrl, Entries: []types.SendMessageBatchRequestEntry{entry}})
	if err != nil {
		print(1234)
	}
	receiveMessageResult, _ = client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: queueUrl})
	message = receiveMessageResult.Messages[0]
	assert.Equal(t, messageBody2, *message.Body)
	spanRecorder.ForceFlush(context.Background())
	spans = spanRecorder.Ended()
	validateTraceparentMessageAttribute(t, message, spans[4])
}

func TestDisableInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	defer os.Setenv("HS_DISABLED", "")

	spanRecorder := getSpanRecorder()
	ctx := context.Background()
	client := initSqsClient(t, ctx)

	createQueueResult := createSqsQueue(t, client, ctx, "test_queue")
	queueUrl := createQueueResult.QueueUrl

	messageBody := "message body"
	_, err := client.SendMessage(ctx, &sqs.SendMessageInput{MessageBody: &messageBody, QueueUrl: queueUrl})
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

	assert.Equal(t, 0, len(spans))
}
