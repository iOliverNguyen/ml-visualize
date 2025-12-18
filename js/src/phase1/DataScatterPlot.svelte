<script lang="ts">
  import type { PointSnapshot } from './types';
  import { tweened } from 'svelte/motion';
  import { cubicOut } from 'svelte/easing';

  interface Props {
    pointDetails: PointSnapshot[];
    currentW: number;
    showResiduals?: boolean;
    highlightedPointIndex?: number;
  }

  let {
    pointDetails,
    currentW,
    showResiduals = true,
    highlightedPointIndex
  }: Props = $props();

  // SVG dimensions
  const width = 500;
  const height = 400;
  const padding = 60;

  // Hover state
  let hoveredIndex = $state<number | null>(null);

  // Animated w for smooth line transitions
  const animatedW = tweened(0, {
    duration: 200,
    easing: cubicOut
  });

  // Update animated w when currentW changes
  $effect(() => {
    animatedW.set(currentW);
  });

  // Calculate data ranges
  const dataRanges = $derived(() => {
    if (pointDetails.length === 0) {
      return { xMin: 0, xMax: 10, yMin: 0, yMax: 20 };
    }

    const xValues = pointDetails.map(p => p.x);
    const yTrueValues = pointDetails.map(p => p.y_true);
    const yPredValues = pointDetails.map(p => p.y_pred);

    const xMin = Math.min(...xValues);
    const xMax = Math.max(...xValues);
    const yMin = Math.min(...yTrueValues, ...yPredValues);
    const yMax = Math.max(...yTrueValues, ...yPredValues);

    // Add 10% padding to ranges
    const xPadding = (xMax - xMin) * 0.1;
    const yPadding = (yMax - yMin) * 0.1;

    return {
      xMin: xMin - xPadding,
      xMax: xMax + xPadding,
      yMin: yMin - yPadding,
      yMax: yMax + yPadding
    };
  });

  // Scaling functions
  const scaleX = $derived((x: number) => {
    const range = dataRanges();
    return padding + ((x - range.xMin) / (range.xMax - range.xMin)) * (width - 2 * padding);
  });

  const scaleY = $derived((y: number) => {
    const range = dataRanges();
    return height - padding - ((y - range.yMin) / (range.yMax - range.yMin)) * (height - 2 * padding);
  });

  // Calculate fitted line points
  const fittedLinePoints = $derived(() => {
    const range = dataRanges();
    const x1 = range.xMin;
    const y1 = $animatedW * x1;
    const x2 = range.xMax;
    const y2 = $animatedW * x2;

    return {
      x1: scaleX(x1),
      y1: scaleY(y1),
      x2: scaleX(x2),
      y2: scaleY(y2)
    };
  });

  // X-axis ticks
  const xTicks = $derived(() => {
    const range = dataRanges();
    const numTicks = 6;
    const step = (range.xMax - range.xMin) / (numTicks - 1);
    return Array.from({ length: numTicks }, (_, i) => range.xMin + i * step);
  });

  // Y-axis ticks
  const yTicks = $derived(() => {
    const range = dataRanges();
    const numTicks = 6;
    const step = (range.yMax - range.yMin) / (numTicks - 1);
    return Array.from({ length: numTicks }, (_, i) => range.yMin + i * step);
  });
</script>

