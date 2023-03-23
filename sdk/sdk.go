package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/go-logr/stdr"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const sdkName = "helios-opentelemetry-sdk"
const customSpanAttr = "hs-custom-span"

var providerSingleton *trace.TracerProvider

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

func WithCollectorInsecure() attribute.KeyValue {
	return attribute.KeyValue{
		Key:   collectorInsecureKey,
		Value: attribute.StringValue("true"),
	}
}

func WithCollectMetrics() attribute.KeyValue {
	return attribute.KeyValue{
		Key:   collectMetricsKey,
		Value: attribute.StringValue("true"),
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

func WithCollectorMetricsPath(collectorMetricsPath string) attribute.KeyValue {
	return attribute.KeyValue{
		Key:   collectorMetricsPathKey,
		Value: attribute.StringValue(collectorMetricsPath),
	}
}

func WithCommitHash(commitHash string) attribute.KeyValue {
	return attribute.KeyValue{
		Key:   commitHashKey,
		Value: attribute.StringValue(commitHash),
	}
}

func WithInstrumentationDisabled() attribute.KeyValue {
	return attribute.KeyValue{
		Key:   instrumentationDisabledKey,
		Value: attribute.StringValue("true"),
	}
}

func WithDebugMode() attribute.KeyValue {
	return attribute.KeyValue{
		Key:   debugKey,
		Value: attribute.StringValue("true"),
	}
}

func WithMetadataOnlyMode() attribute.KeyValue {
	return attribute.KeyValue{
		Key:   metadataOnlyKey,
		Value: attribute.StringValue("true"),
	}
}

func WithObfuscationBlocklistRules(blocklistRules []string) attribute.KeyValue {
	blocklistRulesAsJson, _ := json.Marshal(blocklistRules)
	os.Setenv(hsDataObfuscationBlocklistEnvVar, string(blocklistRulesAsJson))
	return attribute.KeyValue{
		Key:   hsDataObfuscationBlocklistKey,
		Value: attribute.StringSliceValue(blocklistRules),
	}
}

func WithObfuscationAllowlistRules(allowlistRules []string) attribute.KeyValue {
	allowlistRulesAsJson, _ := json.Marshal(allowlistRules)
	os.Setenv(hsDataObfuscationBlocklistEnvVar, string(allowlistRulesAsJson))
	return attribute.KeyValue{
		Key:   hsDataObfuscationAllowlistKey,
		Value: attribute.StringSliceValue(allowlistRules),
	}
}

func WithhmacKey(hMacKey string) attribute.KeyValue {
	os.Setenv(hsDatahMacKeyEnvVar, hMacKey)
	return attribute.KeyValue{
		Key:   hsDatahMacKey,
		Value: attribute.StringValue(hMacKey),
	}
}

func CreateCustomSpan(context context.Context, spanName string, attributes []attribute.KeyValue, callback func()) context.Context {
	if providerSingleton == nil {
		log.Print("Can't create custom span before Initialize is called")
		return nil
	}

	tracer := providerSingleton.Tracer("helios")
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

func getResource(serviceName string, heliosConfig *HeliosConfig) *resource.Resource {
	serviceAttributes := []attribute.KeyValue{semconv.ServiceNameKey.String(serviceName), semconv.TelemetrySDKVersionKey.String(version), semconv.TelemetrySDKNameKey.String(sdkName), semconv.TelemetrySDKLanguageGo}
	if heliosConfig.environment != "" {
		serviceAttributes = append(serviceAttributes, semconv.DeploymentEnvironmentKey.String(heliosConfig.environment))
	}
	if heliosConfig.commitHash != "" {
		serviceAttributes = append(serviceAttributes, semconv.ServiceVersionKey.String(heliosConfig.commitHash))
	}

	return resource.NewWithAttributes(semconv.SchemaURL, serviceAttributes...)
}

func collectMetrics(ctx context.Context, resource *resource.Resource, heliosConfig *HeliosConfig) {
	options := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(heliosConfig.collectorEndpoint),
		otlpmetrichttp.WithURLPath(heliosConfig.collectorPath),
		otlpmetrichttp.WithHeaders(map[string]string{"Authorization": heliosConfig.apiToken}),
	}
	if heliosConfig.collectorInsecure {
		options = append(options, otlpmetrichttp.WithInsecure())
	}

	exporter, err := otlpmetrichttp.New(ctx, options...)
	if err != nil {
		log.Print("Unable to initialize metrics exporter")
		return
	}

	reader := metric.NewPeriodicReader(exporter)
	provider := metric.NewMeterProvider(metric.WithResource(resource), metric.WithReader(reader))
	err = host.Start(host.WithMeterProvider(provider))
	if err != nil {
		log.Print("Unable to start host metrics collection")
	}
}

func Initialize(serviceName string, apiToken string, attrs ...attribute.KeyValue) (*trace.TracerProvider, error) {
	if providerSingleton != nil {
		return providerSingleton, nil
	}

	ctx := context.Background()
	heliosConfig := createHeliosConfig(serviceName, apiToken, attrs...)

	if heliosConfig.instrumentationDisabled {
		os.Setenv("HS_DISABLED", "true")
		return nil, errors.New("helios tracing is disabled")
	}

	var exporter *otlptrace.Exporter
	if heliosConfig.collectorEndpoint != "" {
		options := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(heliosConfig.collectorEndpoint),
			otlptracehttp.WithURLPath(heliosConfig.collectorPath),
			otlptracehttp.WithHeaders(map[string]string{"Authorization": heliosConfig.apiToken}),
		}
		if heliosConfig.collectorInsecure {
			options = append(options, otlptracehttp.WithInsecure())
		}
		var error error
		exporter, error = otlptrace.New(ctx, otlptracehttp.NewClient(options...))
		if error != nil {
			return nil, error
		}
	}

	if heliosConfig.debug {
		stdoutLogger := stdr.New(log.New(os.Stdout, "[opentelemetry-logger] ", 3))
		stdr.SetVerbosity(3)
		otel.SetLogger(stdoutLogger)
	}

	if heliosConfig.metadataOnly {
		os.Setenv("HS_METADATA_ONLY", "true")
	}

	serviceResource := getResource(serviceName, heliosConfig)
	providerParams := []trace.TracerProviderOption{
		trace.WithResource(serviceResource),
		trace.WithSampler(heliosConfig.sampler),
	}

	if exporter != nil {
		providerParams = append(providerParams, trace.WithBatcher(exporter))
	}

	tracerProvider := trace.NewTracerProvider(providerParams...)
	heliosProcessor := HeliosProcessor{}

	if heliosConfig.collectMetrics {
		collectMetrics(ctx, serviceResource, heliosConfig)
	}

	tracerProvider.RegisterSpanProcessor(heliosProcessor)
	otel.SetTracerProvider(tracerProvider)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	_, debugSpan := tracerProvider.Tracer("helios").Start(ctx, "debug span")
	defer debugSpan.End()
	log.Printf("Helios tracing initialized (service: %s, token: %s*****, environment: %s)", serviceName, heliosConfig.apiToken[0:3], heliosConfig.environment)

	// Set singleton
	providerSingleton = tracerProvider
	return tracerProvider, nil
}
