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

type HttpHeaders struct {
	Headers map[string]string `json:"headers"`
}
func heliosEventToCarrier (eventJSON []byte) propagation.TextMapCarrier {
	var headers HttpHeaders
	err := json.Unmarshal(eventJSON, &headers)
	if err != nil {
		return propagation.HeaderCarrier{"": []string{""}}
	} else {
			if val, ok := headers.Headers["Traceparent"]; ok  {
				return propagation.HeaderCarrier{"Traceparent": []string{val}}
			} else if val, ok = headers.Headers["traceparent"]; ok {
				return propagation.HeaderCarrier{"traceparent": []string{val}}
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
