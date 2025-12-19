package neuron

import "math"

// Sigmoid computes the sigmoid activation function: σ(z) = 1 / (1 + e^(-z))
func Sigmoid(z float64) float64 {
	return 1.0 / (1.0 + math.Exp(-z));
}

// SigmoidDerivative computes the derivative of sigmoid: σ'(z) = σ(z) * (1 - σ(z))
func SigmoidDerivative(z float64) float64 {
	s := Sigmoid(z);
	return s * (1.0 - s);
}

// ReLU computes the ReLU activation function: max(0, z)
func ReLU(z float64) float64 {
	return math.Max(0, z);
}

// ReLUDerivative computes the derivative of ReLU: 1 if z > 0, else 0
func ReLUDerivative(z float64) float64 {
	if z > 0 {
		return 1.0;
	}
	return 0.0;
}

// Tanh computes the hyperbolic tangent activation function: tanh(z)
func Tanh(z float64) float64 {
	return math.Tanh(z);
}

// TanhDerivative computes the derivative of tanh: 1 - tanh²(z)
func TanhDerivative(z float64) float64 {
	t := math.Tanh(z);
	return 1.0 - t*t;
}

// IsSaturated checks if the activation function is in saturation zone
// Threshold: |σ'(z)| < 0.01 means gradient is less than 1% of maximum
func IsSaturated(sigmaPrime float64) bool {
	return math.Abs(sigmaPrime) < 0.01;
}

// ApplyActivation applies the specified activation function to z
func ApplyActivation(z float64, activation string) float64 {
	switch activation {
	case "sigmoid":
		return Sigmoid(z);
	case "relu":
		return ReLU(z);
	case "tanh":
		return Tanh(z);
	default:
		return Sigmoid(z); // default to sigmoid
	}
}

// ApplyActivationDerivative applies the derivative of the specified activation function
func ApplyActivationDerivative(z float64, activation string) float64 {
	switch activation {
	case "sigmoid":
		return SigmoidDerivative(z);
	case "relu":
		return ReLUDerivative(z);
	case "tanh":
		return TanhDerivative(z);
	default:
		return SigmoidDerivative(z); // default to sigmoid
	}
}
