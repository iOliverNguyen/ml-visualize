<script lang="ts">
  import { tweened } from 'svelte/motion';
  import { cubicOut } from 'svelte/easing';
  import type { LossGrid, LinearSnapshot } from './types';

  interface Props {
    lossGrid: LossGrid;
    snapshots: LinearSnapshot[];
    currentStep: number;
    fieldDensity?: 'coarse' | 'fine';
    onStepClick?: (step: number) => void;
  }

  let {
    lossGrid,
    snapshots,
    currentStep,
    fieldDensity = 'coarse',
    onStepClick
  }: Props = $props();

  interface GradientVector {
    w1: number;
    w2: number;
    grad_w1: number;
    grad_w2: number;
    magnitude: number;
  }

  // Hover state
  let hoveredStep = $state<number | null>(null);

  // Dimensions
  const width = 800;
  const height = 800;
  const padding = { left: 60, right: 30, top: 30, bottom: 60 };
  const plotWidth = width - padding.left - padding.right;
  const plotHeight = height - padding.top - padding.bottom;

  // Animated current marker
  const markerX = tweened(0, { duration: 150, easing: cubicOut });
  const markerY = tweened(0, { duration: 150, easing: cubicOut });

  // Scaling functions
  function scaleW1(w1: number): number {
    return (
      padding.left +
      ((w1 - lossGrid.w1_min) / (lossGrid.w1_max - lossGrid.w1_min)) *
        plotWidth
    );
  }

  function scaleW2(w2: number): number {
    return (
      height -
      padding.bottom -
      ((w2 - lossGrid.w2_min) / (lossGrid.w2_max - lossGrid.w2_min)) *
        plotHeight
    );
  }

  // Compute gradient vectors from loss grid using finite differences
  let gradientVectors = $derived.by(() => {
    const res = lossGrid.resolution;
    const gridSize = fieldDensity === 'coarse' ? 10 : 20;
    const step = Math.floor(res / gridSize);

    const vectors: GradientVector[] = [];

    // Create lookup map for loss values
    const lossMap = new Map<string, number>();
    lossGrid.points.forEach((p) => {
      const key = `${p.w1.toFixed(6)},${p.w2.toFixed(6)}`;
      lossMap.set(key, p.loss);
    });

    function getLoss(w1: number, w2: number): number | null {
      const key = `${w1.toFixed(6)},${w2.toFixed(6)}`;
      return lossMap.get(key) ?? null;
    }

    // Sample grid points (skip edges to compute finite differences)
    for (let i = step; i < res - step; i += step) {
      for (let j = step; j < res - step; j += step) {
        // Calculate (w1, w2) for this grid position
        const w1 =
          lossGrid.w1_min +
          (i / (res - 1)) * (lossGrid.w1_max - lossGrid.w1_min);
        const w2 =
          lossGrid.w2_min +
          (j / (res - 1)) * (lossGrid.w2_max - lossGrid.w2_min);

        // Neighbor positions for central differences
        const w1_next =
          lossGrid.w1_min +
          ((i + step) / (res - 1)) * (lossGrid.w1_max - lossGrid.w1_min);
        const w1_prev =
          lossGrid.w1_min +
          ((i - step) / (res - 1)) * (lossGrid.w1_max - lossGrid.w1_min);
        const w2_next =
          lossGrid.w2_min +
          ((j + step) / (res - 1)) * (lossGrid.w2_max - lossGrid.w2_min);
        const w2_prev =
          lossGrid.w2_min +
          ((j - step) / (res - 1)) * (lossGrid.w2_max - lossGrid.w2_min);

        // Get loss values
        const loss_w1_next = getLoss(w1_next, w2);
        const loss_w1_prev = getLoss(w1_prev, w2);
        const loss_w2_next = getLoss(w1, w2_next);
        const loss_w2_prev = getLoss(w1, w2_prev);

        // Central difference
        if (
          loss_w1_next !== null &&
          loss_w1_prev !== null &&
          loss_w2_next !== null &&
          loss_w2_prev !== null
        ) {
          const grad_w1 = (loss_w1_next - loss_w1_prev) / (w1_next - w1_prev);
          const grad_w2 = (loss_w2_next - loss_w2_prev) / (w2_next - w2_prev);
          const magnitude = Math.sqrt(grad_w1 ** 2 + grad_w2 ** 2);

          vectors.push({ w1, w2, grad_w1, grad_w2, magnitude });
        }
      }
    }

    return vectors;
  });

  // Find gradient magnitude range for color scaling
  let gradientMagRange = $derived.by(() => {
    const mags = gradientVectors.map((v) => v.magnitude);
    return {
      min: Math.min(...mags),
      max: Math.max(...mags)
    };
  });

  // Map magnitude to color
  function magnitudeToColor(mag: number): string {
    const { min, max } = gradientMagRange;
    const normalized = (mag - min) / (max - min);

    // Blue (low) → Red (high)
    const hue = 240 - normalized * 240; // 240° (blue) to 0° (red)
    return `hsl(${hue}, 80%, 50%)`;
  }

  // Generate arrow rendering data
  function getArrowData(vector: GradientVector, maxLength: number = 30): {
    x1: number;
    y1: number;
    x2: number;
    y2: number;
    color: string;
  } {
    const x1 = scaleW1(vector.w1);
    const y1 = scaleW2(vector.w2);

    // Descent direction = -gradient
    // Normalize and scale
    const scale = maxLength / gradientMagRange.max;
    const dx = -vector.grad_w1 * scale;
    const dy = -vector.grad_w2 * scale;

    // Scale dy for SVG (flip y-axis)
    const x2 = x1 + dx;
    const y2 = y1 - dy;

    const color = magnitudeToColor(vector.magnitude);

    return { x1, y1, x2, y2, color };
  }

  // Update marker animation
  $effect(() => {
    const current = snapshots[currentStep];
    markerX.set(scaleW1(current.w1));
    markerY.set(scaleW2(current.w2));
  });

  // Trajectory path
  let pathData = $derived(
    snapshots
      .map(
        (s, i) => `${i === 0 ? 'M' : 'L'} ${scaleW1(s.w1)} ${scaleW2(s.w2)}`
      )
      .join(' ')
  );

  // Axis ticks
  let w1Ticks = $derived(
    Array.from(
      { length: 5 },
      (_, i) =>
        lossGrid.w1_min +
        (i * (lossGrid.w1_max - lossGrid.w1_min)) / 4
    )
  );
  let w2Ticks = $derived(
    Array.from(
      { length: 5 },
      (_, i) =>
        lossGrid.w2_min +
        (i * (lossGrid.w2_max - lossGrid.w2_min)) / 4
    )
  );
