package sdthrift

import (
	"time"

	"github.com/apache/thrift/lib/go/thrift"
)

type ClientFactory func(t thrift.TTransport, pf thrift.TProtocolFactory) (interface{}, error)

type ClientOptions struct {
	Secure         bool
	SocketTimeout  time.Duration
	ConnectTimeout time.Duration
	TF             thrift.TTransportFactory
	PF             thrift.TProtocolFactory
	CF             ClientFactory
}

type Client interface {
	Connect(addr string, opts ClientOptions) error
	Close() error
	Use(f func(clientIntf interface{}) (interface{}, error)) (interface{}, error)
}
