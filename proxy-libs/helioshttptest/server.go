package helioshttptest

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	origin_httptest "net/http/httptest"

	"github.com/helios/opentelemetry-go-contrib/instrumentation/net/http/otelhttp"
)

type Server struct {
	URL           string
	Listener      net.Listener
	EnableHTTP2   bool
	TLS           *tls.Config
	Config        *http.Server
	wrappedServer *origin_httptest.Server
}

func (s *Server) Start() {
	s.wrappedServer.Start()
}

func (s *Server) StartTLS() {
	s.wrappedServer.StartTLS()
}

func (s *Server) Close() {
	s.wrappedServer.Close()
}

func (s *Server) CloseClientConnections() {
	s.wrappedServer.CloseClientConnections()
}

func (s *Server) Certificate() *x509.Certificate {
	return s.wrappedServer.Certificate()
}

func (s *Server) Client() *http.Client {
	client := s.wrappedServer.Client()
	origTransport := client.Transport
	client.Transport = otelhttp.NewTransport(origTransport)
	return client
}
