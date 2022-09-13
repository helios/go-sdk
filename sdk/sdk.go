package sdk

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const sdkName = "helios-opentelemetry-sdk"
const customSpanAttr = "hs-custom-span"

var providerSingelton *trace.TracerProvider

func WithSamplingRatio(samplingRatio float64) attribute.KeyValue {
	return attribute.KeyValue{
		Key:   samplingRatioKey,
		Value: attribute.Float64Value(samplingRatio),
	}
}

func WithEnvironment(environment string) attribute.KeyValue {
	return attribute.KeyValue{
		Key:   environmentKey,
		Value: attribute.StringValue(environment),
	}
}

func WithCollectorEndpoint(collectorEndpoint string) attribute.KeyValue {
	return attribute.KeyValue{
		Key:   collectorEndpointKey,
		Value: attribute.StringValue(collectorEndpoint),
	}
}

func WithCollectorPath(collectorPath string) attribute.KeyValue {
	return attribute.KeyValue{
		Key:   collectorPathKey,
		Value: attribute.StringValue(collectorPath),
	}
}

func WithCommitHash(commitHash string) attribute.KeyValue {
	return attribute.KeyValue{
		Key:   commitHashKey,
		Value: attribute.StringValue(commitHash),
	}
}

func CreateCustomSpan(context context.Context, spanName string, attributes []attribute.KeyValue, callback func()) context.Context {
	if providerSingelton == nil {
		log.Print("Can't create custom span before Initialize is called")
		return nil
	}

	tracer := providerSingelton.Tracer("helios")
	ctx, span := tracer.Start(context, spanName)
	customSpanAttr := attribute.KeyValue{
		Key:   customSpanAttr,
		Value: attribute.StringValue("true"),
	}
	attributes = append(attributes, customSpanAttr)
	span.SetAttributes(attributes...)
	if callback != nil {
		callback()
	}

	defer span.End()
	return ctx
}

func Initialize(serviceName string, apiToken string, attrs ...attribute.KeyValue) (*trace.TracerProvider, error) {
	if providerSingelton != nil {
		return providerSingelton, nil
	}

	heliosConfig := getHeliosConfig(serviceName, apiToken, attrs...)
	var exporter *otlptrace.Exporter
	if heliosConfig.collectorEndpoint != "" {
		endpoint := otlptracehttp.WithEndpoint(heliosConfig.collectorEndpoint)
		urlPath := otlptracehttp.WithURLPath(heliosConfig.collectorPath)
		headers := otlptracehttp.WithHeaders(map[string]string{"Authorization": heliosConfig.apiToken})
		var error error
		exporter, error = otlptrace.New(context.Background(), otlptracehttp.NewClient(endpoint, headers, urlPath))
		if error != nil {
			return nil, error
		}
	}

	serviceAttributes := []attribute.KeyValue{semconv.ServiceNameKey.String(serviceName), semconv.TelemetrySDKVersionKey.String(version), semconv.TelemetrySDKNameKey.String(sdkName), semconv.TelemetrySDKLanguageGo}
	if heliosConfig.environment != "" {
		serviceAttributes = append(serviceAttributes, semconv.DeploymentEnvironmentKey.String(heliosConfig.environment))
	}
	if heliosConfig.commitHash != "" {
		serviceAttributes = append(serviceAttributes, semconv.ServiceVersionKey.String(heliosConfig.commitHash))
	}

	serviceResource := resource.NewWithAttributes(semconv.SchemaURL, serviceAttributes...)
	providerParams := []trace.TracerProviderOption{
		trace.WithResource(serviceResource),
		trace.WithSampler(heliosConfig.sampler),
	}
	if exporter != nil {
		providerParams = append(providerParams, trace.WithBatcher(exporter))
	}

	tracerProvider := trace.NewTracerProvider(providerParams...)

	otel.SetTracerProvider(tracerProvider)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	log.Printf("Helios tracing initialized (service: %s, token: %s*****, environment: %s)", serviceName, heliosConfig.apiToken[0:3], heliosConfig.environment)

	// Set singleton
	providerSingelton = tracerProvider
	return tracerProvider, nil
}
