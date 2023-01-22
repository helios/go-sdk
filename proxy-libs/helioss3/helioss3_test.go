package helioss3

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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

func assertSpan(t *testing.T, spanRecorder *tracetest.SpanRecorder) []attribute.KeyValue {
	spanRecorder.ForceFlush(context.Background())
	spans := spanRecorder.Ended()
	assert.Equal(t, 1, len(spans))
	span := spans[0]
	assert.Equal(t, "S3", span.Name())
	assert.Equal(t, trace.SpanKindClient, span.SpanKind())
	return span.Attributes()
}

func assertAttributes(t *testing.T, attributes []attribute.KeyValue) {
	for _, attribute := range attributes {
		key := attribute.Key
		value := attribute.Value.AsString()

		switch key {
		case "aws.operation":
			assert.Equal(t, "ListBuckets", value)
		case "aws.service":
			assert.Equal(t, "S3", value)
		}
	}
}

func TestListBuckets(t *testing.T) {
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
	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithCredentialsProvider(newCreds) ,awsConfig.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	s3Client := NewFromConfig(cfg)
	input := &ListBucketsInput{}
	_, err = s3Client.ListBuckets(ctx, input)
	if err != nil {
		log.Fatalf("Got an error retrieving buckets, %v", err)
	}
	attributes := assertSpan(t, spanRecorder)
	fmt.Println(attributes)
	assertAttributes(t, attributes)
}
