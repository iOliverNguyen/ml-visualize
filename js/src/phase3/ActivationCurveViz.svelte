<script lang="ts">
  import { tweened } from 'svelte/motion';
  import { cubicOut } from 'svelte/easing';
  import type { NeuronSnapshot } from './types';

  interface Props {
    snapshot: NeuronSnapshot;
  }

  let { snapshot }: Props = $props();

  // Animated z marker position
  const zPosition = tweened(0, {
    duration: 50,
    easing: cubicOut
  });

  // Update z position when snapshot changes
  $effect(() => {
    zPosition.set(snapshot?.z ?? 0);
  });

  // SVG dimensions
  const width = 600;
  const height = 450; // 200px per panel + 50px gap
  const panelHeight = 200;
  const margin = { top: 20, right: 30, bottom: 40, left: 60 };
  const plotWidth = width - margin.left - margin.right;
  const plotHeight = panelHeight - margin.top - margin.bottom;

  // Z range for visualization
  const zMin = -10;
  const zMax = 10;
  const numPoints = 100;

  // Activation functions
  function sigmoid(z: number): number {
    return 1 / (1 + Math.exp(-z));
  }

  function sigmoidDerivative(z: number): number {
    const s = sigmoid(z);
    return s * (1 - s);
  }

  function relu(z: number): number {
    return Math.max(0, z);
  }

  function reluDerivative(z: number): number {
    return z > 0 ? 1 : 0;
  }

  function tanh(z: number): number {
    return Math.tanh(z);
  }

  function tanhDerivative(z: number): number {
    const t = Math.tanh(z);
    return 1 - t * t;
  }

  // Get activation function based on type
  function getActivationFn(type: string): (z: number) => number {
    switch (type) {
      case 'sigmoid':
        return sigmoid;
      case 'relu':
        return relu;
      case 'tanh':
        return tanh;
      default:
        return sigmoid;
    }
  }

  function getDerivativeFn(type: string): (z: number) => number {
    switch (type) {
      case 'sigmoid':
        return sigmoidDerivative;
      case 'relu':
        return reluDerivative;
      case 'tanh':
        return tanhDerivative;
      default:
        return sigmoidDerivative;
    }
  }

  // Generate curve data
  let curveData = $derived.by(() => {
    const activationFn = getActivationFn(snapshot?.activation ?? 'sigmoid');
    const derivativeFn = getDerivativeFn(snapshot?.activation ?? 'sigmoid');
    const points = [];
    const step = (zMax - zMin) / (numPoints - 1);

    for (let i = 0; i < numPoints; i++) {
      const z = zMin + i * step;
      points.push({
        z,
        activation: activationFn(z),
        derivative: derivativeFn(z)
      });
    }
    return points;
  });

  // Y range for activation (depends on function)
  let yActivationRange = $derived.by(() => {
    if (snapshot?.activation === 'relu') {
      return { min: 0, max: 10 };
    } else if (snapshot?.activation === 'tanh') {
      return { min: -1.2, max: 1.2 };
    } else {
      // sigmoid
      return { min: 0, max: 1 };
    }
  });

  // Y range for derivative
  let yDerivativeRange = $derived.by(() => {
    if (snapshot?.activation === 'relu') {
      return { min: 0, max: 1.2 };
    } else if (snapshot?.activation === 'tanh') {
      return { min: 0, max: 1 };
    } else {
      // sigmoid
      return { min: 0, max: 0.3 };
    }
  });

  // Scale functions
  function scaleX(z: number): number {
    return margin.left + ((z - zMin) / (zMax - zMin)) * plotWidth;
  }

  function scaleYActivation(y: number): number {
    const range = yActivationRange;
    return margin.top + (1 - (y - range.min) / (range.max - range.min)) * plotHeight;
  }

  function scaleYDerivative(y: number): number {
    const range = yDerivativeRange;
    return panelHeight + 50 + margin.top + (1 - (y - range.min) / (range.max - range.min)) * plotHeight;
  }

  // Generate SVG paths
  let activationPath = $derived.by(() => {
    return curveData.map((d, i) => {
      const x = scaleX(d.z);
      const y = scaleYActivation(d.activation);
      return `${i === 0 ? 'M' : 'L'} ${x} ${y}`;
    }).join(' ');
  });

  let derivativePath = $derived.by(() => {
    return curveData.map((d, i) => {
      const x = scaleX(d.z);
      const y = scaleYDerivative(d.derivative);
      return `${i === 0 ? 'M' : 'L'} ${x} ${y}`;
    }).join(' ');
  });

  // Saturation zones (|derivative| < 0.01)
  let saturationZones = $derived.by(() => {
    const zones: { zStart: number; zEnd: number }[] = [];
    let inZone = false;
    let zStart = 0;

    for (const point of curveData) {
      const isSaturated = Math.abs(point.derivative) < 0.01;

      if (isSaturated && !inZone) {
        // Start new zone
        zStart = point.z;
        inZone = true;
      } else if (!isSaturated && inZone) {
        // End zone
        zones.push({ zStart, zEnd: point.z });
        inZone = false;
      }
    }

    // Close any open zone
    if (inZone) {
      zones.push({ zStart, zEnd: zMax });
    }

    return zones;
  });

  // Get activation color
  function getActivationColor(activation: string): string {
    switch (activation) {
      case 'sigmoid':
        return '#3b82f6'; // blue
      case 'relu':
        return '#10b981'; // green
      case 'tanh':
        return '#7c3aed'; // purple
      default:
        return '#6366f1';
    }
  }

  let color = $derived(getActivationColor(snapshot?.activation ?? 'sigmoid'));
