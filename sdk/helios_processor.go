package sdk

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/sdk/trace"
)

const HELIOS_TEST_TRIGGERED_TRACE = "hs-triggered-test"

type HeliosProcessor struct {
	metadataOnlyMode bool
}

func (heliosProcessor HeliosProcessor) OnStart(ctx context.Context, s trace.ReadWriteSpan) {
	spanBaggage := baggage.FromContext(ctx)
	testMember := spanBaggage.Member(HELIOS_TEST_TRIGGERED_TRACE)
	if testMember.Key() == HELIOS_TEST_TRIGGERED_TRACE {
		s.SetAttributes(attribute.String(HELIOS_TEST_TRIGGERED_TRACE, "true"))
	}
}
func (heliosProcessor HeliosProcessor) Shutdown(context.Context) error   { return nil }
func (heliosProcessor HeliosProcessor) ForceFlush(context.Context) error { return nil }
func (heliosProcessor HeliosProcessor) OnEnd(s trace.ReadOnlySpan) {
	if heliosProcessor.metadataOnlyMode {
		newAttrs := []attribute.KeyValue{}
		for _, attr := range s.Attributes() {
			key := attr.Key
			if key == "http.request.body" ||
				key == "http.request.headers" ||
				key == "http.response.body" ||
				key == "http.request.headers" {
				continue
			} else {
				newAttrs = append(newAttrs, attr)
			}
		}
		if len(newAttrs) != len(s.Attributes()) {
			// Can't manipulate ReadOnlySpan
			s.SetAttributes(newAttrs)
		}
	}
}
