package helioshttptest

import (
	"io"
	origin_httptest "net/http/httptest"

	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
	"github.com/helios/opentelemetry-go-contrib/instrumentation/net/http/otelhttp"
)

const DefaultRemoteAddr = origin_httptest.DefaultRemoteAddr

type ResponseRecorder = origin_httptest.ResponseRecorder

func NewRequest(method, target string, body io.Reader) *http.Request {
	return origin_httptest.NewRequest(method, target, body)
}

func NewRecorder() *ResponseRecorder {
	return origin_httptest.NewRecorder()
}

func wrapServer(wrappedServer *origin_httptest.Server) *Server {
	return &Server{
		URL:           wrappedServer.URL,
		Listener:      wrappedServer.Listener,
		EnableHTTP2:   wrappedServer.EnableHTTP2,
		TLS:           wrappedServer.TLS,
		Config:        wrappedServer.Config,
		wrappedServer: wrappedServer,
	}
}

func NewServer(handler http.Handler) *Server {
	wrappedHandler := otelhttp.NewHandler(handler, "test_server")
	wrappedServer := origin_httptest.NewServer(wrappedHandler)
	return wrapServer(wrappedServer)
}

func NewUnstartedServer(handler http.Handler) *Server {
	wrappedHandler := otelhttp.NewHandler(handler, "test_server")
	wrappedServer := origin_httptest.NewUnstartedServer(wrappedHandler)
	return wrapServer(wrappedServer)
}

func NewTLSServer(handler http.Handler) *Server {
	wrappedHandler := otelhttp.NewHandler(handler, "test_server")
	wrappedServer := origin_httptest.NewTLSServer(wrappedHandler)
	return wrapServer(wrappedServer)
}
