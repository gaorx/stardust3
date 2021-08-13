package sdcall

import (
	"reflect"
	"time"

	"github.com/gaorx/stardust3/sderr"
)

func Retry(maxRetries int, f func() error) error {
	if f == nil {
		return nil
	}
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
	if f == nil {
		err = nil
		return
	}

	defer func() {
		if err0 := recover(); err0 != nil {
			err = sderr.ToErr(err0)
		}
	}()
	f()
	return
}

func Time(f func()) time.Duration {
	startAt := time.Now()
	if f != nil {
		f()
	}
	return time.Since(startAt)
}

func Fuse(arrayLike interface{}, f func(int, interface{})) []func() {
	if f == nil {
		f = func(int, interface{}) {}
	}
	if interfaces, ok := arrayLike.([]interface{}); ok {
		if len(interfaces) <= 0 {
			return []func(){}
		}
		funcs := make([]func(), 0, len(interfaces))
		for i, intf := range interfaces {
			i1, intf1 := i, intf
			funcs = append(funcs, func() {
				f(i1, intf1)
			})
		}
		return funcs
	} else {
		arrayVal := reflect.ValueOf(arrayLike)
		arrayTyp := arrayVal.Type()
		if arrayTyp.Kind() != reflect.Slice && arrayTyp.Kind() != reflect.Array {
			panic(&reflect.ValueError{
				Method: "reflect.Value.Index",
				Kind:   arrayTyp.Kind(),
			})
		}
		arrayLen := arrayVal.Len()
		if arrayLen <= 0 {
			return []func(){}
		}
		funcs := make([]func(), 0, arrayLen)
		for i := 0; i < arrayLen; i++ {
			i1, intf1 := i, arrayVal.Index(i).Interface()
			funcs = append(funcs, func() {
				f(i1, intf1)
			})
		}
		return funcs
	}
}
