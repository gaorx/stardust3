package sdkafka

import (
	"fmt"
	"io"

	"github.com/Shopify/sarama"
	"github.com/gaorx/stardust3/sdcall"
	"github.com/gaorx/stardust3/sdchan"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdlog"
)

type consumerImpl struct {
	*Consumer
	Client sarama.Client
}

func newConsumerImpl(consumer *Consumer, client sarama.Client) (*Consumer, *consumerImpl) {
	impl := &consumerImpl{Client: client}
	consumer.impl, impl.Consumer = impl, consumer
	return consumer, impl
}

func (c *consumerImpl) Close() error {
	return c.Client.Close()
}

func (c *consumerImpl) createPartitionOffsetManagersAndConsumers(topic string, partitions []int32, priority MessagePriority, releasing *[]io.Closer) ([]partitionOffsetManagerAndConsumer, error) {
	c_, err := sarama.NewConsumerFromClient(c.Client)
	if err != nil {
		return nil, err
	}
	*releasing = append(*releasing, c_)

	om, err := sarama.NewOffsetManagerFromClient(c.Group, c.Client)
	if err != nil {
		return nil, err
	}
	*releasing = append(*releasing, om)

	makePartitionOffsetManagerAndConsumer := func(partition0 int32) (partitionOffsetManagerAndConsumer, error) {
		pom, err := om.ManagePartition(topic, partition0)
		if err != nil {
			return partitionOffsetManagerAndConsumer{}, err
		}
		*releasing = append(*releasing, pom)
		offset, _ := pom.NextOffset()

		pc, err := c_.ConsumePartition(topic, partition0, offset)
		if err != nil {
			return partitionOffsetManagerAndConsumer{}, err
		}
		*releasing = append(*releasing, pc)

		return partitionOffsetManagerAndConsumer{
			partition:     partition0,
			priority:      priority,
			offsetManager: pom,
			consumer:      pc,
		}, nil
	}

	var r []partitionOffsetManagerAndConsumer
	if len(partitions) == 0 {
		partitions, err = c_.Partitions(topic)
		if err != nil {
			return nil, err
		}
	}
	for _, partition0 := range partitions {
		pomc, err := makePartitionOffsetManagerAndConsumer(partition0)
		if err != nil {
			return nil, err
		}
		r = append(r, pomc)
	}
	if len(r) == 0 {
		return nil, fmt.Errorf("Not found partitions for topic: %s", topic)
	}
	return r, nil
}

func (c *consumerImpl) handleMessage(cmsg *sarama.ConsumerMessage, h ConsumerHandler, p MessagePriority, pom sarama.PartitionOffsetManager) {
	handleConsumerMessage(c.Consumer, cmsg, h, p, func(cmsg *sarama.ConsumerMessage, p MessagePriority) {
		pom.MarkOffset(cmsg.Offset+1, "")
	})
}

func (c *consumerImpl) checkRunArgs(topic string, partitions []int, h ConsumerHandler) error {
	if topic == "" {
		return sderr.New("oo topic")
	}
	if h == nil {
		return sderr.New("no handler")
	}
	return nil
}

func (c *consumerImpl) Run(topic string, partitions []int, h ConsumerHandler) error {
	if err := c.checkRunArgs(topic, partitions, h); err != nil {
		return err
	}

	partitionsI32 := intSliceToInt32Slice(partitions)

	var releasing []io.Closer
	defer closeReleasing(releasing)

	loop := func(_ int, pomci interface{}) {
		pomc := pomci.(partitionOffsetManagerAndConsumer)
		sdlog.WithField("topic", topic).WithField("partition", pomc.partition).Debug("Start consumer Kafka topic")
		pom, pc := pomc.offsetManager, pomc.consumer
		for {
			select {
			case cmsg := <-pc.Messages():
				c.handleMessage(cmsg, h, PN, pom)
			}
		}
	}

	pomcs, err := c.createPartitionOffsetManagersAndConsumers(topic, partitionsI32, PN, &releasing)
	if err != nil {
		return err
	}
	if len(pomcs) == 1 {
		loop(0, pomcs[0])
		return nil
	} else {
		return sdcall.ConcurrentFuse(0, pomcs, loop)
	}
}

