<script lang="ts">
  import type { Snapshot, LayoutMode } from './types';
  import LossPlot from './LossPlot.svelte';
  import ParamPlot from './ParamPlot.svelte';
  import GradientViz from './GradientViz.svelte';

  interface Props {
    snapshots: Snapshot[];
    currentStep: number;
    layoutMode: LayoutMode;
    onStepClick?: (step: number) => void;
  }

  let { snapshots, currentStep, layoutMode, onStepClick }: Props = $props();

  let layoutClass = $derived(`layout-${layoutMode}`);
</script>

<div class="plots {layoutClass}">
  <div class="plot">
    <h2>Loss over Time</h2>
    <LossPlot {snapshots} {currentStep} {onStepClick} />
  </div>

  <div class="plot">
    <h2>Parameter w over Time</h2>
    <ParamPlot {snapshots} {currentStep} {onStepClick} />
  </div>

  <div class="plot">
    <h2>Gradient at Current Step</h2>
    <GradientViz snapshot={snapshots[currentStep]} />
  </div>
</div>

<style>
  /* Base plots container */
  .plots {
    gap: 2rem;
  }

  /* 1 Column - vertical stack */
  .plots.layout-1col {
    display: flex;
    flex-direction: column;
  }

  /* 2 Columns - side by side */
  .plots.layout-2col {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
  }

  /* 3 Columns - three across */
  .plots.layout-3col {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
  }

  /* Auto - responsive */
  .plots.layout-auto {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(600px, 1fr));
  }

  /* Mobile: force single column */
  @media (max-width: 768px) {
    .plots.layout-2col,
    .plots.layout-3col,
    .plots.layout-auto {
      grid-template-columns: 1fr;
    }
  }

  .plot {
    display: grid;
    grid-template-rows: auto 1fr;
    gap: 1rem;
    border: 1px solid #ddd;
    padding: 1.5rem;
    border-radius: 8px;
    background: white;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    align-items: start;
  }

  h2 {
    margin: 0;
    font-size: 1.2rem;
    color: #1a1a1a;
  }
</style>
