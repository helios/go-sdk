package heliossqlx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	_ "modernc.org/sqlite"
)

func getSpanRecorder() *tracetest.SpanRecorder {
	spanRecorder := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return spanRecorder
}

func assertSpan(t *testing.T, spanRecorder *tracetest.SpanRecorder) []attribute.KeyValue {
	spanRecorder.ForceFlush(context.Background())
	spans := spanRecorder.Ended()
	assert.Equal(t, 2, len(spans))
	spanDbConnect := spans[0]
	assert.Equal(t, "db.Connect", spanDbConnect.Name())
	assert.Equal(t, trace.SpanKindClient, spanDbConnect.SpanKind())

	spanDbQuery := spans[1]
	assert.Equal(t, "db.Query", spanDbQuery.Name())
	assert.Equal(t, trace.SpanKindClient, spanDbQuery.SpanKind())
	return spanDbQuery.Attributes()
}

func assertAttributes(t *testing.T, attributes []attribute.KeyValue) {
	for _, attribute := range attributes {
		key := attribute.Key
		value := attribute.Value.AsString()

		switch key {
		case semconv.DBNameKey:
			assert.Equal(t, "test", value)
		case semconv.DBSystemKey:
			assert.Equal(t, semconv.DBSystemSqlite.Value.AsString(), value)
		case semconv.DBStatementKey:
			assert.Equal(t, "SELECT 42", value)	
		}
	}
}

func TestSqlxInstrumentation(t *testing.T) {
	spanRecorder := getSpanRecorder()
	db, err := Open("sqlite", "file::memory:?cache=shared",
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithDBName("test"))
	if err != nil {
		panic(err)
	}
	
	var num int
	if err := db.QueryRowContext(context.Background(), "SELECT 42").Scan(&num); err != nil {
		panic(err)
	}
	attributes := assertSpan(t, spanRecorder)
	assertAttributes(t, attributes)
}