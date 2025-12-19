<script lang="ts">
  import type { Snapshot } from './types';

  interface Props {
    initialSnapshot: Snapshot;
    currentSnapshot: Snapshot;
    finalSnapshot: Snapshot;
  }

  let {
    initialSnapshot,
    currentSnapshot,
    finalSnapshot
  }: Props = $props();

  // Mini scatter plot dimensions
  const width = 180;
  const height = 150;
  const padding = 30;

  // Calculate overall data range across all snapshots
  const dataRange = $derived.by(() => {
    const allPoints = [
      ...initialSnapshot.point_details,
      ...currentSnapshot.point_details,
      ...finalSnapshot.point_details
    ];

    if (allPoints.length === 0) {
      return { xMin: 0, xMax: 10, yMin: 0, yMax: 20 };
    }

    const xValues = allPoints.map(p => p.x);
    const yTrueValues = allPoints.map(p => p.y_true);
    const yPredValues = allPoints.map(p => p.y_pred);

    const xMin = Math.min(...xValues);
    const xMax = Math.max(...xValues);
    const yMin = Math.min(...yTrueValues, ...yPredValues);
    const yMax = Math.max(...yTrueValues, ...yPredValues);

    // Add padding
    const xPadding = (xMax - xMin) * 0.15;
    const yPadding = (yMax - yMin) * 0.15;

    return {
      xMin: xMin - xPadding,
      xMax: xMax + xPadding,
      yMin: yMin - yPadding,
      yMax: yMax + yPadding
    };
  });

  // Scaling functions (shared across all mini plots)
  function scaleX(x: number) {
    const range = dataRange;
    return padding + ((x - range.xMin) / (range.xMax - range.xMin)) * (width - 2 * padding);
  }

  function scaleY(y: number) {
    const range = dataRange;
    return height - padding - ((y - range.yMin) / (range.yMax - range.yMin)) * (height - 2 * padding);
  }

  // Helper to generate fitted line path
  function getFittedLine(w: number) {
    const range = dataRange;
    const x1 = range.xMin;
    const y1 = w * x1;
    const x2 = range.xMax;
    const y2 = w * x2;

    return {
      x1: scaleX(x1),
      y1: scaleY(y1),
      x2: scaleX(x2),
      y2: scaleY(y2)
    };
  }

  // Pre-compute fitted lines for all three snapshots
  const initialLine = $derived(getFittedLine(initialSnapshot.w));
  const currentLine = $derived(getFittedLine(currentSnapshot.w));
  const finalLine = $derived(getFittedLine(finalSnapshot.w));
</script>