</script>

<!-- Density toggle control -->
<div class="controls">
    <label>
        <input type="radio" value="coarse" bind:group={fieldDensity} />
        Coarse (10×10)
    </label>
    <label>
        <input type="radio" value="fine" bind:group={fieldDensity} />
        Fine (20×20)
    </label>
</div>

<svg viewBox="0 0 {width} {height}" width="100%" height="auto" class="plot-svg">
  <!-- Grid lines (light) -->
  {#each w1Ticks as tick}
    <line
      x1={scaleW1(tick)}
      y1={padding.top}
      x2={scaleW1(tick)}
      y2={height - padding.bottom}
      stroke="#f1f5f9"
      stroke-width="1"
    />
  {/each}
  {#each w2Ticks as tick}
    <line
      x1={padding.left}
      y1={scaleW2(tick)}
      x2={width - padding.right}
      y2={scaleW2(tick)}
      stroke="#f1f5f9"
      stroke-width="1"
    />
  {/each}

  <!-- Gradient vector arrows -->
  {#each gradientVectors as vector}
    {@const arrow = getArrowData(vector, 25)}
    <g opacity="0.7">
      <line
        x1={arrow.x1}
        y1={arrow.y1}
        x2={arrow.x2}
        y2={arrow.y2}
        stroke={arrow.color}
        stroke-width="2"
        marker-end="url(#arrowhead)"
      />
    </g>
  {/each}

  <!-- Trajectory path -->
  <path
    d={pathData}
    fill="none"
    stroke="#000"
    stroke-width="3"
    stroke-opacity="0.6"
  />

  <!-- Trajectory markers -->
  {#each snapshots as snapshot, i}
    {#if i === 0 || i === snapshots.length - 1 || i % 10 === 0}
      {@const isStart = i === 0}
      {@const isEnd = i === snapshots.length - 1}

      <circle
        cx={scaleW1(snapshot.w1)}
        cy={scaleW2(snapshot.w2)}
        r={isStart || isEnd ? 7 : 4}
        fill={isStart ? '#16a34a' : isEnd ? '#dc2626' : '#fff'}
        stroke="#000"
        stroke-width="2"
        class="trajectory-marker"
        onmouseenter={() => (hoveredStep = i)}
        onmouseleave={() => (hoveredStep = null)}
        onclick={() => onStepClick?.(i)}
      />
    {/if}
  {/each}

  <!-- Current position marker -->
  <circle
    cx={$markerX}
    cy={$markerY}
    r="10"
    fill="#2563eb"
    stroke="white"
    stroke-width="3"
  />
  <circle cx={$markerX} cy={$markerY} r="15" fill="none" stroke="#2563eb" stroke-width="2" opacity="0.4">
    <animate
      attributeName="r"
      values="15;20;15"
      dur="1.5s"
      repeatCount="indefinite"
    />
  </circle>

  <!-- Axes -->
  <line
    x1={padding.left}
    y1={height - padding.bottom}
    x2={width - padding.right}
    y2={height - padding.bottom}
    stroke="#334155"
    stroke-width="2"
  />
  <line
    x1={padding.left}
    y1={padding.top}
    x2={padding.left}
    y2={height - padding.bottom}
    stroke="#334155"
    stroke-width="2"
  />

  <!-- Axis labels -->
  {#each w1Ticks as tick}
    <text
      x={scaleW1(tick)}
      y={height - padding.bottom + 20}
      text-anchor="middle"
      font-size="12"
      fill="#334155"
    >
      {tick.toFixed(1)}
    </text>
  {/each}
  {#each w2Ticks as tick}
    <text
      x={padding.left - 10}
      y={scaleW2(tick)}
      text-anchor="end"
      dominant-baseline="middle"
      font-size="12"
      fill="#334155"
    >
      {tick.toFixed(1)}
    </text>
  {/each}

  <!-- Axis titles -->
  <text
    x={width / 2}
    y={height - 10}
    text-anchor="middle"
    font-size="14"
    font-weight="600"
    fill="#334155"
  >
    w₁
  </text>
  <text
    x={20}
    y={height / 2}
    text-anchor="middle"
    transform="rotate(-90, 20, {height / 2})"
    font-size="14"
    font-weight="600"
    fill="#334155"
  >
    w₂
  </text>

  <!-- Legend (gradient magnitude color scale) -->
  <g transform="translate({width - 130}, 50)">
    <text x="0" y="-10" font-size="12" fill="#334155" font-weight="600">
      ||∇L||
    </text>
    <!-- Color gradient samples -->
    {#each Array.from({ length: 5 }, (_, i) => i) as i}
      {@const mag = gradientMagRange.min + (i / 4) * (gradientMagRange.max - gradientMagRange.min)}
      <line
        x1="0"
        y1={i * 15}
        x2="30"
        y2={i * 15}
        stroke={magnitudeToColor(mag)}
        stroke-width="4"
      />
      <text
        x="35"
        y={i * 15}
        font-size="10"
        fill="#64748b"
        dominant-baseline="middle"
      >
        {mag.toFixed(1)}
      </text>
    {/each}
  </g>

  <!-- Hover tooltip -->
  {#if hoveredStep !== null}
    {@const s = snapshots[hoveredStep]}
    {@const x = scaleW1(s.w1)}
    {@const y = scaleW2(s.w2)}
    {@const tooltipX = x > width - 180 ? x - 160 : x + 10}
    {@const tooltipY = y > height - 110 ? y - 110 : y + 10}

    <g class="tooltip">
      <rect
        x={tooltipX}
        y={tooltipY}
        width="150"
        height="100"
        fill="white"
        stroke="#cbd5e1"
        stroke-width="1"
        rx="4"
      />
      <text
        x={tooltipX + 10}
        y={tooltipY + 18}
        font-size="12"
        fill="#475569"
        font-weight="600"
      >
        Step {hoveredStep}
      </text>
      <text
        x={tooltipX + 10}
        y={tooltipY + 36}
        font-size="11"
        fill="#64748b"
        font-family="monospace"
      >
        w₁: {s.w1.toFixed(4)}
      </text>
      <text
        x={tooltipX + 10}
        y={tooltipY + 52}
        font-size="11"
        fill="#64748b"
        font-family="monospace"
      >
        w₂: {s.w2.toFixed(4)}
      </text>
      <text
        x={tooltipX + 10}
        y={tooltipY + 68}
        font-size="11"
        fill="#64748b"
        font-family="monospace"
      >
        ||∇L||: {s.gradient_magnitude.toFixed(2)}
      </text>
      <text
        x={tooltipX + 10}
        y={tooltipY + 84}
        font-size="11"
        fill="#64748b"
        font-family="monospace"
      >
        Loss: {s.loss.toFixed(2)}
      </text>
    </g>
  {/if}

  <!-- Arrowhead marker definition -->
  <defs>
    <marker
      id="arrowhead"
      markerWidth="8"
      markerHeight="8"
      refX="6"
      refY="4"
      orient="auto"
    >
      <polygon points="0 0, 8 4, 0 8" fill="currentColor" />
    </marker>
  </defs>
</svg>

<style>
  .plot-svg {
    border: 1px solid #eee;
    border-radius: 4px;
    background: white;
    width: 100%;
    display: block;
  }

  .trajectory-marker {
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .trajectory-marker:hover {
    filter: brightness(1.2);
  }

  .tooltip {
    pointer-events: none;
  }

  .controls {
    display: flex;
    gap: 1rem;
    justify-content: center;
    margin-top: 1rem;
    padding: 0.5rem;
    background: #f8fafc;
    border-radius: 4px;
  }

  .controls label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 14px;
    color: #334155;
    cursor: pointer;
  }
</style>
