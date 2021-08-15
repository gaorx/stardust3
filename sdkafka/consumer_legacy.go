package sdkafka

import (
	"github.com/Shopify/sarama"
	"github.com/gaorx/stardust3/sdcall"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdlog"
	kcg "github.com/wvanbergen/kafka/consumergroup"
)

type legacyConsumerGroupFactory func(topic string) (*kcg.ConsumerGroup, error)

type legacyConsumerImpl struct {
	*Consumer
	consumerGroupFactory legacyConsumerGroupFactory
	consumerGroup        *kcg.ConsumerGroup
}

func newLegacyConsumerImpl(consumer *Consumer, cgFactory legacyConsumerGroupFactory) (*Consumer, *legacyConsumerImpl) {
	impl := &legacyConsumerImpl{consumerGroupFactory: cgFactory}
	consumer.impl, impl.Consumer = impl, consumer
	return consumer, impl
}

func (c *legacyConsumerImpl) Close() error {
	if c.consumerGroup == nil {
		return nil
	}
	return c.consumerGroup.Close()
}

func (c *legacyConsumerImpl) ensureConsumerGroup(topic string) (*kcg.ConsumerGroup, error) {
	if c.consumerGroup == nil {
		cg, err := c.consumerGroupFactory(topic)
		if err != nil {
			return nil, err
		}
		c.consumerGroup = cg
	}
	return c.consumerGroup, nil
}

func (c *legacyConsumerImpl) handleMessage(cmsg *sarama.ConsumerMessage, h ConsumerHandler, p MessagePriority, cg *kcg.ConsumerGroup) {
	handleConsumerMessage(c.Consumer, cmsg, h, p, func(cmsg *sarama.ConsumerMessage, p MessagePriority) {
		_ = cg.CommitUpto(cmsg)
	})
}

func (c *legacyConsumerImpl) checkRunArgs(topic string, partitions []int, h ConsumerHandler) error {
	if topic == "" {
		return sderr.New("no topic")
	}
	if len(partitions) > 0 {
		return sderr.New("not support specified patitions")
	}
	if h == nil {
		return sderr.New("No handler")
	}
	return nil
}

func (c *legacyConsumerImpl) Run(topic string, partitions []int, h ConsumerHandler) error {
	if err := c.checkRunArgs(topic, partitions, h); err != nil {
		return err
	}
	cg, err := c.ensureConsumerGroup(topic)
	if err != nil {
		return err
	}

	logError := func() {
		for err := range cg.Errors() {
			sdlog.WithError(err).Warn("Legacy consumer error")
		}
	}

	process := func() {
		for cmsg := range cg.Messages() {
			c.handleMessage(cmsg, h, PN, cg)
		}
	}
	sdcall.Concurrent(0, []func(){process, logError})
	return nil
}

func (c *legacyConsumerImpl) RunP(topic string, partitions []int, h ConsumerHandler) error {
	if err := c.checkRunArgs(topic, partitions, h); err != nil {
		return err
	}
	return sderr.New("no impl")
}
