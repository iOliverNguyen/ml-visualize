// Phase 2: Linear Regression with 2 Parameters (w1, w2)
// Type definitions matching Go backend JSON structures (snake_case)

export interface DataPoint2D {
  x1: number;
  x2: number;
  y_true: number;
}

export interface PointSnapshot2D {
  x1: number;
  x2: number;
  y_true: number;
  y_pred: number;
  point_loss: number;
  grad_w1: number;
  grad_w2: number;
}

export interface UpdateDetails2D {
  w1_old: number;
  w2_old: number;
  lr: number;
  grad_w1: number;
  grad_w2: number;
  delta_w1: number;
  delta_w2: number;
  w1_new: number;
  w2_new: number;
}

export interface LinearSnapshot {
  step: number;
  w1: number;
  w2: number;
  grad_w1: number;
  grad_w2: number;
  loss: number;
  gradient_magnitude: number;
  gradient_direction: number; // radians [-π, π]
  point_details: PointSnapshot2D[];
  update_components: UpdateDetails2D;
}

export interface LossGridPoint {
  w1: number;
  w2: number;
  loss: number;
}

export interface LossGrid {
  w1_min: number;
  w1_max: number;
  w2_min: number;
  w2_max: number;
  resolution: number;
  points: LossGridPoint[];
}

// Metrics data structure for LinearMetricsPanel
export interface LinearMetricsData {
  initial: LinearSnapshot;
  current: LinearSnapshot;
  final: LinearSnapshot;
  deltaFromStart: {
    w1: number;
    w2: number;
    loss: number;
    gradient_magnitude: number;
    distance: number; // Euclidean distance in parameter space
  };
  deltaToTarget: {
    w1: number;
    w2: number;
    loss: number;
    distance: number;
  };
  progress: number; // percentage
}