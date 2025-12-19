<script lang="ts">
  import { onMount } from 'svelte';
  import { tweened } from 'svelte/motion';
  import { cubicOut } from 'svelte/easing';
  import type { LossGrid, LinearSnapshot } from './types';

  interface Props {
    lossGrid: LossGrid;
    snapshots: LinearSnapshot[];
    currentStep: number;
    colorScale?: 'linear' | 'log';
    onStepClick?: (step: number) => void;
  }

  let {
    lossGrid,
    snapshots,
    currentStep,
    colorScale = 'linear',
    onStepClick
  }: Props = $props();

  // Canvas reference
  let canvasRef: HTMLCanvasElement;

  // Hover state
  let hoveredStep = $state<number | null>(null);

  // Dimensions
  const width = 800;
  const height = 800;
  const padding = { left: 60, right: 100, top: 30, bottom: 60 }; // Extra right padding for legend
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

  // Color mapping function
  function lossToColor(
    loss: number,
    minLoss: number,
    maxLoss: number
  ): string {
    let normalized: number;

    if (colorScale === 'log') {
      // Log scale for wide loss ranges
      const logMin = Math.log10(minLoss + 1);
      const logMax = Math.log10(maxLoss + 1);
      const logLoss = Math.log10(loss + 1);
      normalized = (logLoss - logMin) / (logMax - logMin);
    } else {
      // Linear scale
      normalized = (loss - minLoss) / (maxLoss - minLoss);
    }

    // Clamp to [0, 1]
    normalized = Math.max(0, Math.min(1, normalized));

    // Color gradient: blue (low) → cyan → yellow → red (high)
    if (normalized < 0.33) {
      // Blue to cyan: H 240->180, S 100%, L 50%
      const t = normalized / 0.33; // [0, 1]
      const h = 240 - t * 60; // 240 (blue) to 180 (cyan)
      return `hsl(${h}, 100%, 50%)`;
    } else if (normalized < 0.67) {
      // Cyan to yellow: H 180->60, S 100%, L 50%
      const t = (normalized - 0.33) / 0.34; // [0, 1]
      const h = 180 - t * 120; // 180 (cyan) to 60 (yellow)
      return `hsl(${h}, 100%, 50%)`;
    } else {
      // Yellow to red: H 60->0, S 100%, L 50%->45%
      const t = (normalized - 0.67) / 0.33; // [0, 1]
      const h = 60 - t * 60; // 60 (yellow) to 0 (red)
      const l = 50 - t * 5; // Darken toward red
      return `hsl(${h}, 100%, ${l}%)`;
    }
  }

  // Render heatmap on canvas
  function renderHeatmap() {
    if (!canvasRef) return;

    const ctx = canvasRef.getContext('2d');
    if (!ctx) return;

    // Set canvas resolution
    canvasRef.width = plotWidth;
    canvasRef.height = plotHeight;

    // Find loss range
    const losses = lossGrid.points.map((p) => p.loss);
    const minLoss = Math.min(...losses);
    const maxLoss = Math.max(...losses);

    // Auto-switch to log scale if range is too large
    const shouldUseLog = maxLoss / minLoss > 1000;
    if (shouldUseLog && colorScale === 'linear') {
      console.log('Auto-switching to log scale due to large loss range');
    }

    const effectiveScale = shouldUseLog ? 'log' : colorScale;

    // Resolution
    const res = lossGrid.resolution;
    const cellWidth = plotWidth / res;
    const cellHeight = plotHeight / res;

    // Render each grid cell
    lossGrid.points.forEach((point) => {
      // Calculate grid position
      const gridX = Math.round(
        ((point.w1 - lossGrid.w1_min) /
          (lossGrid.w1_max - lossGrid.w1_min)) *
          (res - 1)
      );
      const gridY = Math.round(
        ((point.w2 - lossGrid.w2_min) /
          (lossGrid.w2_max - lossGrid.w2_min)) *
          (res - 1)
      );

      // Canvas coordinates (flip Y)
      const canvasX = gridX * cellWidth;
      const canvasY = plotHeight - (gridY + 1) * cellHeight;

      ctx.fillStyle = lossToColor(point.loss, minLoss, maxLoss);
      ctx.fillRect(canvasX, canvasY, cellWidth, cellHeight);
    });
  }

  // Lifecycle
  onMount(() => {
    renderHeatmap();
  });

  // Re-render when colorScale changes
  $effect(() => {
    colorScale; // Track dependency
    renderHeatmap();
  });

  // Update marker animation
  $effect(() => {
    const current = snapshots[currentStep];
    markerX.set(scaleW1(current.w1));
    markerY.set(scaleW2(current.w2));
  });

  // Generate trajectory path
  let pathData = $derived(
    snapshots
      .map(
        (s, i) => `${i === 0 ? 'M' : 'L'} ${scaleW1(s.w1)} ${scaleW2(s.w2)}`
      )
      .join(' ')
  );

  // Color legend gradient
  let legendGradient = $derived(() => {
    const losses = lossGrid.points.map((p) => p.loss);
    const minLoss = Math.min(...losses);
    const maxLoss = Math.max(...losses);

    // Create 20 color stops for smooth gradient
    return Array.from({ length: 20 }, (_, i) => {
      const t = i / 19;
      const loss = minLoss + t * (maxLoss - minLoss);
      return lossToColor(loss, minLoss, maxLoss);
    });
  });

  // Loss range for legend
  let lossRange = $derived(() => {
    const losses = lossGrid.points.map((p) => p.loss);
    return {
      min: Math.min(...losses),
      max: Math.max(...losses)
    };
  });

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

