package sdchan

import (
	"reflect"
)

func ReceiveSelect(chs []interface{}) (int, interface{}, bool) {
	if len(chs) == 0 {
		return -1, nil, false
	}
	cl := make([]reflect.SelectCase, 0, len(chs))
	for _, c := range chs {
		cl = append(cl, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c),
		})
	}
	index, v, ok := reflect.Select(cl)
	return index, v.Interface(), ok
}
