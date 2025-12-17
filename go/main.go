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

		// Process each data point
		for _, point := range data {
			// Forward pass
			yPred := Forward(w, point.X)

			// Compute loss
			loss := Loss(yPred, point.YTrue)

			// Compute gradient
			grad := GradW(w, point.X, point.YTrue)

			totalLoss += loss
			totalGrad += grad
		}

		// Average over dataset
		avgLoss := totalLoss / float64(len(data))
		avgGrad := totalGrad / float64(len(data))

		// Create snapshot BEFORE parameter update
		// This captures the state that produced this gradient
		snapshots = append(snapshots, Snapshot{
			Step:  step,
			W:     w,
			GradW: avgGrad,
			Loss:  avgLoss,
		})

		// Print progress every 20 steps
		if step%20 == 0 {
			fmt.Printf("Step %3d: w=%.4f, grad_w=%.4f, loss=%.4f\n",
				step, w, avgGrad, avgLoss)
		}

		// Update parameter using gradient descent
		// w_new = w_old - learning_rate * gradient
		w = w - lr*avgGrad
	}

	fmt.Println()
	fmt.Printf("Training complete! Final w: %.4f\n", w)

	return snapshots
}
