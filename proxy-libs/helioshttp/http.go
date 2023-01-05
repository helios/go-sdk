package helioshttp

import (
	"context"
	"io"
	realHttp "net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var InstrumentedSymbols = [...]string{"Handle", "Client", "DefaultClient", "Get", "Post", "Head"}

var DefaultClient = &Client{}

func Handle(pattern string, handler realHttp.Handler) {
	handler = otelhttp.NewHandler(handler, pattern)
	realHttp.Handle(pattern, handler)
}

func Get(url string) (resp *realHttp.Response, err error) {
	ctx := context.Background()
	req, err := realHttp.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	return DefaultClient.Do(req)
}

func Post(url, contentType string, body io.Reader) (resp *realHttp.Response, err error) {
	ctx := context.Background()
	req, err := realHttp.NewRequestWithContext(ctx, "POST", url, body)
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		return nil, err
	}

	return DefaultClient.Do(req)
}

func Head(url string) (resp *realHttp.Response, err error) {
	ctx := context.Background()
	req, err := realHttp.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return DefaultClient.Do(req)
}
