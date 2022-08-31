package sdk

import (
	"context"
	"os"
	"strconv"

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
const defaultCollectorEndpoint = "collector.heliosphere.io:443"
const defaultCollectorPath = "traces"
const environmentEnvVar = "HS_ENVIRONMENT"
const samplingRatioEnvVar = "HS_SAMPLING_RATIO"
const collectorEndpointEnvVar = "HS_COLLECTOR_ENDPOINT"
const collectorPathEnvVar = "HS_COLLECTOR_PATH"

var providerSingelton *trace.TracerProvider

func getSampler() trace.Sampler {
	samplingRatio := os.Getenv(samplingRatioEnvVar)
	if samplingRatio == "" {
		return trace.AlwaysSample()
	}

	res, err := strconv.ParseFloat(samplingRatio, 64)
	if err != nil {
		return trace.AlwaysSample()
	}

	return trace.TraceIDRatioBased(res)
}

func getCollectorEndpoint() string {
	collectorEndpoint := os.Getenv(collectorEndpointEnvVar)
	if collectorEndpoint == "" {
		return defaultCollectorEndpoint
	}

	return collectorEndpoint
}

func getCollectorPath() string {
	collectorPath := os.Getenv(collectorPathEnvVar)
	if collectorPath == "" {
		return defaultCollectorPath
	}

	return collectorPath
}

func Initialize(serviceName string, apiToken string) (*trace.TracerProvider, error) {
	if providerSingelton != nil {
		return providerSingelton, nil
	}

	collectorEndpoint := getCollectorEndpoint()
	collectorPath := getCollectorPath()
	endpoint := otlptracehttp.WithEndpoint(collectorEndpoint)
	urlPath := otlptracehttp.WithURLPath(collectorPath)
	headers := otlptracehttp.WithHeaders(map[string]string{"Authorization": apiToken})
	exporter, error := otlptrace.New(context.Background(), otlptracehttp.NewClient(endpoint, headers, urlPath))

	if error != nil {
		return nil, error
	}

	serviceAttributes := []attribute.KeyValue{semconv.ServiceNameKey.String(serviceName), semconv.TelemetrySDKVersionKey.String(version), semconv.TelemetrySDKNameKey.String(sdkName), semconv.TelemetrySDKLanguageGo}
	if os.Getenv(environmentEnvVar) != "" {
		serviceAttributes = append(serviceAttributes, semconv.DeploymentEnvironmentKey.String(os.Getenv(environmentEnvVar)))
	}

	serviceResource := resource.NewWithAttributes(semconv.SchemaURL, serviceAttributes...)
	sampler := getSampler()

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(serviceResource),
		trace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tracerProvider)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	// Set singleton
	providerSingelton = tracerProvider
	return tracerProvider, nil
}
