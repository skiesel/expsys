package tables

import (
	"math"
)

func StdDevAndVariance(values []float64) (stddev, variance float64) {
	average := 0.
	for i := range values {
		average += values[i]
	}
	average /= float64(len(values))

	for i := range values {
		diff := values[i] - average
		variance += diff * diff
	}

	variance /= float64(len(values))
	stddev = math.Sqrt(variance)
	return
}