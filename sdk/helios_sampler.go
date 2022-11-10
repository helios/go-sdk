package sdk

import (
	"fmt"

	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type HeliosSampler struct {
	sampler       trace.Sampler
	samplingRatio float64
}

func NewHeliosSampler(samplingRatio float64) HeliosSampler {
	var sampler trace.Sampler
	switch {
	case samplingRatio <= 0:
		samplingRatio = 0
		sampler = trace.NeverSample()
	case samplingRatio >= 1:
		samplingRatio = 1
		sampler = trace.AlwaysSample()
	default:
		sampler = trace.TraceIDRatioBased(samplingRatio)
	}

	return HeliosSampler{
		sampler:       sampler,
		samplingRatio: samplingRatio,
	}
}

func (hs HeliosSampler) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	parentSpanContext := oteltrace.SpanContextFromContext(parameters.ParentContext)
	if parentSpanContext.IsSampled() {
		// Make the sampling consistent
		return trace.SamplingResult{
			Decision:   trace.RecordAndSample,
			Tracestate: parentSpanContext.TraceState(),
		}
	}

	return hs.sampler.ShouldSample(parameters)
}

func (hs HeliosSampler) Description() string {
	return fmt.Sprintf("HeliosSampler(%.4f)", hs.samplingRatio)
}
