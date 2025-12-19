<script lang="ts">
  import type { NeuronSnapshot } from './types';

  interface Props {
    snapshot: NeuronSnapshot;
  }

  let { snapshot }: Props = $props();

  // SVG dimensions
  const width = 700;
  const height = 300;

  // Node positions
  const inputNodes = [
    { id: 'x1', x: 100, y: 100, label: 'x₁' },
    { id: 'x2', x: 100, y: 200, label: 'x₂' }
  ];
  const biasNode = { id: 'b', x: 100, y: 150, label: 'bias' };
  const sumNode = { id: 'sum', x: 300, y: 150, label: 'Σ' };
  const activationNode = { id: 'activation', x: 450, y: 150, label: snapshot?.activation ?? 'sigmoid' };
  const outputNode = { id: 'output', x: 600, y: 150, label: 'ŷ' };

  // Get average input values from first point for visualization
  let avgX1 = $derived.by(() => {
    if (snapshot?.point_details?.length > 0) {
      const sum = snapshot.point_details.reduce((acc, p) => acc + p.x[0], 0);
      return sum / snapshot.point_details.length;
    }
    return 0;
  });

  let avgX2 = $derived.by(() => {
    if (snapshot?.point_details?.length > 0) {
      const sum = snapshot.point_details.reduce((acc, p) => acc + p.x[1], 0);
      return sum / snapshot.point_details.length;
    }
    return 0;
  });

  // Edge styling functions
  function getEdgeColor(weight: number): string {
    return weight >= 0 ? '#3b82f6' : '#dc2626'; // blue for positive, red for negative
  }

  function getEdgeWidth(weight: number): number {
    const absWeight = Math.abs(weight);
    return Math.max(2, Math.min(8, absWeight * 3));
  }

  function getNodeOpacity(value: number, maxValue: number): number {
    return 0.3 + (0.7 * Math.abs(value) / Math.max(maxValue, 1));
  }

  function formatNumber(n: number): string {
    return n.toFixed(3);
  }

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
</script>

