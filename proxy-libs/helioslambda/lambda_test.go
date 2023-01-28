package helioslambda

import (
	"context"
	"reflect"
	"testing"

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
)

const response = "hello world"

func validateResults(t *testing.T, resp []reflect.Value) {
	assert.Len(t, resp, 2)
	assert.Equal(t, response, resp[0].Interface())
	assert.Nil(t, resp[1].Interface())

	spans := exporter.GetSpans()
	assert.Len(t, spans, 2)
	lambdaSpan := spans[1]
	assert.Equal(t, traceId, lambdaSpan.SpanContext.TraceID().String())
	assert.Equal(t, parentSpanId, lambdaSpan.Parent.SpanID().String())
	customSpan := spans[0]
	assert.Equal(t, traceId, customSpan.SpanContext.TraceID().String())
	assert.Equal(t, lambdaSpan.SpanContext.SpanID().String(), customSpan.Parent.SpanID().String())
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
	validateResults(t, resp)
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
	validateResults(t, resp)
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
	validateResults(t, resp)
}
