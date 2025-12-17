package core

// Forward computes the prediction: y_pred = w * x
func Forward(w, x float64) float64 {
	return w * x
}

// Loss computes the mean squared error: (y_pred - y_true)^2
func Loss(yPred, yTrue float64) float64 {
	diff := yPred - yTrue
	return diff * diff
}

// GradW computes the gradient of loss with respect to w
// Derivation:
//   loss = (w*x - y_true)^2
//   d(loss)/dw = 2 * (w*x - y_true) * d(w*x)/dw
//              = 2 * (w*x - y_true) * x
func GradW(w, x, yTrue float64) float64 {
	yPred := Forward(w, x)
	return 2 * (yPred - yTrue) * x
}
