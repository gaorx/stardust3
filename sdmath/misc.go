package sdmath

// 归一化, 将value从[min, max]区间映射到[newMin, newMax]区间,不做参数检查
func Normalize(min, max, value, newMin, newMax float64) float64 {
	if value > max {
		value = max
	} else if value < min {
		value = min
	}

	return (value-min)/(max-min)*newMax + newMin
}
