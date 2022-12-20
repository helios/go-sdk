package heliossarama

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
)

type consumerGroupWrapper struct {
	consumerGroup ConsumerGroup
}

func (wrapper consumerGroupWrapper) Consume(ctx context.Context, topics []string, handler ConsumerGroupHandler) error {
	handler = otelsarama.WrapConsumerGroupHandler(handler)
	return wrapper.consumerGroup.Consume(ctx, topics, handler)
}

func (wrapper consumerGroupWrapper) Errors() <-chan error {
	return wrapper.consumerGroup.Errors()
}

func (wrapper consumerGroupWrapper) Close() error {
	return wrapper.consumerGroup.Close()
}

func (wrapper consumerGroupWrapper) Pause(partitions map[string][]int32) {
	wrapper.consumerGroup.Pause(partitions)
}

func (wrapper consumerGroupWrapper) Resume(partitions map[string][]int32) {
	wrapper.consumerGroup.Resume(partitions)
}

func (wrapper consumerGroupWrapper) PauseAll() {
	wrapper.consumerGroup.PauseAll()
}

func (wrapper consumerGroupWrapper) ResumeAll() {
	wrapper.consumerGroup.ResumeAll()
}
