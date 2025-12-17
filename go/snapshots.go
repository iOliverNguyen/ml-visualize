package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// PointSnapshot captures per-point breakdown for pedagogical inspection
type PointSnapshot struct {
	X         float64 `json:"x"`
	YTrue     float64 `json:"y_true"`
	YPred     float64 `json:"y_pred"`
	PointLoss float64 `json:"point_loss"` // (y_pred - y_true)^2
	PointGrad float64 `json:"point_grad"` // 2(y_pred - y_true)*x
}

// UpdateDetails captures parameter update breakdown for pedagogy
type UpdateDetails struct {
	WOld   float64 `json:"w_old"`
	LR     float64 `json:"lr"`
	GradW  float64 `json:"grad_w"`
	DeltaW float64 `json:"delta_w"` // -lr * grad_w
	WNew   float64 `json:"w_new"`
}

// Snapshot captures the complete state at one training step
type Snapshot struct {
	Step  int     `json:"step"`
	W     float64 `json:"w"`
	GradW float64 `json:"grad_w"`
	Loss  float64 `json:"loss"`

	// NEW: Per-point details for inspection
	PointDetails []PointSnapshot `json:"point_details"`

	// NEW: Update breakdown for pedagogy
	UpdateComponents UpdateDetails `json:"update_components"`
}

// WriteSnapshots marshals snapshots to JSON and writes to file
func WriteSnapshots(snapshots []Snapshot, filePath string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(snapshots, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
