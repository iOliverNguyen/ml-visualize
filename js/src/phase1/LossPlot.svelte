<script lang="ts">
  import { tweened } from 'svelte/motion'
  import { cubicOut } from 'svelte/easing'
  import type { Snapshot } from './types'

  interface Props {
    snapshots: Snapshot[]
    currentStep: number
    onStepClick?: (step: number) => void
  }

  let { snapshots, currentStep, onStepClick }: Props = $props()

  // Hover state
  let hoveredStep = $state<number | null>(null)

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

  <!-- Interactive data point circles -->
  {#each snapshots as snapshot, i}
    <circle
      cx={scaleX(i)}
      cy={scaleY(snapshot.loss)}
      r={i === currentStep ? 6 : hoveredStep === i ? 5 : 3}
      fill={i === currentStep ? '#dc2626' : hoveredStep === i ? '#3b82f6' : '#94a3b8'}
      class="data-point"
      onmouseenter={() => (hoveredStep = i)}
      onmouseleave={() => (hoveredStep = null)}
      onclick={() => onStepClick?.(i)}
    />
  {/each}

  <!-- Hover tooltip -->
  {#if hoveredStep !== null}
    {@const x = scaleX(hoveredStep)}
    {@const y = scaleY(snapshots[hoveredStep].loss)}
    {@const tooltipX = x > width - 150 ? x - 140 : x + 10}
    {@const tooltipY = y > height - 60 ? y - 60 : y + 10}

    <g class="tooltip">
      <rect
        x={tooltipX}
        y={tooltipY}
        width="130"
        height="50"
        fill="white"
        stroke="#cbd5e1"
        stroke-width="1"
        rx="4"
      />
      <text
        x={tooltipX + 10}
        y={tooltipY + 20}
        font-size="12"
        fill="#475569"
        font-weight="600"
      >
        Step {hoveredStep}
      </text>
      <text
        x={tooltipX + 10}
        y={tooltipY + 38}
        font-size="11"
        fill="#64748b"
        font-family="monospace"
      >
        Loss: {snapshots[hoveredStep].loss.toFixed(4)}
      </text>
    </g>
  {/if}

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
