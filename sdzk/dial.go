package sdzk

import (
	"reflect"

	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdtime"
	"github.com/go-zookeeper/zk"
)

type Conn = zk.Conn
type Event = zk.Event

type Client struct {
	*Conn
	EC <-chan Event
}

type Address struct {
	Servers          []string `json:"servers" toml:"servers"`
	SessionTimeoutMs int64    `json:"session_timeout" toml:"session_timeout"`
}

type DialOptions = func(*zk.Conn)

var (
	WithDialer            = zk.WithDialer
	WithEventCallback     = zk.WithEventCallback
	WithHostProvider      = zk.WithHostProvider
	WithLogger            = zk.WithLogger
	WithLogInfo           = zk.WithLogInfo
	WithMaxBufferSize     = zk.WithMaxBufferSize
	WithMaxConnBufferSize = zk.WithMaxConnBufferSize
)

func Dial(addr Address, opts ...DialOptions) (*Client, error) {
	if len(addr.Servers) == 0 {
		return nil, sderr.New("no servers")
	}
	if addr.SessionTimeoutMs <= 0 {
		addr.SessionTimeoutMs = 60 * 1000 // 1分钟
	}

	argVals := []reflect.Value{
		reflect.ValueOf(addr.Servers),
		reflect.ValueOf(sdtime.Millis(addr.SessionTimeoutMs)),
	}
	for _, opt := range opts {
		if opt != nil {
			argVals = append(argVals, reflect.ValueOf(opt))
		}
	}
	connReturnVals := reflect.ValueOf(zk.Connect).Call(argVals)
	conn := connReturnVals[0].Interface().(*zk.Conn)
	ec := connReturnVals[1].Interface().(<-chan zk.Event)
	errVal := connReturnVals[2]
	if !errVal.IsNil() {
		return nil, sderr.WithStack(errVal.Interface().(error))
	}
	return &Client{Conn: conn, EC: ec}, nil
}

func (c *Client) Close() error {
	c.Conn.Close()
	return nil
}
