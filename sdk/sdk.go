package sdk

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const sdkName = "helios-opentelemetry-sdk"
const collectorEndpoint = "collector.heliosphere.io:443"
const collectorPath = "traces"

func Initialize(serviceName string, apiToken string) (*trace.TracerProvider, error) {
	endpoint := otlptracehttp.WithEndpoint(collectorEndpoint)
	headers := otlptracehttp.WithHeaders(map[string]string{"Authorization": apiToken})
	urlPath := otlptracehttp.WithURLPath(collectorPath)
	exporter, error := otlptrace.New(context.Background(), otlptracehttp.NewClient(endpoint, headers, urlPath))

	if error != nil {
		return nil, error
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName), semconv.TelemetrySDKVersionKey.String(version), semconv.TelemetrySDKNameKey.String(sdkName))),
		trace.WithSampler(trace.AlwaysSample()),
	)

	otel.SetTracerProvider(tracerProvider)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return tracerProvider, nil
}
