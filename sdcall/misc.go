package sdcall

import (
	"github.com/gaorx/stardust3/sderr"
)

func Retry(maxRetries int, f func() error) error {
	err0 := f()
	if err0 == nil {
		return nil
	}
	err := sderr.Append(err0)
	for i := 1; i <= maxRetries; i++ {
		err0 := f()
		if err0 != nil {
			err = sderr.Append(err, err0)
		} else {
			return nil
		}
	}
	return err
}

func Safe(f func()) (err error) {
	defer func() {
		if err0 := recover(); err0 != nil {
			err = sderr.ToErr(err0)
		}
	}()

	f()
	return
}
