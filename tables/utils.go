package tables

import (
	"math"
)

func MeanStdDevVariance(values []float64) (mean, stddev, variance float64) {
	for i := range values {
		mean += values[i]
	}
	mean /= float64(len(values))

	for i := range values {
		diff := values[i] - mean
		variance += diff * diff
	}

	variance /= float64(len(values))
	stddev = math.Sqrt(variance)
	return
}

func Sum(values []float64) (sum float64) {
	for i := range values {
		sum += values[i]
	}
	return
}