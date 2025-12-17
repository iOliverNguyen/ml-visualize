<script lang="ts">
  import { tweened } from 'svelte/motion'
  import { cubicOut } from 'svelte/easing'
  import type { Snapshot } from './types'

  interface Props {
    snapshots: Snapshot[]
    currentStep: number
  }

  let { snapshots, currentStep }: Props = $props()

  // Dimensions and padding
  const width = 800
  const height = 300
  const padding = { left: 40, right: 30, top: 30, bottom: 50 }
  const plotWidth = width - padding.left - padding.right
  const plotHeight = height - padding.top - padding.bottom

  // Scaling calculations
  let losses = $derived(snapshots.map(s => s.loss))
  let minLoss = $derived(Math.min(...losses))
  let maxLoss = $derived(Math.max(...losses))
  let lossRange = $derived(maxLoss - minLoss || 1)

  // Animated current marker position
  const markerX = tweened(0, { duration: 150, easing: cubicOut })
  const markerY = tweened(0, { duration: 150, easing: cubicOut })

  function scaleX(index: number): number {
    return padding.left + (index / (snapshots.length - 1)) * plotWidth
  }

  function scaleY(loss: number): number {
    return height - padding.bottom - ((loss - minLoss) / lossRange) * plotHeight
  }

  // Update animated marker when currentStep changes
  $effect(() => {
    markerX.set(scaleX(currentStep))
    markerY.set(scaleY(snapshots[currentStep].loss))
  })

  // Generate path data for the line
  let pathData = $derived(
    snapshots
      .map((s, i) => `${i === 0 ? 'M' : 'L'} ${scaleX(i)} ${scaleY(s.loss)}`)
      .join(' ')
  )
</script>

<svg viewBox="0 0 {width} {height}" width="100%" height="auto" class="plot-svg">
  <!-- Line path -->
  <path
    d={pathData}
    fill="none"
    stroke="#2563eb"
    stroke-width="2"
  />

  <!-- Animated current step marker -->
  <circle
    cx={$markerX}
    cy={$markerY}
    r="6"
    fill="#dc2626"
  />

  <!-- Axes labels -->
  <text
    x={width / 2}
    y={height - 10}
    text-anchor="middle"
    font-size="14"
    fill="#333"
  >
    Step
  </text>

  <text
    x="15"
    y={height / 2}
    text-anchor="middle"
    font-size="14"
    fill="#333"
    transform="rotate(-90 15 {height / 2})"
  >
    Loss
  </text>

  <!-- Current value display -->
  <text
    x={width - 10}
    y="20"
    text-anchor="end"
    font-family="monospace"
    font-size="12"
    fill="#dc2626"
  >
    Loss: {snapshots[currentStep].loss.toFixed(4)}
  </text>
</svg>

<style>
  .plot-svg {
    border: 1px solid #eee;
    border-radius: 4px;
    background: white;
    max-height: 400px;
    width: 100%;
    display: block;
  }
</style>
