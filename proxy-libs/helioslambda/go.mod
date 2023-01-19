module github.com/helios/go-sdk/proxy-libs/helioslambda

go 1.18

require (
	github.com/aws/aws-lambda-go v1.35.0
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda v0.37.0
)

require go.opentelemetry.io/otel/trace v1.11.2 // indirect

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
)
