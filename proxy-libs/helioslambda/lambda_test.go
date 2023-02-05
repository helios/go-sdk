package helioslambda

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var (
	traceId               = "83d8d6c5347593d092e9409f4978bd51"
	parentSpanId          = "6f2a23d2d1e9159c"
	tracingHeader         = "00" + "-" + traceId + "-" + parentSpanId + "-" + "01"
	traceCarrier          = map[string]string{"traceparent": tracingHeader}
	testApiGatewayEvent   = apiGatewayEvent{Headers: traceCarrier}
	testEventBridgeEvent1 = eventBridgeEvent{Detail: traceCarrier}
	testEventBridgeEvent2 = eventBridgeEvent{TraceHeader: tracingHeader}
	exporter              = tracetest.NewInMemoryExporter()
	provider              = trace.NewTracerProvider(trace.WithBatcher(exporter))
	testSqsMessage        = events.SQSMessage{MessageId: "1234", MessageAttributes: map[string]events.SQSMessageAttribute{"traceparent": {DataType: "String", StringValue: &tracingHeader}}}
	testSqsRecord         = events.SQSEvent{Records: []events.SQSMessage{testSqsMessage}}
)

const response = "hello world"

func assertPayloads(t *testing.T, span tracetest.SpanStub, expectedEvent string) {
	foundRes := false
	foundEvent := false
	for _, attr := range span.Attributes {
		if attr.Key == "faas.res" {
			foundRes = true
			assert.Equal(t, response, attr.Value.AsString())
		} else if attr.Key == "faas.event" {
			foundEvent = true
			assert.Equal(t, expectedEvent, attr.Value.AsString())
		}
	}

	assert.True(t, foundRes)
	assert.True(t, foundEvent)
}

func validateResults(t *testing.T, resp []reflect.Value, expectedEvent string) {
	assert.Len(t, resp, 2)
	assert.Equal(t, response, resp[0].Interface())
	assert.Nil(t, resp[1].Interface())

	spans := exporter.GetSpans()
	assert.Len(t, spans, 2)
	lambdaSpan := spans[1]
	assert.Equal(t, traceId, lambdaSpan.SpanContext.TraceID().String())
	assert.Equal(t, parentSpanId, lambdaSpan.Parent.SpanID().String())
	assertPayloads(t, lambdaSpan, expectedEvent)
	customSpan := spans[0]
	assert.Equal(t, traceId, customSpan.SpanContext.TraceID().String())
	assert.Equal(t, lambdaSpan.SpanContext.SpanID().String(), customSpan.Parent.SpanID().String())
}

func validateSqsTestResults(t *testing.T, resp []reflect.Value) {
	assert.Len(t, resp, 2)
	assert.Equal(t, response, resp[0].Interface())
	assert.Nil(t, resp[1].Interface())

	spans := exporter.GetSpans()
	assert.Len(t, spans, 3)
	lambdaSqsHandlerSpan := spans[1]
	assert.Equal(t, traceId, lambdaSqsHandlerSpan.SpanContext.TraceID().String())
	assert.Equal(t, parentSpanId, lambdaSqsHandlerSpan.Parent.SpanID().String())
	customSpan := spans[0]
	assert.Equal(t, traceId, customSpan.SpanContext.TraceID().String())
	assert.Equal(t, lambdaSqsHandlerSpan.SpanContext.SpanID().String(), customSpan.Parent.SpanID().String())
}

func TestApiGatewayContextPropagation(t *testing.T) {
	ctx := context.Background()
	exporter.Reset()
	otel.SetTracerProvider(provider)

	customerHandler := func(lambdaContext context.Context, event apiGatewayEvent) (string, error) {
		_, customSpan := provider.Tracer("test").Start(lambdaContext, "custom_span")
		customSpan.End()
		return response, nil
	}

	wrapped := instrumentHandler(customerHandler)

	wrappedCallable := reflect.ValueOf(wrapped)
	resp := wrappedCallable.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(testApiGatewayEvent)})
	rawEvent, _ := json.Marshal(testApiGatewayEvent)
	validateResults(t, resp, string(rawEvent))
}

