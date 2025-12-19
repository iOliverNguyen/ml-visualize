package neuron

import "math"

// Forward computes the forward pass for a single data point
// z = w1*x1 + w2*x2 + b
// a = σ(z)
func Forward(x []float64, params NeuronParams, activation string) (z float64, a float64) {
	// Compute pre-activation: z = w·x + b
	z = 0.0;
	for i := 0; i < len(x); i++ {
		z += params.W[i] * x[i];
	}
	z += params.B;

	// Apply activation function: a = σ(z)
	a = ApplyActivation(z, activation);

	return z, a;
}

// Loss computes the MSE loss for a single prediction: L = (a - y_true)²
func Loss(a float64, yTrue float64) float64 {
	diff := a - yTrue;
	return diff * diff;
}

// ComputeAvgLoss computes the average loss across the dataset
func ComputeAvgLoss(dataset []DataPoint2DNeuron, params NeuronParams, activation string) float64 {
	totalLoss := 0.0;
	for _, point := range dataset {
		_, a := Forward(point.X, params, activation);
		totalLoss += Loss(a, point.Y);
	}
	return totalLoss / float64(len(dataset));
}

// ComputeAvgZ computes the average pre-activation across the dataset
func ComputeAvgZ(dataset []DataPoint2DNeuron, params NeuronParams) float64 {
	totalZ := 0.0;
	for _, point := range dataset {
		z := 0.0;
		for i := 0; i < len(point.X); i++ {
			z += params.W[i] * point.X[i];
		}
		z += params.B;
		totalZ += z;
	}
	return totalZ / float64(len(dataset));
}

// ComputeAvgA computes the average post-activation across the dataset
func ComputeAvgA(dataset []DataPoint2DNeuron, params NeuronParams, activation string) float64 {
	totalA := 0.0;
	for _, point := range dataset {
		_, a := Forward(point.X, params, activation);
		totalA += a;
	}
	return totalA / float64(len(dataset));
}

// ComputeAvgDerivative computes the average σ'(z) across the dataset
func ComputeAvgDerivative(dataset []DataPoint2DNeuron, params NeuronParams, activation string) float64 {
	totalDeriv := 0.0;
	for _, point := range dataset {
		z, _ := Forward(point.X, params, activation);
		deriv := ApplyActivationDerivative(z, activation);
		totalDeriv += deriv;
	}
	return totalDeriv / float64(len(dataset));
}

// GradientMagnitude computes ||∇L|| = sqrt(sum of squared gradients)
func GradientMagnitude(grads NeuronGrads) float64 {
	sumSq := 0.0;
	for _, g := range grads.GradW {
		sumSq += g * g;
	}
	sumSq += grads.GradB * grads.GradB;
	return math.Sqrt(sumSq);
}
