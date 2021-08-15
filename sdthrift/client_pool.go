package sdthrift

import (
	"sync"

	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdlog"
	"github.com/gaorx/stardust3/sdsync"
)

type Connector func(addr string) (Client, error)

type Pool struct {
	mtx     sync.Mutex
	clients map[string]Client
}

func NewPool() *Pool {
	return &Pool{
		clients: map[string]Client{},
	}
}

func (p *Pool) Open(addr string, connector Connector) (Client, error) {
	if addr == "" {
		return nil, sderr.New("addr is nil")
	}
	if connector == nil {
		return nil, sderr.New("no connector")
	}

	var err error
	var c Client
	sdsync.Lock(&p.mtx, func() {
		c0, ok := p.clients[addr]
		if ok {
			c, err = c0, nil
			return
		}
		c0, err0 := connector(addr)
		if err0 != nil {
			c, err = nil, err0
			return
		}
		if c0 == nil {
			c, err = nil, sderr.New("nil client from connector")
			return
		}
		p.clients[addr] = c0
		c, err = c0, nil
	})
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return c, nil
}

func (p *Pool) Close() error {
	sdsync.Lock(&p.mtx, func() {
		for addr, c := range p.clients {
			_ = c.Close()
			sdlog.WithField("addr", addr).Debug("Close thrift client")
		}
	})
	return nil
}
