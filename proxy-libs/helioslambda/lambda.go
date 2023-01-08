package helioslambda

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
)

var InstrumentedSymbols = [...]string{"Start", "StartWithContext", "StartWithOptions"}

func Start(handler interface{}) {
	lambda.Start(otellambda.InstrumentHandler(handler))
}

func StartWithOptions(handler interface{}, options ...lambda.Option) {
	lambda.StartWithOptions(otellambda.InstrumentHandler(handler), options...)
}

func StartWithContext(ctx context.Context, handler interface{}) {
	lambda.StartWithContext(ctx, otellambda.InstrumentHandler(handler))
}

func StartHandler(handler lambda.Handler) {
	lambda.StartHandler(handler)
}
