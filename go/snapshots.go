package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Snapshot captures the complete state at one training step
type Snapshot struct {
	Step  int     `json:"step"`
	W     float64 `json:"w"`
	GradW float64 `json:"grad_w"`
	Loss  float64 `json:"loss"`
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