<div class="contour-container">
  <!-- Canvas for heatmap (positioned absolutely) -->
  <canvas
    bind:this={canvasRef}
    style="
      position: absolute;
      left: {padding.left}px;
      top: {padding.top}px;
      width: {plotWidth}px;
      height: {plotHeight}px;
    "
  />

  <!-- SVG overlay for trajectory and axes -->
  <svg
    viewBox="0 0 {width} {height}"
    width="100%"
    height="auto"
    class="overlay-svg"
  >
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

    <!-- Trajectory path with glow effect -->
    <path
      d={pathData}
      fill="none"
      stroke="white"
      stroke-width="4"
      stroke-opacity="0.9"
      filter="url(#glow)"
    />
    <path
      d={pathData}
      fill="none"
      stroke="#000"
      stroke-width="2"
      stroke-opacity="0.7"
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
          fill={isStart ? '#16a34a' : isEnd ? '#dc2626' : 'white'}
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
      filter="url(#glow)"
    />

    <!-- Color legend (vertical gradient bar) -->
    <defs>
      <linearGradient id="legend-gradient" x1="0%" y1="100%" x2="0%" y2="0%">
        {#each legendGradient() as color, i}
          <stop
            offset="{(i / (legendGradient().length - 1)) * 100}%"
            stop-color={color}
          />
        {/each}
      </linearGradient>

      <!-- Glow filter for trajectory -->
      <filter id="glow">
        <feGaussianBlur stdDeviation="2" result="coloredBlur" />
        <feMerge>
          <feMergeNode in="coloredBlur" />
          <feMergeNode in="SourceGraphic" />
        </feMerge>
      </filter>
    </defs>

    <!-- Legend bar -->
    <rect
      x={width - padding.right + 20}
      y={padding.top}
      width="30"
      height={plotHeight}
      fill="url(#legend-gradient)"
      stroke="#334155"
      stroke-width="1"
    />

    <!-- Legend labels -->
    <text
      x={width - padding.right + 55}
      y={padding.top + 5}
      font-size="11"
      fill="#334155"
      dominant-baseline="hanging"
    >
      {lossRange().max.toFixed(1)}
    </text>
    <text
      x={width - padding.right + 55}
      y={height - padding.bottom - 5}
      font-size="11"
      fill="#334155"
      dominant-baseline="baseline"
    >
      {lossRange().min.toFixed(2)}
    </text>
    <text
      x={width - padding.right + 55}
      y={height / 2}
      font-size="10"
      fill="#64748b"
      dominant-baseline="middle"
    >
      Loss
    </text>

    <!-- Hover tooltip -->
    {#if hoveredStep !== null}
      {@const s = snapshots[hoveredStep]}
      {@const x = scaleW1(s.w1)}
      {@const y = scaleW2(s.w2)}
      {@const tooltipX = x > width - 180 ? x - 140 : x + 15}
      {@const tooltipY = y > height - 90 ? y - 85 : y + 15}

      <g class="tooltip">
        <rect
          x={tooltipX}
          y={tooltipY}
          width="125"
          height="75"
          fill="rgba(255, 255, 255, 0.95)"
          stroke="#000"
          stroke-width="2"
          rx="4"
        />
        <text
          x={tooltipX + 10}
          y={tooltipY + 18}
          font-size="12"
          fill="#000"
          font-weight="600"
        >
          Step {hoveredStep}
        </text>
        <text
          x={tooltipX + 10}
          y={tooltipY + 36}
          font-size="11"
          fill="#334155"
          font-family="monospace"
        >
          w₁: {s.w1.toFixed(3)}
        </text>
        <text
          x={tooltipX + 10}
          y={tooltipY + 52}
          font-size="11"
          fill="#334155"
          font-family="monospace"
        >
          w₂: {s.w2.toFixed(3)}
        </text>
        <text
          x={tooltipX + 10}
          y={tooltipY + 68}
          font-size="11"
          fill="#334155"
          font-family="monospace"
        >
          L: {s.loss.toFixed(2)}
        </text>
      </g>
    {/if}
  </svg>
</div>

<style>
  .contour-container {
    position: relative;
    width: 100%;
    background: white;
    border: 1px solid #eee;
    border-radius: 4px;
  }

  canvas {
    image-rendering: pixelated; /* Sharp grid cells */
  }

  .overlay-svg {
    position: relative;
    width: 100%;
    display: block;
    pointer-events: none;
  }

  .trajectory-marker {
    cursor: pointer;
    pointer-events: all;
    transition: all 0.15s ease;
  }

  .trajectory-marker:hover {
    filter: brightness(1.2);
  }

  .tooltip {
    pointer-events: none;
  }
</style>
