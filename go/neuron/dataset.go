package neuron

import (
	"math/rand"
)

// GenerateLinearDataset2D generates a 2D linear dataset: y = w_true·x + b_true + noise
func GenerateLinearDataset2D(numPoints int, wTrue []float64, bTrue float64, noiseStdDev float64, xRange [2][2]float64, seed int64) []DataPoint2DNeuron {
	rng := rand.New(rand.NewSource(seed));
	dataset := make([]DataPoint2DNeuron, numPoints);

	for i := 0; i < numPoints; i++ {
		// Generate random x values in specified ranges
		x1 := xRange[0][0] + rng.Float64()*(xRange[0][1]-xRange[0][0]);
		x2 := xRange[1][0] + rng.Float64()*(xRange[1][1]-xRange[1][0]);

		// Compute true y value
		y := wTrue[0]*x1 + wTrue[1]*x2 + bTrue;

		// Add Gaussian noise
		noise := rng.NormFloat64() * noiseStdDev;
		y += noise;

		dataset[i] = DataPoint2DNeuron{
			X: []float64{x1, x2},
			Y: y,
		};
	}

	return dataset;
}

// GenerateNonlinearDataset2D generates a dataset suitable for activation function testing
// Uses a nonlinear target function to make activation functions meaningful
func GenerateNonlinearDataset2D(numPoints int, targetFunc func([]float64) float64, noiseStdDev float64, xRange [2][2]float64, seed int64) []DataPoint2DNeuron {
	rng := rand.New(rand.NewSource(seed));
	dataset := make([]DataPoint2DNeuron, numPoints);

	for i := 0; i < numPoints; i++ {
		// Generate random x values in specified ranges
		x1 := xRange[0][0] + rng.Float64()*(xRange[0][1]-xRange[0][0]);
		x2 := xRange[1][0] + rng.Float64()*(xRange[1][1]-xRange[1][0]);

		x := []float64{x1, x2};

		// Compute target value using nonlinear function
		y := targetFunc(x);

		// Add Gaussian noise
		noise := rng.NormFloat64() * noiseStdDev;
		y += noise;

		dataset[i] = DataPoint2DNeuron{
			X: x,
			Y: y,
		};
	}

	return dataset;
}

// SimpleLinearTarget creates a simple linear target function
func SimpleLinearTarget(wTrue []float64, bTrue float64) func([]float64) float64 {
	return func(x []float64) float64 {
		y := bTrue;
		for i := 0; i < len(x); i++ {
			y += wTrue[i] * x[i];
		}
		return y;
	};
}
