<script lang="ts">
  import type { NeuronSnapshot } from './types';

  interface NeuronMetricsData {
    current: NeuronSnapshot;
    initial: NeuronSnapshot;
    final: NeuronSnapshot;
    deltaParams: {
      w1: number;
      w2: number;
      b: number;
    };
    deltaLoss: number;
    progress: number;
  }

  interface Props {
    metrics: NeuronMetricsData;
  }

  let { metrics }: Props = $props();

  function formatNumber(n: number): string {
    return n.toFixed(4);
  }

  function formatPercent(n: number): string {
    return n.toFixed(2) + '%';
  }

  function getActivationColor(activation: string): string {
    switch (activation) {
      case 'sigmoid':
        return '#3b82f6'; // blue
      case 'relu':
        return '#10b981'; // green
      case 'tanh':
        return '#7c3aed'; // purple
      default:
        return '#6366f1'; // indigo
    }
  }

  function getActivationBg(activation: string): string {
    switch (activation) {
      case 'sigmoid':
        return '#eff6ff'; // blue-50
      case 'relu':
        return '#d1fae5'; // green-50
      case 'tanh':
        return '#f5f3ff'; // purple-50
      default:
        return '#eef2ff'; // indigo-50
    }
  }
</script>


<div class="metrics-panel">
  <div class="metric-card initial">
    <h3>Initial (Step 0)</h3>
    <div class="metric-row">
      <span class="label">w₁₀:</span>
      <span class="value">{formatNumber(metrics?.initial?.params?.w?.[0] ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">w₂₀:</span>
      <span class="value">{formatNumber(metrics?.initial?.params?.w?.[1] ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">b₀:</span>
      <span class="value">{formatNumber(metrics?.initial?.params?.b ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">loss₀:</span>
      <span class="value">{formatNumber(metrics?.initial?.loss ?? 0)}</span>
    </div>
  </div>

  <div class="metric-card current" style="--activation-color: {getActivationColor(metrics?.current?.activation ?? 'unknown')}; --activation-bg: {getActivationBg(metrics?.current?.activation ?? 'unknown')}">
    <h3>Current (Step {metrics?.current?.step ?? 0})</h3>
    <div class="metric-row">
      <span class="label">w₁:</span>
      <span class="value">{formatNumber(metrics?.current?.params?.w?.[0] ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">w₂:</span>
      <span class="value">{formatNumber(metrics?.current?.params?.w?.[1] ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">b:</span>
      <span class="value">{formatNumber(metrics?.current?.params?.b ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">z:</span>
      <span class="value">{formatNumber(metrics?.current?.z ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">a:</span>
      <span class="value">{formatNumber(metrics?.current?.a ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">σ'(z):</span>
      <span class="value">{formatNumber(metrics?.current?.local_derivative ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">loss:</span>
      <span class="value">{formatNumber(metrics?.current?.loss ?? 0)}</span>
    </div>
  </div>

  <div class="metric-card final">
    <h3>Final (Step {metrics?.final?.step ?? 0})</h3>
    <div class="metric-row">
      <span class="label">w₁:</span>
      <span class="value">{formatNumber(metrics?.final?.params?.w?.[0] ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">w₂:</span>
      <span class="value">{formatNumber(metrics?.final?.params?.w?.[1] ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">b:</span>
      <span class="value">{formatNumber(metrics?.final?.params?.b ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">loss:</span>
      <span class="value">{formatNumber(metrics?.final?.loss ?? 0)}</span>
    </div>
  </div>

  <div class="metric-card delta">
    <h3>Progress</h3>
    <div class="metric-row">
      <span class="label">Δw₁:</span>
      <span class="value">{formatNumber(metrics?.deltaParams?.w1 ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Δw₂:</span>
      <span class="value">{formatNumber(metrics?.deltaParams?.w2 ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Δb:</span>
      <span class="value">{formatNumber(metrics?.deltaParams?.b ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Δloss:</span>
      <span class="value">{formatNumber(metrics?.deltaLoss ?? 0)}</span>
    </div>
    <div class="metric-row">
      <span class="label">Progress:</span>
      <span class="value">{formatPercent(metrics?.progress ?? 0)}</span>
    </div>
  </div>

  <div class="metric-card saturation" class:in-saturation={metrics?.current?.in_saturation_zone ?? false}>
    <h3>Saturation Status</h3>
    <div class="saturation-indicator">
      {#if metrics?.current?.in_saturation_zone}
        <span class="status-icon warning">⚠️</span>
        <span class="status-text">In Saturation Zone</span>
      {:else}
        <span class="status-icon healthy">✓</span>
        <span class="status-text">Active Region</span>
      {/if}
    </div>
    <div class="metric-row">
      <span class="label">Activation:</span>
      <span class="value">{metrics?.current?.activation}</span>
    </div>
    <div class="metric-row">
      <span class="label">|σ'(z)|:</span>
      <span class="value">{formatNumber(Math.abs(metrics?.current?.local_derivative ?? 0))}</span>
    </div>
    <div class="saturation-explanation">
      {#if metrics?.current?.in_saturation_zone}
        <p class="warning-text">Gradient magnitude is very small (&lt; 0.01). This can cause vanishing gradients!</p>
      {:else}
        <p class="healthy-text">Neuron is in active learning region.</p>
      {/if}
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
    border-color: var(--activation-color, #dc2626);
    background: var(--activation-bg, #fef2f2);
  }

  .metric-card.final {
    border-color: #16a34a;
    background: #f0fdf4;
  }

  .metric-card.delta {
    border-color: #7c3aed;
    background: #faf5ff;
  }

  .metric-card.saturation {
    border-color: #64748b;
    background: #f8fafc;
  }

  .metric-card.saturation.in-saturation {
    border-color: #dc2626;
    background: #fef2f2;
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
    color: var(--activation-color, #dc2626);
  }

  .final h3 {
    color: #16a34a;
  }

  .delta h3 {
    color: #7c3aed;
  }

  .saturation h3 {
    color: #64748b;
  }

  .saturation.in-saturation h3 {
    color: #dc2626;
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

  .saturation-indicator {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 1rem;
    margin: 0.5rem 0;
    border-radius: 6px;
    background: white;
  }

  .status-icon {
    font-size: 1.5rem;
  }

  .status-text {
    font-weight: 600;
    font-size: 1rem;
  }

  .status-icon.warning + .status-text {
    color: #dc2626;
  }

  .status-icon.healthy + .status-text {
    color: #16a34a;
  }

  .saturation-explanation {
    margin-top: 0.75rem;
    padding: 0.75rem;
    border-radius: 4px;
    font-size: 0.85rem;
    line-height: 1.4;
  }

  .warning-text {
    margin: 0;
    color: #991b1b;
    background: #fee2e2;
    padding: 0.5rem;
    border-radius: 4px;
  }

  .healthy-text {
    margin: 0;
    color: #166534;
    background: #dcfce7;
    padding: 0.5rem;
    border-radius: 4px;
  }

  @media (max-width: 768px) {
    .metrics-panel {
      grid-template-columns: 1fr;
    }
  }
</style>
