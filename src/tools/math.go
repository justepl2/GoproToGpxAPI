package tools

import "math"

func TruncateFloat(f float64) float64 {
	return math.Floor(f*100) / 100
}
