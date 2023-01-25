package helioshttp

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const responseBody = "hello1234"

// Simulate wrapped client with xray

type ApiClient struct {
	client *http.Client
	ctx    context.Context
}

func XrayClient(c *http.Client) *http.Client {
	if c == nil {
		c = http.DefaultClient
	}
	transport := c.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &http.Client{
		Transport:     RoundTripper(transport),
		CheckRedirect: c.CheckRedirect,
		Jar:           c.Jar,
		Timeout:       c.Timeout,
	}
}

func XrayRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return rt
}

func newWrappedClient() *http.Client {
	c := DefaultClient
	return &http.Client{
		Transport:     c.Transport,
		CheckRedirect: c.CheckRedirect,
		Jar:           c.Jar,
		Timeout:       c.Timeout,
	}
}

func NewApiClientWithContext(ctx context.Context) *ApiClient {
	return &ApiClient{
		client: XrayClient(newWrappedClient()),
		ctx:    ctx,
	}
}

func getHello(responseWriter ResponseWriter, request *Request) {
	io.WriteString(responseWriter, responseBody)
}

func validateAttributes(t *testing.T, attrs []attribute.KeyValue, path string) {
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "GET", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, path, value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		}
	}
}

func TestHttpInstrumentation(t *testing.T) {
	ctx := context.Background()
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	Handle("/test", HandlerFunc(getHello))
	Handle("/test2", HandlerFunc(getHello))
	go func() {
		ListenAndServe(":8000", nil)
	}()

	res, _ := Get("http://localhost:8000/test")
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, responseBody, string(body))
	sr.ForceFlush(ctx)
	spans := sr.Ended()
	assert.Equal(t, 2, len(spans))
	serverSpan := spans[0]
	assert.Equal(t, trace.SpanKind(2), serverSpan.SpanKind())
	validateAttributes(t, serverSpan.Attributes(), "/test")
	clientSpan := spans[1]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.False(t, clientSpan.Parent().HasTraceID())
	validateAttributes(t, clientSpan.Attributes(), "/test")
	assert.Equal(t, serverSpan.Parent().SpanID(), clientSpan.SpanContext().SpanID())

	client := NewApiClientWithContext(ctx)
	req, _ := NewRequestWithContext(ctx, "GET", "http://localhost:8000/test2", nil)
	res, _ = client.client.Do(req)
	res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)
	sr.ForceFlush(ctx)
	spans = sr.Ended()
	assert.Equal(t, 4, len(spans))
	xrayClientSpan := spans[3]
	validateAttributes(t, xrayClientSpan.Attributes(), "/test2")
}
