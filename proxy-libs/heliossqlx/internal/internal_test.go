package internal

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	// "github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/helios/go-sdk/proxy-libs/heliossqlx"
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

func assertAttributes(t *testing.T, attributes []attribute.KeyValue) bool {
	dbNameExists := false
	dbSystemExists := false
	dbStatementExists := false
	for _, attribute := range attributes {
		key := attribute.Key
		value := attribute.Value.AsString()

		switch key {
		case semconv.DBNameKey:
			dbNameExists = true
			assert.Equal(t, "test", value)
		case semconv.DBSystemKey:
			dbSystemExists = true
			assert.Equal(t, semconv.DBSystemSqlite.Value.AsString(), value)
		case semconv.DBStatementKey:
			dbStatementExists = true
			assert.Equal(t, "SELECT 42", value)
		}
	}
	return dbNameExists && dbSystemExists && dbStatementExists
}

func getDBConnection() *heliossqlx.DB {
	db, err := heliossqlx.Open("sqlite", "file::memory:?cache=shared",
		otelsql.WithDBName("test"))
	if err != nil {
		panic(err)
	}

	return db
}

func TestSqlxInstrumentation(t *testing.T) {
	spanRecorder := getSpanRecorder()
	db := getDBConnection()

	var num int
	if err := db.QueryRowContext(context.Background(), "SELECT 42").Scan(&num); err != nil {
		panic(err)
	}
	attributes := assertSpan(t, spanRecorder)
	assert.True(t, assertAttributes(t, attributes))
}

// func TestDisableInstrumentation(t *testing.T) {
// 	os.Setenv("HS_DISABLED", "true")
// 	defer os.Setenv("HS_DISABLED", "")
// 	spanRecorder := getSpanRecorder()
// 	db := getDBConnection()

// 	var num int
// 	if err := db.QueryRowContext(context.Background(), "SELECT 42").Scan(&num); err != nil {
// 		panic(err)
// 	}
// 	spanRecorder.ForceFlush(context.Background())
// 	spans := spanRecorder.Ended()
// 	assert.Equal(t, 0, len(spans))
// }
