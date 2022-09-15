package sdk

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/sdk/trace"
)

const HELIOS_TEST_TRIGGERED_TRACE = "hs-triggered-test"

type HeliosProcessor struct{}

func (heliosProcessor HeliosProcessor) OnStart(ctx context.Context, s trace.ReadWriteSpan) {
	spanBaggage := baggage.FromContext(ctx)
	testMember := spanBaggage.Member(HELIOS_TEST_TRIGGERED_TRACE)
	if testMember.Key() == HELIOS_TEST_TRIGGERED_TRACE {
		s.SetAttributes(attribute.String(HELIOS_TEST_TRIGGERED_TRACE, "true"))
	}
}
func (heliosProcessor HeliosProcessor) Shutdown(context.Context) error   { return nil }
func (heliosProcessor HeliosProcessor) ForceFlush(context.Context) error { return nil }
func (heliosProcessor HeliosProcessor) OnEnd(s trace.ReadOnlySpan)       {}
