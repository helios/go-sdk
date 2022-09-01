package sdk

import (
	"context"

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

func Initialize(serviceName string, apiToken string, attrs ...attribute.KeyValue) (*trace.TracerProvider, error) {
	if providerSingelton != nil {
		return providerSingelton, nil
	}

	heliosConfig := getHeliosConfig(serviceName, apiToken, attrs...)
	endpoint := otlptracehttp.WithEndpoint(heliosConfig.collectorEndpoint)
	urlPath := otlptracehttp.WithURLPath(heliosConfig.collectorPath)
	headers := otlptracehttp.WithHeaders(map[string]string{"Authorization": heliosConfig.apiToken})
	exporter, error := otlptrace.New(context.Background(), otlptracehttp.NewClient(endpoint, headers, urlPath))

	if error != nil {
		return nil, error
	}

	serviceAttributes := []attribute.KeyValue{semconv.ServiceNameKey.String(serviceName), semconv.TelemetrySDKVersionKey.String(version), semconv.TelemetrySDKNameKey.String(sdkName), semconv.TelemetrySDKLanguageGo}
	if heliosConfig.environment != "" {
		serviceAttributes = append(serviceAttributes, semconv.DeploymentEnvironmentKey.String(heliosConfig.environment))
	}
	if heliosConfig.commitHash != "" {
		serviceAttributes = append(serviceAttributes, semconv.ServiceVersionKey.String(heliosConfig.commitHash))
	}

	serviceResource := resource.NewWithAttributes(semconv.SchemaURL, serviceAttributes...)

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(serviceResource),
		trace.WithSampler(heliosConfig.sampler),
	)

	otel.SetTracerProvider(tracerProvider)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	// Set singleton
	providerSingelton = tracerProvider
	return tracerProvider, nil
}