package sdkafka

import (
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdjson"
)

// MessageEncoder

type MessageEncoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte) (interface{}, error)
}

var (
	NilEncoder  MessageEncoder = &nilEncoder{}
	TextEncoder MessageEncoder = &textEncoder{}
)

// Nil encoder

type nilEncoder struct {
	Encoder func(v interface{}) ([]byte, error)
	Decoder func(data []byte) (interface{}, error)
}

func (enc *nilEncoder) Encode(v interface{}) ([]byte, error) {
	data, ok := v.([]byte)
	if !ok {
		return nil, sderr.New("the value is not byte slice")
	}
	return data, nil
}

func (enc *nilEncoder) Decode(data []byte) (interface{}, error) {
	return data, nil
}

// Function encoder

type FuncEncoder struct {
	Encoder func(v interface{}) ([]byte, error)
	Decoder func(data []byte) (interface{}, error)
}

func (enc *FuncEncoder) Encode(v interface{}) ([]byte, error) {
	if enc.Encoder == nil {
		return nil, sderr.New("nil encoder")
	}
	r, err := enc.Encoder(v)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func (enc *FuncEncoder) Decode(data []byte) (interface{}, error) {
	if enc.Decoder == nil {
		return nil, sderr.New("nil decoder")
	}
	r, err := enc.Decoder(data)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

// text encoder

type textEncoder struct {
}

func (enc *textEncoder) Encode(v interface{}) ([]byte, error) {
	s, ok := v.(string)
	if !ok {
		return nil, sderr.New("the value is not string")
	}
	return []byte(s), nil
}

func (enc *textEncoder) Decode(data []byte) (interface{}, error) {
	return string(data), nil
}

// JsonEncoder

type JsonEncoder struct {
}

func NewJsonEncoder() JsonEncoder {
	return JsonEncoder{}
}

func (enc *JsonEncoder) Encode(v interface{}) ([]byte, error) {
	r, err := sdjson.Marshal(v)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func (enc *JsonEncoder) Decode(data []byte) (interface{}, error) {
	r, err := sdjson.UnmarshalValue(data)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r.Interface, nil
}
