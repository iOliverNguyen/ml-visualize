package linear

// DataPoint2D represents a training example with two input features
type DataPoint2D struct {
	X1    float64 `json:"x1"`
	X2    float64 `json:"x2"`
	YTrue float64 `json:"y_true"`
}

// PointSnapshot2D captures per-point breakdown at a training step
type PointSnapshot2D struct {
	X1        float64 `json:"x1"`
	X2        float64 `json:"x2"`
	YTrue     float64 `json:"y_true"`
	YPred     float64 `json:"y_pred"`
	PointLoss float64 `json:"point_loss"`
	GradW1    float64 `json:"grad_w1"`
	GradW2    float64 `json:"grad_w2"`
}

// UpdateDetails2D captures parameter update breakdown
type UpdateDetails2D struct {
	W1Old   float64 `json:"w1_old"`
	W2Old   float64 `json:"w2_old"`
	LR      float64 `json:"lr"`
	GradW1  float64 `json:"grad_w1"`
	GradW2  float64 `json:"grad_w2"`
	DeltaW1 float64 `json:"delta_w1"`
	DeltaW2 float64 `json:"delta_w2"`
	W1New   float64 `json:"w1_new"`
	W2New   float64 `json:"w2_new"`
}

// LinearSnapshot captures complete state at one training step
type LinearSnapshot struct {
	Step              int                 `json:"step"`
	W1                float64             `json:"w1"`
	W2                float64             `json:"w2"`
	GradW1            float64             `json:"grad_w1"`
	GradW2            float64             `json:"grad_w2"`
	Loss              float64             `json:"loss"`
	GradientMagnitude float64             `json:"gradient_magnitude"`
	GradientDirection float64             `json:"gradient_direction"`
	PointDetails      []PointSnapshot2D   `json:"point_details"`
	UpdateComponents  UpdateDetails2D     `json:"update_components"`
}
