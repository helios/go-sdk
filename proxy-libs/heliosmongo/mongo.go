package heliosmongo

import (
	"context"

	originalMongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var InstrumentedSymbols = [...]string{"Connect", "NewClient"}

func setMonitor(opts []*options.ClientOptions) []*options.ClientOptions {
	monitor := otelmongo.NewMonitor()

	for _, opt := range opts {
		if opt != nil {
			opt.Monitor = monitor
		}
	}

	return opts
}

func Connect(ctx context.Context, opts ...*options.ClientOptions) (*originalMongo.Client, error) {
	return originalMongo.Connect(ctx, setMonitor(opts)...)
}

func NewClient(opts ...*options.ClientOptions) (*originalMongo.Client, error) {
	return originalMongo.NewClient(setMonitor(opts)...)
}
