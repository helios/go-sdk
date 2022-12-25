package helioschi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	exportsExtractor "github.com/helios/go-instrumentor/exports_extractor"
	"github.com/stretchr/testify/assert"
)

const requestBody = "{\"id\":123,\"name\":\"Lior Govrin\",\"role\":\"Software Engineer\""

func validateAttributes(attrs []attribute.KeyValue, t *testing.T) {
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "POST", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/users/abcd1234", value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		} else if key == semconv.HTTPRouteKey {
			assert.Equal(t, "/users/{id}", value.Value.AsString())
		} else if key == "http.request.body" {
			assert.Equal(t, requestBody, value.Value.AsString())
		}
	}
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func TestInstrumentation(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	r := NewRouter()

	r.HandleFunc("/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var user User
		decoder.Decode(&user)

		id := URLParam(r, "id")
		name := "test"
		reply := fmt.Sprintf("user %s (id %s)", name, id)
		w.Write(([]byte)(reply))
	}))

	go func() {
		http.ListenAndServe(":3333", r)
	}()

	res, _ := http.Post("http://localhost:3333/users/abcd1234", "application/json", bytes.NewBuffer([]byte(requestBody)))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "user test (id abcd1234)", string(body))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 1, len(spans))
	serverSpan := spans[0]
	assert.Equal(t, "/users/{id}", serverSpan.Name())
	assert.Equal(t, serverSpan.SpanKind(), trace.SpanKindServer)
	validateAttributes(serverSpan.Attributes(), t)
}

func TestInterfaceMatch(t *testing.T) {
	originalRepository := exportsExtractor.CloneGitRepository("https://github.com/go-chi/chi", "v5.0.8")
	originalExports := exportsExtractor.ExtractExports(originalRepository, "chi")
	os.RemoveAll(originalRepository)
	sort.Slice(originalExports, func(i, j int) bool {
		return originalExports[i].Name < originalExports[j].Name
	})

	srcDir, _ := filepath.Abs(".")
	heliosExports := exportsExtractor.ExtractExports(srcDir, "helioschi")
	sort.Slice(heliosExports, func(i, j int) bool {
		return heliosExports[i].Name < heliosExports[j].Name
	})

	assert.EqualValues(t, originalExports, heliosExports)
}
