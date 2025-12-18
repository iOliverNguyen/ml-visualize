package linear

// LossGridPoint represents loss at a specific (w1, w2) point
type LossGridPoint struct {
	W1   float64 `json:"w1"`
	W2   float64 `json:"w2"`
	Loss float64 `json:"loss"`
}

// LossGrid represents a 2D grid of loss values for visualization
type LossGrid struct {
	W1Min      float64          `json:"w1_min"`
	W1Max      float64          `json:"w1_max"`
	W2Min      float64          `json:"w2_min"`
	W2Max      float64          `json:"w2_max"`
	Resolution int              `json:"resolution"`
	Points     []LossGridPoint  `json:"points"`
}

// ComputeLossGrid generates a grid of loss values for contour plotting
func ComputeLossGrid(data []DataPoint2D, w1Min, w1Max, w2Min, w2Max float64, resolution int) LossGrid {
	w1Step := (w1Max - w1Min) / float64(resolution-1)
	w2Step := (w2Max - w2Min) / float64(resolution-1)

	points := make([]LossGridPoint, 0, resolution*resolution)

	for i := 0; i < resolution; i++ {
		w1 := w1Min + float64(i)*w1Step
		for j := 0; j < resolution; j++ {
			w2 := w2Min + float64(j)*w2Step

			// Compute loss at (w1, w2)
			totalLoss := 0.0
			for _, point := range data {
				yPred := Forward(w1, w2, point.X1, point.X2)
				totalLoss += Loss(yPred, point.YTrue)
			}
			avgLoss := totalLoss / float64(len(data))

			points = append(points, LossGridPoint{
				W1:   w1,
				W2:   w2,
				Loss: avgLoss,
			})
		}
	}

	return LossGrid{
		W1Min:      w1Min,
		W1Max:      w1Max,
		W2Min:      w2Min,
		W2Max:      w2Max,
		Resolution: resolution,
		Points:     points,
	}
}
