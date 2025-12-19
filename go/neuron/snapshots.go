package neuron

// NeuronParams represents the parameters of the neuron: w = [w1, w2], b
type NeuronParams struct {
	W []float64 `json:"w"` // weights [w1, w2]
	B float64   `json:"b"` // bias
}

// NeuronGrads represents the gradients with respect to parameters
type NeuronGrads struct {
	GradW []float64 `json:"grad_w"` // [∂L/∂w1, ∂L/∂w2]
	GradB float64   `json:"grad_b"` // ∂L/∂b
}

// DataPoint2DNeuron represents a single training data point with 2 features
type DataPoint2DNeuron struct {
	X []float64 `json:"x"` // [x1, x2]
	Y float64   `json:"y"` // target value
}

// PointSnapshotNeuron represents the forward and backward pass for a single data point
type PointSnapshotNeuron struct {
	Index        int       `json:"index"`         // point index in dataset
	X            []float64 `json:"x"`             // [x1, x2]
	YTrue        float64   `json:"y_true"`        // target
	Z            float64   `json:"z"`             // pre-activation: z = w·x + b
	A            float64   `json:"a"`             // post-activation: a = σ(z)
	Loss         float64   `json:"loss"`          // (a - y_true)²
	DLda         float64   `json:"dL_da"`         // ∂L/∂a = 2(a - y_true)
	DaDz         float64   `json:"da_dz"`         // ∂a/∂z = σ'(z)
	DLdz         float64   `json:"dL_dz"`         // ∂L/∂z = ∂L/∂a × ∂a/∂z
	DzDw         []float64 `json:"dz_dw"`         // [∂z/∂w1, ∂z/∂w2] = [x1, x2]
	DLdw         []float64 `json:"dL_dw"`         // [∂L/∂w1, ∂L/∂w2]
	DLdb         float64   `json:"dL_db"`         // ∂L/∂b
	InSaturation bool      `json:"in_saturation"` // whether σ'(z) < 0.01
}

// UpdateDetailsNeuron contains the details of the parameter update for this step
type UpdateDetailsNeuron struct {
	LearningRate  float64   `json:"learning_rate"`
	GradMagnitude float64   `json:"gradient_magnitude"` // ||∇L||
	UpdateW       []float64 `json:"update_w"`           // -lr × grad_w
	UpdateB       float64   `json:"update_b"`           // -lr × grad_b
	StepSize      float64   `json:"step_size"`          // ||update|| magnitude
}

// ChainRuleComponent represents one parameter's chain rule breakdown
type ChainRuleComponent struct {
	ParamName string  `json:"param_name"` // "w1", "w2", "b"
	DLda      float64 `json:"dL_da"`      // ∂L/∂a
	DaDz      float64 `json:"da_dz"`      // ∂a/∂z = σ'(z)
	DzDparam  float64 `json:"dz_dparam"`  // ∂z/∂w = x, ∂z/∂b = 1
	DLdparam  float64 `json:"dL_dparam"`  // ∂L/∂param = product
}

// ChainRuleViz contains the chain rule breakdown for all parameters
type ChainRuleViz struct {
	Components []ChainRuleComponent `json:"components"` // [w1, w2, b]
}

// NeuronSnapshot represents the state at a single training step
type NeuronSnapshot struct {
	Step               int                   `json:"step"`
	Params             NeuronParams          `json:"params"`
	Grads              NeuronGrads           `json:"grads"`
	Z                  float64               `json:"z"`                    // avg pre-activation across dataset
	A                  float64               `json:"a"`                    // avg post-activation
	DLdz               float64               `json:"dL_dz"`                // avg ∂L/∂z
	DLda               float64               `json:"dL_da"`                // avg ∂L/∂a
	LocalDerivative    float64               `json:"local_derivative"`     // avg σ'(z)
	Activation         string                `json:"activation"`           // "sigmoid", "relu", "tanh"
	InSaturationZone   bool                  `json:"in_saturation_zone"`   // whether avg |σ'(z)| < 0.01
	Loss               float64               `json:"loss"`                 // avg loss across dataset
	PointDetails       []PointSnapshotNeuron `json:"point_details"`        // per-point breakdown
	UpdateComponents   UpdateDetailsNeuron   `json:"update_components"`    // update details
	ChainRuleBreakdown ChainRuleViz          `json:"chain_rule_breakdown"` // chain rule for each param
}

// NeuronTrainingCase represents a complete training case with metadata
type NeuronTrainingCase struct {
	CaseID      string           `json:"case_id"`
	Description string           `json:"description"`
	Category    string           `json:"category"`    // "saturation", "optimal", "dying-relu", "comparison"
	Activation  string           `json:"activation"`  // "sigmoid", "relu", "tanh"
	Dataset     []DataPoint2DNeuron `json:"dataset"`
	InitParams  NeuronParams     `json:"init_params"`
	FinalParams NeuronParams     `json:"final_params"`
	Config      TrainingConfig   `json:"config"`
	Snapshots   []NeuronSnapshot `json:"snapshots"`
}

// TrainingConfig contains the hyperparameters for training
type TrainingConfig struct {
	LearningRate float64 `json:"learning_rate"`
	NumSteps     int     `json:"num_steps"`
	Activation   string  `json:"activation"`
}
