<script lang="ts">
  import type { LinearSnapshot } from './types';

  interface Props {
    snapshot: LinearSnapshot;
    selectedPointIndex: number | null;
    onSelectPoint: (index: number | null) => void;
  }

  let { snapshot, selectedPointIndex, onSelectPoint }: Props = $props();

  // SVG dimensions
  const width = 700;
  const height = 500;
  const margin = { top: 40, right: 40, bottom: 60, left: 60 };
  const plotWidth = width - margin.left - margin.right;
  const plotHeight = height - margin.top - margin.bottom;

  // Hover state
  let hoveredIndex = $state<number | null>(null);

  // Compute scales for x1 and x2
  let x1Range = $derived.by(() => {
    const values = snapshot.point_details.map(p => p.x1);
    const min = Math.min(...values);
    const max = Math.max(...values);
    const padding = (max - min) * 0.1 || 1;
    return { min: min - padding, max: max + padding };
  });

  let x2Range = $derived.by(() => {
    const values = snapshot.point_details.map(p => p.x2);
    const min = Math.min(...values);
    const max = Math.max(...values);
    const padding = (max - min) * 0.1 || 1;
    return { min: min - padding, max: max + padding };
  });

  // Loss range for color scale
  let lossRange = $derived.by(() => {
    const losses = snapshot.point_details.map(p => p.point_loss);
    return { min: Math.min(...losses), max: Math.max(...losses) };
  });

  // Gradient magnitude range for size scale
  let gradMagRange = $derived.by(() => {
    const mags = snapshot.point_details.map(p => Math.sqrt(p.grad_w1 ** 2 + p.grad_w2 ** 2));
    return { min: Math.min(...mags), max: Math.max(...mags) };
  });

  // Scale functions
  function scaleX1(x1: number): number {
    const range = x1Range;
    return margin.left + ((x1 - range.min) / (range.max - range.min)) * plotWidth;
  }

  function scaleX2(x2: number): number {
    const range = x2Range;
    return margin.top + plotHeight - ((x2 - range.min) / (range.max - range.min)) * plotHeight;
  }

  // Color scale: blue (low loss) to red (high loss)
  function getLossColor(loss: number): string {
    const range = lossRange;
    const normalized = (loss - range.min) / (range.max - range.min || 1);

    if (normalized < 0.33) {
      // Blue to cyan
      const t = normalized / 0.33;
      return `rgb(${Math.round(59 + (14 - 59) * t)}, ${Math.round(130 + (165 - 130) * t)}, ${Math.round(246 + (255 - 246) * t)})`;
    } else if (normalized < 0.67) {
      // Cyan to yellow
      const t = (normalized - 0.33) / 0.34;
      return `rgb(${Math.round(14 + (253 - 14) * t)}, ${Math.round(165 + (224 - 165) * t)}, ${Math.round(255 + (71 - 255) * t)})`;
    } else {
      // Yellow to red
      const t = (normalized - 0.67) / 0.33;
      return `rgb(${Math.round(253 + (220 - 253) * t)}, ${Math.round(224 + (38 - 224) * t)}, ${Math.round(71 + (38 - 71) * t)})`;
    }
  }

  // Size scale: larger for higher gradient magnitude
  function getPointRadius(gradMag: number): number {
    const range = gradMagRange;
    const normalized = (gradMag - range.min) / (range.max - range.min || 1);
    return 5 + normalized * 5; // 5-10px radius
  }

  // Format number for display
  function formatNumber(n: number): string {
    return Math.abs(n) < 100 ? n.toFixed(3) : n.toFixed(1);
  }

  // Handle point click
  function handlePointClick(index: number) {
    if (selectedPointIndex === index) {
      onSelectPoint(null); // Deselect if clicking same point
    } else {
      onSelectPoint(index);
    }
  }

  // Grid lines for x1 axis
  let x1Ticks = $derived.by(() => {
    const range = x1Range;
    const span = range.max - range.min;
    const step = Math.pow(10, Math.floor(Math.log10(span / 5)));
    const start = Math.ceil(range.min / step) * step;
    const ticks = [];
    for (let val = start; val <= range.max; val += step) {
      ticks.push(val);
    }
    return ticks;
  });

  // Grid lines for x2 axis
  let x2Ticks = $derived.by(() => {
    const range = x2Range;
    const span = range.max - range.min;
    const step = Math.pow(10, Math.floor(Math.log10(span / 5)));
    const start = Math.ceil(range.min / step) * step;
    const ticks = [];
    for (let val = start; val <= range.max; val += step) {
      ticks.push(val);
    }
    return ticks;
  });
