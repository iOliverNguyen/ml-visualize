package linear

// TrainingConfig2D holds configuration for 2-parameter training
type TrainingConfig2D struct {
	W1Init float64 `json:"w1_init"`
	W2Init float64 `json:"w2_init"`
	LR     float64 `json:"lr"`
	Steps  int     `json:"steps"`
}

// RunTraining performs gradient descent training and returns snapshots
func RunTraining(data []DataPoint2D, config TrainingConfig2D) []LinearSnapshot {
	w1, w2 := config.W1Init, config.W2Init
	lr := config.LR
	steps := config.Steps

	snapshots := make([]LinearSnapshot, 0, steps)

	for step := 0; step < steps; step++ {
		totalLoss := 0.0
		totalGradW1 := 0.0
		totalGradW2 := 0.0
		pointDetails := make([]PointSnapshot2D, 0, len(data))

		// Compute per-point values
		for _, point := range data {
			yPred := Forward(w1, w2, point.X1, point.X2)
			pointLoss := Loss(yPred, point.YTrue)
			gradW1 := GradW1(w1, w2, point.X1, point.X2, point.YTrue)
			gradW2 := GradW2(w1, w2, point.X1, point.X2, point.YTrue)

			pointDetails = append(pointDetails, PointSnapshot2D{
				X1:        point.X1,
				X2:        point.X2,
				YTrue:     point.YTrue,
				YPred:     yPred,
				PointLoss: pointLoss,
				GradW1:    gradW1,
				GradW2:    gradW2,
			})

			totalLoss += pointLoss
			totalGradW1 += gradW1
			totalGradW2 += gradW2
		}

		// Average over dataset
		n := float64(len(data))
		avgLoss := totalLoss / n
		avgGradW1 := totalGradW1 / n
		avgGradW2 := totalGradW2 / n

		gradMag := GradientMagnitude(avgGradW1, avgGradW2)
		gradDir := GradientDirection(avgGradW1, avgGradW2)

		// Compute updates
		deltaW1 := -lr * avgGradW1
		deltaW2 := -lr * avgGradW2
		w1New := w1 + deltaW1
		w2New := w2 + deltaW2

		// Create snapshot
		snapshots = append(snapshots, LinearSnapshot{
			Step:              step,
			W1:                w1,
			W2:                w2,
			GradW1:            avgGradW1,
			GradW2:            avgGradW2,
			Loss:              avgLoss,
			GradientMagnitude: gradMag,
			GradientDirection: gradDir,
			PointDetails:      pointDetails,
			UpdateComponents: UpdateDetails2D{
				W1Old:   w1,
				W2Old:   w2,
				LR:      lr,
				GradW1:  avgGradW1,
				GradW2:  avgGradW2,
				DeltaW1: deltaW1,
				DeltaW2: deltaW2,
				W1New:   w1New,
				W2New:   w2New,
			},
		})

		// Update parameters
		w1, w2 = w1New, w2New
	}

	return snapshots
}
