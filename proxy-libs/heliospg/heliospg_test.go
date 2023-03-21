package heliospg

import (
	"context"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type User struct {
	Id     int64
	Name   string
	Emails []string
}

func createSchema(db *DB) error {
	models := []interface{}{
		(*User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func getSpanRecorder() (*tracetest.SpanRecorder, *sdktrace.TracerProvider) {
	spanRecorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
	otel.SetTracerProvider(provider)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return spanRecorder, provider
}

func ConnectToDb() *DB {
	db := Connect(&Options{
		User:     "postgres",
		Password: "postgres"})
	err := createSchema(db)
	if err != nil {
		panic(err)
	}
	return db
}

func insertUser(db *DB) {
	user1 := &User{
		Name:   "admin",
		Emails: []string{"admin1@admin", "admin2@admin"},
	}
	_, err := db.Model(user1).Insert()
	if err != nil {
		panic(err)
	}
}

func assertSpan(t *testing.T, spanRecorder *tracetest.SpanRecorder) []attribute.KeyValue {
	spanRecorder.ForceFlush(context.Background())
	spans := spanRecorder.Ended()
	assert.Equal(t, 1, len(spans))
	span := spans[0]
	assert.Equal(t, "INSERT", span.Name())
	assert.Equal(t, trace.SpanKindInternal, span.SpanKind())
	return span.Attributes()
}

func assertAttributes(t *testing.T, attributes []attribute.KeyValue) {
	assert.Contains(t, attributes, attribute.String("db.name", "postgres"))
	assert.Contains(t, attributes, attribute.String("db.system", "postgresql"))
	assert.Contains(t, attributes, attribute.String("db.statement", "INSERT INTO \"users\" (\"id\", \"name\", \"emails\") VALUES (?, ?, ?)"))
}

func TestConnectInstrumentation(t *testing.T) {
	spanRecorder, provider := getSpanRecorder()
	tracer := provider.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "custom")
	defer span.End()
	db := ConnectToDb()

	insertUser(db.WithContext(ctx))
	attributes := assertSpan(t, spanRecorder)
	assertAttributes(t, attributes)
	defer db.Close()
}
