package sdredis

import (
	"fmt"

	"github.com/gaorx/stardust3/sderr"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	redis.Cmdable
}

var _ redis.Cmdable = (*Client)(nil)

func (c *Client) String() string {
	if c.Cmdable == nil {
		return "Nil"
	}
	c1, ok := c.Cmdable.(fmt.Stringer)
	if !ok {
		return "Redis"
	}
	return c1.String()
}

func (c *Client) Close() error {
	type closeable interface {
		Close() error
	}
	if c.Cmdable == nil {
		return sderr.New("nil cmdable")
	}
	closable, ok := c.Cmdable.(closeable)
	if !ok {
		return sderr.New("not closable")
	}
	return closable.Close()
}
