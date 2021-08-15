package sdrange

func Int(start, stop, step int) []int {
	if start == stop {
		return nil
	}
	a := make([]int, 0, 4)
	for i := start; i < stop; i += step {
		a = append(a, i)
	}
	return a
}

func Int64(start, stop, step int64) []int64 {
	if start == stop {
		return nil
	}
	a := make([]int64, 0, 4)
	for i := start; i < stop; i += step {
		a = append(a, i)
	}
	return a
}
