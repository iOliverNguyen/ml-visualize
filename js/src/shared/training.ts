// Client-side gradient descent training for Phase 1
// Ported from Go implementation in /go/

// ============================================================================
// Data Structures
// ============================================================================

export interface DataPoint {
  x: number;
  y_true: number;  // Changed from yTrue to match Go JSON format
}

export interface TrainingConfig {
  wInit?: number; // Initial weight (default: 0)
  learningRate: number;
  maxSteps: number;
  data: DataPoint[];
}

export interface PointSnapshot {
  x: number;
  y_true: number;      // Changed from yTrue
  y_pred: number;      // Changed from yPred
  point_loss: number;  // Changed from pointLoss
  point_grad: number;  // Changed from pointGrad
}

export interface UpdateDetails {
  w_old: number;    // Changed from wOld
  lr: number;
  grad_w: number;   // Changed from gradW
  delta_w: number;  // Changed from deltaW
  w_new: number;    // Changed from wNew
}

export interface Snapshot {
  step: number;
  w: number;
  grad_w: number;                    // Changed from gradW
  loss: number;
  point_details: PointSnapshot[];    // Changed from pointDetails
  update_components: UpdateDetails;  // Changed from updateComponents
}

export interface DataGenConfig {
  numPoints: number;
  xMin: number;
  xMax: number;
  trueSlope: number;
  noiseLevel: number;
  seed?: number;
}

// ============================================================================
// Core Training Functions
// ============================================================================

/**
 * Forward pass: y_pred = w * x
 */
function forward(w: number, x: number): number {
  return w * x;
}

/**
 * Compute loss for a single point: (y_pred - y_true)^2
 */
function pointLoss(y_pred: number, y_true: number): number {
  const diff = y_pred - y_true;
  return diff * diff;
}

/**
 * Compute gradient with respect to w for a single point
 * Derivation:
 *   loss = (w*x - y_true)^2
 *   d(loss)/dw = 2 * (w*x - y_true) * x
 */
function gradW(w: number, x: number, y_true: number): number {
  const y_pred = forward(w, x);
  return 2 * (y_pred - y_true) * x;
}

/**
 * Compute total loss across all data points (mean squared error)
 */
function computeLoss(data: DataPoint[], weight: number): number {
  let totalLoss = 0;
  for (const point of data) {
    const y_pred = forward(weight, point.x);
    totalLoss += pointLoss(y_pred, point.y_true);
  }
  return totalLoss / data.length;
}

/**
 * Compute gradient across all data points (mean gradient)
 */
function computeGradient(data: DataPoint[], weight: number): number {
  let gradient = 0;
  for (const point of data) {
    gradient += gradW(weight, point.x, point.y_true);
  }
  return gradient / data.length;
}

/**
 * Create a snapshot of the current training state
 */
function createSnapshot(
  step: number,
  weight: number,
  data: DataPoint[],
  config: TrainingConfig
): Snapshot {
  const loss = computeLoss(data, weight);
  const gradient = computeGradient(data, weight);

  // Compute per-point details
  const point_details: PointSnapshot[] = data.map(point => {
    const y_pred = forward(weight, point.x);
    return {
      x: point.x,
      y_true: point.y_true,
      y_pred,
      point_loss: pointLoss(y_pred, point.y_true),
      point_grad: gradW(weight, point.x, point.y_true)
    };
  });

  // Compute update components
  const delta_w = -config.learningRate * gradient;
  const w_new = weight + delta_w;

  const update_components: UpdateDetails = {
    w_old: weight,
    lr: config.learningRate,
    grad_w: gradient,
    delta_w,
    w_new
  };

  return {
    step,
    w: weight,
    grad_w: gradient,
    loss,
    point_details,
    update_components
  };
}

/**
 * Train a linear model using gradient descent
 * Returns snapshots of each training step
 */
export function trainModel(config: TrainingConfig): Snapshot[] {
  const snapshots: Snapshot[] = [];
  let weight = config.wInit ?? 0; // Start at wInit (default: 0)

  for (let step = 0; step < config.maxSteps; step++) {
    // Create snapshot before update
    const snapshot = createSnapshot(step, weight, config.data, config);
    snapshots.push(snapshot);

    // Update weight using gradient descent
    weight -= config.learningRate * snapshot.grad_w;
  }

  return snapshots;
}

// ============================================================================
// Data Generation
// ============================================================================

/**
 * Generate random training data
 * Data follows: y = trueSlope * x + noise
 *
 * Note: This generates evenly-spaced x values (not random x),
 * matching the Go implementation
 */
export function generateRandomData(config: DataGenConfig): DataPoint[] {
  // Validate config
  if (config.numPoints <= 0) {
    throw new Error(`num_points must be positive, got ${config.numPoints}`);
  }
  if (config.xMax <= config.xMin) {
    throw new Error('x_max must be greater than x_min');
  }
  if (config.noiseLevel < 0) {
    throw new Error(`noise_level must be non-negative, got ${config.noiseLevel}`);
  }

  // Use seed for reproducibility if provided
  const rng = config.seed ? seedRandom(config.seed) : Math.random;

  const data: DataPoint[] = [];
  const xRange = config.xMax - config.xMin;

  for (let i = 0; i < config.numPoints; i++) {
    // Evenly spaced x values (matches Go implementation)
    const x = config.xMin + (xRange * i / (config.numPoints - 1));

    // y = slope * x + noise
    // Uniform noise in [-noiseLevel, +noiseLevel]
    const noise = (rng() * 2 - 1) * config.noiseLevel;
    const y = config.trueSlope * x + noise;

    data.push({ x, y_true: y });
  }

  return data;
}

/**
 * Simple seeded random number generator for reproducibility
 * Uses linear congruential generator (LCG)
 */
function seedRandom(seed: number): () => number {
  let state = seed;
  return () => {
    state = (state * 9301 + 49297) % 233280;
    return state / 233280;
  };
}

/**
 * Validate that a dataset is suitable for training
 */
export function validateDataset(data: DataPoint[]): string | null {
  if (data.length === 0) {
    return 'Dataset is empty';
  }

  for (let i = 0; i < data.length; i++) {
    const point = data[i];
    if (!Number.isFinite(point.x)) {
      return `Point ${i} has invalid X value: ${point.x}`;
    }
    if (!Number.isFinite(point.y_true)) {
      return `Point ${i} has invalid y_true value: ${point.y_true}`;
    }
  }

  return null;
}
