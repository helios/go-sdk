package helioshttp

import (
	"io"
	realHttp "net/http"
	"net/url"
	"reflect"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	realClient    realHttp.Client
	initialized   bool
	Transport     RoundTripper
	CheckRedirect func(req *Request, via []*Request) error
	Jar           CookieJar
	Timeout       time.Duration
}

func copyClientProxyToReal(from *Client, to *realHttp.Client) {
	if !from.initialized {
		if from.Transport == nil {
			to.Transport = otelhttp.NewTransport(realHttp.DefaultTransport)
		} else {
			// Check if the existing transport is the wrapped OpenTelemetry transport
			transportType := reflect.ValueOf(from.Transport).Elem()
			found := transportType.FieldByName("tracer")
			if found == (reflect.Value{}) {
				to.Transport = otelhttp.NewTransport(from.Transport)
			} else {
				to.Transport = from.Transport
			}
		}

		from.initialized = true
	}

	to.CheckRedirect = from.CheckRedirect
	to.Jar = from.Jar
	to.Timeout = from.Timeout
}

func copyClientRealToProxy(from *realHttp.Client, to *Client) {
	// The only field ever modified by net/http is the CookieJar
	to.Jar = from.Jar
}

func (c *Client) CloseIdleConnections() {
	copyClientProxyToReal(c, &c.realClient)
	c.realClient.CloseIdleConnections()
}

func (c *Client) Do(req *Request) (resp *Response, err error) {
	copyClientProxyToReal(c, &c.realClient)
	resp, err = c.realClient.Do(req)
	copyClientRealToProxy(&c.realClient, c)

	return resp, err
}

func (c *Client) Get(url string) (resp *Response, err error) {
	copyClientProxyToReal(c, &c.realClient)
	resp, err = c.realClient.Get(url)
	copyClientRealToProxy(&c.realClient, c)

	return resp, err
}

func (c *Client) Head(url string) (resp *Response, err error) {
	copyClientProxyToReal(c, &c.realClient)
	resp, err = c.realClient.Head(url)
	copyClientRealToProxy(&c.realClient, c)

	return resp, err
}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *Response, err error) {
	copyClientProxyToReal(c, &c.realClient)
	resp, err = c.realClient.Post(url, contentType, body)
	copyClientRealToProxy(&c.realClient, c)

	return resp, err
}

func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error) {
	copyClientProxyToReal(c, &c.realClient)
	resp, err = c.realClient.PostForm(url, data)
	copyClientRealToProxy(&c.realClient, c)

	return resp, err
}
