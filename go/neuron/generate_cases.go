package neuron

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GenerateAllCases generates all pre-computed Phase 3 cases
func GenerateAllCases(outputDir string) error {
	fmt.Println("Generating Phase 3 cases...");

	cases := []struct {
		CaseID      string
		Description string
		Category    string
		Generator   func() NeuronTrainingCase
	}{
		{
			CaseID:      "sigmoid-vanishing",
			Description: "Large initialization leads to saturation and vanishing gradients",
			Category:    "saturation",
			Generator:   generateSigmoidVanishing,
		},
		{
			CaseID:      "sigmoid-optimal",
			Description: "Small initialization keeps sigmoid in active region",
			Category:    "optimal",
			Generator:   generateSigmoidOptimal,
		},
		{
			CaseID:      "relu-dying",
			Description: "Negative initialization causes ReLU to die (outputs zero forever)",
			Category:    "dying-relu",
			Generator:   generateReLUDying,
		},
		{
			CaseID:      "relu-optimal",
			Description: "Positive initialization allows ReLU to converge quickly",
			Category:    "optimal",
			Generator:   generateReLUOptimal,
		},
		{
			CaseID:      "tanh-saturation",
			Description: "Large initialization pushes tanh into saturation zones",
			Category:    "saturation",
			Generator:   generateTanhSaturation,
		},
		{
			CaseID:      "tanh-optimal",
			Description: "Centered initialization keeps tanh active and converging",
			Category:    "optimal",
			Generator:   generateTanhOptimal,
		},
		{
			CaseID:      "activation-comparison",
			Description: "Same data and initialization, different activation functions",
			Category:    "comparison",
			Generator:   generateActivationComparison,
		},
		{
			CaseID:      "lr-saturation-interaction",
			Description: "High learning rate causes oscillations into saturation",
			Category:    "saturation",
			Generator:   generateLRSaturation,
		},
	};

	for _, caseSpec := range cases {
		fmt.Printf("  Generating case: %s\n", caseSpec.CaseID);

		trainingCase := caseSpec.Generator();
		trainingCase.CaseID = caseSpec.CaseID;
		trainingCase.Description = caseSpec.Description;
		trainingCase.Category = caseSpec.Category;

		// Create output directory
		caseDir := filepath.Join(outputDir, caseSpec.CaseID);
		err := os.MkdirAll(caseDir, 0755);
		if err != nil {
			return fmt.Errorf("failed to create directory for case %s: %w", caseSpec.CaseID, err);
		}

		// Write snapshots.json
		snapshotsPath := filepath.Join(caseDir, "snapshots.json");
		err = writeJSON(snapshotsPath, trainingCase);
		if err != nil {
			return fmt.Errorf("failed to write snapshots for case %s: %w", caseSpec.CaseID, err);
		}

		fmt.Printf("    ✓ Wrote %s (%.1f KB)\n", snapshotsPath, float64(getFileSize(snapshotsPath))/1024.0);
	}

	fmt.Printf("✓ Generated %d Phase 3 cases\n", len(cases));
	return nil;
}

// Helper function to write JSON
func writeJSON(path string, data interface{}) error {
	file, err := os.Create(path);
	if err != nil {
		return err;
	}
	defer file.Close();

	encoder := json.NewEncoder(file);
	encoder.SetIndent("", "  ");
	return encoder.Encode(data);
}

// Helper function to get file size
func getFileSize(path string) int64 {
	info, err := os.Stat(path);
	if err != nil {
		return 0;
	}
	return info.Size();
}

// Case 1: Sigmoid Vanishing - Large init → saturation → gradients vanish
func generateSigmoidVanishing() NeuronTrainingCase {
	// Large initial weights push z into extreme regions
	initParams := NeuronParams{
		W: []float64{5.0, 5.0}, // Large weights
		B: 2.0,                  // Positive bias pushes z higher
	};

	// Simple linear target
	wTrue := []float64{0.5, 0.8};
	bTrue := 0.3;
	dataset := GenerateLinearDataset2D(50, wTrue, bTrue, 0.1, [2][2]float64{{-1, 1}, {-1, 1}}, 100);

	config := TrainingConfig{
		LearningRate: 0.01,
		NumSteps:     200,
		Activation:   "sigmoid",
	};

	snapshots := Train(dataset, initParams, config);

	return NeuronTrainingCase{
		Dataset:     dataset,
		InitParams:  initParams,
		FinalParams: snapshots[len(snapshots)-1].Params,
		Config:      config,
		Snapshots:   snapshots,
	};
}

// Case 2: Sigmoid Optimal - Small init → active region → smooth convergence
func generateSigmoidOptimal() NeuronTrainingCase {
	// Small initial weights keep z in active region
	initParams := NeuronParams{
		W: []float64{0.1, 0.1},
		B: 0.0,
	};

	wTrue := []float64{0.5, 0.8};
	bTrue := 0.3;
	dataset := GenerateLinearDataset2D(50, wTrue, bTrue, 0.1, [2][2]float64{{-1, 1}, {-1, 1}}, 101);

	config := TrainingConfig{
		LearningRate: 0.1, // Can use higher LR with small init
		NumSteps:     200,
		Activation:   "sigmoid",
	};

	snapshots := Train(dataset, initParams, config);

	return NeuronTrainingCase{
		Dataset:     dataset,
		InitParams:  initParams,
		FinalParams: snapshots[len(snapshots)-1].Params,
		Config:      config,
		Snapshots:   snapshots,
	};
}

