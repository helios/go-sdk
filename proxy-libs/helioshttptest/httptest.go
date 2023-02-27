package helioshttptest

import (
	"io"
	"net/http"
	origin_httptest "net/http/httptest"
)

const DefaultRemoteAddr = origin_httptest.DefaultRemoteAddr

type Server = origin_httptest.Server
type ResponseRecorder = origin_httptest.ResponseRecorder

func NewRequest(method, target string, body io.Reader) *http.Request {
	return origin_httptest.NewRequest(method, target, body)
}

func NewRecorder() *ResponseRecorder {
	return origin_httptest.NewRecorder()
}

func NewServer(handler http.Handler) *Server {
	return origin_httptest.NewServer(handler)
}

func NewUnstartedServer(handler http.Handler) *Server {
	return origin_httptest.NewUnstartedServer(handler)
}

func NewTLSServer(handler http.Handler) *Server {
	return origin_httptest.NewTLSServer(handler)
}
