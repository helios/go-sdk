package sdk

import (
	"context"
	"fmt"
	"reflect"

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

func (heliosProcessor HeliosProcessor) Shutdown(context.Context) error { return nil }

func (heliosProcessor HeliosProcessor) ForceFlush(context.Context) error { return nil }

func getSpanAttributeByName(span trace.ReadOnlySpan, attrName string) *reflect.Value {
	reflectionValue := reflect.ValueOf(span)
	if reflectionValue.Kind() == reflect.Ptr {
		reflectionValue = reflectionValue.Elem()
	}

	attrs := reflectionValue.FieldByName("attributes")
	for i := 0; i < attrs.Len(); i++ {
		attr := attrs.Index(i)
		key := attr.FieldByName("Key").String()
		if key == attrName {
			return &attr
		}
	}

	return nil
}

func (heliosProcessor HeliosProcessor) OnEnd(s trace.ReadOnlySpan) {
	name := s.Name()
	if name == "sampled1" {
		value := getSpanAttributeByName(s, "key1")
		if value != nil {
			bla := value.FieldByName("Value").FieldByName("Value").String()
			fmt.Printf("bla: %v\n", bla)
			value.FieldByName("Value").SetString("value2")
		}
	}
}
