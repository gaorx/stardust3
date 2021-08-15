package sdsort

import (
	"reflect"
	"sort"

	"github.com/gaorx/stardust3/sderr"
)

func Interface(arr interface{}, comparator Comparator) {
	if arr == nil {
		return
	}
	arrVal := reflect.ValueOf(arr)
	if arrVal.Kind() != reflect.Slice && arrVal.Kind() != reflect.Array {
		panic(sderr.WithStack(sderr.ErrIllegalType))
	}
	if arrVal.Len() <= 1 {
		return
	}
	sort.Sort(sortable{arrVal, comparator})
}

type sortable struct {
	arrayValues reflect.Value
	comparator  Comparator
}

func (s sortable) Len() int {
	arr := s.arrayValues
	return arr.Len()
}

func (s sortable) Swap(i, j int) {
	arr := s.arrayValues
	iv, jv := arr.Index(i).Interface(), arr.Index(j).Interface()
	arr.Index(i).Set(reflect.ValueOf(jv))
	arr.Index(j).Set(reflect.ValueOf(iv))
}

func (s sortable) Less(i, j int) bool {
	arr := s.arrayValues
	return s.comparator(arr.Index(i).Interface(), arr.Index(j).Interface()) < 0
}
