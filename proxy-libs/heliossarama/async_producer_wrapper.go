package heliossarama

import (
	"context"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
)

type asyncProducerWrapper struct {
	asyncProducer sarama.AsyncProducer
}

func (wrapper asyncProducerWrapper) AsyncClose() {
	wrapper.asyncProducer.AsyncClose()
}

func (wrapper asyncProducerWrapper) Close() error {
	return wrapper.asyncProducer.Close()
}

func (wrapper asyncProducerWrapper) Input() chan<- *sarama.ProducerMessage {
	channel := make(chan *sarama.ProducerMessage, 1)

	go func() {
		producerMessage := <-channel
		carrier := otelsarama.NewProducerMessageCarrier(producerMessage)
		otel.GetTextMapPropagator().Inject(context.Background(), carrier)
		wrapper.asyncProducer.Input() <- producerMessage
		close(channel)
	}()

	return channel
}

func (wrapper asyncProducerWrapper) Successes() <-chan *sarama.ProducerMessage {
	return wrapper.asyncProducer.Successes()
}

func (wrapper asyncProducerWrapper) Errors() <-chan *sarama.ProducerError {
	return wrapper.asyncProducer.Errors()
}

func (wrapper asyncProducerWrapper) IsTransactional() bool {
	return wrapper.asyncProducer.IsTransactional()
}

func (wrapper asyncProducerWrapper) TxnStatus() sarama.ProducerTxnStatusFlag {
	return wrapper.asyncProducer.TxnStatus()
}

func (wrapper asyncProducerWrapper) BeginTxn() error {
	return wrapper.asyncProducer.BeginTxn()
}

func (wrapper asyncProducerWrapper) CommitTxn() error {
	return wrapper.asyncProducer.CommitTxn()
}

func (wrapper asyncProducerWrapper) AbortTxn() error {
	return wrapper.asyncProducer.AbortTxn()
}

func (wrapper asyncProducerWrapper) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	return wrapper.asyncProducer.AddOffsetsToTxn(offsets, groupId)
}

func (wrapper asyncProducerWrapper) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	return wrapper.asyncProducer.AddMessageToTxn(msg, groupId, metadata)
}
