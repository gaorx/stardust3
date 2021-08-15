package sdkafka

import (
	"github.com/Shopify/sarama"
	"github.com/gaorx/stardust3/sderr"
)

type SyncProducer struct {
	SyncProducer   sarama.SyncProducer
	MessageEncoder MessageEncoder
}

type AsyncProducer struct {
	AsyncProducer  sarama.AsyncProducer
	MessageEncoder MessageEncoder
}

// SyncProducer

func NewSyncProducer(addrs []string, config *sarama.Config, msgEnc MessageEncoder) (*SyncProducer, error) {
	if msgEnc == nil {
		return nil, sderr.New("nil message encoder")
	}
	p, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}
	return &SyncProducer{
		SyncProducer:   p,
		MessageEncoder: msgEnc,
	}, nil
}

func NewSyncProducerFromSarama(p sarama.SyncProducer, msgEnc MessageEncoder) (*SyncProducer, error) {
	if p == nil {
		return nil, sderr.New("nil sync producer")
	}
	if msgEnc == nil {
		return nil, sderr.New("nil message encoder")
	}
	return &SyncProducer{
		SyncProducer:   p,
		MessageEncoder: msgEnc,
	}, nil
}

func (p *SyncProducer) Close() error {
	return p.SyncProducer.Close()
}

func (p *SyncProducer) Send(msg *Message) error {
	pmsg, err := encodeProducerMessage(p.MessageEncoder, msg)
	if err != nil {
		return err
	}
	_, _, err = p.SyncProducer.SendMessage(pmsg)
	//fmt.Printf("****Send: %s %s\n", msg.Topic, freejson.MarshalString(msg.Value, ""))
	return err
}

func (p *SyncProducer) SendMulti(msgs []*Message) error {
	nMsgs := len(msgs)
	if nMsgs == 0 {
		return nil
	}

	pmsgs := make([]*sarama.ProducerMessage, 0, nMsgs)
	for _, msg := range msgs {
		pmsg, err := encodeProducerMessage(p.MessageEncoder, msg)
		if err != nil {
			return err
		}
		pmsgs = append(pmsgs, pmsg)
	}
	return p.SyncProducer.SendMessages(pmsgs)
}

// AsyncProducer

func NewAsyncProducer(addrs []string, config *sarama.Config, msgEnc MessageEncoder) (*AsyncProducer, error) {
	if msgEnc == nil {
		return nil, sderr.New("nil message encoder")
	}
	p, err := sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}
	return &AsyncProducer{
		AsyncProducer:  p,
		MessageEncoder: msgEnc,
	}, nil
}

func NewAsyncProducerFromSarama(p sarama.AsyncProducer, msgEnc MessageEncoder) (*AsyncProducer, error) {
	if p == nil {
		return nil, sderr.New("nil async producer")
	}
	if msgEnc == nil {
		return nil, sderr.New("nil message encoder")
	}
	return &AsyncProducer{
		AsyncProducer:  p,
		MessageEncoder: msgEnc,
	}, nil
}

func (p *AsyncProducer) Close() error {
	return p.AsyncProducer.Close()
}

func (p *AsyncProducer) Send(msg *Message) error {
	pmsg, err := encodeProducerMessage(p.MessageEncoder, msg)
	if err != nil {
		return err
	}

	p.AsyncProducer.Input() <- pmsg

	return nil
}

func (p *AsyncProducer) SendMulti(msgs []*Message) error {
	nMsgs := len(msgs)
	if nMsgs == 0 {
		return nil
	}

	pmsgs := make([]*sarama.ProducerMessage, 0, nMsgs)
	for _, msg := range msgs {
		pmsg, err := encodeProducerMessage(p.MessageEncoder, msg)
		if err != nil {
			return err
		}
		pmsgs = append(pmsgs, pmsg)
	}

	for _, pmsg := range pmsgs {
		p.AsyncProducer.Input() <- pmsg
	}

	return nil
}
