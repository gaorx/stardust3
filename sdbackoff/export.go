package sdbackoff

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gaorx/stardust3/sderr"
)

type (
	BackOff = backoff.BackOff
	Ticker  = backoff.Ticker
)

var (
	Stop = backoff.Stop
)

func Const(d time.Duration) BackOff {
	if d > 0 {
		return backoff.NewConstantBackOff(d)
	} else {
		return &backoff.ZeroBackOff{}
	}
}

type ExponentialOptions struct {
	InitialInterval     time.Duration
	RandomizationFactor float64
	Multiplier          float64
	MaxInterval         time.Duration
	MaxElapsedTime      time.Duration
}

func Exponential(opts *ExponentialOptions) BackOff {
	if opts == nil {
		opts = &ExponentialOptions{}
	}
	if opts.InitialInterval <= 0 {
		opts.InitialInterval = backoff.DefaultInitialInterval
	}
	if opts.RandomizationFactor <= 0.0 {
		opts.RandomizationFactor = backoff.DefaultRandomizationFactor
	}
	if opts.Multiplier <= 0.0 {
		opts.Multiplier = backoff.DefaultMultiplier
	}
	if opts.MaxInterval <= 0 {
		opts.MaxInterval = backoff.DefaultMaxInterval
	}
	if opts.MaxElapsedTime <= 0 {
		opts.MaxElapsedTime = backoff.DefaultMaxElapsedTime
	}
	b := &backoff.ExponentialBackOff{
		InitialInterval:     opts.InitialInterval,
		RandomizationFactor: opts.RandomizationFactor,
		Multiplier:          opts.Multiplier,
		MaxInterval:         opts.MaxInterval,
		MaxElapsedTime:      opts.MaxElapsedTime,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	return b
}

func StopBack() BackOff {
	return &backoff.StopBackOff{}
}

func Zero() BackOff {
	return &backoff.ZeroBackOff{}
}

// Ticker

func TickerOf(b BackOff) *Ticker {
	return backoff.NewTicker(b)
}

// Retry

func Retry(b BackOff, f func() error) error {
	if b == nil {
		return sderr.New("nil backoff")
	}
	if f == nil {
		return sderr.New("nil action")
	}
	err := backoff.Retry(f, b)
	return sderr.WithStack(err)
}