func (c *consumerImpl) RunP(topic string, partitions []int, h ConsumerHandler) error {
	if err := c.checkRunArgs(topic, partitions, h); err != nil {
		return err
	}

	partitionsI32 := intSliceToInt32Slice(partitions)

	var releasing []io.Closer
	defer closeReleasing(releasing)

	loop := func(_ int, ppomci interface{}) {
		ppomc := ppomci.(*partitionPriorityOffsetManagerAndConsumer)
		topicStatus := fmt.Sprintf("%s(PN:%s,P1:%s,P2:%s,P3:%s,P4:%s)",
			topic,
			boolToStr(ppomc.omPN != nil && ppomc.cPN != nil),
			boolToStr(ppomc.omP1 != nil && ppomc.cP1 != nil),
			boolToStr(ppomc.omP2 != nil && ppomc.cP2 != nil),
			boolToStr(ppomc.omP3 != nil && ppomc.cP3 != nil),
			boolToStr(ppomc.omP4 != nil && ppomc.cP4 != nil),
		)
		sdlog.WithField("topic", topicStatus).
			WithField("paritition", ppomc.partition).
			Debug("Start consumer Kafka priority topic")
		recvChans := make([]interface{}, 5, 5)
		for {
			// P4
			if ppomc.cP4 != nil {
				select {
				case cmsg := <-ppomc.cP4.Messages():
					c.handleMessage(cmsg, h, P4, ppomc.omP4)
					continue
				default:
				}
			}

			// P3
			if ppomc.cP3 != nil {
				select {
				case cmsg := <-ppomc.cP3.Messages():
					c.handleMessage(cmsg, h, P3, ppomc.omP3)
					continue
				default:
				}
			}

			// P2
			if ppomc.cP2 != nil {
				select {
				case cmsg := <-ppomc.cP2.Messages():
					c.handleMessage(cmsg, h, P2, ppomc.omP2)
					continue
				default:
				}
			}

			// P1
			if ppomc.cP1 != nil {
				select {
				case cmsg := <-ppomc.cP1.Messages():
					c.handleMessage(cmsg, h, P1, ppomc.omP1)
					continue
				default:
				}
			}

			// Default
			if ppomc.cP4 != nil {
				recvChans[0] = ppomc.cP4.Messages()
			} else {
				recvChans[0] = nil
			}
			if ppomc.cP3 != nil {
				recvChans[1] = ppomc.cP3.Messages()
			} else {
				recvChans[1] = nil
			}
			if ppomc.cP2 != nil {
				recvChans[2] = ppomc.cP2.Messages()
			} else {
				recvChans[2] = nil
			}
			if ppomc.cP1 != nil {
				recvChans[3] = ppomc.cP1.Messages()
			} else {
				recvChans[3] = nil
			}
			if ppomc.cPN != nil {
				recvChans[4] = ppomc.cPN.Messages()
			} else {
				recvChans[4] = nil
			}

			index, cmsgi, ok := sdchan.ReceiveSelect(recvChans)
			if !ok {
				continue
			}
			cmsg, ok := cmsgi.(*sarama.ConsumerMessage)
			if !ok {
				continue
			}
			switch index {
			case 0:
				c.handleMessage(cmsg, h, P4, ppomc.omP4)
			case 1:
				c.handleMessage(cmsg, h, P3, ppomc.omP3)
			case 2:
				c.handleMessage(cmsg, h, P2, ppomc.omP2)
			case 3:
				c.handleMessage(cmsg, h, P1, ppomc.omP1)
			case 4:
				c.handleMessage(cmsg, h, PN, ppomc.omPN)
			}

			// 上面的代码约等于下面的形式，只不过对于ppomc.cPx为nil不敏感
			//select {
			//case cmsg := <-ppomc.cP4.Messages():
			//	c.handleMessage(cmsg, h, ppomc.omP4, P4)
			//case cmsg := <-ppomc.cP3.Messages():
			//	c.handleMessage(cmsg, h, ppomc.omP3, P3)
			//case cmsg := <-ppomc.cP2.Messages():
			//	c.handleMessage(cmsg, h, ppomc.omP2, P2)
			//case cmsg := <-ppomc.cP1.Messages():
			//	c.handleMessage(cmsg, h, ppomc.omP1, P1)
			//case cmsg := <-ppomc.cPN.Messages():
			//	c.handleMessage(cmsg, h, ppomc.omPN, PN)
			//}
		}
	}

	pomcsPN, err := c.createPartitionOffsetManagersAndConsumers(topicOfPriority(topic, PN), partitionsI32, PN, &releasing)
	if err != nil && err != sarama.ErrUnknownTopicOrPartition {
		return err
	}

	pomcsP1, err := c.createPartitionOffsetManagersAndConsumers(topicOfPriority(topic, P1), partitionsI32, P1, &releasing)
	if err != nil && err != sarama.ErrUnknownTopicOrPartition {
		return err
	}

	pomcsP2, err := c.createPartitionOffsetManagersAndConsumers(topicOfPriority(topic, P2), partitionsI32, P2, &releasing)
	if err != nil && err != sarama.ErrUnknownTopicOrPartition {
		return err
	}

	pomcsP3, err := c.createPartitionOffsetManagersAndConsumers(topicOfPriority(topic, P3), partitionsI32, P3, &releasing)
	if err != nil && err != sarama.ErrUnknownTopicOrPartition {
		return err
	}

	pomcsP4, err := c.createPartitionOffsetManagersAndConsumers(topicOfPriority(topic, P4), partitionsI32, P4, &releasing)
	if err != nil && err != sarama.ErrUnknownTopicOrPartition {
		return err
	}

	if len(pomcsPN) == 0 && len(pomcsP1) == 0 && len(pomcsP2) == 0 && len(pomcsP3) == 0 && len(pomcsP4) == 0 {
		return fmt.Errorf("Not found topic (%s | %s.1 | %s.2 | %s.3 | %s.4) ", topic, topic, topic, topic, topic)
	}

	ppomcs := groupOffsetManagerAndConsumerForPriority(pomcsPN, pomcsP1, pomcsP2, pomcsP3, pomcsP4)
	switch len(ppomcs) {
	case 0:
		return sderr.New("group error")
	case 1:
		loop(0, ppomcs[0])
		return nil
	default:
		return sdcall.ConcurrentFuse(0, ppomcs, loop)
	}
}

