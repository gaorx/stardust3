package sdcall

import (
	"reflect"
	"sync"

	"github.com/gaorx/stardust3/sderr"
)

func Concurrent(funcs []func()) {
	if len(funcs) == 0 {
		return
	}
	var wg sync.WaitGroup
	for _, action := range funcs {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(action)
	}
	wg.Wait()
}

func ConcurrentFor(arr interface{}, f func(elem interface{})) error {
	arr1 := reflect.ValueOf(arr)
	if arr1.Kind() != reflect.Slice && arr1.Kind() != reflect.Array {
		return sderr.New("the param(arr) is not array or slice")
	}
	if arr1.Len() == 0 {
		return nil
	}
	concurrentFor(arr1, f)
	return nil
}

func concurrentFor(arr1 reflect.Value, f func(elem interface{})) {
	var wg sync.WaitGroup
	for i := 0; i < arr1.Len(); i++ {
		elem := arr1.Index(i)
		wg.Add(1)
		go func(elem interface{}) {
			defer wg.Done()
			f(elem)
		}(elem.Interface())
	}
	wg.Wait()
}
