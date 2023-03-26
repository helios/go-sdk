package helioshttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/helios/opentelemetry-go-contrib/instrumentation/net/http/otelhttp"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const expectedRequestBody = "{\"id\":123,\"name\":\"Lior Govrin\",\"role\":\"Software Engineer\"}"
const expectedObfuscatedRequestBody = "{\"id\":123,\"name\":\"dac02c19\",\"role\":\"Software Engineer\"}"
const expectedResponseBody = "hello1234"
const expectedCssResponseBody = "body {\n  font-size: 16px;\n  font-family: 'Lato', sans-serif;\n  background-size: cover;\n  background-color: black;\n}\n"
const expectedObfuscatedResponseBody = "87468e56"

func init() {
	blocklistRules, _ := json.Marshal([]string{"$.name"})
	os.Setenv("HS_DATA_OBFUSCATION_HMAC_KEY", "12345")
	os.Setenv("HS_DATA_OBFUSCATION_BLOCKLIST", string(blocklistRules))
}

func getHello(responseWriter ResponseWriter, request *Request) {
	body, _ := io.ReadAll(request.Body)
	if string(body) != expectedRequestBody {
		log.Fatal("Invalid request body")
	}
	io.WriteString(responseWriter, expectedResponseBody)
}

func getCss(responseWriter ResponseWriter, request *Request) {
	responseWriter.Header().Set("content-type", "text/css")
	responseWriter.Header().Add("foo", "bar1")
	io.WriteString(responseWriter, expectedCssResponseBody)
}

func validateAttributes(attrs []attribute.KeyValue, path string, metadataOnly bool, t *testing.T) {
	requestBodyFound := false
	requestHeadersFound := false
	responseBodyFound := false
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "POST", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/"+path, value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		} else if key == "http.response.body" {
			responseBodyFound = true
			assert.Equal(t, expectedObfuscatedResponseBody, value.Value.AsString())
		} else if key == "http.request.headers" {
			requestHeadersFound = true
			headers := map[string][]string{}
			json.Unmarshal([]byte(value.Value.AsString()), &headers)
			assert.Equal(t, "application/json", headers["Content-Type"][0])
		} else if key == "http.request.body" {
			requestBodyFound = true
			assert.Equal(t, expectedObfuscatedRequestBody, value.Value.AsString())
		}
	}

	assert.Equal(t, metadataOnly, !requestBodyFound)
	assert.Equal(t, metadataOnly, !requestHeadersFound)
	assert.Equal(t, metadataOnly, !responseBodyFound)
}

func validateStaticContentAttributes(attrs []attribute.KeyValue, metadataOnly bool, t *testing.T) {
	requestHeadersFound := false
	for _, value := range attrs {
		key := value.Key
		assert.NotEqualValues(t, key, "http.request.body")
		assert.NotEqualValues(t, key, "http.response.body")

		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "GET", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/style.css", value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		} else if key == "http.request.headers" {
			requestHeadersFound = true
			fmt.Println("requestHeadersFound = true", requestHeadersFound)
			headers := map[string][]string{}
			json.Unmarshal([]byte(value.Value.AsString()), &headers)
			if headers["Content-Type"] != nil {
				assert.Equal(t, "text/css", headers["Content-Type"][0])
			}
		}
	}

	assert.EqualValues(t, metadataOnly, !requestHeadersFound)
}

func staticContentTestHelper(t *testing.T, port int, metadataOnly bool) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	Handle("/style.css", HandlerFunc(getCss))
	go func() {
		ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	res, _ := Get(fmt.Sprintf("http://localhost:%d/style.css", port))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, expectedCssResponseBody, string(body))

	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 2, len(spans))
	
	serverSpan := spans[0]
	assert.Equal(t, trace.SpanKind(2), serverSpan.SpanKind())
	validateStaticContentAttributes(serverSpan.Attributes(), metadataOnly, t)
	
	clientSpan := spans[1]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.False(t, clientSpan.Parent().HasTraceID())
	validateStaticContentAttributes(clientSpan.Attributes(), metadataOnly, t)

	assert.Equal(t, serverSpan.Parent().SpanID(), clientSpan.SpanContext().SpanID())
	assert.Equal(t, res.Header.Get("traceresponse"), fmt.Sprintf("00-%s-%s-01", serverSpan.SpanContext().TraceID().String(), serverSpan.SpanContext().SpanID().String()))
}

