package sdkafka

import (
	"io"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gaorx/stardust3/sderr"
	kcg "github.com/wvanbergen/kafka/consumergroup"
)

type ConsumerHandler func(msg *Message) error

type Consumer struct {
	impl           consumer
	MessageEncoder MessageEncoder
	Group          string
	MessageLock    ConsumerMessageLock
}

type consumer interface {
	io.Closer
	Run(topic string, partitions []int, h ConsumerHandler) error
	RunP(topic string, partitions []int, h ConsumerHandler) error
}

// default

func NewConsumer(addrs []string, config *sarama.Config, msgEnc MessageEncoder, group string) (*Consumer, error) {
	if msgEnc == nil {
		return nil, sderr.New("nil message encoder")
	}
	client, err := sarama.NewClient(addrs, config)
	if err != nil {
		return nil, err
	}
	consumer, _ := newConsumerImpl(&Consumer{
		MessageEncoder: msgEnc,
		Group:          group,
	}, client)
	return consumer, nil
}

func NewConsumerFromSarama(client sarama.Client, msgEnc MessageEncoder, group string) (*Consumer, error) {
	if msgEnc == nil {
		return nil, sderr.New("nil message encoder")
	}
	consumer, _ := newConsumerImpl(&Consumer{
		MessageEncoder: msgEnc,
		Group:          group,
	}, client)
	return consumer, nil
}

// legacy

func NewLegacyConsumer(zkAddrs []string, zkRoot string, config *sarama.Config, msgEnc MessageEncoder, group string) (*Consumer, error) {
	if msgEnc == nil {
		return nil, sderr.New("nil message encoder")
	}
	if len(zkAddrs) == 0 {
		return nil, sderr.New("no addresses")
	}
	if zkRoot == "" {
		return nil, sderr.New("no root")
	}
	cgConfig := kcg.NewConfig()
	if config != nil {
		cgConfig.Config = config
	}
	cgConfig.Offsets.Initial = sarama.OffsetNewest
	cgConfig.Offsets.ProcessingTimeout = time.Second * 20
	cgConfig.Offsets.CommitInterval = time.Second * 10
	cgConfig.Offsets.ResetOffsets = false
	cgConfig.Zookeeper.Chroot = zkRoot

	cgf := func(topic string) (*kcg.ConsumerGroup, error) {
		return kcg.JoinConsumerGroup(group, []string{topic}, zkAddrs, cgConfig)
	}

	consumer, _ := newLegacyConsumerImpl(&Consumer{
		MessageEncoder: msgEnc,
		Group:          group,
	}, cgf)
	return consumer, nil
}

func (c *Consumer) Close() error {
	return c.impl.Close()
}

func (c *Consumer) Run(topic string, partitions []int, h ConsumerHandler) error {
	return c.impl.Run(topic, partitions, h)
}

func (c *Consumer) RunP(topic string, partitions []int, h ConsumerHandler) error {
	return c.impl.RunP(topic, partitions, h)
}
