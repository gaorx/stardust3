package sdcall

import (
	"sync"

	"github.com/gaorx/stardust3/sderr"
)

func Concurrent(concurrency int, funcs []func()) error {
	nFuncs := len(funcs)
	if nFuncs == 0 {
		return nil
	}
	if concurrency <= 0 {
		var wg sync.WaitGroup
		for _, f := range funcs {
			wg.Add(1)
			go func(f func()) {
				defer wg.Done()
				Safe(f)
			}(f)
		}
		wg.Wait()
		return nil
	} else {
		if concurrency > nFuncs {
			concurrency = nFuncs
		}
		pool, err := NewPool(concurrency, &PoolOptions{
			PreAlloc: true,
		})
		if err != nil {
			return sderr.WithStack(err)
		}
		defer pool.Close()
		var wg sync.WaitGroup
		for _, f := range funcs {
			f1 := f
			wg.Add(1)
			err := pool.Submit(func() {
				defer wg.Done()
				Safe(f1)
			})
			if err != nil {
				return sderr.WithStack(err)
			}
		}
		wg.Wait()
		return nil
	}
}
