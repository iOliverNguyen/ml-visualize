<script lang="ts">
  import type { MetricsData } from './types';

  interface Props {
    metrics: MetricsData;
  }

  let { metrics }: Props = $props();

  function formatNumber(n: number): string {
    return n.toFixed(4);
  }

  function formatPercent(n: number): string {
    return n.toFixed(2) + '%';
  }
</script>

<div class="metrics-panel">
  <div class="metric-card initial">
    <h3>Initial (Step 0)</h3>
    <div class="metric-row">
      <span class="label">w₀:</span>
      <span class="value">{formatNumber(metrics.initial.w)}</span>
    </div>
    <div class="metric-row">
      <span class="label">grad_w₀:</span>
      <span class="value">{formatNumber(metrics.initial.grad_w)}</span>
    </div>
    <div class="metric-row">
      <span class="label">loss₀:</span>
      <span class="value">{formatNumber(metrics.initial.loss)}</span>
    </div>
  </div>

  <div class="metric-card current">
    <h3>Current</h3>
    <div class="metric-row">
      <span class="label">w:</span>
      <span class="value">{formatNumber(metrics.current.w)}</span>
    </div>
    <div class="metric-row">
      <span class="label">grad_w:</span>
      <span class="value">{formatNumber(metrics.current.grad_w)}</span>
    </div>
    <div class="metric-row">
      <span class="label">loss:</span>
      <span class="value">{formatNumber(metrics.current.loss)}</span>
    </div>
  </div>

  <div class="metric-card final">
    <h3>Final (Target)</h3>
    <div class="metric-row">
      <span class="label">w_final:</span>
      <span class="value">{formatNumber(metrics.final.w)}</span>
    </div>
    <div class="metric-row">
      <span class="label">grad_w_final:</span>
      <span class="value">{formatNumber(metrics.final.grad_w)}</span>
    </div>
    <div class="metric-row">
      <span class="label">loss_final:</span>
      <span class="value">{formatNumber(metrics.final.loss)}</span>
    </div>
  </div>

  <div class="metric-card delta">
    <h3>Progress</h3>
    <div class="metric-row">
      <span class="label">Δw:</span>
      <span class="value">{formatNumber(metrics.deltaFromStart.w)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Δloss:</span>
      <span class="value">{formatNumber(metrics.deltaFromStart.loss)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Progress:</span>
      <span class="value">{formatPercent(metrics.deltaFromStart.progress)}</span>
    </div>
  </div>

  <div class="metric-card delta">
    <h3>Remaining to Target</h3>
    <div class="metric-row">
      <span class="label">Δw to target:</span>
      <span class="value">{formatNumber(metrics.deltaToTarget.w)}</span>
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
