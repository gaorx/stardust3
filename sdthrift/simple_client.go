package sdthrift

import (
	"crypto/tls"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gaorx/stardust3/sdbackoff"
	"github.com/gaorx/stardust3/sdcall"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdtime"
)

type SimpleClient struct {
	addr    string
	opts    ClientOptions
	backOff sdbackoff.BackOff
}

func NewSimpleClient() *SimpleClient {
	return &SimpleClient{}
}

func (c *SimpleClient) Connect(addr string, opts ClientOptions) error {
	if addr == "" {
		return sderr.New("addr is empty")
	}
	if opts.TF == nil {
		return sderr.New("nil transport factory")
	}
	if opts.PF == nil {
		return sderr.New("nil protocol factory")
	}
	if opts.CF == nil {
		return sderr.New("nil client factory")
	}
	c.addr = addr
	c.opts = opts
	c.backOff = sdbackoff.Synchronized(sdbackoff.Exponential(&sdbackoff.ExponentialOptions{
		InitialInterval:     sdtime.Millis(5),
		RandomizationFactor: 0.0,
		Multiplier:          1.8,
		MaxInterval:         sdtime.Millis(40),
	}))
	return nil
}

func (c *SimpleClient) Close() error {
	return nil
}

func (c *SimpleClient) Use(f func(clientIntf interface{}) (interface{}, error)) (interface{}, error) {
	addr, tf, pf, cf := c.addr, c.opts.TF, c.opts.PF, c.opts.CF
	if cf == nil {
		return nil, sderr.New("nil client factory")
	}

	// make transport
	var err error
	var t thrift.TTransport
	if c.opts.Secure {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		var t1 *thrift.TSSLSocket
		t1, err = thrift.NewTSSLSocketTimeout(addr, cfg, c.opts.ConnectTimeout, c.opts.SocketTimeout)
		if err != nil {
			return nil, sderr.WithStack(err)
		}
		t = t1
	} else {
		var t1 *thrift.TSocket
		t1, err = thrift.NewTSocketTimeout(addr, c.opts.ConnectTimeout, c.opts.SocketTimeout)
		if err != nil {
			return nil, sderr.WithStack(err)
		}
		t = t1
	}
	t, err = tf.GetTransport(t)
	if err != nil {
		return nil, sderr.WithStack(err)
	}

	// 开启Transport
	err = t.Open()
	if err != nil {
		// 使用指数增长的sleep避免不停重连造成CPU过大消耗
		t := c.backOff.NextBackOff()
		if t != sdbackoff.Stop {
			time.Sleep(t)
		}
		return nil, sderr.WithStack(err)
	} else {
		// 如果成功的话将指数backOff重置xzl
		c.backOff.Reset()
	}
	defer func() {
		_ = t.Close()
	}()

	// make thrift client
	thriftClient, err := cf(t, pf)
	if err != nil {
		return nil, err
	}

	// use thrift client
	var r interface{}
	safeErr := sdcall.Safe(func() {
		r, err = f(thriftClient)
	})
	if safeErr != nil {
		return nil, sderr.WithStack(safeErr)
	}
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}
