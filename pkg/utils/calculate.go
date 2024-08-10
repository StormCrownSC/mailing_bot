package utils

import "math"

func Round(value float64, decimal int) float64 {
	power10 := math.Pow10(decimal)
	return math.Round(value*power10) / power10
}
