package linear

import (
	"math/rand"
)

// DataGenConfig2D holds configuration for 2D dataset generation
type DataGenConfig2D struct {
	NumPoints  int     `json:"num_points"`
	X1Min      float64 `json:"x1_min"`
	X1Max      float64 `json:"x1_max"`
	X2Min      float64 `json:"x2_min"`
	X2Max      float64 `json:"x2_max"`
	TrueW1     float64 `json:"true_w1"`
	TrueW2     float64 `json:"true_w2"`
	NoiseLevel float64 `json:"noise_level"`
	Seed       int64   `json:"seed"`
}

// GenerateRandomData creates a synthetic 2D linear dataset
// Data follows: y = w1*x1 + w2*x2 + noise
func GenerateRandomData(config DataGenConfig2D) []DataPoint2D {
	rng := rand.New(rand.NewSource(config.Seed))
	data := make([]DataPoint2D, config.NumPoints)

	x1Range := config.X1Max - config.X1Min
	x2Range := config.X2Max - config.X2Min

	for i := 0; i < config.NumPoints; i++ {
		// Generate random x1, x2 in specified ranges
		x1 := config.X1Min + rng.Float64()*x1Range
		x2 := config.X2Min + rng.Float64()*x2Range

		// Compute true y value
		yTrue := config.TrueW1*x1 + config.TrueW2*x2

		// Add Gaussian noise
		noise := (rng.Float64()*2 - 1) * config.NoiseLevel
		y := yTrue + noise

		data[i] = DataPoint2D{
			X1:    x1,
			X2:    x2,
			YTrue: y,
		}
	}

	return data
}
