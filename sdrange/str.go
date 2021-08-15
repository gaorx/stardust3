package sdrange

import (
	"fmt"
)

func StrByFormat(start, stop, step int, layout string) []string {
	return StrByFunc(start, stop, step, func(i int) string {
		return fmt.Sprintf(layout, i)
	})
}

func StrByFunc(start, stop, step int, formatter func(i int) string) []string {
	if start == stop {
		return nil
	}
	a := make([]string, 0, 4)
	for i := start; i < stop; i += step {
		a = append(a, formatter(i))
	}
	return a
}
