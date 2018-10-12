package util

import (
	"math"
)

// 整数指数求幂
func Power(x float64, n int) float64 {
	ans := 1.0

	for n != 0 {
		ans *= x
		n--
	}
	return ans
}

// golang 标准库的求幂
func Powerf(x, n float64) float64 {
	return math.Pow(x, n)
}
