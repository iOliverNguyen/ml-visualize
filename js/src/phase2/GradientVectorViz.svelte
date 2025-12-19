<script lang="ts">
  import { tweened } from 'svelte/motion';
  import { cubicOut } from 'svelte/easing';
  import type { LinearSnapshot } from './types';

  interface Props {
    snapshot: LinearSnapshot;
    mode?: 'gradient' | 'both';
  }

  let { snapshot, mode = 'both' }: Props = $props();

  const width = 800;
  const height = 600;
  const centerX = width / 2;
  const centerY = height / 2;

  // Scale factor to fit arrows in viewport
  const arrowScale = 80; // pixels per unit magnitude
  const maxLength = 250; // max arrow length in pixels

  // Animated gradient arrow properties
  const gradEndX = tweened(centerX, { duration: 200, easing: cubicOut });
  const gradEndY = tweened(centerY, { duration: 200, easing: cubicOut });

  // Animated update arrow properties
  const updateEndX = tweened(centerX, { duration: 200, easing: cubicOut });
  const updateEndY = tweened(centerY, { duration: 200, easing: cubicOut });

  // Update animations when snapshot changes
  $effect(() => {
    const magnitude = snapshot.gradient_magnitude;
    const direction = snapshot.gradient_direction;

    // Gradient vector arrow (pointing uphill in loss landscape)
    let gradLength = Math.min(magnitude * arrowScale, maxLength);

    // Convert polar to Cartesian (SVG: Y increases downward, so negate Y component)
    const gEndX = centerX + gradLength * Math.cos(direction);
    const gEndY = centerY - gradLength * Math.sin(direction); // Negate for SVG coords

    gradEndX.set(gEndX);
    gradEndY.set(gEndY);

    // Update vector arrow (pointing downhill - descent direction)
    const updateMagnitude = snapshot.update_components.lr * magnitude;
    let updateLength = Math.min(updateMagnitude * arrowScale, maxLength);

    // Update direction is opposite of gradient (add π radians)
    const updateDir = direction + Math.PI;
    const uEndX = centerX + updateLength * Math.cos(updateDir);
    const uEndY = centerY - updateLength * Math.sin(updateDir); // Negate for SVG coords

    updateEndX.set(uEndX);
    updateEndY.set(uEndY);
  });

  // Calculate arrow head polygon points
  function arrowHead(endX: number, endY: number, angle: number): string {
    const size = 12;
    const tipX = endX;
    const tipY = endY;

    // Left base point (rotated 150° from direction)
    const leftAngle = angle + (Math.PI * 5) / 6;
    const leftX = endX - size * Math.cos(leftAngle);
    const leftY = endY + size * Math.sin(leftAngle); // Note: + because SVG Y is flipped

    // Right base point (rotated -150° from direction)
    const rightAngle = angle - (Math.PI * 5) / 6;
    const rightX = endX - size * Math.cos(rightAngle);
    const rightY = endY + size * Math.sin(rightAngle); // Note: + because SVG Y is flipped

    return `${tipX},${tipY} ${leftX},${leftY} ${rightX},${rightY}`;
  }

  let gradArrowHead = $derived(
    arrowHead($gradEndX, $gradEndY, snapshot.gradient_direction)
  );
  let updateArrowHead = $derived(
    arrowHead($updateEndX, $updateEndY, snapshot.gradient_direction + Math.PI)
  );
</script>

<svg viewBox="0 0 {width} {height}" width="100%" height="auto" class="vector-viz">
  <!-- Coordinate axes -->
  <line
    x1="0"
    y1={centerY}
    x2={width}
    y2={centerY}
    stroke="#e0e0e0"
    stroke-width="2"
    stroke-dasharray="5,5"
  />
  <line
    x1={centerX}
    y1="0"
    x2={centerX}
    y2={height}
    stroke="#e0e0e0"
    stroke-width="2"
    stroke-dasharray="5,5"
  />

  <!-- Axis labels -->
  <text x={width - 30} y={centerY - 10} font-size="14" fill="#666">w₁</text>
  <text x={centerX + 10} y="20" font-size="14" fill="#666">w₂</text>

  {#if mode === 'both'}
    <!-- Update vector (green, solid) - drawn first so gradient is on top -->
    <line
      x1={centerX}
      y1={centerY}
      x2={$updateEndX}
      y2={$updateEndY}
      stroke="#10b981"
      stroke-width="3"
    />

    <!-- Update arrowhead -->
    <polygon points={updateArrowHead} fill="#10b981" />
  {/if}

  <!-- Gradient vector (blue, dashed) -->
  <line
    x1={centerX}
    y1={centerY}
    x2={$gradEndX}
    y2={$gradEndY}
    stroke="#3b82f6"
    stroke-width="3"
    stroke-dasharray="8,4"
  />

  <!-- Gradient arrowhead -->
  <polygon points={gradArrowHead} fill="#3b82f6" />

  <!-- Origin marker -->
  <circle cx={centerX} cy={centerY} r="4" fill="#333" />

  <!-- Legend -->
  <g transform="translate(20, 30)">
    <line
      x1="0"
      y1="0"
      x2="40"
      y2="0"
      stroke="#3b82f6"
      stroke-width="3"
      stroke-dasharray="8,4"
    />
    <text x="50" y="5" font-size="14" fill="#333">∇L (gradient)</text>

    {#if mode === 'both'}
      <line
        x1="0"
        y1="25"
        x2="40"
        y2="25"
        stroke="#10b981"
        stroke-width="3"
      />
      <text x="50" y="30" font-size="14" fill="#333">-α∇L (update)</text>
    {/if}
  </g>

  <!-- Vector component labels -->
  <text
    x="20"
    y={height - 90}
    font-family="monospace"
    font-size="14"
    fill="#333"
  >
    ∇L = ({snapshot.grad_w1.toFixed(3)}, {snapshot.grad_w2.toFixed(3)})
  </text>

  <text
    x="20"
    y={height - 65}
    font-family="monospace"
    font-size="14"
    fill="#333"
  >
    ||∇L|| = {snapshot.gradient_magnitude.toFixed(3)}
  </text>

  <text
    x="20"
    y={height - 40}
    font-family="monospace"
    font-size="14"
    fill="#333"
  >
    ∠∇L = {((snapshot.gradient_direction * 180) / Math.PI).toFixed(1)}°
  </text>

  {#if mode === 'both'}
    <text
      x={width - 300}
      y={height - 65}
      font-family="monospace"
      font-size="14"
      fill="#10b981"
    >
      Δw = ({snapshot.update_components.delta_w1.toFixed(3)}, {snapshot.update_components.delta_w2.toFixed(
        3
      )})
    </text>

    <text
      x={width - 300}
      y={height - 40}
      font-family="monospace"
      font-size="14"
      fill="#10b981"
    >
      α = {snapshot.update_components.lr.toFixed(4)}
    </text>
  {/if}
</svg>

<style>
  .vector-viz {
    border: 1px solid #eee;
    border-radius: 8px;
    background: white;
    max-height: 600px;
    width: 100%;
    display: block;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  }

  text {
    user-select: none;
  }

  @media (max-width: 768px) {
    .vector-viz {
      max-height: 400px;
    }
  }
</style>
