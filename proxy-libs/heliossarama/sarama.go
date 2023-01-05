package heliossarama

import (
	"context"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
)

var InstrumentedSymbols = [...]string{"NewAsyncProducer", "NewAsyncProducerFromClient", "NewConsumerGroup", "NewConsumerGroupFromClient", "NewSyncProducer", "NewSyncProducerFromClient"}

func NewAsyncProducer(addrs []string, conf *sarama.Config) (sarama.AsyncProducer, error) {
	asyncProducer, error := sarama.NewAsyncProducer(addrs, conf)

	if asyncProducer != nil {
		asyncProducer = asyncProducerWrapper{asyncProducer}
		asyncProducer = otelsarama.WrapAsyncProducer(conf, asyncProducer)
	}

	return asyncProducer, error
}

func NewAsyncProducerFromClient(client sarama.Client) (sarama.AsyncProducer, error) {
	asyncProducer, error := sarama.NewAsyncProducerFromClient(client)

	if asyncProducer != nil {
		asyncProducer = asyncProducerWrapper{asyncProducer}
		asyncProducer = otelsarama.WrapAsyncProducer(client.Config(), asyncProducer)
	}

	return asyncProducer, error
}

func NewConsumerGroup(addrs []string, groupID string, config *sarama.Config) (sarama.ConsumerGroup, error) {
	consumerGroup, error := sarama.NewConsumerGroup(addrs, groupID, config)

	if consumerGroup != nil {
		consumerGroup = consumerGroupWrapper{consumerGroup}
	}

	return consumerGroup, error
}

func NewConsumerGroupFromClient(groupID string, client sarama.Client) (sarama.ConsumerGroup, error) {
	consumerGroup, error := sarama.NewConsumerGroupFromClient(groupID, client)

	if consumerGroup != nil {
		consumerGroup = consumerGroupWrapper{consumerGroup}
	}

	return consumerGroup, error
}

func NewSyncProducer(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
	syncProducer, error := sarama.NewSyncProducer(addrs, config)

	if syncProducer != nil {
		syncProducer = otelsarama.WrapSyncProducer(config, syncProducer)
	}

	return syncProducer, error
}

func NewSyncProducerFromClient(client sarama.Client) (sarama.SyncProducer, error) {
	syncProducer, error := sarama.NewSyncProducerFromClient(client)

	if syncProducer != nil {
		syncProducer = otelsarama.WrapSyncProducer(client.Config(), syncProducer)
	}

	return syncProducer, error
}

func InjectContextToMessage(ctx context.Context, message *sarama.ProducerMessage) {
	carrier := otelsarama.NewProducerMessageCarrier(message)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}
