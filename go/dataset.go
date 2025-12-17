package core

import (
	"fmt"
	"math"
	"math/rand"
)

// DataPoint represents a single (x, y_true) training example
type DataPoint struct {
	X     float64 `json:"x"`
	YTrue float64 `json:"y_true"`
}

// DataGenConfig configures random data generation
type DataGenConfig struct {
	NumPoints  int     `json:"num_points"`
	XMin       float64 `json:"x_min"`
	XMax       float64 `json:"x_max"`
	TrueSlope  float64 `json:"true_slope"`
	NoiseLevel float64 `json:"noise_level"`
	Seed       int64   `json:"seed"`
}

// DatasetMetadata contains metadata about a dataset
type DatasetMetadata struct {
	Name      string  `json:"name"`
	Source    string  `json:"source"` // "hardcoded", "random", "custom"
	NumPoints int     `json:"num_points"`
	TrueSlope float64 `json:"true_slope,omitempty"`
}

// Dataset wraps data points with metadata
type Dataset struct {
	Points   []DataPoint     `json:"points"`
	Metadata DatasetMetadata `json:"metadata"`
}

// GetDataset returns hardcoded training data
// Data points lie approximately on the line y = 2x with small noise
func GetDataset() []DataPoint {
	return []DataPoint{
		{X: 1.0, YTrue: 2.1},
		{X: 2.0, YTrue: 3.9},
		{X: 3.0, YTrue: 6.2},
		{X: 4.0, YTrue: 7.8},
		{X: 5.0, YTrue: 10.1},
		{X: 6.0, YTrue: 11.9},
		{X: 7.0, YTrue: 14.2},
		{X: 8.0, YTrue: 15.8},
		{X: 9.0, YTrue: 18.1},
		{X: 10.0, YTrue: 19.9},
	}
}

// GenerateRandomData creates random training data with configurable parameters
// Data follows: y = trueSlope * x + noise
func GenerateRandomData(config DataGenConfig) ([]DataPoint, error) {
	// Validate config
	if config.NumPoints <= 0 {
		return nil, fmt.Errorf("num_points must be positive, got %d", config.NumPoints)
	}
	if config.XMax <= config.XMin {
		return nil, fmt.Errorf("x_max must be greater than x_min")
	}
	if config.NoiseLevel < 0 {
		return nil, fmt.Errorf("noise_level must be non-negative, got %f", config.NoiseLevel)
	}

	// Set up random number generator
	rng := rand.New(rand.NewSource(config.Seed))
	if config.Seed == 0 {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	// Generate data points
	points := make([]DataPoint, config.NumPoints)
	xRange := config.XMax - config.XMin

	for i := 0; i < config.NumPoints; i++ {
		// Evenly spaced x values
		x := config.XMin + (xRange * float64(i) / float64(config.NumPoints-1))

		// y = slope * x + noise
		noise := (rng.Float64()*2 - 1) * config.NoiseLevel // Uniform noise in [-noiseLevel, +noiseLevel]
		y := config.TrueSlope*x + noise

		points[i] = DataPoint{
			X:     x,
			YTrue: y,
		}
	}

	return points, nil
}

// ValidateDataset checks if a dataset is valid for training
func ValidateDataset(data []DataPoint) error {
	if len(data) == 0 {
		return fmt.Errorf("dataset is empty")
	}

	for i, point := range data {
		if math.IsNaN(point.X) || math.IsInf(point.X, 0) {
			return fmt.Errorf("point %d has invalid X value: %f", i, point.X)
		}
		if math.IsNaN(point.YTrue) || math.IsInf(point.YTrue, 0) {
			return fmt.Errorf("point %d has invalid YTrue value: %f", i, point.YTrue)
		}
	}

	return nil
}
