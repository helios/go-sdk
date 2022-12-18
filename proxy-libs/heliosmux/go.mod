module github.com/helios/go-sdk/proxy-libs/heliosmux

go 1.12

require (
	github.com/gorilla/mux v1.8.0
	github.com/helios/go-instrumentor/exports_extractor v0.0.9
	github.com/stretchr/testify v1.8.1
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.37.0
	go.opentelemetry.io/otel v1.11.2
	go.opentelemetry.io/otel/sdk v1.11.2
	go.opentelemetry.io/otel/trace v1.11.2
)