// Case 3: ReLU Dying - Negative init → z < 0 → ReLU outputs zero forever
func generateReLUDying() NeuronTrainingCase {
	// Negative initialization causes z < 0 for all points
	initParams := NeuronParams{
		W: []float64{-2.0, -2.0}, // Negative weights
		B: -5.0,                   // Large negative bias
	};

	wTrue := []float64{0.5, 0.8};
	bTrue := 0.3;
	dataset := GenerateLinearDataset2D(50, wTrue, bTrue, 0.1, [2][2]float64{{-1, 1}, {-1, 1}}, 102);

	config := TrainingConfig{
		LearningRate: 0.01,
		NumSteps:     200,
		Activation:   "relu",
	};

	snapshots := Train(dataset, initParams, config);

	return NeuronTrainingCase{
		Dataset:     dataset,
		InitParams:  initParams,
		FinalParams: snapshots[len(snapshots)-1].Params,
		Config:      config,
		Snapshots:   snapshots,
	};
}

// Case 4: ReLU Optimal - Positive init → active → fast convergence
func generateReLUOptimal() NeuronTrainingCase {
	// Positive initialization keeps ReLU active
	initParams := NeuronParams{
		W: []float64{0.1, 0.1},
		B: 0.5, // Positive bias helps keep z > 0
	};

	wTrue := []float64{0.5, 0.8};
	bTrue := 0.3;
	dataset := GenerateLinearDataset2D(50, wTrue, bTrue, 0.1, [2][2]float64{{-1, 1}, {-1, 1}}, 103);

	config := TrainingConfig{
		LearningRate: 0.1,
		NumSteps:     200,
		Activation:   "relu",
	};

	snapshots := Train(dataset, initParams, config);

	return NeuronTrainingCase{
		Dataset:     dataset,
		InitParams:  initParams,
		FinalParams: snapshots[len(snapshots)-1].Params,
		Config:      config,
		Snapshots:   snapshots,
	};
}

// Case 5: Tanh Saturation - Large init → saturation (similar to sigmoid)
func generateTanhSaturation() NeuronTrainingCase {
	// Large initial weights push tanh into saturation
	initParams := NeuronParams{
		W: []float64{4.0, 4.0},
		B: 1.0,
	};

	wTrue := []float64{0.5, 0.8};
	bTrue := 0.3;
	dataset := GenerateLinearDataset2D(50, wTrue, bTrue, 0.1, [2][2]float64{{-1, 1}, {-1, 1}}, 104);

	config := TrainingConfig{
		LearningRate: 0.01,
		NumSteps:     200,
		Activation:   "tanh",
	};

	snapshots := Train(dataset, initParams, config);

	return NeuronTrainingCase{
		Dataset:     dataset,
		InitParams:  initParams,
		FinalParams: snapshots[len(snapshots)-1].Params,
		Config:      config,
		Snapshots:   snapshots,
	};
}

// Case 6: Tanh Optimal - Centered init → good convergence
func generateTanhOptimal() NeuronTrainingCase {
	// Small centered initialization
	initParams := NeuronParams{
		W: []float64{0.1, 0.1},
		B: 0.0,
	};

	wTrue := []float64{0.5, 0.8};
	bTrue := 0.3;
	dataset := GenerateLinearDataset2D(50, wTrue, bTrue, 0.1, [2][2]float64{{-1, 1}, {-1, 1}}, 105);

	config := TrainingConfig{
		LearningRate: 0.1,
		NumSteps:     200,
		Activation:   "tanh",
	};

	snapshots := Train(dataset, initParams, config);

	return NeuronTrainingCase{
		Dataset:     dataset,
		InitParams:  initParams,
		FinalParams: snapshots[len(snapshots)-1].Params,
		Config:      config,
		Snapshots:   snapshots,
	};
}

// Case 7: Activation Comparison - Same data/init, different activations
// Note: This returns sigmoid version; frontend will compare with other cases
func generateActivationComparison() NeuronTrainingCase {
	// Same initialization for all activations
	initParams := NeuronParams{
		W: []float64{0.2, 0.2},
		B: 0.1,
	};

	wTrue := []float64{0.5, 0.8};
	bTrue := 0.3;
	dataset := GenerateLinearDataset2D(50, wTrue, bTrue, 0.1, [2][2]float64{{-1, 1}, {-1, 1}}, 106);

	config := TrainingConfig{
		LearningRate: 0.05,
		NumSteps:     200,
		Activation:   "sigmoid", // This case uses sigmoid; compare with relu/tanh cases
	};

	snapshots := Train(dataset, initParams, config);

	return NeuronTrainingCase{
		Dataset:     dataset,
		InitParams:  initParams,
		FinalParams: snapshots[len(snapshots)-1].Params,
		Config:      config,
		Snapshots:   snapshots,
	};
}

// Case 8: LR-Saturation Interaction - High LR oscillates into saturation
func generateLRSaturation() NeuronTrainingCase {
	// Moderate initialization with high learning rate
	initParams := NeuronParams{
		W: []float64{1.0, 1.0},
		B: 0.5,
	};

	wTrue := []float64{0.5, 0.8};
	bTrue := 0.3;
	dataset := GenerateLinearDataset2D(50, wTrue, bTrue, 0.1, [2][2]float64{{-1, 1}, {-1, 1}}, 107);

	config := TrainingConfig{
		LearningRate: 0.5, // Very high LR causes oscillations
		NumSteps:     200,
		Activation:   "sigmoid",
	};

	snapshots := Train(dataset, initParams, config);

	return NeuronTrainingCase{
		Dataset:     dataset,
		InitParams:  initParams,
		FinalParams: snapshots[len(snapshots)-1].Params,
		Config:      config,
		Snapshots:   snapshots,
	};
}
