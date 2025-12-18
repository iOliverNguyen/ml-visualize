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
				Steps:  200,
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
				Steps:  100,
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
				Steps:  100,
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
				Steps:  150,
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
				Steps:  200,
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
				Steps:  150,
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
				Steps:  200,
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

func generateCase(outputDir string, caseConfig CaseConfig2D) error {
	// Create case directory
	caseDir := filepath.Join(outputDir, caseConfig.ID)
	if err := os.MkdirAll(caseDir, 0755); err != nil {
		return fmt.Errorf("failed to create case directory: %w", err)
	}

	// Generate dataset
	data := GenerateRandomData(caseConfig.DataConfig)

	// Run training
	snapshots := RunTraining(data, caseConfig.TrainConfig)

	// Compute loss grid
	lossGrid := ComputeLossGrid(data, -1.0, 4.0, -1.0, 4.0, 50)

	// Write snapshots
	snapshotsPath := filepath.Join(caseDir, "snapshots.json")
	if err := writeJSON(snapshotsPath, snapshots); err != nil {
		return fmt.Errorf("failed to write snapshots: %w", err)
	}

	// Write loss grid
	lossGridPath := filepath.Join(caseDir, "loss_grid.json")
	if err := writeJSON(lossGridPath, lossGrid); err != nil {
		return fmt.Errorf("failed to write loss grid: %w", err)
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
