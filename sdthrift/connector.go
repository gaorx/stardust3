package sdthrift

import (
	"github.com/gaorx/stardust3/sderr"
)

func SimpleConnector(connOpts ClientOptions) Connector {
	return func(addr string) (Client, error) {
		c := NewSimpleClient()
		err := c.Connect(addr, connOpts)
		if err != nil {
			return nil, sderr.WithStack(err)
		}
		return c, nil
	}
}

func ReuseConnConnector(connOpts ClientOptions, initCap, maxCap int) Connector {
	return func(addr string) (Client, error) {
		c := NewReuseConnClient(initCap, maxCap)
		err := c.Connect(addr, connOpts)
		if err != nil {
			return nil, sderr.WithStack(err)
		}
		return c, nil
	}
}
