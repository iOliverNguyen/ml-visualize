<script lang="ts">
  import Controls from '../shared/Controls.svelte';
  import LinearMetricsPanel from './LinearMetricsPanel.svelte';
  import GradientVectorViz from './GradientVectorViz.svelte';
  import ParameterSpace2D from './ParameterSpace2D.svelte';
  import GradientField from './GradientField.svelte';
  import LossContour from './LossContour.svelte';
  import type { LinearSnapshot, LossGrid, LinearMetricsData } from './types';

  let snapshots = $state<LinearSnapshot[]>([]);
  let lossGrid = $state<LossGrid | null>(null);
  let currentStep = $state(0);
  let playing = $state(false);
  let intervalId = $state<number | null>(null);
  let loading = $state(true);
  let error = $state('');

  // Reactive derived values
  let snapshot = $derived(snapshots[currentStep] || null);

  // Compute comprehensive metrics
  let metrics = $derived<LinearMetricsData>(() => {
    if (snapshots.length === 0) {
      return {
        initial: {} as LinearSnapshot,
        current: {} as LinearSnapshot,
        final: {} as LinearSnapshot,
        deltaFromStart: {
          w1: 0,
          w2: 0,
          loss: 0,
          gradient_magnitude: 0,
          distance: 0
        },
        deltaToTarget: {
          w1: 0,
          w2: 0,
          loss: 0,
          distance: 0
        },
        progress: 0
      };
    }

    const initial = snapshots[0];
    const current = snapshots[currentStep];
    const final = snapshots[snapshots.length - 1];

    const deltaW1 = current.w1 - initial.w1;
    const deltaW2 = current.w2 - initial.w2;
    const distanceFromStart = Math.sqrt(deltaW1 ** 2 + deltaW2 ** 2);

    const deltaW1ToTarget = final.w1 - current.w1;
    const deltaW2ToTarget = final.w2 - current.w2;
    const distanceToTarget = Math.sqrt(
      deltaW1ToTarget ** 2 + deltaW2ToTarget ** 2
    );

    return {
      initial,
      current,
      final,
      deltaFromStart: {
        w1: deltaW1,
        w2: deltaW2,
        loss: current.loss - initial.loss,
        gradient_magnitude:
          current.gradient_magnitude - initial.gradient_magnitude,
        distance: distanceFromStart
      },
      deltaToTarget: {
        w1: deltaW1ToTarget,
        w2: deltaW2ToTarget,
        loss: final.loss - current.loss,
        distance: distanceToTarget
      },
      progress:
        snapshots.length > 1 ? (currentStep / (snapshots.length - 1)) * 100 : 0
    };
  });

  // Load data on mount
  $effect(() => {
    async function loadData() {
      try {
        const snapshotsResp = await fetch(
          '/cases-phase2/lr-optimal/snapshots.json'
        );
        if (!snapshotsResp.ok) {
          throw new Error('Failed to load snapshots');
        }
        snapshots = await snapshotsResp.json();

        const gridResp = await fetch('/cases-phase2/lr-optimal/loss_grid.json');
        if (!gridResp.ok) {
          throw new Error('Failed to load loss grid');
        }
        lossGrid = await gridResp.json();

        loading = false;
      } catch (e) {
        console.error('Failed to load Phase 2 data:', e);
        error = e instanceof Error ? e.message : 'Unknown error';
        loading = false;
      }
    }
    loadData();
  });

  // Control functions
  function stepForward() {
    if (currentStep < snapshots.length - 1) {
      currentStep++;
    }
  }

  function stepBackward() {
    if (currentStep > 0) {
      currentStep--;
    }
  }

  function setStep(step: number) {
    currentStep = step;
  }

  function reset() {
    currentStep = 0;
    if (playing) {
      pause();
    }
  }

  function jumpToFinal() {
    currentStep = snapshots.length - 1;
    if (playing) {
      pause();
    }
  }

  function togglePlayback() {
    if (playing) {
      pause();
    } else {
      play();
    }
  }

  function play() {
    playing = true;
    intervalId = window.setInterval(() => {
      if (currentStep < snapshots.length - 1) {
        currentStep++;
      } else {
        pause();
      }
    }, 100);
  }

  function pause() {
    playing = false;
    if (intervalId !== null) {
      clearInterval(intervalId);
      intervalId = null;
    }
  }
</script>

<div class="phase2-app">
  <header class="app-header">
    <h1>Phase 2: 2-Parameter Gradient Descent</h1>
    <p class="subtitle">Visualizing y = w₁·x₁ + w₂·x₂ optimization</p>
  </header>

  {#if loading}
    <div class="loading">Loading Phase 2 data...</div>
  {:else if error}
    <div class="error">Error: {error}</div>
  {:else if !lossGrid || snapshots.length === 0}
    <div class="error">No data loaded</div>
  {:else}
    <Controls
      currentStep={currentStep}
      totalSteps={snapshots.length}
      playing={playing}
      onstepBack={stepBackward}
      onstepForward={stepForward}
      onplayPause={togglePlayback}
      onsliderChange={setStep}
      onreset={reset}
      onjumpToFinal={jumpToFinal}
    />

    <LinearMetricsPanel metrics={metrics()} />

    <div class="visualizations">
      <div class="viz-section">
        <h2>Parameter Space Trajectory</h2>
        <ParameterSpace2D
          snapshots={snapshots}
          currentStep={currentStep}
          lossGrid={lossGrid}
          onStepClick={setStep}
        />
      </div>

      <div class="viz-section">
        <h2>Loss Contour Surface</h2>
        <LossContour
          lossGrid={lossGrid}
          snapshots={snapshots}
          currentStep={currentStep}
          onStepClick={setStep}
        />
      </div>

      <div class="viz-section">
        <h2>Gradient Vector Field</h2>
        <GradientField
          lossGrid={lossGrid}
          snapshots={snapshots}
          currentStep={currentStep}
          onStepClick={setStep}
        />
      </div>

      <div class="viz-section">
        <h2>Gradient Vector</h2>
        {#if snapshot}
          <GradientVectorViz snapshot={snapshot} />
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .phase2-app {
    padding: 2rem;
    max-width: 1800px;
    margin: 0 auto;
  }

  .app-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .app-header h1 {
    margin: 0 0 0.5rem 0;
    font-size: 2rem;
    color: #2563eb;
  }

  .subtitle {
    margin: 0;
    font-size: 1.1rem;
    color: #64748b;
  }

  .loading,
  .error {
    text-align: center;
    padding: 3rem;
    font-size: 1.2rem;
  }

  .error {
    color: #dc2626;
  }

  .visualizations {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(700px, 1fr));
    gap: 2rem;
    margin-top: 2rem;
  }

  .viz-section {
    background: #f8fafc;
    padding: 1.5rem;
    border-radius: 8px;
    border: 1px solid #e2e8f0;
  }

  .viz-section h2 {
    margin: 0 0 1rem 0;
    font-size: 1.25rem;
    color: #334155;
  }

  @media (max-width: 1400px) {
    .visualizations {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 768px) {
    .phase2-app {
      padding: 1rem;
    }

    .app-header h1 {
      font-size: 1.5rem;
    }

    .subtitle {
      font-size: 1rem;
    }
  }
</style>
