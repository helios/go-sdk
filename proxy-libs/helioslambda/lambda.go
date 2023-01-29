package helioslambda

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

var InstrumentedSymbols = [...]string{"Start", "StartWithContext", "StartWithOptions"}

type apiGatewayEvent struct {
	Headers map[string]string `json:"headers"`
}

type eventBridgeEvent struct {
	Detail      map[string]string `json:"detail"`
	TraceHeader string            `json:"trace-header"`
}

func heliosEventToCarrier(eventJSON []byte) propagation.TextMapCarrier {
	const traceParentHeader = "Traceparent"
	const lowerCaseTraceParentHeader = "traceparent"

	// Try API Gateway context propagation
	var headers apiGatewayEvent
	err := json.Unmarshal(eventJSON, &headers)
	if err == nil && headers.Headers != nil {
		if val, ok := headers.Headers[traceParentHeader]; ok {
			return propagation.HeaderCarrier{traceParentHeader: []string{val}}
		} else if val, ok = headers.Headers[lowerCaseTraceParentHeader]; ok {
			return propagation.HeaderCarrier{traceParentHeader: []string{val}}
		}
	}

	// Try EventBridge context propagation
	var payload eventBridgeEvent
	err = json.Unmarshal(eventJSON, &payload)
	if err == nil {
		if payload.Detail != nil {
			if val, ok := payload.Detail[traceParentHeader]; ok {
				return propagation.HeaderCarrier{traceParentHeader: []string{val}}
			} else if val, ok = payload.Detail[lowerCaseTraceParentHeader]; ok {
				return propagation.HeaderCarrier{traceParentHeader: []string{val}}
			}
		}

		if payload.TraceHeader != "" {
			return propagation.HeaderCarrier{traceParentHeader: []string{payload.TraceHeader}}
		}
	}

	return propagation.HeaderCarrier{"": []string{""}}
}

func instrumentHandler(handler interface{}) interface{} {
	provider := otel.GetTracerProvider()

	options := []otellambda.Option{}
	castProvider, success := provider.(*trace.TracerProvider)
	if success {
		options = append(options, otellambda.WithFlusher(castProvider),
			otellambda.WithEventToCarrier(heliosEventToCarrier),
			otellambda.WithPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})))
	}
	return otellambda.InstrumentHandler(handler, options...)
}

func Start(handler interface{}) {
	lambda.Start(instrumentHandler(handler))
}

func StartWithOptions(handler interface{}, options ...Option) {
	lambda.StartWithOptions(instrumentHandler(handler), options...)
}

func StartWithContext(ctx context.Context, handler interface{}) {
	lambda.StartWithContext(ctx, instrumentHandler(handler))
}

func StartHandlerWithContext(ctx context.Context, handler Handler) {
	lambda.StartWithOptions(handler, WithContext(ctx))
}

func StartHandler(handler Handler) {
	lambda.StartHandler(handler)
}

type Handler = lambda.Handler
type Function = lambda.Function
type Option = lambda.Option

func WithContext(ctx context.Context) Option {
	return lambda.WithContext(ctx)
}

func WithSetEscapeHTML(escapeHTML bool) Option {
	return lambda.WithSetEscapeHTML(escapeHTML)
}

func WithSetIndent(prefix, indent string) Option {
	return lambda.WithSetIndent(prefix, indent)
}

func WithEnableSIGTERM(callbacks ...func()) Option {
	return lambda.WithEnableSIGTERM(callbacks...)
}

func NewHandler(handlerFunc interface{}) Handler {
	return lambda.NewHandler(handlerFunc)
}

func NewHandlerWithOptions(handlerFunc interface{}, options ...Option) Handler {
	return lambda.NewHandlerWithOptions(handlerFunc, options...)
}

func NewFunction(handler Handler) *Function {
	return lambda.NewFunction(handler)
}
