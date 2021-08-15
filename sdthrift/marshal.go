package sdthrift

import (
	"bytes"
	"context"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/gaorx/stardust3/sderr"
)

var (
	compactPF thrift.TProtocolFactory = thrift.NewTCompactProtocolFactory()
	binaryPF  thrift.TProtocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
)

var (
	nilStructErr = sderr.Sentinel("nil struct")
)

func Marshal(pf thrift.TProtocolFactory, v thrift.TStruct) ([]byte, error) {
	if v == nil {
		return nil, sderr.WithStack(nilStructErr)
	}
	buf := thrift.TMemoryBuffer{Buffer: &bytes.Buffer{}}
	p := pf.GetProtocol(&buf)
	err := v.Write(context.Background(), p)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	err = buf.Flush(context.Background())
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return buf.Bytes(), nil
}

func Unmarshal(pf thrift.TProtocolFactory, data []byte, v thrift.TStruct) error {
	if v == nil {
		return sderr.WithStack(nilStructErr)
	}
	buf := thrift.TMemoryBuffer{Buffer: bytes.NewBuffer(data)}
	p := pf.GetProtocol(&buf)
	return sderr.WithStack(v.Read(context.Background(), p))
}

func MarshalCompact(v thrift.TStruct) ([]byte, error) {
	return Marshal(compactPF, v)
}

func UnmarshalCompact(data []byte, v thrift.TStruct) error {
	return Unmarshal(compactPF, data, v)
}

func MarshalBinary(v thrift.TStruct) ([]byte, error) {
	return Marshal(binaryPF, v)
}

func UnmarshalBinary(data []byte, v thrift.TStruct) error {
	return Unmarshal(binaryPF, data, v)
}