</script>

<div class="scatter-container">
  <h3>Training Data Points (x₁ vs x₂)</h3>
  <p class="description">Click a point to view detailed analysis. Click again to deselect.</p>

  <div class="scatter-layout">
    <div class="scatter-plot-section">
      <svg {width} {height}>
    <!-- Grid lines -->
    <g class="grid">
      {#each x1Ticks as tick}
        <line
          x1={scaleX1(tick)}
          y1={margin.top}
          x2={scaleX1(tick)}
          y2={margin.top + plotHeight}
          stroke="#e2e8f0"
          stroke-width="1"
        />
      {/each}
      {#each x2Ticks as tick}
        <line
          x1={margin.left}
          y1={scaleX2(tick)}
          x2={margin.left + plotWidth}
          y2={scaleX2(tick)}
          stroke="#e2e8f0"
          stroke-width="1"
        />
      {/each}
    </g>

    <!-- Axes -->
    <line
      x1={margin.left}
      y1={margin.top + plotHeight}
      x2={margin.left + plotWidth}
      y2={margin.top + plotHeight}
      stroke="#64748b"
      stroke-width="2"
    />
    <line
      x1={margin.left}
      y1={margin.top}
      x2={margin.left}
      y2={margin.top + plotHeight}
      stroke="#64748b"
      stroke-width="2"
    />

    <!-- Axis ticks and labels -->
    {#each x1Ticks as tick}
      <line
        x1={scaleX1(tick)}
        y1={margin.top + plotHeight}
        x2={scaleX1(tick)}
        y2={margin.top + plotHeight + 6}
        stroke="#64748b"
        stroke-width="2"
      />
      <text
        x={scaleX1(tick)}
        y={margin.top + plotHeight + 20}
        text-anchor="middle"
        class="axis-label"
      >
        {formatNumber(tick)}
      </text>
    {/each}

    {#each x2Ticks as tick}
      <line
        x1={margin.left - 6}
        y1={scaleX2(tick)}
        x2={margin.left}
        y2={scaleX2(tick)}
        stroke="#64748b"
        stroke-width="2"
      />
      <text
        x={margin.left - 10}
        y={scaleX2(tick)}
        text-anchor="end"
        dominant-baseline="middle"
        class="axis-label"
      >
        {formatNumber(tick)}
      </text>
    {/each}

    <!-- Axis titles -->
    <text
      x={margin.left + plotWidth / 2}
      y={height - 10}
      text-anchor="middle"
      class="axis-title"
    >
      x₁ (Feature 1)
    </text>
    <text
      x={15}
      y={margin.top + plotHeight / 2}
      text-anchor="middle"
      dominant-baseline="middle"
      transform="rotate(-90, 15, {margin.top + plotHeight / 2})"
      class="axis-title"
    >
      x₂ (Feature 2)
    </text>

    <!-- Data points -->
    {#each snapshot.point_details as point, i}
      {@const gradMag = Math.sqrt(point.grad_w1 ** 2 + point.grad_w2 ** 2)}
      {@const isSelected = i === selectedPointIndex}
      {@const isHovered = i === hoveredIndex}

      <circle
        cx={scaleX1(point.x1)}
        cy={scaleX2(point.x2)}
        r={getPointRadius(gradMag)}
        fill={getLossColor(point.point_loss)}
        stroke={isSelected ? '#2563eb' : isHovered ? '#64748b' : 'white'}
        stroke-width={isSelected ? 3 : isHovered ? 2 : 1}
        opacity={isSelected || isHovered ? 1 : 0.8}
        class="data-point"
        onclick={() => handlePointClick(i)}
        onmouseenter={() => (hoveredIndex = i)}
        onmouseleave={() => (hoveredIndex = null)}
      />
    {/each}
      </svg>

      <!-- Legend -->
      <div class="legend">
        <div class="legend-item">
          <div class="legend-color-gradient"></div>
          <span class="legend-label">Loss: Low → High</span>
        </div>
        <div class="legend-item">
          <div class="legend-size-demo">
            <div class="size-dot small"></div>
            <div class="size-dot large"></div>
          </div>
          <span class="legend-label">Size: Gradient Magnitude</span>
        </div>
      </div>
    </div>

    <!-- Details Panel with Click-based Selection -->
    {#if selectedPointIndex !== null}
      {@const point = snapshot.point_details[selectedPointIndex]}
      {@const error = point.y_pred - point.y_true}
      {@const errorPercent = point.y_true !== 0 ? (error / point.y_true) * 100 : 0}
      {@const maxVal = Math.max(Math.abs(point.y_pred), Math.abs(point.y_true), 1)}
      {@const scale = 100 / maxVal}
      {@const predBar = Math.abs(point.y_pred) * scale}
      {@const trueBar = Math.abs(point.y_true) * scale}
      {@const lossRangeValues = snapshot.point_details.map(p => p.point_loss)}
      {@const maxLoss = Math.max(...lossRangeValues)}
      {@const lossPercent = (point.point_loss / maxLoss) * 100}
      {@const maxGrad = Math.max(Math.abs(point.grad_w1), Math.abs(point.grad_w2), 0.1)}
      {@const gradScale = 80 / maxGrad}
      {@const w1Length = Math.abs(point.grad_w1) * gradScale}
      {@const w1Start = point.grad_w1 > 0 ? 150 : 150 - w1Length}
      {@const w2Length = Math.abs(point.grad_w2) * gradScale}
      {@const w2Start = point.grad_w2 > 0 ? 150 : 150 - w2Length}
      {@const gradMag = Math.sqrt(point.grad_w1 ** 2 + point.grad_w2 ** 2)}

      <div class="point-details-panel">
        <!-- Header with close button -->
        <div class="panel-header">
          <h4>Point {selectedPointIndex + 1} Analysis</h4>
          <button
            class="close-button"
            onclick={() => onSelectPoint(null)}
            aria-label="Close details"
            type="button"
          >
            ×
          </button>
        </div>

        <!-- Section 1: Input Features -->
        <div class="details-section">
          <h5>Input Features</h5>
          <div class="feature-row">
            <span class="feature-label">x₁:</span>
            <span class="feature-value">{formatNumber(point.x1)}</span>
          </div>
          <div class="feature-row">
            <span class="feature-label">x₂:</span>
            <span class="feature-value">{formatNumber(point.x2)}</span>
          </div>
        </div>

        <!-- Section 2: Prediction vs Truth Visualization -->
        <div class="details-section">
          <h5>Prediction vs True Value</h5>
          <svg viewBox="0 0 300 120" class="comparison-viz">
            <!-- Y_pred bar -->
            <g transform="translate(50, 20)">
              <rect x="0" y="0" width={predBar} height="20" fill="#3b82f6" rx="2"/>
              <text x={predBar + 5} y="15" font-size="12" fill="#334155">
                ŷ = {formatNumber(point.y_pred)}
              </text>
            </g>

            <!-- Y_true bar -->
            <g transform="translate(50, 50)">
              <rect x="0" y="0" width={trueBar} height="20" fill="#10b981" rx="2"/>
              <text x={trueBar + 5} y="15" font-size="12" fill="#334155">
                y = {formatNumber(point.y_true)}
              </text>
            </g>

            <!-- Error indicator -->
            <g transform="translate(50, 85)">
              <text x="0" y="0" font-size="11" fill="#64748b">
                Error: {formatNumber(error)} ({errorPercent.toFixed(1)}%)
              </text>
            </g>
          </svg>
        </div>

        <!-- Section 3: Loss Visualization -->
        <div class="details-section">
          <h5>Point Loss</h5>
          <div class="loss-display">
            <div class="loss-bar-container">
              <div
                class="loss-bar"
                style="width: {lossPercent}%; background: {getLossColor(point.point_loss)}"
              ></div>
            </div>
            <span class="loss-value">{formatNumber(point.point_loss)}</span>
          </div>
        </div>

        <!-- Section 4: Gradient Components -->
        <div class="details-section">
          <h5>Gradient Components</h5>
          <svg viewBox="0 0 300 100" class="gradient-viz">
            <!-- Center line -->
            <line x1="150" y1="20" x2="150" y2="80" stroke="#e2e8f0" stroke-width="2"/>

            <!-- Grad w1 bar -->
            <g transform="translate(0, 30)">
              <rect
                x={w1Start}
                y="0"
                width={w1Length}
                height="15"
                fill="#ef4444"
                opacity="0.7"
                rx="2"
              />
              <text x="10" y="12" font-size="11" fill="#334155">∂L/∂w₁:</text>
              <text x="220" y="12" font-size="11" fill="#334155">{formatNumber(point.grad_w1)}</text>
            </g>

            <!-- Grad w2 bar -->
            <g transform="translate(0, 55)">
              <rect
                x={w2Start}
                y="0"
                width={w2Length}
                height="15"
                fill="#8b5cf6"
                opacity="0.7"
                rx="2"
              />
              <text x="10" y="12" font-size="11" fill="#334155">∂L/∂w₂:</text>
              <text x="220" y="12" font-size="11" fill="#334155">{formatNumber(point.grad_w2)}</text>
            </g>
          </svg>

          <div class="gradient-magnitude">
            <span class="magnitude-label">Magnitude:</span>
            <span class="magnitude-value">{formatNumber(gradMag)}</span>
          </div>
        </div>
      </div>
    {:else}
      <div class="point-details-panel placeholder">
        <div class="placeholder-content">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <circle cx="12" cy="12" r="10" stroke-width="2"/>
            <path d="M12 16v-4m0-4h.01" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <p class="details-hint">Click a point to see detailed analysis</p>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .scatter-container {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
  }

  h3 {
    margin: 0 0 0.5rem 0;
    font-size: 1.25rem;
    color: #0f172a;
  }

  .description {
    margin: 0 0 1rem 0;
    font-size: 0.9rem;
    color: #64748b;
  }

  svg {
    display: block;
    margin: 0 auto;
    background: white;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
  }

  .axis-label {
    font-size: 12px;
    fill: #64748b;
    font-family: 'SF Mono', 'Monaco', 'Courier New', monospace;
  }

  .axis-title {
    font-size: 14px;
    fill: #334155;
    font-weight: 600;
  }

  .data-point {
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .data-point:hover {
    filter: brightness(1.1);
  }

  /* Main layout container - uses flexbox for smart positioning */
  .scatter-layout {
    display: flex;
    gap: 1.5rem;
    align-items: flex-start;
    flex-wrap: wrap;
  }

  .scatter-plot-section {
    flex: 1 1 700px;
    min-width: 0;
  }

  /* Details panel - positioned right or wraps below */
  .point-details-panel {
    flex: 0 0 350px;
    min-width: 300px;
    max-width: 400px;
    background: white;
    border: 1px solid #cbd5e1;
    border-radius: 8px;
    padding: 0;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
    max-height: 600px;
    overflow-y: auto;
  }

  .point-details-panel.placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 300px;
    border-style: dashed;
    background: #f8fafc;
  }

  .placeholder-content {
    text-align: center;
    color: #94a3b8;
  }

  .placeholder-content svg {
    margin: 0 auto 1rem;
    color: #cbd5e1;
  }

  .details-hint {
    font-size: 0.9rem;
    color: #64748b;
    margin: 0;
  }

  /* Panel header with close button */
  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.25rem;
    border-bottom: 2px solid #e2e8f0;
    background: #f8fafc;
    border-radius: 8px 8px 0 0;
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .panel-header h4 {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: #0f172a;
  }

  .close-button {
    background: transparent;
    border: none;
    font-size: 1.5rem;
    line-height: 1;
    cursor: pointer;
    color: #64748b;
    padding: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    transition: all 0.15s ease;
  }

  .close-button:hover {
    background: #e2e8f0;
    color: #0f172a;
  }

  /* Details sections */
  .details-section {
    padding: 1rem 1.25rem;
    border-bottom: 1px solid #f1f5f9;
  }

  .details-section:last-child {
    border-bottom: none;
  }

  .details-section h5 {
    margin: 0 0 0.75rem 0;
    font-size: 0.85rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: #64748b;
  }

  /* Feature rows */
  .feature-row {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem 0;
    font-family: 'SF Mono', 'Monaco', 'Courier New', monospace;
    font-size: 0.875rem;
  }

  .feature-label {
    font-weight: 600;
    color: #64748b;
  }

  .feature-value {
    color: #0f172a;
    font-weight: 500;
  }

  /* Visualizations */
  .comparison-viz,
  .gradient-viz {
    width: 100%;
    height: auto;
    margin: 0.5rem 0;
  }

  /* Loss display */
  .loss-display {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-top: 0.5rem;
  }

  .loss-bar-container {
    flex: 1;
    height: 24px;
    background: #f1f5f9;
    border-radius: 4px;
    overflow: hidden;
    border: 1px solid #e2e8f0;
  }

  .loss-bar {
    height: 100%;
    transition: width 0.3s ease;
    border-radius: 3px;
  }

  .loss-value {
    font-family: 'SF Mono', 'Monaco', 'Courier New', monospace;
    font-size: 0.875rem;
    font-weight: 600;
    color: #0f172a;
    min-width: 60px;
    text-align: right;
  }

  /* Gradient magnitude display */
  .gradient-magnitude {
    display: flex;
    justify-content: space-between;
    margin-top: 0.75rem;
    padding: 0.5rem;
    background: #f8fafc;
    border-radius: 4px;
    font-family: 'SF Mono', 'Monaco', 'Courier New', monospace;
    font-size: 0.875rem;
  }

  .magnitude-label {
    font-weight: 600;
    color: #64748b;
  }

  .magnitude-value {
    color: #0f172a;
    font-weight: 500;
  }

  /* Responsive behavior */
  @media (max-width: 1200px) {
    .scatter-layout {
      flex-direction: column;
    }

    .point-details-panel {
      flex: 1 1 auto;
      max-width: 100%;
      width: 100%;
    }
  }

  @media (max-width: 768px) {
    .point-details-panel {
      min-width: 0;
      width: 100%;
    }

    .details-section {
      padding: 0.75rem 1rem;
    }
  }

  .legend {
    display: flex;
    justify-content: center;
    gap: 2rem;
    margin-top: 1rem;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .legend-color-gradient {
    width: 100px;
    height: 16px;
    background: linear-gradient(to right, #3b82f6, #0ea5e9, #fde047, #dc2626);
    border: 1px solid #cbd5e1;
    border-radius: 4px;
  }

  .legend-size-demo {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .size-dot {
    border-radius: 50%;
    background: #64748b;
    border: 1px solid white;
  }

  .size-dot.small {
    width: 10px;
    height: 10px;
  }

  .size-dot.large {
    width: 20px;
    height: 20px;
  }

  .legend-label {
    font-size: 0.85rem;
    color: #64748b;
  }
</style>