func sendRequestAndValidate(t *testing.T, port int, path string, metadataOnly bool, sr *tracetest.SpanRecorder) {
	go func() {
		ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	res, _ := Post(fmt.Sprintf("http://localhost:%d/%s", port, path), "application/json", bytes.NewBuffer([]byte(expectedRequestBody)))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, expectedResponseBody, string(body))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 2, len(spans))
	serverSpan := spans[0]
	assert.Equal(t, trace.SpanKind(2), serverSpan.SpanKind())
	validateAttributes(serverSpan.Attributes(), path, metadataOnly, t)
	clientSpan := spans[1]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.False(t, clientSpan.Parent().HasTraceID())
	validateAttributes(clientSpan.Attributes(), path, metadataOnly, t)
	assert.Equal(t, serverSpan.Parent().SpanID(), clientSpan.SpanContext().SpanID())
	assert.Equal(t, res.Header.Get("traceresponse"), fmt.Sprintf("00-%s-%s-01", serverSpan.SpanContext().TraceID().String(), serverSpan.SpanContext().SpanID().String()))

	// Send again
	res, _ = Post(fmt.Sprintf("http://localhost:%d/%s", port, path), "application/json", bytes.NewBuffer([]byte(expectedRequestBody)))
	body, _ = io.ReadAll(res.Body)
	assert.Equal(t, expectedResponseBody, string(body))
	sr.ForceFlush(context.Background())
	spans = sr.Ended()
	serverSpan = spans[2]
	assert.Equal(t, trace.SpanKind(2), serverSpan.SpanKind())
	validateAttributes(serverSpan.Attributes(), path, metadataOnly, t)
	clientSpan = spans[3]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.False(t, clientSpan.Parent().HasTraceID())
	validateAttributes(clientSpan.Attributes(), path, metadataOnly, t)
}

func setupSpanRecording() *tracetest.SpanRecorder {
	resetClient()
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return sr
}

func resetClient() {
	otelhttp.DefaultClient = &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	DefaultClient = &Client{}
}

func testHelper(t *testing.T, port int, path string, metadataOnly bool) {
	sr := setupSpanRecording()
	Handle("/"+path, HandlerFunc(getHello))
	sendRequestAndValidate(t, port, path, metadataOnly, sr)
}

func TestServerInstrumentationWithSkippedContent(t *testing.T) {
	staticContentTestHelper(t, 8003, false)
}

func TestDisableInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	defer os.Setenv("HS_DISABLED", "")

	sr := setupSpanRecording()
	port := 8004
	path := "test4"
	HandleFunc("/"+path, getHello)

	go func() {
		ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	res, _ := Post(fmt.Sprintf("http://localhost:%d/%s", port, path), "application/json", bytes.NewBuffer([]byte(expectedRequestBody)))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, expectedResponseBody, string(body))
	fmt.Println("AFTER RESPONSE ASSERT: ", string(body))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 0, len(spans))
	fmt.Println("AFTER SPANS ASSERTION: ", len(spans))
}

func TestServerInstrumentation(t *testing.T) {
	testHelper(t, 8000, "test1", false)
}

func TestClientInstrumentation(t *testing.T) {
	sr := setupSpanRecording()
	client := &Client{}
	realClient := client.GetOriginHttpClient()
	_,_ = realClient.Get("google.com")
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	
	assert.Equal(t, 1, len(spans))
	clientSpan := spans[0]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.Contains(t, clientSpan.Attributes(), attribute.String("http.url", "google.com"))
}

func TestServerInstrumentationMetadataOnly(t *testing.T) {
	os.Setenv("HS_METADATA_ONLY", "true")
	testHelper(t, 8001, "test2", true)
}

func TestHandleFunc(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	port := 8002
	path := "test3"
	HandleFunc("/"+path, getHello)
	sendRequestAndValidate(t, port, path, true, sr)
}
