<script lang="ts">
  import type { LinearMetricsData } from './types';

  interface Props {
    metrics: LinearMetricsData;
  }

  let { metrics }: Props = $props();

  function formatNumber(n: number): string {
    return n.toFixed(4);
  }

  function formatPercent(n: number): string {
    return n.toFixed(2) + '%';
  }

  function formatDegrees(radians: number): string {
    const degrees = (radians * 180) / Math.PI;
    return degrees.toFixed(1) + '°';
  }
</script>

<div class="metrics-panel">
  <div class="metric-card initial">
    <h3>Initial (Step 0)</h3>
    <div class="metric-row">
      <span class="label">w₁₀:</span>
      <span class="value">{formatNumber(metrics.initial.w1)}</span>
    </div>
    <div class="metric-row">
      <span class="label">w₂₀:</span>
      <span class="value">{formatNumber(metrics.initial.w2)}</span>
    </div>
    <div class="metric-row">
      <span class="label">||∇L||₀:</span>
      <span class="value">{formatNumber(metrics.initial.gradient_magnitude)}</span>
    </div>
    <div class="metric-row">
      <span class="label">loss₀:</span>
      <span class="value">{formatNumber(metrics.initial.loss)}</span>
    </div>
  </div>

  <div class="metric-card current">
    <h3>Current (Step {metrics.current.step})</h3>
    <div class="metric-row">
      <span class="label">w₁:</span>
      <span class="value">{formatNumber(metrics.current.w1)}</span>
    </div>
    <div class="metric-row">
      <span class="label">w₂:</span>
      <span class="value">{formatNumber(metrics.current.w2)}</span>
    </div>
    <div class="metric-row">
      <span class="label">||∇L||:</span>
      <span class="value">{formatNumber(metrics.current.gradient_magnitude)}</span>
    </div>
    <div class="metric-row">
      <span class="label">∠∇L:</span>
      <span class="value">{formatDegrees(metrics.current.gradient_direction)}</span>
    </div>
    <div class="metric-row">
      <span class="label">loss:</span>
      <span class="value">{formatNumber(metrics.current.loss)}</span>
    </div>
  </div>

  <div class="metric-card final">
    <h3>Final (Target)</h3>
    <div class="metric-row">
      <span class="label">w₁_final:</span>
      <span class="value">{formatNumber(metrics.final.w1)}</span>
    </div>
    <div class="metric-row">
      <span class="label">w₂_final:</span>
      <span class="value">{formatNumber(metrics.final.w2)}</span>
    </div>
    <div class="metric-row">
      <span class="label">||∇L||_final:</span>
      <span class="value">{formatNumber(metrics.final.gradient_magnitude)}</span>
    </div>
    <div class="metric-row">
      <span class="label">loss_final:</span>
      <span class="value">{formatNumber(metrics.final.loss)}</span>
    </div>
  </div>

  <div class="metric-card delta">
    <h3>Progress from Start</h3>
    <div class="metric-row">
      <span class="label">Δw₁:</span>
      <span class="value">{formatNumber(metrics.deltaFromStart.w1)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Δw₂:</span>
      <span class="value">{formatNumber(metrics.deltaFromStart.w2)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Distance:</span>
      <span class="value">{formatNumber(metrics.deltaFromStart.distance)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Δloss:</span>
      <span class="value">{formatNumber(metrics.deltaFromStart.loss)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Progress:</span>
      <span class="value">{formatPercent(metrics.progress)}</span>
    </div>
  </div>

  <div class="metric-card delta">
    <h3>Remaining to Target</h3>
    <div class="metric-row">
      <span class="label">Δw₁ to target:</span>
      <span class="value">{formatNumber(metrics.deltaToTarget.w1)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Δw₂ to target:</span>
      <span class="value">{formatNumber(metrics.deltaToTarget.w2)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Distance:</span>
      <span class="value">{formatNumber(metrics.deltaToTarget.distance)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Δloss to target:</span>
      <span class="value">{formatNumber(metrics.deltaToTarget.loss)}</span>
    </div>
  </div>
</div>

<style>
  .metrics-panel {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 1rem;
    margin: 2rem 0;
  }

  .metric-card {
    border: 2px solid;
    border-radius: 8px;
    padding: 1rem;
    background: white;
  }

  .metric-card.initial {
    border-color: #2563eb;
    background: #eff6ff;
  }

  .metric-card.current {
    border-color: #dc2626;
    background: #fef2f2;
  }

  .metric-card.final {
    border-color: #16a34a;
    background: #f0fdf4;
  }

  .metric-card.delta {
    border-color: #7c3aed;
    background: #faf5ff;
  }

  h3 {
    margin: 0 0 0.75rem 0;
    font-size: 1rem;
    font-weight: bold;
  }

  .initial h3 {
    color: #2563eb;
  }

  .current h3 {
    color: #dc2626;
  }

  .final h3 {
    color: #16a34a;
  }

  .delta h3 {
    color: #7c3aed;
  }

  .metric-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.25rem 0;
    font-size: 0.9rem;
  }

  .label {
    font-weight: 500;
    color: #555;
  }

  .value {
    font-family: 'Courier New', monospace;
    font-weight: bold;
    color: #1a1a1a;
  }

  @media (max-width: 768px) {
    .metrics-panel {
      grid-template-columns: 1fr;
    }
  }
</style>