<div class="fit-comparison-container">
  <h3 class="comparison-title">Training Progress</h3>

  <div class="comparison-grid">
    <!-- Initial Fit -->
    <div class="comparison-item">
      <div class="item-header">
        <span class="item-title">Initial</span>
        <span class="item-step">Step {initialSnapshot.step}</span>
      </div>

      <svg {width} {height} class="mini-plot">
        <!-- Axes -->
        <line
          x1={padding}
          y1={height - padding}
          x2={width - padding}
          y2={height - padding}
          stroke="#cbd5e1"
          stroke-width="1"
        />
        <line
          x1={padding}
          y1={padding}
          x2={padding}
          y2={height - padding}
          stroke="#cbd5e1"
          stroke-width="1"
        />

        <!-- Fitted line -->
        <line
          x1={initialLine.x1}
          y1={initialLine.y1}
          x2={initialLine.x2}
          y2={initialLine.y2}
          stroke="#94a3b8"
          stroke-width="2"
          opacity="0.7"
        />

        <!-- Data points -->
        {#each initialSnapshot.point_details as point}
          <circle
            cx={scaleX(point.x)}
            cy={scaleY(point.y_true)}
            r="3"
            fill="#3b82f6"
            opacity="0.8"
          />
          <circle
            cx={scaleX(point.x)}
            cy={scaleY(point.y_pred)}
            r="2"
            fill="#ef4444"
            opacity="0.6"
          />
        {/each}
      </svg>

      <div class="item-stats">
        <div class="stat">w = {initialSnapshot.w.toFixed(3)}</div>
        <div class="stat loss">Loss = {initialSnapshot.loss.toFixed(2)}</div>
      </div>
    </div>

    <!-- Current Fit -->
    <div class="comparison-item current">
      <div class="item-header">
        <span class="item-title">Current</span>
        <span class="item-step">Step {currentSnapshot.step}</span>
      </div>

      <svg {width} {height} class="mini-plot">
        <!-- Axes -->
        <line
          x1={padding}
          y1={height - padding}
          x2={width - padding}
          y2={height - padding}
          stroke="#cbd5e1"
          stroke-width="1"
        />
        <line
          x1={padding}
          y1={padding}
          x2={padding}
          y2={height - padding}
          stroke="#cbd5e1"
          stroke-width="1"
        />

        <!-- Fitted line -->
        <line
          x1={currentLine.x1}
          y1={currentLine.y1}
          x2={currentLine.x2}
          y2={currentLine.y2}
          stroke="#3b82f6"
          stroke-width="2.5"
          opacity="0.8"
        />

        <!-- Data points -->
        {#each currentSnapshot.point_details as point}
          <circle
            cx={scaleX(point.x)}
            cy={scaleY(point.y_true)}
            r="3.5"
            fill="#3b82f6"
            opacity="0.9"
          />
          <circle
            cx={scaleX(point.x)}
            cy={scaleY(point.y_pred)}
            r="2.5"
            fill="#ef4444"
            opacity="0.7"
          />
        {/each}
      </svg>

      <div class="item-stats">
        <div class="stat">w = {currentSnapshot.w.toFixed(3)}</div>
        <div class="stat loss">Loss = {currentSnapshot.loss.toFixed(2)}</div>
      </div>
    </div>

    <!-- Final Fit -->
    <div class="comparison-item">
      <div class="item-header">
        <span class="item-title">Final</span>
        <span class="item-step">Step {finalSnapshot.step}</span>
      </div>

      <svg {width} {height} class="mini-plot">
        <!-- Axes -->
        <line
          x1={padding}
          y1={height - padding}
          x2={width - padding}
          y2={height - padding}
          stroke="#cbd5e1"
          stroke-width="1"
        />
        <line
          x1={padding}
          y1={padding}
          x2={padding}
          y2={height - padding}
          stroke="#cbd5e1"
          stroke-width="1"
        />

        <!-- Fitted line -->
        <line
          x1={finalLine.x1}
          y1={finalLine.y1}
          x2={finalLine.x2}
          y2={finalLine.y2}
          stroke="#10b981"
          stroke-width="2.5"
          opacity="0.8"
        />

        <!-- Data points -->
        {#each finalSnapshot.point_details as point}
          <circle
            cx={scaleX(point.x)}
            cy={scaleY(point.y_true)}
            r="3.5"
            fill="#3b82f6"
            opacity="0.9"
          />
          <circle
            cx={scaleX(point.x)}
            cy={scaleY(point.y_pred)}
            r="2.5"
            fill="#ef4444"
            opacity="0.7"
          />
        {/each}
      </svg>

      <div class="item-stats">
        <div class="stat">w = {finalSnapshot.w.toFixed(3)}</div>
        <div class="stat loss">Loss = {finalSnapshot.loss.toFixed(2)}</div>
      </div>
    </div>
  </div>

  <div class="progress-indicator">
    <div class="progress-arrow">→</div>
    <p class="progress-text">
      Watch how the fitted line (colored) moves closer to the data points as training progresses
    </p>
  </div>
</div>

<style>
  .fit-comparison-container {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
    margin-top: 1rem;
  }

  .comparison-title {
    margin: 0 0 1rem 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: #0f172a;
  }

  .comparison-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
  }

  .comparison-item {
    background: #f8fafc;
    padding: 1rem;
    border-radius: 8px;
    border: 2px solid transparent;
    transition: all 0.2s ease;
  }

  .comparison-item.current {
    border-color: #3b82f6;
    background: #eff6ff;
  }

  .item-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .item-title {
    font-weight: 600;
    font-size: 0.9rem;
    color: #0f172a;
  }

  .item-step {
    font-size: 0.75rem;
    color: #64748b;
    font-family: 'SF Mono', monospace;
  }

  .mini-plot {
    display: block;
    margin: 0 auto;
  }

  .item-stats {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    margin-top: 0.5rem;
  }

  .stat {
    font-size: 0.85rem;
    font-family: 'SF Mono', monospace;
    color: #334155;
  }

  .stat.loss {
    color: #ef4444;
    font-weight: 600;
  }

  .progress-indicator {
    margin-top: 1rem;
    text-align: center;
    padding-top: 1rem;
    border-top: 1px solid #e2e8f0;
  }

  .progress-arrow {
    font-size: 2rem;
    color: #10b981;
    margin-bottom: 0.5rem;
  }

  .progress-text {
    margin: 0;
    font-size: 0.85rem;
    color: #64748b;
    line-height: 1.5;
  }

  @media (max-width: 768px) {
    .comparison-grid {
      grid-template-columns: 1fr;
    }

    .comparison-item {
      max-width: 280px;
      margin: 0 auto;
    }
  }
</style>
