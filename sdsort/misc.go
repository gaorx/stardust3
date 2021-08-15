package sdsort

import (
	"github.com/gaorx/stardust3/sderr"
)

func Reverse(c Comparator) Comparator {
	if c == nil {
		panic(sderr.New("nil comparator"))
	}
	return func(a interface{}, b interface{}) int {
		return -c(a, b) // or c(b, a)
	}
}
