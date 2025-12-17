<script lang="ts">
  import type { PointSnapshot } from './types';

  interface Props {
    dataset: PointSnapshot[];
    currentW: number;
    initialW: number;
    targetW: number;
    currentLoss: number;
  }

  let {
    dataset,
    currentW,
    initialW,
    targetW,
    currentLoss
  }: Props = $props();

  // Calculate dataset statistics
  const stats = $derived(() => {
    if (dataset.length === 0) {
      return {
        numPoints: 0,
        xMin: 0,
        xMax: 0,
        yMin: 0,
        yMax: 0
      };
    }

    const xValues = dataset.map(p => p.x);
    const yValues = dataset.map(p => p.y_true);

    return {
      numPoints: dataset.length,
      xMin: Math.min(...xValues),
      xMax: Math.max(...xValues),
      yMin: Math.min(...yValues),
      yMax: Math.max(...yValues)
    };
  });

  // Calculate progress toward target
  const progress = $derived(() => {
    const range = Math.abs(targetW - initialW);
    if (range === 0) return 100;
    const current = Math.abs(currentW - initialW);
    return Math.min(100, (current / range) * 100);
  });

  // Determine if we're close to target
  const isClose = $derived(() => {
    return Math.abs(currentW - targetW) < 0.1;
  });
</script>

<div class="dataset-info-panel">
  <h3 class="panel-title">
    <span class="icon">🎯</span>
    Training Objective
  </h3>

  <div class="objective-statement">
    <p class="main-text">
      We're fitting a line <strong>y = w*x</strong> to data
    </p>
    <p class="sub-text">
      Finding the best parameter <strong>w</strong> that minimizes prediction errors
    </p>
  </div>

  <div class="stats-grid">
    <div class="stat-card">
      <div class="stat-label">Dataset</div>
      <div class="stat-value">{stats().numPoints} points</div>
    </div>

    <div class="stat-card">
      <div class="stat-label">X range</div>
      <div class="stat-value">
        {stats().xMin.toFixed(1)} → {stats().xMax.toFixed(1)}
      </div>
    </div>

    <div class="stat-card">
      <div class="stat-label">Y range</div>
      <div class="stat-value">
        {stats().yMin.toFixed(1)} → {stats().yMax.toFixed(1)}
      </div>
    </div>

    <div class="stat-card">
      <div class="stat-label">Current Loss</div>
      <div class="stat-value loss-value">
        {currentLoss.toFixed(3)}
      </div>
    </div>
  </div>

  <div class="parameter-section">
    <div class="param-row">
      <span class="param-label">Target w:</span>
      <span class="param-value target">≈ {targetW.toFixed(2)}</span>
    </div>

    <div class="param-row">
      <span class="param-label">Initial w:</span>
      <span class="param-value">{initialW.toFixed(2)}</span>
    </div>

    <div class="param-row">
      <span class="param-label">Current w:</span>
      <span class="param-value current" class:close={isClose()}>
        {currentW.toFixed(3)}
        {#if isClose()}
          <span class="check-mark">✓</span>
        {/if}
      </span>
    </div>

    <div class="progress-bar-container">
      <div class="progress-bar-bg">
        <div
          class="progress-bar-fill"
          style="width: {progress()}%"
        ></div>
      </div>
      <div class="progress-label">
        {progress().toFixed(0)}% to target
      </div>
    </div>
  </div>

  <div class="explanation">
    <div class="explain-item">
      <strong>Loss:</strong> Measures how far predictions are from actual data
    </div>
    <div class="explain-item">
      <strong>Gradient:</strong> Direction to adjust w to reduce loss
    </div>
    <div class="explain-item">
      <strong>Training:</strong> Iteratively update w to minimize loss
    </div>
  </div>
</div>

<style>
  .dataset-info-panel {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 1.5rem;
    border-radius: 12px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }

  .panel-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 0 0 1rem 0;
    font-size: 1.25rem;
    font-weight: 700;
  }

  .icon {
    font-size: 1.5rem;
  }

  .objective-statement {
    background: rgba(255, 255, 255, 0.15);
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1rem;
    backdrop-filter: blur(10px);
  }

  .main-text {
    margin: 0 0 0.5rem 0;
    font-size: 1.05rem;
    line-height: 1.5;
  }

  .sub-text {
    margin: 0;
    font-size: 0.9rem;
    opacity: 0.9;
    line-height: 1.4;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .stat-card {
    background: rgba(255, 255, 255, 0.15);
    padding: 0.75rem;
    border-radius: 6px;
    backdrop-filter: blur(10px);
  }

  .stat-label {
    font-size: 0.75rem;
    opacity: 0.9;
    margin-bottom: 0.25rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .stat-value {
    font-size: 1rem;
    font-weight: 600;
    font-family: 'SF Mono', monospace;
  }

  .loss-value {
    color: #fbbf24;
  }

  .parameter-section {
    background: rgba(255, 255, 255, 0.15);
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1rem;
    backdrop-filter: blur(10px);
  }

  .param-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 0.5rem 0;
  }

  .param-label {
    font-size: 0.9rem;
    font-weight: 500;
  }

  .param-value {
    font-size: 1.1rem;
    font-weight: 700;
    font-family: 'SF Mono', monospace;
  }

  .param-value.target {
    color: #fbbf24;
  }

  .param-value.current {
    color: #60a5fa;
  }

  .param-value.current.close {
    color: #34d399;
  }

  .check-mark {
    margin-left: 0.25rem;
    font-size: 1.2rem;
  }

  .progress-bar-container {
    margin-top: 1rem;
  }

  .progress-bar-bg {
    height: 8px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 4px;
    overflow: hidden;
  }

  .progress-bar-fill {
    height: 100%;
    background: linear-gradient(90deg, #60a5fa 0%, #34d399 100%);
    transition: width 0.3s ease;
    border-radius: 4px;
  }

  .progress-label {
    margin-top: 0.5rem;
    font-size: 0.85rem;
    text-align: center;
    opacity: 0.9;
  }

  .explanation {
    background: rgba(255, 255, 255, 0.1);
    padding: 1rem;
    border-radius: 8px;
    backdrop-filter: blur(10px);
  }

  .explain-item {
    margin: 0.5rem 0;
    font-size: 0.85rem;
    line-height: 1.5;
    opacity: 0.95;
  }

  .explain-item strong {
    font-weight: 600;
  }

  @media (max-width: 640px) {
    .dataset-info-panel {
      padding: 1rem;
    }

    .stats-grid {
      grid-template-columns: 1fr;
    }

    .param-row {
      flex-direction: column;
      align-items: flex-start;
      gap: 0.25rem;
    }
  }
</style>
