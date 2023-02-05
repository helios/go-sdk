package helioslambda

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/helios/opentelemetry-go-contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

var InstrumentedSymbols = [...]string{"Start", "StartWithContext", "StartWithOptions"}
var propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

const otellambdaTracerName = "go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
const traceParentHeader = "Traceparent"
const lowerCaseTraceParentHeader = "traceparent"

type eventBridgeEvent struct {
	Detail      map[string]string `json:"detail"`
	TraceHeader string            `json:"trace-header"`
}

type sqsMessageCarrier struct {
	messageAttrs map[string]events.SQSMessageAttribute
}

func (c sqsMessageCarrier) Get(key string) string {
	if c.messageAttrs == nil {
		return ""
	}

	for attrKey, val := range c.messageAttrs {
		if attrKey == key {
			return *val.StringValue
		}
	}

	return ""
}

func (c sqsMessageCarrier) Set(key, val string) {
	dataType := "String"
	c.messageAttrs[key] = events.SQSMessageAttribute{DataType: dataType, StringValue: &val}
}

func (c sqsMessageCarrier) Keys() []string {
	result := []string{}
	for key := range c.messageAttrs {
		result = append(result, key)
	}

	return result
}

func extractSpanNameFromQueueName(eventSourceArn string) string {
	if eventSourceArn != "" {
		queueNameIndex := strings.LastIndex(eventSourceArn, ":")
		if queueNameIndex != -1 {
			return fmt.Sprintf("process %s", eventSourceArn[queueNameIndex+1:])
		}
	}
	return "process SQS"
}

func extractContextFromSqsMessage(ctx context.Context, attributes map[string]events.SQSMessageAttribute) context.Context {
	return propagator.Extract(ctx, sqsMessageCarrier{attributes})
}

func extractContextFromEventBridgeSqsMessage(ctx context.Context, record events.SQSMessage) context.Context {
	var messageBody map[string]interface{}
	err := json.Unmarshal([]byte(record.Body), &messageBody)
	if err != nil {
		return ctx
	}
	if val, ok := messageBody["detail"]; ok {
		var detail map[string]interface{}
		detail = val.(map[string]interface{})
		if traceparentVal, ok := detail[lowerCaseTraceParentHeader]; ok {
			return propagator.Extract(ctx, propagation.MapCarrier{lowerCaseTraceParentHeader: traceparentVal.(string)})
		}
	}
	return ctx
}

func HandleRecord(ctx context.Context, record events.SQSMessage, handleRecordHelper func(ctx context.Context, message events.SQSMessage) (any, error)) (any, error) {
	var recordCtx context.Context
	if _, ok := record.MessageAttributes[lowerCaseTraceParentHeader]; ok {
		recordCtx = extractContextFromSqsMessage(ctx, record.MessageAttributes)
	} else {
		recordCtx = extractContextFromEventBridgeSqsMessage(ctx, record)
	}
	tp := otel.GetTracerProvider()
	messageId := attribute.KeyValue{
		Key:   semconv.MessageIDKey,
		Value: attribute.StringValue(record.MessageId),
	}
	messagingSystem := attribute.KeyValue{
		Key:   semconv.MessagingSystemKey,
		Value: attribute.StringValue("aws.sqs"),
	}
	messagingPayload := attribute.KeyValue{
		Key:   "faas.event",
		Value: attribute.StringValue(record.Body),
	}
	spanName := extractSpanNameFromQueueName(record.EventSource)
	updatedCtx, span := tp.Tracer(otellambdaTracerName).Start(recordCtx, spanName, trace.WithAttributes(messageId, messagingSystem, messagingPayload))
	defer span.End()
	return handleRecordHelper(updatedCtx, record)
}

func heliosEventToCarrier(eventJSON []byte) propagation.TextMapCarrier {
	// Try API Gateway context propagation
	var apiGatewayEvent events.APIGatewayV2HTTPRequest
	err := json.Unmarshal(eventJSON, &apiGatewayEvent)
	if err == nil && apiGatewayEvent.Headers != nil {
		if val, ok := apiGatewayEvent.Headers[traceParentHeader]; ok {
			return propagation.HeaderCarrier{traceParentHeader: []string{val}}
		} else if val, ok = apiGatewayEvent.Headers[lowerCaseTraceParentHeader]; ok {
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
	castProvider, success := provider.(*sdkTrace.TracerProvider)
	if success {
		options = append(options, otellambda.WithFlusher(castProvider),
			otellambda.WithEventToCarrier(heliosEventToCarrier),
			otellambda.WithPropagator(propagator))
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
