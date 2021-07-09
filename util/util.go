package util

import "math"

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func RoundCredits(credits int) int {
	return Max(int(math.Round(float64(credits)/10)*10), 100)
}