func TestApiGatewayContextPropagationWithObfuscation(t *testing.T) {
	ctx := context.Background()
	exporter.Reset()
	otel.SetTracerProvider(provider)

	blocklistRules,_ := json.Marshal([]string{"$.headers.*"})
	obfuscatedExpectedPayload := "{\"headers\":{\"traceparent\":\"4161107f\"}}"

	os.Setenv("HS_DATA_OBFUSCATION_HMAC_KEY", "12345")
	os.Setenv("HS_DATA_OBFUSCATION_BLOCKLIST", string(blocklistRules))

	customerHandler := func(lambdaContext context.Context, event apiGatewayEvent) (string, error) {
		_, customSpan := provider.Tracer("test").Start(lambdaContext, "custom_span")
		customSpan.End()
		return response, nil
	}

	wrapped := instrumentHandler(customerHandler)

	wrappedCallable := reflect.ValueOf(wrapped)
	resp := wrappedCallable.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(testApiGatewayEvent)})
	validateResults(t, resp, string(obfuscatedExpectedPayload))
}

func TestEventbridgeContextPropagationInDetail(t *testing.T) {
	exporter.Reset()
	ctx := context.Background()
	otel.SetTracerProvider(provider)

	customerHandler := func(lambdaContext context.Context, event eventBridgeEvent) (string, error) {
		_, customSpan := provider.Tracer("test").Start(lambdaContext, "custom_span")
		customSpan.End()
		return response, nil
	}

	wrapped := instrumentHandler(customerHandler)

	wrappedCallable := reflect.ValueOf(wrapped)
	resp := wrappedCallable.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(testEventBridgeEvent1)})
	rawEvent, _ := json.Marshal(testEventBridgeEvent1)
	validateResults(t, resp, string(rawEvent))
}

func TestEventbridgeContextPropagationInTraceHeader(t *testing.T) {
	exporter.Reset()
	ctx := context.Background()
	otel.SetTracerProvider(provider)

	customerHandler := func(lambdaContext context.Context, event eventBridgeEvent) (string, error) {
		_, customSpan := provider.Tracer("test").Start(lambdaContext, "custom_span")
		customSpan.End()
		return response, nil
	}

	wrapped := instrumentHandler(customerHandler)

	wrappedCallable := reflect.ValueOf(wrapped)
	resp := wrappedCallable.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(testEventBridgeEvent2)})
	rawEvent, _ := json.Marshal(testEventBridgeEvent2)
	validateResults(t, resp, string(rawEvent))
}

func TestSqsContextPropagationInMessageAttribute(t *testing.T) {
	exporter.Reset()
	ctx := context.Background()
	otel.SetTracerProvider(provider)

	innerMethod := func(lambdaContext context.Context, event events.SQSMessage) (any, error) {
		_, customSpan := provider.Tracer("test").Start(lambdaContext, "custom_span")
		customSpan.End()
		return response, nil
	}

	newHandler := func(lambdaContext context.Context, event events.SQSEvent) (any, error) {
		var returnVal any
		for _, record := range event.Records {
			returnVal, _ = HandleRecord(lambdaContext, record, innerMethod)
		}
		return returnVal, nil
	}

	wrapped := instrumentHandler(newHandler)

	wrappedCallable := reflect.ValueOf(wrapped)
	resp := wrappedCallable.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(testSqsRecord)})
	validateSqsTestResults(t, resp)
}

func TestSqsEventBridgeContextPropagaion(t *testing.T) {
	jsonFile, err := os.Open("sqsMessage.json")
	if err != nil {
		assert.Fail(t, "could not open json file")
	}
	defer jsonFile.Close()
	var record events.SQSEvent
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &record)
	if err != nil {
		assert.Fail(t, "could not parse json file")
	}

	exporter.Reset()
	ctx := context.Background()
	otel.SetTracerProvider(provider)

	innerMethod := func(lambdaContext context.Context, event events.SQSMessage) (any, error) {
		_, customSpan := provider.Tracer("test").Start(lambdaContext, "custom_span")
		customSpan.End()
		return response, nil
	}

	newHandler := func(lambdaContext context.Context, event events.SQSEvent) (any, error) {
		var returnVal any
		for _, record := range event.Records {
			returnVal, _ = HandleRecord(lambdaContext, record, innerMethod)
		}
		return returnVal, nil
	}

	wrapped := instrumentHandler(newHandler)

	wrappedCallable := reflect.ValueOf(wrapped)
	resp := wrappedCallable.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(record)})
	validateSqsTestResults(t, resp)
}
