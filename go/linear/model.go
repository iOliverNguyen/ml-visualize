package linear

import "math"

// Forward computes the prediction: y_pred = w1*x1 + w2*x2
func Forward(w1, w2, x1, x2 float64) float64 {
	return w1*x1 + w2*x2
}

// Loss computes the mean squared error: (y_pred - y_true)^2
func Loss(yPred, yTrue float64) float64 {
	diff := yPred - yTrue
	return diff * diff
}

// GradW1 computes the gradient of loss with respect to w1
// Derivation:
//   loss = (w1*x1 + w2*x2 - y_true)^2
//   ∂(loss)/∂w1 = 2 * (w1*x1 + w2*x2 - y_true) * ∂(w1*x1 + w2*x2)/∂w1
//                = 2 * (y_pred - y_true) * x1
func GradW1(w1, w2, x1, x2, yTrue float64) float64 {
	yPred := Forward(w1, w2, x1, x2)
	return 2 * (yPred - yTrue) * x1
}

// GradW2 computes the gradient of loss with respect to w2
// Derivation:
//   loss = (w1*x1 + w2*x2 - y_true)^2
//   ∂(loss)/∂w2 = 2 * (w1*x1 + w2*x2 - y_true) * ∂(w1*x1 + w2*x2)/∂w2
//                = 2 * (y_pred - y_true) * x2
func GradW2(w1, w2, x1, x2, yTrue float64) float64 {
	yPred := Forward(w1, w2, x1, x2)
	return 2 * (yPred - yTrue) * x2
}

// GradientMagnitude computes the L2 norm of the gradient vector
// ||∇L|| = sqrt(grad_w1² + grad_w2²)
func GradientMagnitude(gradW1, gradW2 float64) float64 {
	return math.Sqrt(gradW1*gradW1 + gradW2*gradW2)
}

// GradientDirection computes the angle of the gradient vector in radians
// θ = atan2(grad_w2, grad_w1)
// Returns angle in range [-π, π]
func GradientDirection(gradW1, gradW2 float64) float64 {
	return math.Atan2(gradW2, gradW1)
}
