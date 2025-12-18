<script lang="ts">
  import { tweened } from 'svelte/motion'
  import { cubicOut } from 'svelte/easing'
  import type { Snapshot } from './types'

  interface Props {
    snapshot: Snapshot
  }

  let { snapshot }: Props = $props()

  const width = 800
  const height = 200
  const centerX = width / 2
  const centerY = height / 2

  // Animated arrow properties
  const animatedLength = tweened(0, { duration: 200, easing: cubicOut })
  const animatedDirection = tweened(0, { duration: 200, easing: cubicOut })

  // Update animations when snapshot changes
  $effect(() => {
    const length = Math.min(Math.abs(snapshot.grad_w) * 50, 200)
    const direction = snapshot.grad_w > 0 ? -1 : 1

    animatedLength.set(length)
    animatedDirection.set(direction)
  })

  // Calculated end position for arrow
  let arrowEndX = $derived(centerX + $animatedDirection * $animatedLength)
</script>

<svg viewBox="0 0 {width} {height}" width="100%" height="auto" class="plot-svg">
  <!-- Center line -->
  <line
    x1={centerX}
    y1="0"
    x2={centerX}
    y2={height}
    stroke="#e0e0e0"
    stroke-width="1"
  />

  <!-- Gradient arrow line -->
  <line
    x1={centerX}
    y1={centerY}
    x2={arrowEndX}
    y2={centerY}
    stroke="#7c3aed"
    stroke-width="3"
  />

  <!-- Arrow head -->
  <polygon
    points={`
      ${arrowEndX},${centerY}
      ${arrowEndX - $animatedDirection * 10},${centerY - 10}
      ${arrowEndX - $animatedDirection * 10},${centerY + 10}
    `}
    fill="#7c3aed"
  />

  <!-- Value labels -->
  <text
    x="20"
    y="30"
    font-family="monospace"
    font-size="16"
    fill="#333"
  >
    grad_w = {snapshot.grad_w.toFixed(4)}
  </text>

  <text
    x="20"
    y="55"
    font-family="monospace"
    font-size="16"
    fill="#333"
  >
    w = {snapshot.w.toFixed(4)}
  </text>

  <!-- Direction indicator -->
  <text
    x={centerX}
    y={height - 20}
    text-anchor="middle"
    font-size="14"
    fill="#666"
  >
    {#if snapshot.grad_w > 0}
      ← w decreases
    {:else if snapshot.grad_w < 0}
      w increases →
    {:else}
      w stays (grad = 0)
    {/if}
  </text>
</svg>

<style>
  .plot-svg {
    border: 1px solid #eee;
    border-radius: 4px;
    background: white;
    max-height: 300px;
    width: 100%;
    display: block;
  }
</style>
