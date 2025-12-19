// Client-side gradient descent training for Phase 1
// Ported from Go implementation in /go/

// ============================================================================
// Data Structures
// ============================================================================

export interface DataPoint {
  x: number;
  yTrue: number;
}

export interface TrainingConfig {
  wInit?: number; // Initial weight (default: 0)
  learningRate: number;
  maxSteps: number;
  data: DataPoint[];
}

export interface PointSnapshot {
  x: number;
  yTrue: number;
  yPred: number;
  pointLoss: number; // (yPred - yTrue)^2
  pointGrad: number; // 2 * (yPred - yTrue) * x
}

export interface UpdateDetails {
  wOld: number;
  lr: number;
  gradW: number;
  deltaW: number; // -lr * gradW
  wNew: number;
}

export interface Snapshot {
  step: number;
  w: number;
  gradW: number;
  loss: number;
  pointDetails: PointSnapshot[];
  updateComponents: UpdateDetails;
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
function pointLoss(yPred: number, yTrue: number): number {
  const diff = yPred - yTrue;
  return diff * diff;
}

/**
 * Compute gradient with respect to w for a single point
 * Derivation:
 *   loss = (w*x - y_true)^2
 *   d(loss)/dw = 2 * (w*x - y_true) * x
 */
function gradW(w: number, x: number, yTrue: number): number {
  const yPred = forward(w, x);
  return 2 * (yPred - yTrue) * x;
}

/**
 * Compute total loss across all data points (mean squared error)
 */
function computeLoss(data: DataPoint[], weight: number): number {
  let totalLoss = 0;
  for (const point of data) {
    const yPred = forward(weight, point.x);
    totalLoss += pointLoss(yPred, point.yTrue);
  }
  return totalLoss / data.length;
}

/**
 * Compute gradient across all data points (mean gradient)
 */
function computeGradient(data: DataPoint[], weight: number): number {
  let gradient = 0;
  for (const point of data) {
    gradient += gradW(weight, point.x, point.yTrue);
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
  const pointDetails: PointSnapshot[] = data.map(point => {
    const yPred = forward(weight, point.x);
    return {
      x: point.x,
      yTrue: point.yTrue,
      yPred,
      pointLoss: pointLoss(yPred, point.yTrue),
      pointGrad: gradW(weight, point.x, point.yTrue)
    };
  });

  // Compute update components
  const deltaW = -config.learningRate * gradient;
  const wNew = weight + deltaW;

  const updateComponents: UpdateDetails = {
    wOld: weight,
    lr: config.learningRate,
    gradW: gradient,
    deltaW,
    wNew
  };

  return {
    step,
    w: weight,
    gradW: gradient,
    loss,
    pointDetails,
    updateComponents
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
    weight -= config.learningRate * snapshot.gradW;
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

    data.push({ x, yTrue: y });
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
    if (!Number.isFinite(point.yTrue)) {
      return `Point ${i} has invalid YTrue value: ${point.yTrue}`;
    }
  }

  return null;
}
