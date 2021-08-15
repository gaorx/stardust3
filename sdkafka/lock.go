package sdkafka

import (
	"fmt"

	"github.com/gaorx/stardust3/sdtime"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type ConsumerMessageLock interface {
	TryLockMessage(kafkaTopic, group string, partition int32, offset int64) (bool, error)
}

// impl

type redisLock struct {
	timeoutSec int64
	c          redis.Cmdable
}

func NewRedisLock(c redis.Cmdable, timeoutSec int64) ConsumerMessageLock {
	return &redisLock{
		timeoutSec: timeoutSec,
		c:          c,
	}
}

func (l *redisLock) TryLockMessage(kafkaTopic, group string, partition int32, offset int64) (bool, error) {
	key := fmt.Sprintf("SDCMR:%s:%s:%d:%d", kafkaTopic, group, partition, offset)
	return l.c.SetNX(context.Background(), key, sdtime.NowUnixS(), sdtime.Seconds(l.timeoutSec)).Result()
}