func groupOffsetManagerAndConsumerForPriority(pomcsArr ...[]partitionOffsetManagerAndConsumer) []*partitionPriorityOffsetManagerAndConsumer {
	m := make(map[int32]*partitionPriorityOffsetManagerAndConsumer, 16)
	for _, pomcs := range pomcsArr {
		for _, pomc := range pomcs {
			v, ok := m[pomc.partition]
			if !ok {
				v = &partitionPriorityOffsetManagerAndConsumer{
					partition: pomc.partition,
				}
				m[pomc.partition] = v
			}
			switch pomc.priority {
			case PN:
				v.omPN, v.cPN = pomc.offsetManager, pomc.consumer
			case P1:
				v.omP1, v.cP1 = pomc.offsetManager, pomc.consumer
			case P2:
				v.omP2, v.cP2 = pomc.offsetManager, pomc.consumer
			case P3:
				v.omP3, v.cP3 = pomc.offsetManager, pomc.consumer
			case P4:
				v.omP4, v.cP4 = pomc.offsetManager, pomc.consumer
			}
		}
	}
	var r []*partitionPriorityOffsetManagerAndConsumer
	for _, ppomc := range m {
		r = append(r, ppomc)
	}
	// TODO: 应该根据partition来个排序
	return r
}

type partitionOffsetManagerAndConsumer struct {
	partition     int32
	priority      MessagePriority
	offsetManager sarama.PartitionOffsetManager
	consumer      sarama.PartitionConsumer
}

type partitionPriorityOffsetManagerAndConsumer struct {
	partition                    int32
	omPN, omP1, omP2, omP3, omP4 sarama.PartitionOffsetManager
	cPN, cP1, cP2, cP3, cP4      sarama.PartitionConsumer
}