<div class="neuron-diagram">
  <h2>Neuron Architecture</h2>
  <p class="subtitle">Forward pass: z = w₁·x₁ + w₂·x₂ + b → a = {snapshot?.activation ?? 'sigmoid'}(z) → ŷ</p>

  <svg {width} {height} class="diagram">
    <!-- Edges from inputs to sum -->
    <!-- x1 to sum -->
    <line
      x1={inputNodes[0].x + 30}
      y1={inputNodes[0].y}
      x2={sumNode.x - 30}
      y2={sumNode.y}
      stroke={getEdgeColor(snapshot?.params?.w?.[0] ?? 0)}
      stroke-width={getEdgeWidth(snapshot?.params?.w?.[0] ?? 0)}
      opacity="0.7"
    />
    <text
      x={(inputNodes[0].x + sumNode.x) / 2}
      y={(inputNodes[0].y + sumNode.y) / 2 - 10}
      class="edge-label"
      fill={getEdgeColor(snapshot?.params?.w?.[0] ?? 0)}
    >
      w₁={formatNumber(snapshot?.params?.w?.[0] ?? 0)}
    </text>

    <!-- x2 to sum -->
    <line
      x1={inputNodes[1].x + 30}
      y1={inputNodes[1].y}
      x2={sumNode.x - 30}
      y2={sumNode.y}
      stroke={getEdgeColor(snapshot?.params?.w?.[1] ?? 0)}
      stroke-width={getEdgeWidth(snapshot?.params?.w?.[1] ?? 0)}
      opacity="0.7"
    />
    <text
      x={(inputNodes[1].x + sumNode.x) / 2}
      y={(inputNodes[1].y + sumNode.y) / 2 + 20}
      class="edge-label"
      fill={getEdgeColor(snapshot?.params?.w?.[1] ?? 0)}
    >
      w₂={formatNumber(snapshot?.params?.w?.[1] ?? 0)}
    </text>

    <!-- bias to sum -->
    <line
      x1={biasNode.x + 25}
      y1={biasNode.y}
      x2={sumNode.x - 30}
      y2={sumNode.y}
      stroke={getEdgeColor(snapshot?.params?.b ?? 0)}
      stroke-width={getEdgeWidth(snapshot?.params?.b ?? 0)}
      opacity="0.5"
      stroke-dasharray="4,4"
    />
    <text
      x={(biasNode.x + sumNode.x) / 2}
      y={biasNode.y - 10}
      class="edge-label small"
      fill={getEdgeColor(snapshot?.params?.b ?? 0)}
    >
      b={formatNumber(snapshot?.params?.b ?? 0)}
    </text>

    <!-- Edge from sum to activation -->
    <line
      x1={sumNode.x + 30}
      y1={sumNode.y}
      x2={activationNode.x - 50}
      y2={activationNode.y}
      stroke="#64748b"
      stroke-width="3"
      opacity="0.7"
    />
    <text
      x={(sumNode.x + activationNode.x) / 2}
      y={sumNode.y - 15}
      class="edge-label"
      fill="#64748b"
    >
      z={formatNumber(snapshot?.z ?? 0)}
    </text>

    <!-- Edge from activation to output -->
    <line
      x1={activationNode.x + 50}
      y1={activationNode.y}
      x2={outputNode.x - 30}
      y2={outputNode.y}
      stroke={getActivationColor(snapshot?.activation ?? 'sigmoid')}
      stroke-width="3"
      opacity="0.7"
    />
    <text
      x={(activationNode.x + outputNode.x) / 2}
      y={activationNode.y - 15}
      class="edge-label"
      fill={getActivationColor(snapshot?.activation ?? 'sigmoid')}
    >
      a={formatNumber(snapshot?.a ?? 0)}
    </text>

    <!-- Input nodes -->
    {#each inputNodes as node, i}
      <g>
        <circle
          cx={node.x}
          cy={node.y}
          r="25"
          fill="#e0f2fe"
          stroke="#0284c7"
          stroke-width="2"
          opacity={getNodeOpacity(i === 0 ? avgX1 : avgX2, 2)}
        />
        <text x={node.x} y={node.y - 5} class="node-label" text-anchor="middle">{node.label}</text>
        <text x={node.x} y={node.y + 10} class="node-value" text-anchor="middle">
          {formatNumber(i === 0 ? avgX1 : avgX2)}
        </text>
      </g>
    {/each}

    <!-- Bias node -->
    <g>
      <rect
        x={biasNode.x - 20}
        y={biasNode.y - 15}
        width="40"
        height="30"
        rx="4"
        fill="#fef3c7"
        stroke="#f59e0b"
        stroke-width="2"
        opacity="0.8"
      />
      <text x={biasNode.x} y={biasNode.y + 5} class="node-label small" text-anchor="middle">
        {biasNode.label}
      </text>
    </g>

    <!-- Sum node -->
    <g>
      <rect
        x={sumNode.x - 30}
        y={sumNode.y - 30}
        width="60"
        height="60"
        rx="8"
        fill="#f1f5f9"
        stroke="#64748b"
        stroke-width="3"
      />
      <text x={sumNode.x} y={sumNode.y + 5} class="node-label large" text-anchor="middle">
        {sumNode.label}
      </text>
    </g>

    <!-- Activation node -->
    <g>
      <rect
        x={activationNode.x - 50}
        y={activationNode.y - 30}
        width="100"
        height="60"
        rx="12"
        fill={getActivationColor(snapshot?.activation ?? 'sigmoid')}
        fill-opacity={getNodeOpacity(snapshot?.a ?? 0, 1)}
        stroke={getActivationColor(snapshot?.activation ?? 'sigmoid')}
        stroke-width="3"
      />
      <text x={activationNode.x} y={activationNode.y + 5} class="node-label activation" text-anchor="middle">
        {activationNode.label}()
      </text>
      {#if snapshot?.in_saturation_zone}
        <text x={activationNode.x} y={activationNode.y + 50} class="warning-label" text-anchor="middle">
          ⚠️ Saturated
        </text>
      {/if}
    </g>

    <!-- Output node -->
    <g>
      <circle
        cx={outputNode.x}
        cy={outputNode.y}
        r="25"
        fill="#dcfce7"
        stroke="#16a34a"
        stroke-width="2"
        opacity={getNodeOpacity(snapshot?.a ?? 0, 1)}
      />
      <text x={outputNode.x} y={outputNode.y + 5} class="node-label" text-anchor="middle">
        {outputNode.label}
      </text>
    </g>
  </svg>

  <!-- Legend -->
  <div class="legend">
    <div class="legend-section">
      <h4>Edge Colors:</h4>
      <div class="legend-item">
        <div class="legend-color" style="background: #3b82f6;"></div>
        <span>Positive weight</span>
      </div>
      <div class="legend-item">
        <div class="legend-color" style="background: #dc2626;"></div>
        <span>Negative weight</span>
      </div>
    </div>
    <div class="legend-section">
      <h4>Edge Width:</h4>
      <span>Thickness ∝ |weight magnitude|</span>
    </div>
    <div class="legend-section">
      <h4>Node Opacity:</h4>
      <span>Opacity ∝ |activation value|</span>
    </div>
  </div>
</div>

<style>
  .neuron-diagram {
    background: white;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    padding: 1.5rem;
    margin-top: 2rem;
  }

  h2 {
    margin: 0 0 0.5rem 0;
    font-size: 1.25rem;
    color: #1e293b;
    text-align: center;
  }

  .subtitle {
    text-align: center;
    color: #64748b;
    font-size: 0.95rem;
    margin: 0 0 1.5rem 0;
    font-style: italic;
  }

  .diagram {
    display: block;
    margin: 0 auto;
  }

  .node-label {
    font-size: 14px;
    font-weight: bold;
    fill: #1e293b;
    pointer-events: none;
  }

  .node-label.small {
    font-size: 11px;
  }

  .node-label.large {
    font-size: 24px;
  }

  .node-label.activation {
    font-size: 16px;
    fill: white;
  }

  .node-value {
    font-size: 10px;
    font-family: 'Courier New', monospace;
    fill: #475569;
    pointer-events: none;
  }

  .edge-label {
    font-size: 12px;
    font-weight: 600;
    font-family: 'Courier New', monospace;
    pointer-events: none;
  }

  .edge-label.small {
    font-size: 10px;
  }

  .warning-label {
    font-size: 11px;
    fill: #dc2626;
    font-weight: bold;
  }

  .legend {
    display: flex;
    justify-content: space-around;
    margin-top: 1.5rem;
    padding: 1rem;
    background: #f8fafc;
    border-radius: 6px;
    font-size: 0.9rem;
  }

  .legend-section {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .legend-section h4 {
    margin: 0;
    font-size: 0.85rem;
    color: #64748b;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .legend-color {
    width: 30px;
    height: 4px;
    border-radius: 2px;
  }

  @media (max-width: 768px) {
    .diagram {
      width: 100%;
      height: auto;
    }

    .legend {
      flex-direction: column;
      gap: 1rem;
    }
  }
</style>
