package heliospg

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type User struct {
    Id     int64
    Name   string
    Emails []string
}

// createSchema creates database schema for User and Story models.
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

func getSpanRecorder() *tracetest.SpanRecorder {
	spanRecorder := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return spanRecorder
}

func ConnectToDb() *DB {
	db := Connect(&Options{
        User: "postgres",
    })
	defer db.Close()
	err := createSchema(db)
    if err != nil {
        panic(err)
    }
	return db
}

func insertUser(db *DB, id int, name string, role string) {
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
	// assert.Equal(t, "users.insert", span.Name())
	assert.Equal(t, trace.SpanKindClient, span.SpanKind())
	return span.Attributes()
}

func assertAttributes(t *testing.T, attributes []attribute.KeyValue, id int, name string, role string) {
	for _, attribute := range attributes {
		key := attribute.Key
		value := attribute.Value.AsString()

		switch key {
		case semconv.DBNameKey:
			assert.Contains(t, "test", value)
		case semconv.DBOperationKey:
			assert.Contains(t, "insert", value)
		case semconv.DBStatementKey:
			assert.Contains(t, value, fmt.Sprintf("\"id\":%v", id))
			assert.Contains(t, value, fmt.Sprintf("\"name\":\"%v\"", name))
			assert.Contains(t, value, fmt.Sprintf("\"role\":\"%v\"", role))
		case semconv.DBSystemKey:
			assert.Equal(t, "mongodb", value)
		}
	}
}

func TestConnectInstrumentation(t *testing.T) {
	spanRecorder := getSpanRecorder()
	db := ConnectToDb()

	insertUser(db, 12345, "Lior Govrin", "Software Engineer")
	assertSpan(t, spanRecorder)
	// assertAttributes(t, attributes, 12345, "Lior Govrin", "Software Engineer")
}
