<script lang="ts">
  import { tweened } from 'svelte/motion';
  import { cubicOut } from 'svelte/easing';
  import type { LinearSnapshot, LossGrid } from './types';

  interface Props {
    snapshots: LinearSnapshot[];
    currentStep: number;
    lossGrid: LossGrid;
    onStepClick?: (step: number) => void;
  }

  let { snapshots, currentStep, lossGrid, onStepClick }: Props = $props();

  // Hover state
  let hoveredStep = $state<number | null>(null);

  // SVG dimensions (square for accurate 2D visualization)
  const width = 800;
  const height = 800;
  const padding = { left: 60, right: 40, top: 40, bottom: 60 };
  const plotWidth = width - padding.left - padding.right;
  const plotHeight = height - padding.top - padding.bottom;

  // Use loss grid bounds as fixed axis ranges
  const w1Min = lossGrid.w1_min;
  const w1Max = lossGrid.w1_max;
  const w2Min = lossGrid.w2_min;
  const w2Max = lossGrid.w2_max;

  // Scaling functions
  function scaleW1(w1: number): number {
    return padding.left + ((w1 - w1Min) / (w1Max - w1Min)) * plotWidth;
  }

  function scaleW2(w2: number): number {
    // Flip Y axis (w2 increases upward)
    return (
      height - padding.bottom - ((w2 - w2Min) / (w2Max - w2Min)) * plotHeight
    );
  }

  // Animated current marker
  const markerX = tweened(0, { duration: 150, easing: cubicOut });
  const markerY = tweened(0, { duration: 150, easing: cubicOut });

  $effect(() => {
    const current = snapshots[currentStep];
    markerX.set(scaleW1(current.w1));
    markerY.set(scaleW2(current.w2));
  });

  // Generate trajectory path
  let pathData = $derived(
    snapshots
      .map(
        (s, i) =>
          `${i === 0 ? 'M' : 'L'} ${scaleW1(s.w1)} ${scaleW2(s.w2)}`
      )
      .join(' ')
  );

  // Axis ticks
  let w1Ticks = $derived(
    Array.from({ length: 6 }, (_, i) => w1Min + (i * (w1Max - w1Min)) / 5)
  );
  let w2Ticks = $derived(
    Array.from({ length: 6 }, (_, i) => w2Min + (i * (w2Max - w2Min)) / 5)
  );
</script>

<svg viewBox="0 0 {width} {height}" width="100%" height="auto" class="plot-svg">
  <!-- Grid lines -->
  {#each w1Ticks as tick}
    <line
      x1={scaleW1(tick)}
      y1={padding.top}
      x2={scaleW1(tick)}
      y2={height - padding.bottom}
      stroke="#e2e8f0"
      stroke-width="1"
    />
  {/each}
  {#each w2Ticks as tick}
    <line
      x1={padding.left}
      y1={scaleW2(tick)}
      x2={width - padding.right}
      y2={scaleW2(tick)}
      stroke="#e2e8f0"
      stroke-width="1"
    />
  {/each}

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

  <!-- Axis ticks and labels (w1 - horizontal) -->
  {#each w1Ticks as tick}
    <line
      x1={scaleW1(tick)}
      y1={height - padding.bottom}
      x2={scaleW1(tick)}
      y2={height - padding.bottom + 5}
      stroke="#94a3b8"
      stroke-width="2"
    />
    <text
      x={scaleW1(tick)}
      y={height - padding.bottom + 20}
      text-anchor="middle"
      font-size="12"
      fill="#64748b"
    >
      {tick.toFixed(1)}
    </text>
  {/each}

  <!-- Axis ticks and labels (w2 - vertical) -->
  {#each w2Ticks as tick}
    <line
      x1={padding.left - 5}
      y1={scaleW2(tick)}
      x2={padding.left}
      y2={scaleW2(tick)}
      stroke="#94a3b8"
      stroke-width="2"
    />
    <text
      x={padding.left - 10}
      y={scaleW2(tick)}
      text-anchor="end"
      dominant-baseline="middle"
      font-size="12"
      fill="#64748b"
    >
      {tick.toFixed(1)}
    </text>
  {/each}

  <!-- Axis labels -->
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

  <!-- Trajectory path -->
  <path
    d={pathData}
    fill="none"
    stroke="#64748b"
    stroke-width="2"
    stroke-opacity="0.5"
  />

  <!-- Parameter points -->
  {#each snapshots as snapshot, i}
    {@const isStart = i === 0}
    {@const isEnd = i === snapshots.length - 1}
    {@const isCurrent = i === currentStep}
    {@const isHovered = i === hoveredStep}

    <circle
      cx={scaleW1(snapshot.w1)}
      cy={scaleW2(snapshot.w2)}
      r={isStart || isEnd ? 8 : isCurrent ? 6 : isHovered ? 5 : 3}
      fill={isStart ? '#16a34a' : isEnd ? '#dc2626' : isCurrent ? '#2563eb' : isHovered ? '#3b82f6' : '#94a3b8'}
      stroke="white"
      stroke-width={isStart || isEnd ? 2 : 1}
      class="data-point"
      onmouseenter={() => (hoveredStep = i)}
      onmouseleave={() => (hoveredStep = null)}
      onclick={() => onStepClick?.(i)}
    />
  {/each}

  <!-- Animated current position marker (pulsing ring) -->
  <circle cx={$markerX} cy={$markerY} r="12" fill="none" stroke="#2563eb" stroke-width="2" opacity="0.6">
    <animate
      attributeName="r"
      values="12;16;12"
      dur="1.5s"
      repeatCount="indefinite"
    />
    <animate
      attributeName="opacity"
      values="0.6;0.2;0.6"
      dur="1.5s"
      repeatCount="indefinite"
    />
  </circle>

  <!-- Hover tooltip -->
  {#if hoveredStep !== null}
    {@const s = snapshots[hoveredStep]}
    {@const x = scaleW1(s.w1)}
    {@const y = scaleW2(s.w2)}
    {@const tooltipX = x > width - 170 ? x - 160 : x + 10}
    {@const tooltipY = y > height - 90 ? y - 90 : y + 10}

    <g class="tooltip">
      <rect
        x={tooltipX}
        y={tooltipY}
        width="150"
        height="80"
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
        Loss: {s.loss.toFixed(4)}
      </text>
    </g>
  {/if}

  <!-- Legend -->
  <g transform="translate(20, 60)">
    <circle cx="0" cy="0" r="6" fill="#16a34a" stroke="white" stroke-width="2" />
    <text x="12" y="4" font-size="12" fill="#334155">Start</text>

    <circle cx="0" cy="20" r="6" fill="#2563eb" stroke="white" stroke-width="1" />
    <text x="12" y="24" font-size="12" fill="#334155">Current</text>

    <circle cx="0" cy="40" r="6" fill="#dc2626" stroke="white" stroke-width="2" />
    <text x="12" y="44" font-size="12" fill="#334155">End</text>
  </g>
</svg>

<style>
  .plot-svg {
    border: 1px solid #eee;
    border-radius: 4px;
    background: white;
    width: 100%;
    display: block;
  }

  .data-point {
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .data-point:hover {
    filter: brightness(1.2);
  }

  .tooltip {
    pointer-events: none;
  }
</style>
