package adamqp

import (
	"github.com/gaorx/stardust3/sderr"
	"github.com/streadway/amqp"
)

type ChannelConn struct {
	Chan *amqp.Channel
	Conn *amqp.Connection
}

func (cc *ChannelConn) Close() error {
	var chanErr, connErr error
	if cc.Chan != nil {
		chanErr = cc.Chan.Close()
		cc.Chan = nil
	}
	if cc.Conn != nil {
		connErr = cc.Conn.Close()
		cc.Conn = nil
	}
	if chanErr != nil {
		return sderr.WithStack(chanErr)
	} else if connErr != nil {
		return sderr.WithStack(connErr)
	} else {
		return nil
	}
}
