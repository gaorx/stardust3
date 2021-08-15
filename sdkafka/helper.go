package sdkafka

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/gaorx/stardust3/sdencoding"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdlog"
)

type consumerOffsetSaver func(cmsg *sarama.ConsumerMessage, p MessagePriority)

var (
	errNilMsg      = sderr.Sentinel("the message is nil")
	errNilMsgValue = sderr.Sentinel("the message value is nil")
)

func encodeProducerMessage(msgEnc MessageEncoder, msg *Message) (*sarama.ProducerMessage, error) {
	if msg == nil {
		return nil, errNilMsg
	}
	if msg.Value == nil {
		return nil, errNilMsgValue
	}
	data, err := msgEnc.Encode(msg.Value)
	if err != nil {
		return nil, err
	}
	return &sarama.ProducerMessage{
		Topic: msg.KafkaTopic(),
		Value: sarama.ByteEncoder(data),
	}, nil
}

func decodeConsumerMessage(msgEnc MessageEncoder, cmsg *sarama.ConsumerMessage, p MessagePriority) (*Message, error) {
	v, err := msgEnc.Decode(cmsg.Value)
	if err != nil {
		return nil, err
	}
	if p == PN {
		return &Message{
			Topic: cmsg.Topic,
			Value: v,
		}, nil
	} else {
		topic := strings.TrimSuffix(cmsg.Topic, "."+strconv.Itoa(int(p)))
		return &Message{
			Topic:    topic,
			Value:    v,
			Priority: p,
		}, nil
	}
}

func topicOfPriority(topic string, priority MessagePriority) string {
	if priority == PN {
		return topic
	}
	if priority < 0 || priority > PriorityNum {
		panic("Illegal message priority")
	}
	return fmt.Sprintf("%s.%d", topic, priority)
}

func handleConsumerMessage(c *Consumer, cmsg *sarama.ConsumerMessage, h ConsumerHandler, p MessagePriority, offsetSaver consumerOffsetSaver) {
	defer func() {
		err := sderr.ToErr(recover())
		if err != nil {
			sdlog.WithError(err).Error("Handle message painc")
		}
	}()

	// lock
	msgLock := c.MessageLock
	if msgLock != nil {
		ok, err := msgLock.TryLockMessage(cmsg.Topic, c.Group, cmsg.Partition, cmsg.Offset)
		if err != nil {
			sdlog.WithError(err).Warn("Kafka consumer try lock message error")
		} else {
			if !ok {
				return
			}
		}
	}

	// Save offset
	offsetSaver(cmsg, p)

	// Decode
	msg, err := decodeConsumerMessage(c.MessageEncoder, cmsg, p)
	if err != nil {
		sdlog.
			WithError(err).
			WithField("topic", msg.Topic).
			WithField("priority", msg.Priority).
			WithField("Data", sdencoding.Base64Url.EncodeStr(cmsg.Value)).
			Error("Decode message error")
		return
	}

	// Handle
	err = h(msg)
	if err != nil {
		sdlog.
			WithError(err).
			WithField("topic", msg.Topic).
			WithField("priority", msg.Priority).
			WithField("Data", sdencoding.Base64Url.EncodeStr(cmsg.Value)).
			Error("Handle message error")
	} else {
		//slog.
		//	WithField("topic", msg.Topic).
		//	WithField("priority", msg.Priority).
		//	WithField("Data", base64x.EncodeString(cmsg.Value)).
		//	Debug("Handle message")
	}
}

func closeReleasing(releasing []io.Closer) {
	for _, closable := range releasing {
		_ = closable.Close()
	}
}

func boolToStr(b bool) string {
	if b {
		return "y"
	} else {
		return "n"
	}
}

func intSliceToInt32Slice(a []int) []int32 {
	r := make([]int32, 0, len(a))
	for _, e := range a {
		r = append(r, int32(e))
	}
	return r
}
