package core

import (
	"fmt"
)

func RunTraining() []Snapshot {
	// Load dataset
	data := GetDataset()

	// Training hyperparameters
	w := 0.0     // Initial parameter
	lr := 0.01   // Learning rate
	steps := 100 // Number of training steps

	// Storage for all training snapshots
	snapshots := []Snapshot{}

	fmt.Println("Starting training...")
	fmt.Printf("Initial w: %.4f\n", w)
	fmt.Printf("Learning rate: %.4f\n", lr)
	fmt.Printf("Dataset size: %d\n", len(data))
	fmt.Println()

	// Training loop - explicit and imperative (no Model or Trainer abstraction)
	for step := 0; step < steps; step++ {
		totalLoss := 0.0
		totalGrad := 0.0
		pointDetails := make([]PointSnapshot, 0, len(data))

		// Process each data point
		for _, point := range data {
			// Forward pass
			yPred := Forward(w, point.X)

			// Compute loss
			pointLoss := Loss(yPred, point.YTrue)

			// Compute gradient
			pointGrad := GradW(w, point.X, point.YTrue)

			// Store per-point details for pedagogical inspection
			pointDetails = append(pointDetails, PointSnapshot{
				X:         point.X,
				YTrue:     point.YTrue,
				YPred:     yPred,
				PointLoss: pointLoss,
				PointGrad: pointGrad,
			})

			totalLoss += pointLoss
			totalGrad += pointGrad
		}

		// Average over dataset
		avgLoss := totalLoss / float64(len(data))
		avgGrad := totalGrad / float64(len(data))

		// Compute w_new before creating snapshot
		deltaW := -lr * avgGrad
		wNew := w + deltaW

		// Create snapshot BEFORE parameter update
		// This captures the state that produced this gradient
		snapshots = append(snapshots, Snapshot{
			Step:  step,
			W:     w,
			GradW: avgGrad,
			Loss:  avgLoss,
			PointDetails: pointDetails,
			UpdateComponents: UpdateDetails{
				WOld:   w,
				LR:     lr,
				GradW:  avgGrad,
				DeltaW: deltaW,
				WNew:   wNew,
			},
		})

		// Print progress every 20 steps
		if step%20 == 0 {
			fmt.Printf("Step %3d: w=%.4f, grad_w=%.4f, loss=%.4f\n",
				step, w, avgGrad, avgLoss)
		}

		// Update parameter using gradient descent
		w = wNew
	}

	fmt.Println()
	fmt.Printf("Training complete! Final w: %.4f\n", w)

	return snapshots
}
