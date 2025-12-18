package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CaseConfig defines a training scenario for the case library
type CaseConfig struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Emoji       string         `json:"emoji"`
	Category    string         `json:"category"`
	DataConfig  DataGenConfig  `json:"data_config"`
	Training    TrainingConfig `json:"training_config"`
	Insights    []string       `json:"insights"`
}

// CaseManifest contains metadata for all cases
type CaseManifest struct {
	Version string       `json:"version"`
	Cases   []CaseConfig `json:"cases"`
}

// GenerateCases generates 8 pre-computed training cases for the case library
func GenerateCases() error {
	// Define the 8 beginner cases
	cases := []CaseConfig{
		// Foundational cases
		{
			ID:          "perfect-start",
			Name:        "Perfect Start",
			Description: "Clean data, good learning rate, ideal initialization. Everything goes smoothly!",
			Emoji:       "✓",
			Category:    "foundational",
			DataConfig: DataGenConfig{
				NumPoints:  10,
				XMin:       1.0,
				XMax:       10.0,
				TrueSlope:  2.0,
				NoiseLevel: 0.1,
				Seed:       42,
			},
			Training: TrainingConfig{
				WInit: 0.0,
				LR:    0.001,
				Steps: 100,
			},
			Insights: []string{
				"Loss decreases smoothly every step",
				"Line quickly moves toward the data points",
				"Final w ≈ 2.0 matches the true slope",
			},
		},
		{
			ID:          "noisy-but-ok",
			Name:        "Noisy But OK",
			Description: "Moderate noise in the data. Training still works, but loss won't reach zero.",
			Emoji:       "✓",
			Category:    "foundational",
			DataConfig: DataGenConfig{
				NumPoints:  10,
				XMin:       1.0,
				XMax:       10.0,
				TrueSlope:  2.0,
				NoiseLevel: 0.5,
				Seed:       123,
			},
			Training: TrainingConfig{
				WInit: 0.0,
				LR:    0.001,
				Steps: 100,
			},
			Insights: []string{
				"Loss decreases but stabilizes above zero",
				"Noise prevents perfect fit",
				"This is normal in real-world data",
			},
		},
		{
			ID:          "very-noisy",
			Name:        "Very Noisy",
			Description: "High noise makes training harder. Loss decreases slowly and stays high.",
			Emoji:       "⚠",
			Category:    "foundational",
			DataConfig: DataGenConfig{
				NumPoints:  10,
				XMin:       1.0,
				XMax:       10.0,
				TrueSlope:  2.0,
				NoiseLevel: 2.0,
				Seed:       456,
			},
			Training: TrainingConfig{
				WInit: 0.0,
				LR:    0.001,
				Steps: 100,
			},
			Insights: []string{
				"Loss decreases but remains high",
				"Line struggles to find pattern in noisy data",
				"May need more data or different model",
			},
		},

		// Learning rate cases
		{
			ID:          "lr-too-slow",
			Name:        "Learning Too Slow",
			Description: "Learning rate is too small. Training barely moves!",
			Emoji:       "😴",
			Category:    "learning-rate",
			DataConfig: DataGenConfig{
				NumPoints:  10,
				XMin:       1.0,
				XMax:       10.0,
				TrueSlope:  2.0,
				NoiseLevel: 0.1,
				Seed:       42,
			},
			Training: TrainingConfig{
				WInit: 0.0,
				LR:    0.00001,
				Steps: 200,
			},
			Insights: []string{
				"Barely moving after 200 steps",
				"Loss decreases extremely slowly",
				"Would need thousands of steps to converge",
			},
		},
		{
			ID:          "lr-just-right",
			Name:        "Learning Just Right",
			Description: "Perfect learning rate. Fast convergence without overshooting.",
			Emoji:       "✓",
			Category:    "learning-rate",
			DataConfig: DataGenConfig{
				NumPoints:  10,
				XMin:       1.0,
				XMax:       10.0,
				TrueSlope:  2.0,
				NoiseLevel: 0.1,
				Seed:       42,
			},
			Training: TrainingConfig{
				WInit: 0.5,
				LR:    0.002,
				Steps: 100,
			},
			Insights: []string{
				"Converges in ~50 steps",
				"Smooth, steady improvement",
				"This is what good training looks like",
			},
		},
		{
			ID:          "lr-too-fast",
			Name:        "Learning Too Fast",
			Description: "Learning rate is too large. Watch the line bounce around!",
			Emoji:       "⚠",
			Category:    "learning-rate",
			DataConfig: DataGenConfig{
				NumPoints:  10,
				XMin:       1.0,
				XMax:       10.0,
				TrueSlope:  2.0,
				NoiseLevel: 0.1,
				Seed:       42,
			},
			Training: TrainingConfig{
				WInit: 0.8,
				LR:    0.008,
				Steps: 100,
			},
			Insights: []string{
				"Converges quickly but overshoots",
				"Shows bouncing behavior in first steps",
				"Eventually stabilizes but less smoothly",
			},
		},

		// Initialization cases
		{
			ID:          "start-at-zero",
			Name:        "Start at Zero",
			Description: "Initialize w=0. Line starts completely flat.",
			Emoji:       "✓",
			Category:    "initialization",
			DataConfig: DataGenConfig{
				NumPoints:  10,
				XMin:       1.0,
				XMax:       10.0,
				TrueSlope:  2.0,
				NoiseLevel: 0.1,
				Seed:       42,
			},
			Training: TrainingConfig{
				WInit: 0.0,
				LR:    0.001,
				Steps: 100,
			},
			Insights: []string{
				"Starts with flat line (w=0)",
				"Gradually tilts upward toward data",
				"Common initialization choice",
			},
		},
		{
			ID:          "start-far-away",
			Name:        "Start Far Away",
			Description: "Initialize w=-3. Line starts pointing the wrong direction!",
			Emoji:       "⚠",
			Category:    "initialization",
			DataConfig: DataGenConfig{
				NumPoints:  10,
				XMin:       1.0,
				XMax:       10.0,
				TrueSlope:  2.0,
				NoiseLevel: 0.1,
				Seed:       42,
			},
			Training: TrainingConfig{
				WInit: -3.0,
				LR:    0.001,
				Steps: 150,
			},
			Insights: []string{
				"Starts with negative slope",
				"Takes longer to reach target",
				"Bad initialization wastes training time",
			},
		},
	};

	// Create output directory
	outputDir := "js/public/cases";
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err);
	}

	// Generate snapshots for each case
	for _, caseConfig := range cases {
		fmt.Printf("Generating case: %s (%s)\n", caseConfig.Name, caseConfig.ID);

		// Generate data
		data, err := GenerateRandomData(caseConfig.DataConfig);
		if err != nil {
			return fmt.Errorf("failed to generate data for case %s: %w", caseConfig.ID, err);
		}

		// Run training
		snapshots := RunTrainingWithDataset(data, caseConfig.Training);

		// Create case directory
		caseDir := filepath.Join(outputDir, caseConfig.ID);
		if err := os.MkdirAll(caseDir, 0755); err != nil {
			return fmt.Errorf("failed to create case directory: %w", err);
		}

		// Save snapshots
		snapshotsPath := filepath.Join(caseDir, "snapshots.json");
		snapshotsFile, err := os.Create(snapshotsPath);
		if err != nil {
			return fmt.Errorf("failed to create snapshots file: %w", err);
		}

		encoder := json.NewEncoder(snapshotsFile);
		encoder.SetIndent("", "  ");
		if err := encoder.Encode(snapshots); err != nil {
			snapshotsFile.Close();
			return fmt.Errorf("failed to encode snapshots: %w", err);
		}
		snapshotsFile.Close();

		fmt.Printf("  ✓ Saved %d snapshots to %s\n", len(snapshots), snapshotsPath);
	}

	// Create manifest
	manifest := CaseManifest{
		Version: "1.0",
		Cases:   cases,
	};

	manifestPath := filepath.Join(outputDir, "manifest.json");
	manifestFile, err := os.Create(manifestPath);
	if err != nil {
		return fmt.Errorf("failed to create manifest file: %w", err);
	}
	defer manifestFile.Close();

	encoder := json.NewEncoder(manifestFile);
	encoder.SetIndent("", "  ");
	if err := encoder.Encode(manifest); err != nil {
		return fmt.Errorf("failed to encode manifest: %w", err);
	}

	fmt.Printf("\n✓ Generated %d cases\n", len(cases));
	fmt.Printf("✓ Saved manifest to %s\n", manifestPath);

	return nil;
}
