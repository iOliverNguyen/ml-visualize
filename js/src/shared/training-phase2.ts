// Client-side gradient descent training for Phase 2 (2 parameters)
// Ported from Go implementation in /go/linear/

// ============================================================================
// Data Structures
// ============================================================================

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
  gradient_direction: number;
  point_details: PointSnapshot2D[];
  update_components: UpdateDetails2D;
}

export interface TrainingConfig2D {
  w1_init?: number; // Initial weight 1 (default: 0)
  w2_init?: number; // Initial weight 2 (default: 0)
  lr: number;
  max_steps: number;
  data: DataPoint2D[];
}

export interface DataGenConfig2D {
  num_points: number;
  x1_min: number;
  x1_max: number;
  x2_min: number;
  x2_max: number;
  true_w1: number;
  true_w2: number;
  noise_level: number;
  seed?: number;
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

export interface LossGridConfig {
  w1_min: number;
  w1_max: number;
  w2_min: number;
  w2_max: number;
  resolution: number;
}

// ============================================================================
// Core Training Functions
// ============================================================================

/**
 * Forward pass: y_pred = w1*x1 + w2*x2
 */
function forward(w1: number, w2: number, x1: number, x2: number): number {
  return w1 * x1 + w2 * x2;
}

/**
 * Compute loss for a single point: (y_pred - y_true)^2
 */
function pointLoss(y_pred: number, y_true: number): number {
  const diff = y_pred - y_true;
  return diff * diff;
}

/**
 * Compute gradient with respect to w1 for a single point
 * Derivation:
 *   loss = (w1*x1 + w2*x2 - y_true)^2
 *   ∂(loss)/∂w1 = 2 * (y_pred - y_true) * x1
 */
function gradW1(w1: number, w2: number, x1: number, x2: number, y_true: number): number {
  const y_pred = forward(w1, w2, x1, x2);
  return 2 * (y_pred - y_true) * x1;
}

/**
 * Compute gradient with respect to w2 for a single point
 * Derivation:
 *   loss = (w1*x1 + w2*x2 - y_true)^2
 *   ∂(loss)/∂w2 = 2 * (y_pred - y_true) * x2
 */
function gradW2(w1: number, w2: number, x1: number, x2: number, y_true: number): number {
  const y_pred = forward(w1, w2, x1, x2);
  return 2 * (y_pred - y_true) * x2;
}

/**
 * Compute gradient magnitude (L2 norm)
 * ||∇L|| = sqrt(grad_w1² + grad_w2²)
 */
function gradientMagnitude(grad_w1: number, grad_w2: number): number {
  return Math.sqrt(grad_w1 * grad_w1 + grad_w2 * grad_w2);
}

/**
 * Compute gradient direction (angle in radians)
 * θ = atan2(grad_w2, grad_w1)
 * Returns angle in range [-π, π]
 */
function gradientDirection(grad_w1: number, grad_w2: number): number {
  return Math.atan2(grad_w2, grad_w1);
}

/**
 * Compute total loss across all data points (mean squared error)
 */
function computeLoss2D(data: DataPoint2D[], w1: number, w2: number): number {
  let totalLoss = 0;
  for (const point of data) {
    const y_pred = forward(w1, w2, point.x1, point.x2);
    totalLoss += pointLoss(y_pred, point.y_true);
  }
  return totalLoss / data.length;
}

/**
 * Train a 2-parameter linear model using gradient descent
 * Returns snapshots of each training step
 */
export function trainModel2D(config: TrainingConfig2D): LinearSnapshot[] {
  const snapshots: LinearSnapshot[] = [];
  let w1 = config.w1_init ?? 0;
  let w2 = config.w2_init ?? 0;
  const lr = config.lr;

  for (let step = 0; step < config.max_steps; step++) {
    let totalLoss = 0;
    let totalGradW1 = 0;
    let totalGradW2 = 0;
    const pointDetails: PointSnapshot2D[] = [];

    // Compute per-point values
    for (const point of config.data) {
      const y_pred = forward(w1, w2, point.x1, point.x2);
      const loss = pointLoss(y_pred, point.y_true);
      const gw1 = gradW1(w1, w2, point.x1, point.x2, point.y_true);
      const gw2 = gradW2(w1, w2, point.x1, point.x2, point.y_true);

      pointDetails.push({
        x1: point.x1,
        x2: point.x2,
        y_true: point.y_true,
        y_pred,
        point_loss: loss,
        grad_w1: gw1,
        grad_w2: gw2
      });

      totalLoss += loss;
      totalGradW1 += gw1;
      totalGradW2 += gw2;
    }

    // Average over dataset
    const n = config.data.length;
    const avgLoss = totalLoss / n;
    const avgGradW1 = totalGradW1 / n;
    const avgGradW2 = totalGradW2 / n;

    const gradMag = gradientMagnitude(avgGradW1, avgGradW2);
    const gradDir = gradientDirection(avgGradW1, avgGradW2);

    // Compute updates
    const deltaW1 = -lr * avgGradW1;
    const deltaW2 = -lr * avgGradW2;
    const w1New = w1 + deltaW1;
    const w2New = w2 + deltaW2;

    // Create snapshot
    snapshots.push({
      step,
      w1,
      w2,
      grad_w1: avgGradW1,
      grad_w2: avgGradW2,
      loss: avgLoss,
      gradient_magnitude: gradMag,
      gradient_direction: gradDir,
      point_details: pointDetails,
      update_components: {
        w1_old: w1,
        w2_old: w2,
        lr,
        grad_w1: avgGradW1,
        grad_w2: avgGradW2,
        delta_w1: deltaW1,
        delta_w2: deltaW2,
        w1_new: w1New,
        w2_new: w2New
      }
    });

    // Update parameters
    w1 = w1New;
    w2 = w2New;
  }

  return snapshots;
}

// ============================================================================
// Data Generation
// ============================================================================

/**
 * Generate random 2D training data
 * Data follows: y = trueW1*x1 + trueW2*x2 + noise
 */
export function generateRandomData2D(config: DataGenConfig2D): DataPoint2D[] {
  // Validate config
  if (config.num_points <= 0) {
    throw new Error(`num_points must be positive, got ${config.num_points}`);
  }
  if (config.x1_max <= config.x1_min) {
    throw new Error('x1_max must be greater than x1_min');
  }
  if (config.x2_max <= config.x2_min) {
    throw new Error('x2_max must be greater than x2_min');
  }
  if (config.noise_level < 0) {
    throw new Error(`noise_level must be non-negative, got ${config.noise_level}`);
  }

  // Use seed for reproducibility if provided
  const rng = config.seed ? seedRandom(config.seed) : Math.random;

  const data: DataPoint2D[] = [];
  const x1Range = config.x1_max - config.x1_min;
  const x2Range = config.x2_max - config.x2_min;

  for (let i = 0; i < config.num_points; i++) {
    // Generate random x1, x2 in specified ranges
    const x1 = config.x1_min + rng() * x1Range;
    const x2 = config.x2_min + rng() * x2Range;

    // Compute true y value
    const yTrue = config.true_w1 * x1 + config.true_w2 * x2;

    // Add uniform noise in [-noise_level, +noise_level]
    const noise = (rng() * 2 - 1) * config.noise_level;
    const y = yTrue + noise;

    data.push({ x1, x2, y_true: y });
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
 * Validate that a 2D dataset is suitable for training
 */
export function validateDataset2D(data: DataPoint2D[]): string | null {
  if (data.length === 0) {
    return 'Dataset is empty';
  }

  for (let i = 0; i < data.length; i++) {
    const point = data[i];
    if (!Number.isFinite(point.x1)) {
      return `Point ${i} has invalid x1 value: ${point.x1}`;
    }
    if (!Number.isFinite(point.x2)) {
      return `Point ${i} has invalid x2 value: ${point.x2}`;
    }
    if (!Number.isFinite(point.y_true)) {
      return `Point ${i} has invalid y_true value: ${point.y_true}`;
    }
  }

  return null;
}

// ============================================================================
// Loss Grid Computation
// ============================================================================

/**
 * Compute a 2D grid of loss values for contour visualization
 * This computes loss at resolution×resolution points in (w1, w2) space
 */
export function computeLossGrid(data: DataPoint2D[], config: LossGridConfig): LossGrid {
  const { w1_min, w1_max, w2_min, w2_max, resolution } = config;

  const w1Step = (w1_max - w1_min) / (resolution - 1);
  const w2Step = (w2_max - w2_min) / (resolution - 1);

  const points: LossGridPoint[] = [];

  for (let i = 0; i < resolution; i++) {
    const w1 = w1_min + i * w1Step;
    for (let j = 0; j < resolution; j++) {
      const w2 = w2_min + j * w2Step;

      // Compute loss at (w1, w2)
      const loss = computeLoss2D(data, w1, w2);

      points.push({ w1, w2, loss });
    }
  }

  return {
    w1_min,
    w1_max,
    w2_min,
    w2_max,
    resolution,
    points
  };
}
