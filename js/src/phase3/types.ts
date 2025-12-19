// Phase 3 TypeScript interfaces for Single Neuron with Activation Functions

export interface NeuronParams {
  w: number[];  // weights [w1, w2]
  b: number;    // bias
}

export interface NeuronGrads {
  grad_w: number[];  // [∂L/∂w1, ∂L/∂w2]
  grad_b: number;    // ∂L/∂b
}

export interface DataPoint2DNeuron {
  x: number[];  // [x1, x2]
  y: number;    // target value
}

export interface PointSnapshotNeuron {
  index: number;
  x: number[];           // [x1, x2]
  y_true: number;        // target
  z: number;             // pre-activation: z = w·x + b
  a: number;             // post-activation: a = σ(z)
  loss: number;          // (a - y_true)²
  dL_da: number;         // ∂L/∂a = 2(a - y_true)
  da_dz: number;         // ∂a/∂z = σ'(z)
  dL_dz: number;         // ∂L/∂z = ∂L/∂a × ∂a/∂z
  dz_dw: number[];       // [∂z/∂w1, ∂z/∂w2] = [x1, x2]
  dL_dw: number[];       // [∂L/∂w1, ∂L/∂w2]
  dL_db: number;         // ∂L/∂b
  in_saturation: boolean; // whether σ'(z) < 0.01
}

export interface UpdateDetailsNeuron {
  learning_rate: number;
  gradient_magnitude: number; // ||∇L||
  update_w: number[];         // -lr × grad_w
  update_b: number;           // -lr × grad_b
  step_size: number;          // ||update|| magnitude
}

export interface ChainRuleComponent {
  param_name: string;   // "w1", "w2", "b"
  dL_da: number;        // ∂L/∂a
  da_dz: number;        // ∂a/∂z = σ'(z)
  dz_dparam: number;    // ∂z/∂w = x, ∂z/∂b = 1
  dL_dparam: number;    // ∂L/∂param = product
}

export interface ChainRuleBreakdown {
  components: ChainRuleComponent[];
}

export interface NeuronSnapshot {
  step: number;
  params: NeuronParams;
  grads: NeuronGrads;
  z: number;                    // avg pre-activation across dataset
  a: number;                    // avg post-activation
  dL_dz: number;                // avg ∂L/∂z
  dL_da: number;                // avg ∂L/∂a
  local_derivative: number;     // avg σ'(z)
  activation: "sigmoid" | "relu" | "tanh";
  in_saturation_zone: boolean;  // whether avg |σ'(z)| < 0.01
  loss: number;                 // avg loss across dataset
  point_details: PointSnapshotNeuron[];
  update_components: UpdateDetailsNeuron;
  chain_rule_breakdown: ChainRuleBreakdown;
}

export interface TrainingConfig {
  learning_rate: number;
  num_steps: number;
  activation: string;
}

export interface NeuronTrainingCase {
  case_id: string;
  description: string;
  category: string;           // "saturation", "optimal", "dying-relu", "comparison"
  activation: string;         // "sigmoid", "relu", "tanh"
  dataset: DataPoint2DNeuron[];
  init_params: NeuronParams;
  final_params: NeuronParams;
  config: TrainingConfig;
  snapshots: NeuronSnapshot[];
}
