package linear

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CaseConfig2D represents metadata for a Phase 2 case
type CaseConfig2D struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Emoji       string           `json:"emoji"`
	Category    string           `json:"category"`
	DataConfig  DataGenConfig2D  `json:"data_config"`
	TrainConfig TrainingConfig2D `json:"training_config"`
	Insights    []string         `json:"insights"`
}

// CaseManifest2D represents the manifest of all Phase 2 cases
type CaseManifest2D struct {
	Version string         `json:"version"`
	Cases   []CaseConfig2D `json:"cases"`
}

// GenerateCases2D generates all pre-computed Phase 2 cases
func GenerateCases2D(outputDir string) error {
	cases := []CaseConfig2D{
		{
			ID:          "lr-small",
			Name:        "Learning Rate Too Small",
			Description: "Tiny steps make convergence painfully slow. Watch gradient descent crawl toward the optimum.",
			Emoji:       "🐌",
			Category:    "learning-rate",
			DataConfig: DataGenConfig2D{
				NumPoints:  20,
				X1Min:      0.0,
				X1Max:      5.0,
				X2Min:      0.0,
				X2Max:      5.0,
				TrueW1:     2.0,
				TrueW2:     1.5,
				NoiseLevel: 0.5,
				Seed:       42,
			},
			TrainConfig: TrainingConfig2D{
				W1Init: 0.0,
				W2Init: 0.0,
				LR:     0.0001,
				MaxSteps: 200,
			},
			Insights: []string{
				"Tiny steps: learning rate = 0.0001",
				"Converges very slowly, never reaches optimum",
				"Gradient magnitude stays large even after 200 steps",
			},
		},
		{
			ID:          "lr-optimal",
			Name:        "Optimal Learning Rate",
			Description: "The Goldilocks zone: smooth, efficient steps toward convergence. This is what you want.",
			Emoji:       "✓",
			Category:    "learning-rate",
			DataConfig: DataGenConfig2D{
				NumPoints:  20,
				X1Min:      0.0,
				X1Max:      5.0,
				X2Min:      0.0,
				X2Max:      5.0,
				TrueW1:     2.0,
				TrueW2:     1.5,
				NoiseLevel: 0.5,
				Seed:       42,
			},
			TrainConfig: TrainingConfig2D{
				W1Init: 0.0,
				W2Init: 0.0,
				LR:     0.01,
				MaxSteps: 100,
			},
			Insights: []string{
				"Smooth convergence in ~50 steps",
				"Gradient magnitude decreases steadily",
				"Direct path to optimum with minimal oscillation",
			},
		},
		{
			ID:          "lr-large",
			Name:        "Learning Rate Too Large",
			Description: "Big steps overshoot the target. Watch the zigzag pattern as optimization bounces around.",
			Emoji:       "↔️",
			Category:    "learning-rate",
			DataConfig: DataGenConfig2D{
				NumPoints:  20,
				X1Min:      0.0,
				X1Max:      5.0,
				X2Min:      0.0,
				X2Max:      5.0,
				TrueW1:     2.0,
				TrueW2:     1.5,
				NoiseLevel: 0.5,
				Seed:       42,
			},
			TrainConfig: TrainingConfig2D{
				W1Init: 0.0,
				W2Init: 0.0,
				LR:     0.1,
				MaxSteps: 100,
			},
			Insights: []string{
				"Overshooting causes zigzag pattern",
				"Eventually converges but inefficiently",
				"Large oscillations in parameter space",
			},
		},
		{
			ID:          "anisotropic-easy",
			Name:        "Anisotropic Loss Surface (Mild)",
			Description: "Different scales in x1 vs x2 create an elliptical loss surface. Faster progress in one direction.",
			Emoji:       "⬭",
			Category:    "anisotropy",
			DataConfig: DataGenConfig2D{
				NumPoints:  20,
				X1Min:      0.0,
				X1Max:      1.0,
				X2Min:      0.0,
				X2Max:      10.0,
				TrueW1:     2.0,
				TrueW2:     1.5,
				NoiseLevel: 0.5,
				Seed:       42,
			},
			TrainConfig: TrainingConfig2D{
				W1Init: 0.0,
				W2Init: 0.0,
				LR:     0.01,
				MaxSteps: 150,
			},
			Insights: []string{
				"Elliptical contours due to different x1/x2 scales",
				"Faster movement in w2 direction",
				"Still converges with standard learning rate",
			},
		},
		{
			ID:          "anisotropic-hard",
			Name:        "Anisotropic Loss Surface (Extreme)",
			Description: "Extreme scale difference creates a narrow valley. Optimization struggles with zigzag motion.",
			Emoji:       "⬬",
			Category:    "anisotropy",
			DataConfig: DataGenConfig2D{
				NumPoints:  20,
				X1Min:      0.0,
				X1Max:      1.0,
				X2Min:      0.0,
				X2Max:      20.0,
				TrueW1:     2.0,
				TrueW2:     0.5,
				NoiseLevel: 0.3,
				Seed:       42,
			},
			TrainConfig: TrainingConfig2D{
				W1Init: 0.0,
				W2Init: 0.0,
				LR:     0.008,
				MaxSteps: 200,
			},
			Insights: []string{
				"Very elongated ellipse creates narrow valley",
				"Zigzag pattern even with reduced learning rate",
				"Demonstrates why feature scaling matters",
			},
		},
		{
			ID:          "saddle-point",
			Name:        "Near-Saddle Geometry",
			Description: "Starting near a saddle-like region shows curved, non-direct trajectory to optimum.",
			Emoji:       "〰️",
			Category:    "geometry",
			DataConfig: DataGenConfig2D{
				NumPoints:  20,
				X1Min:      -2.0,
				X1Max:      3.0,
				X2Min:      -2.0,
				X2Max:      3.0,
				TrueW1:     2.0,
				TrueW2:     1.5,
				NoiseLevel: 0.5,
				Seed:       42,
			},
			TrainConfig: TrainingConfig2D{
				W1Init: -1.0,
				W2Init: -1.0,
				LR:     0.01,
				MaxSteps: 150,
			},
			Insights: []string{
				"Gradient direction changes rapidly",
				"Curved trajectory through parameter space",
				"Starting position affects convergence path",
			},
		},
		{
			ID:          "zigzag-convergence",
			Name:        "Zigzag with High LR + Anisotropy",
			Description: "Combined effect: high learning rate meets anisotropic surface. Maximum zigzag demonstration.",
			Emoji:       "⚡",
			Category:    "geometry",
			DataConfig: DataGenConfig2D{
				NumPoints:  20,
				X1Min:      0.0,
				X1Max:      1.0,
				X2Min:      0.0,
				X2Max:      15.0,
				TrueW1:     2.0,
				TrueW2:     0.8,
				NoiseLevel: 0.4,
				Seed:       42,
			},
			TrainConfig: TrainingConfig2D{
				W1Init: 3.0,
				W2Init: -1.5,
				LR:     0.03,
				MaxSteps: 200,
			},
			Insights: []string{
				"Bouncing back and forth dramatically",
				"High LR + anisotropy = worst case scenario",
				"Slow progress despite high learning rate",
			},
		},
	}

	// Create manifest
	manifest := CaseManifest2D{
		Version: "1.0",
		Cases:   cases,
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write manifest
	manifestPath := filepath.Join(outputDir, "manifest.json")
	if err := writeJSON(manifestPath, manifest); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	fmt.Printf("Generated Phase 2 manifest at %s\n", manifestPath)

	// Generate each case
	for _, caseConfig := range cases {
		if err := generateCase(outputDir, caseConfig); err != nil {
			return fmt.Errorf("failed to generate case %s: %w", caseConfig.ID, err)
		}
		fmt.Printf("Generated case: %s\n", caseConfig.ID)
	}

	return nil
}

// CaseConfigJSON represents the minimal config file for client-side training
type CaseConfigJSON struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	DataConfig     DataGenConfig2D  `json:"data_config"`
	TrainingConfig TrainingConfig2D `json:"training_config"`
	LossGridConfig LossGridConfig   `json:"loss_grid_config"`
}

// LossGridConfig holds parameters for loss grid computation
type LossGridConfig struct {
	W1Min      float64 `json:"w1_min"`
	W1Max      float64 `json:"w1_max"`
	W2Min      float64 `json:"w2_min"`
	W2Max      float64 `json:"w2_max"`
	Resolution int     `json:"resolution"`
}

func generateCase(outputDir string, caseConfig CaseConfig2D) error {
	// Create case directory
	caseDir := filepath.Join(outputDir, caseConfig.ID)
	if err := os.MkdirAll(caseDir, 0755); err != nil {
		return fmt.Errorf("failed to create case directory: %w", err)
	}

	// Create minimal config file for client-side training
	config := CaseConfigJSON{
		Name:           caseConfig.Name,
		Description:    caseConfig.Description,
		DataConfig:     caseConfig.DataConfig,
		TrainingConfig: caseConfig.TrainConfig,
		LossGridConfig: LossGridConfig{
			W1Min:      -1.0,
			W1Max:      4.0,
			W2Min:      -1.0,
			W2Max:      4.0,
			Resolution: 50,
		},
	}

	// Write config file (~5KB)
	configPath := filepath.Join(caseDir, "config.json")
	if err := writeJSON(configPath, config); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func writeJSON(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