</script>

<div class="activation-viz">
  <h2>Activation Function: {snapshot?.activation ?? 'sigmoid'}</h2>

  <svg {width} {height} class="chart">
    <!-- Activation Panel -->
    <g class="activation-panel">
      <text x={margin.left / 2} y={margin.top + plotHeight / 2} class="axis-label" transform="rotate(-90, {margin.left / 2}, {margin.top + plotHeight / 2})">
        a = {snapshot?.activation ?? 'sigmoid'}(z)
      </text>

      <!-- Grid lines -->
      {#each [0, 0.25, 0.5, 0.75, 1] as yVal}
        {#if scaleYActivation(yVal) >= margin.top && scaleYActivation(yVal) <= margin.top + plotHeight}
          <line x1={margin.left} y1={scaleYActivation(yVal)} x2={width - margin.right} y2={scaleYActivation(yVal)} stroke="#e5e7eb" stroke-width="1" stroke-dasharray="2,2" />
          <text x={margin.left - 5} y={scaleYActivation(yVal) + 4} text-anchor="end" class="tick-label">{yVal.toFixed(2)}</text>
        {/if}
      {/each}

      <!-- Activation curve -->
      <path d={activationPath} stroke={color} stroke-width="3" fill="none" />

      <!-- Current z marker -->
      <line x1={scaleX($zPosition)} y1={margin.top} x2={scaleX($zPosition)} y2={margin.top + plotHeight} stroke={color} stroke-width="2" stroke-dasharray="4,4" opacity="0.7" />

      <!-- Current point circle -->
      <circle cx={scaleX($zPosition)} cy={scaleYActivation(getActivationFn(snapshot?.activation ?? 'sigmoid')($zPosition))} r="6" fill={color} stroke="white" stroke-width="2" />
    </g>

    <!-- Derivative Panel -->
    <g class="derivative-panel">
      <text x={margin.left / 2} y={panelHeight + 50 + margin.top + plotHeight / 2} class="axis-label" transform="rotate(-90, {margin.left / 2}, {panelHeight + 50 + margin.top + plotHeight / 2})">
        σ'(z)
      </text>

      <!-- Grid lines -->
      {#each [0, 0.1, 0.2, 0.3] as yVal}
        {#if scaleYDerivative(yVal) >= panelHeight + 50 + margin.top && scaleYDerivative(yVal) <= panelHeight + 50 + margin.top + plotHeight}
          <line x1={margin.left} y1={scaleYDerivative(yVal)} x2={width - margin.right} y2={scaleYDerivative(yVal)} stroke="#e5e7eb" stroke-width="1" stroke-dasharray="2,2" />
          <text x={margin.left - 5} y={scaleYDerivative(yVal) + 4} text-anchor="end" class="tick-label">{yVal.toFixed(2)}</text>
        {/if}
      {/each}

      <!-- Saturation zones -->
      {#each saturationZones as zone}
        <rect
          x={scaleX(zone.zStart)}
          y={panelHeight + 50 + margin.top}
          width={scaleX(zone.zEnd) - scaleX(zone.zStart)}
          height={plotHeight}
          fill="#dc2626"
          opacity="0.1"
        />
      {/each}

      <!-- Derivative curve -->
      <path d={derivativePath} stroke={color} stroke-width="3" fill="none" />

      <!-- Current z marker -->
      <line x1={scaleX($zPosition)} y1={panelHeight + 50 + margin.top} x2={scaleX($zPosition)} y2={panelHeight + 50 + margin.top + plotHeight} stroke={color} stroke-width="2" stroke-dasharray="4,4" opacity="0.7" />

      <!-- Current point circle -->
      <circle cx={scaleX($zPosition)} cy={scaleYDerivative(getDerivativeFn(snapshot?.activation ?? 'sigmoid')($zPosition))} r="6" fill={color} stroke="white" stroke-width="2" />

      <!-- Saturation threshold line -->
      <line x1={margin.left} y1={scaleYDerivative(0.01)} x2={width - margin.right} y2={scaleYDerivative(0.01)} stroke="#dc2626" stroke-width="1" stroke-dasharray="5,5" opacity="0.5" />
      <text x={width - margin.right + 5} y={scaleYDerivative(0.01) + 4} class="tick-label" fill="#dc2626">0.01</text>
    </g>

    <!-- X axis (shared) -->
    <g class="x-axis">
      <line x1={margin.left} y1={height - margin.bottom} x2={width - margin.right} y2={height - margin.bottom} stroke="#94a3b8" stroke-width="2" />

      {#each [-10, -5, 0, 5, 10] as zVal}
        <line x1={scaleX(zVal)} y1={height - margin.bottom} x2={scaleX(zVal)} y2={height - margin.bottom + 6} stroke="#94a3b8" stroke-width="2" />
        <text x={scaleX(zVal)} y={height - margin.bottom + 20} text-anchor="middle" class="tick-label">{zVal}</text>
      {/each}

      <text x={width / 2} y={height - 5} text-anchor="middle" class="axis-label">z (pre-activation)</text>
    </g>
  </svg>

  <!-- Status display -->
  <div class="status-panel" class:saturated={snapshot?.in_saturation_zone ?? false}>
    <div class="status-row">
      <span class="label">Current z:</span>
      <span class="value">{snapshot?.z?.toFixed(4) ?? '0.0000'}</span>
    </div>
    <div class="status-row">
      <span class="label">Current σ'(z):</span>
      <span class="value">{snapshot?.local_derivative?.toFixed(4) ?? '0.0000'}</span>
    </div>
    {#if snapshot?.in_saturation_zone}
      <div class="warning">
        ⚠️ Saturation Zone: |σ'(z)| &lt; 0.01 — Vanishing gradient!
      </div>
    {:else}
      <div class="healthy">
        ✓ Active Region: Healthy gradient flow
      </div>
    {/if}
  </div>
</div>

<style>
  .activation-viz {
    background: white;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    padding: 1.5rem;
  }

  h2 {
    margin: 0 0 1rem 0;
    font-size: 1.25rem;
    color: #1e293b;
    text-align: center;
  }

  .chart {
    display: block;
    margin: 0 auto;
  }

  .axis-label {
    font-size: 12px;
    font-weight: 600;
    fill: #64748b;
  }

  .tick-label {
    font-size: 11px;
    fill: #94a3b8;
  }

  .status-panel {
    margin-top: 1rem;
    padding: 1rem;
    background: #f8fafc;
    border-radius: 6px;
    border: 2px solid #e2e8f0;
  }

  .status-panel.saturated {
    background: #fef2f2;
    border-color: #dc2626;
  }

  .status-row {
    display: flex;
    justify-content: space-between;
    padding: 0.25rem 0;
  }

  .status-row .label {
    font-weight: 500;
    color: #64748b;
  }

  .status-row .value {
    font-family: 'Courier New', monospace;
    font-weight: bold;
    color: #1e293b;
  }

  .warning {
    margin-top: 0.75rem;
    padding: 0.75rem;
    background: #fee2e2;
    border: 1px solid #dc2626;
    border-radius: 4px;
    color: #991b1b;
    font-weight: 600;
    font-size: 0.9rem;
  }

  .healthy {
    margin-top: 0.75rem;
    padding: 0.75rem;
    background: #dcfce7;
    border: 1px solid #16a34a;
    border-radius: 4px;
    color: #166534;
    font-weight: 600;
    font-size: 0.9rem;
  }

  @media (max-width: 768px) {
    .chart {
      width: 100%;
      height: auto;
    }
  }
</style>