<div class="scatter-plot-container">
  <h3 class="plot-title">Data Fit Visualization</h3>

  <svg {width} {height} class="scatter-plot">
    <!-- Axes -->
    <line
      x1={padding}
      y1={height - padding}
      x2={width - padding}
      y2={height - padding}
      stroke="#94a3b8"
      stroke-width="2"
    />
    <line
      x1={padding}
      y1={padding}
      x2={padding}
      y2={height - padding}
      stroke="#94a3b8"
      stroke-width="2"
    />

    <!-- X-axis ticks and labels -->
    {#each xTicks() as tick}
      <line
        x1={scaleX(tick)}
        y1={height - padding}
        x2={scaleX(tick)}
        y2={height - padding + 5}
        stroke="#94a3b8"
        stroke-width="1"
      />
      <text
        x={scaleX(tick)}
        y={height - padding + 20}
        text-anchor="middle"
        font-size="12"
        fill="#64748b"
      >
        {tick.toFixed(1)}
      </text>
    {/each}

    <!-- Y-axis ticks and labels -->
    {#each yTicks() as tick}
      <line
        x1={padding - 5}
        y1={scaleY(tick)}
        x2={padding}
        y2={scaleY(tick)}
        stroke="#94a3b8"
        stroke-width="1"
      />
      <text
        x={padding - 10}
        y={scaleY(tick)}
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
      Input (x)
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
      Output (y)
    </text>

    <!-- Residual lines (errors) -->
    {#if showResiduals}
      {#each pointDetails as point, i}
        <line
          x1={scaleX(point.x)}
          y1={scaleY(point.y_true)}
          x2={scaleX(point.x)}
          y2={scaleY(point.y_pred)}
          stroke="#cbd5e1"
          stroke-width={hoveredIndex === i ? 3 : 1.5}
          stroke-dasharray="4,2"
          opacity={hoveredIndex === i ? 1 : 0.6}
        />
      {/each}
    {/if}

    <!-- Fitted line -->
    <line
      x1={fittedLinePoints().x1}
      y1={fittedLinePoints().y1}
      x2={fittedLinePoints().x2}
      y2={fittedLinePoints().y2}
      stroke="#10b981"
      stroke-width="3"
      opacity="0.8"
    />

    <!-- Actual data points (y_true) -->
    {#each pointDetails as point, i}
      <circle
        cx={scaleX(point.x)}
        cy={scaleY(point.y_true)}
        r={hoveredIndex === i || highlightedPointIndex === i ? 7 : 5}
        fill="#3b82f6"
        stroke="white"
        stroke-width="2"
        opacity={hoveredIndex === i || highlightedPointIndex === i ? 1 : 0.8}
        style="cursor: pointer;"
        onmouseenter={() => hoveredIndex = i}
        onmouseleave={() => hoveredIndex = null}
      />
    {/each}

    <!-- Predicted data points (y_pred) -->
    {#each pointDetails as point, i}
      <circle
        cx={scaleX(point.x)}
        cy={scaleY(point.y_pred)}
        r={hoveredIndex === i || highlightedPointIndex === i ? 6 : 4}
        fill="#ef4444"
        stroke="white"
        stroke-width="2"
        opacity={hoveredIndex === i || highlightedPointIndex === i ? 1 : 0.7}
        style="cursor: pointer;"
        onmouseenter={() => hoveredIndex = i}
        onmouseleave={() => hoveredIndex = null}
      />
    {/each}

    <!-- Equation label -->
    <text
      x={width - padding - 10}
      y={padding + 20}
      text-anchor="end"
      font-size="16"
      font-weight="600"
      fill="#10b981"
    >
      y = {currentW.toFixed(3)}*x
    </text>
  </svg>

  <!-- Legend -->
  <div class="legend">
    <div class="legend-item">
      <div class="legend-marker" style="background: #3b82f6;"></div>
      <span>Actual data (y_true)</span>
    </div>
    <div class="legend-item">
      <div class="legend-marker" style="background: #ef4444;"></div>
      <span>Predictions (y_pred)</span>
    </div>
    <div class="legend-item">
      <div class="legend-line" style="background: #10b981;"></div>
      <span>Fitted line (y = w*x)</span>
    </div>
    {#if showResiduals}
      <div class="legend-item">
        <div class="legend-residual"></div>
        <span>Residual errors</span>
      </div>
    {/if}
  </div>

  <!-- Hover tooltip -->
  {#if hoveredIndex !== null && pointDetails[hoveredIndex]}
    {@const point = pointDetails[hoveredIndex]}
    <div class="tooltip">
      <div class="tooltip-row">
        <span class="tooltip-label">x:</span>
        <span class="tooltip-value">{point.x.toFixed(2)}</span>
      </div>
      <div class="tooltip-row">
        <span class="tooltip-label">y_true:</span>
        <span class="tooltip-value">{point.y_true.toFixed(2)}</span>
      </div>
      <div class="tooltip-row">
        <span class="tooltip-label">y_pred:</span>
        <span class="tooltip-value">{point.y_pred.toFixed(2)}</span>
      </div>
      <div class="tooltip-row">
        <span class="tooltip-label">Error:</span>
        <span class="tooltip-value">{Math.abs(point.y_true - point.y_pred).toFixed(3)}</span>
      </div>
    </div>
  {/if}
</div>

<style>
  .scatter-plot-container {
    position: relative;
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  }

  .plot-title {
    margin: 0 0 1rem 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: #0f172a;
  }

  .scatter-plot {
    display: block;
    margin: 0 auto;
  }

  .legend {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem;
    justify-content: center;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid #e2e8f0;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.85rem;
    color: #64748b;
  }

  .legend-marker {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    border: 2px solid white;
  }

  .legend-line {
    width: 24px;
    height: 3px;
    border-radius: 2px;
  }

  .legend-residual {
    width: 24px;
    height: 2px;
    background: repeating-linear-gradient(
      to right,
      #cbd5e1 0px,
      #cbd5e1 4px,
      transparent 4px,
      transparent 6px
    );
  }

  .tooltip {
    position: absolute;
    top: 1rem;
    right: 1rem;
    background: rgba(15, 23, 42, 0.95);
    color: white;
    padding: 0.75rem;
    border-radius: 6px;
    font-size: 0.85rem;
    pointer-events: none;
    z-index: 10;
  }

  .tooltip-row {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    margin: 0.25rem 0;
  }

  .tooltip-label {
    font-weight: 600;
    color: #94a3b8;
  }

  .tooltip-value {
    font-family: 'SF Mono', monospace;
    color: white;
  }

  @media (max-width: 640px) {
    .scatter-plot-container {
      padding: 1rem;
    }

    .scatter-plot {
      width: 100%;
      height: auto;
    }

    .legend {
      flex-direction: column;
      align-items: flex-start;
      gap: 0.5rem;
    }
  }
</style>
