package sdk

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-logr/stdr"
	ec2Detector "go.opentelemetry.io/contrib/detectors/aws/ec2"
	ecsDetector "go.opentelemetry.io/contrib/detectors/aws/ecs"
	eksDetector "go.opentelemetry.io/contrib/detectors/aws/eks"
	lambdaDetector "go.opentelemetry.io/contrib/detectors/aws/lambda"
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

func WithCollectorInsecure() attribute.KeyValue {
	return attribute.KeyValue{
		Key:   collectorInsecureKey,
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

func WithCommitHash(commitHash string) attribute.KeyValue {
	return attribute.KeyValue{
		Key:   commitHashKey,
		Value: attribute.StringValue(commitHash),
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

func getServerResource(serviceName string, heliosConfig HeliosConfig) *resource.Resource {
	serviceAttributes := []attribute.KeyValue{semconv.ServiceNameKey.String(serviceName), semconv.TelemetrySDKVersionKey.String(version), semconv.TelemetrySDKNameKey.String(sdkName), semconv.TelemetrySDKLanguageGo}
	if heliosConfig.environment != "" {
		serviceAttributes = append(serviceAttributes, semconv.DeploymentEnvironmentKey.String(heliosConfig.environment))
	}
	if heliosConfig.commitHash != "" {
		serviceAttributes = append(serviceAttributes, semconv.ServiceVersionKey.String(heliosConfig.commitHash))
	}

	return resource.NewWithAttributes(semconv.SchemaURL, serviceAttributes...)
}

func getResourceOptions(ctx context.Context, serviceName string, heliosConfig HeliosConfig) []trace.TracerProviderOption {
	resources := []trace.TracerProviderOption{trace.WithResource(getServerResource(serviceName, heliosConfig))}

	lambdaDetector := lambdaDetector.NewResourceDetector()
	res, err := lambdaDetector.Detect(ctx)
	if err == nil {
		resources = append(resources, trace.WithResource(res))
	}

	eksDetector := eksDetector.NewResourceDetector()
	res, err = eksDetector.Detect(ctx)
	if err == nil {
		resources = append(resources, trace.WithResource(res))
	}

	ecsDetector := ecsDetector.NewResourceDetector()
	res, err = ecsDetector.Detect(ctx)
	if err == nil {
		resources = append(resources, trace.WithResource(res))
	}

	ec2Detector := ec2Detector.NewResourceDetector()
	res, err = ec2Detector.Detect(ctx)
	if err == nil {
		resources = append(resources, trace.WithResource(res))
	}

	return resources
}

func Initialize(serviceName string, apiToken string, attrs ...attribute.KeyValue) (*trace.TracerProvider, error) {
	if providerSingelton != nil {
		return providerSingelton, nil
	}

	ctx := context.Background()
	heliosConfig := createHeliosConfig(serviceName, apiToken, attrs...)
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

	resourceOptions := getResourceOptions(ctx, serviceName, *heliosConfig)
	providerParams := []trace.TracerProviderOption{
		trace.WithSampler(heliosConfig.sampler),
	}
	providerParams = append(providerParams, resourceOptions...)

	if exporter != nil {
		providerParams = append(providerParams, trace.WithBatcher(exporter))
	}

	tracerProvider := trace.NewTracerProvider(providerParams...)
	heliosProcessor := HeliosProcessor{}

	tracerProvider.RegisterSpanProcessor(heliosProcessor)
	otel.SetTracerProvider(tracerProvider)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	_, debugSpan := tracerProvider.Tracer("helios").Start(ctx, "debug span")
	defer debugSpan.End()
	log.Printf("Helios tracing initialized (service: %s, token: %s*****, environment: %s)", serviceName, heliosConfig.apiToken[0:3], heliosConfig.environment)

	// Set singleton
	providerSingelton = tracerProvider
	return tracerProvider, nil
}
