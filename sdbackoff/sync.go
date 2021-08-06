package sdbackoff

import (
	"sync"
	"time"

	"github.com/gaorx/stardust3/sderr"
)

type syncBackOff struct {
	backOff BackOff
	mtx     sync.Mutex
}

func Synchronized(b BackOff) BackOff {
	if b == nil {
		panic(sderr.New("nil backoff"))
	}
	return &syncBackOff{backOff: b}
}

func (b *syncBackOff) NextBackOff() time.Duration {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	return b.backOff.NextBackOff()
}

func (b *syncBackOff) Reset() {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.backOff.Reset()
}
