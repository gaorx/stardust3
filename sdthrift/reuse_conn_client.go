package sdthrift

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	connpool "github.com/fatih/pool"
	"github.com/gaorx/stardust3/sdbackoff"
	"github.com/gaorx/stardust3/sdcall"
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdtime"
)

type ReuseConnClient struct {
	p               connpool.Pool
	initCap, maxCap int
	opts            ClientOptions
	backOff         sdbackoff.BackOff
}

func NewReuseConnClient(initCap, maxCap int) *ReuseConnClient {
	if initCap < 0 {
		initCap = 0
	}
	if maxCap < 0 {
		maxCap = 1
	}
	backOff := sdbackoff.Synchronized(sdbackoff.Exponential(&sdbackoff.ExponentialOptions{
		InitialInterval:     sdtime.Millis(5),
		RandomizationFactor: 0.0,
		Multiplier:          1.8,
		MaxInterval:         sdtime.Millis(40),
	}))
	return &ReuseConnClient{
		initCap: initCap,
		maxCap:  maxCap,
		backOff: backOff,
	}
}

func makeConn(addr string, opts ClientOptions, backOff sdbackoff.BackOff) (net.Conn, error) {
	var err error
	var conn net.Conn
	if opts.Secure {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		conn, err = tls.Dial("tcp", addr, cfg)
	} else {
		var tcpAddr *net.TCPAddr
		tcpAddr, err = net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return nil, err
		}
		conn, err = net.DialTimeout(tcpAddr.Network(), tcpAddr.String(), opts.ConnectTimeout)
	}
	// 使用指数增长的sleep避免不停重连造成CPU过大消耗
	if err != nil {
		t := backOff.NextBackOff()
		if t != sdbackoff.Stop {
			time.Sleep(t)
		}
	} else {
		backOff.Reset()
	}
	return conn, err
}

func (c *ReuseConnClient) Connect(addr string, opts ClientOptions) error {
	if c.p != nil {
		return sderr.New("opened")
	}
	p, err := connpool.NewChannelPool(c.initCap, c.maxCap, func() (net.Conn, error) {
		return makeConn(addr, opts, c.backOff)
	})
	if err != nil {
		return err
	}
	c.p = p
	c.opts = opts
	return nil
}

func (c *ReuseConnClient) Close() error {
	if c.p != nil {
		c.p.Close()
	}
	return nil
}

func (c *ReuseConnClient) NumConn() int {
	return c.p.Len()
}

func (c *ReuseConnClient) Use(f func(clientIntf interface{}) (interface{}, error)) (interface{}, error) {
	tf, pf, cf := c.opts.TF, c.opts.PF, c.opts.CF
	if cf == nil {
		return nil, sderr.New("no client factory")
	}

	var err error
	var t thrift.TTransport = &connTransport{
		p:       c.p,
		timeout: c.opts.SocketTimeout,
	}
	t, err = tf.GetTransport(t)
	if err != nil {
		return nil, sderr.WithStack(err)
	}

	// open transport
	err = t.Open()
	if err != nil {
		return nil, err
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

var (
	errNotOpen = sderr.Sentinel("not open")
)

type connTransport struct {
	p       connpool.Pool
	conn    net.Conn
	timeout time.Duration
}

func (t *connTransport) Open() error {
	conn, err := t.p.Get()
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *connTransport) Close() error {
	if !t.IsOpen() {
		return sderr.WithStack(errNotOpen)
	}
	err := t.conn.Close()
	t.conn = nil
	return err
}

func (t *connTransport) IsOpen() bool {
	return t.conn != nil
}

func (t *connTransport) Read(p []byte) (int, error) {
	if !t.IsOpen() {
		return 0, sderr.WithStack(errNotOpen)
	}
	err := t.setDeadline(true, false)
	if err != nil {
		return 0, sderr.WithStack(err)
	}
	n, err := t.conn.Read(p)
	if err != nil {
		markUnusable(t.conn)
	}
	return n, err
}

func (t *connTransport) Write(p []byte) (int, error) {
	if !t.IsOpen() {
		return 0, sderr.WithStack(errNotOpen)
	}
	err := t.setDeadline(false, true)
	if err != nil {
		return 0, sderr.WithStack(err)
	}
	n, err := t.conn.Write(p)
	if err != nil {
		markUnusable(t.conn)
	}
	return n, err
}

func markUnusable(conn net.Conn) {
	if pconn, ok := conn.(*connpool.PoolConn); ok {
		pconn.MarkUnusable()
	}
}

func (t *connTransport) Flush(ctx context.Context) error {
	return nil
}

func (t *connTransport) RemainingBytes() uint64 {
	const maxSize = ^uint64(0)
	return maxSize
}

func (t *connTransport) setDeadline(read, write bool) error {
	var deadline time.Time
	if t.timeout > 0 {
		deadline = time.Now().Add(t.timeout)
	}
	if read && write {
		return t.conn.SetDeadline(deadline)
	} else if read {
		return t.conn.SetReadDeadline(deadline)
	} else if write {
		return t.conn.SetWriteDeadline(deadline)
	} else {
		return nil
	}
}
