package core

// DataPoint represents a single (x, y_true) training example
type DataPoint struct {
	X     float64
	YTrue float64
}

// GetDataset returns hardcoded training data
// Data points lie approximately on the line y = 2x with small noise
func GetDataset() []DataPoint {
	return []DataPoint{
		{X: 1.0, YTrue: 2.1},
		{X: 2.0, YTrue: 3.9},
		{X: 3.0, YTrue: 6.2},
		{X: 4.0, YTrue: 7.8},
		{X: 5.0, YTrue: 10.1},
		{X: 6.0, YTrue: 11.9},
		{X: 7.0, YTrue: 14.2},
		{X: 8.0, YTrue: 15.8},
		{X: 9.0, YTrue: 18.1},
		{X: 10.0, YTrue: 19.9},
	}
}
