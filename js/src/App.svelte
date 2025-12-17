<script lang="ts">
  import Controls from './Controls.svelte';
  import Plots from './Plots.svelte';
  import MetricsPanel from './MetricsPanel.svelte';
  import LayoutSelector from './LayoutSelector.svelte';
  import StepBreakdown from './StepBreakdown.svelte';
  import FormulaOverlay from './FormulaOverlay.svelte';
  import type { Snapshot, LayoutMode, MetricsData } from './types';

  let snapshots = $state<Snapshot[]>([]);
  let currentStep = $state(0);
  let playing = $state(false);
  let intervalId = $state<number | null>(null);
  let loading = $state(true);
  let error = $state('');

  // Layout mode with localStorage persistence
  let layoutMode = $state<LayoutMode>(
    (typeof localStorage !== 'undefined' ? localStorage.getItem('layoutMode') : null) as LayoutMode || '1col'
  );

  function handleLayoutChange(mode: LayoutMode) {
    layoutMode = mode;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('layoutMode', mode);
    }
  }

  // Reactive derived value - explicit dependency on snapshots and currentStep
  let snapshot = $derived(snapshots[currentStep]);

  // Compute comprehensive metrics
  let metrics = $derived<MetricsData>({
    initial: {
      w: snapshots[0]?.w ?? 0,
      grad_w: snapshots[0]?.grad_w ?? 0,
      loss: snapshots[0]?.loss ?? 0
    },
    current: {
      w: snapshots[currentStep]?.w ?? 0,
      grad_w: snapshots[currentStep]?.grad_w ?? 0,
      loss: snapshots[currentStep]?.loss ?? 0
    },
    final: {
      w: snapshots[snapshots.length - 1]?.w ?? 0,
      grad_w: snapshots[snapshots.length - 1]?.grad_w ?? 0,
      loss: snapshots[snapshots.length - 1]?.loss ?? 0
    },
    deltaFromStart: {
      w: (snapshots[currentStep]?.w ?? 0) - (snapshots[0]?.w ?? 0),
      loss: (snapshots[currentStep]?.loss ?? 0) - (snapshots[0]?.loss ?? 0),
      progress: snapshots.length > 1
        ? (currentStep / (snapshots.length - 1)) * 100
        : 0
    },
    deltaToTarget: {
      w: (snapshots[snapshots.length - 1]?.w ?? 0) - (snapshots[currentStep]?.w ?? 0),
      loss: (snapshots[snapshots.length - 1]?.loss ?? 0) - (snapshots[currentStep]?.loss ?? 0)
    }
  });

  // Load snapshots on mount
  $effect(() => {
    async function loadSnapshots() {
      try {
        // Try server mode first (Go server running)
        const response = await fetch('/api/snapshots');
        if (response.ok) {
          snapshots = await response.json();
          loading = false;
          return;
        }
      } catch (e) {
        // Server not available, try static file
      }

      // Fall back to batch mode (static file)
      try {
        const response = await fetch('/snapshots.json');
        if (!response.ok) {
          throw new Error('Snapshots not found');
        }
        snapshots = await response.json();
        loading = false;
      } catch (e) {
        error = 'Failed to load snapshots';
        loading = false;
      }
    }
    loadSnapshots();
  })

  function stepForward() {
    if (currentStep < snapshots.length - 1) {
      currentStep++
    }
  }

  function stepBackward() {
    if (currentStep > 0) {
      currentStep--
    }
  }

  function setStep(step: number) {
    currentStep = step
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
      pause()
    } else {
      play()
    }
  }

  function play() {
    playing = true
    intervalId = window.setInterval(() => {
      if (currentStep < snapshots.length - 1) {
        currentStep++
      } else {
        pause()
      }
    }, 100) // 10 steps per second
  }

  function pause() {
    playing = false
    if (intervalId !== null) {
      clearInterval(intervalId)
      intervalId = null
    }
  }
</script>

<main>
  <h1>Scalar Gradient Descent: y = w*x</h1>

  {#if loading}
    <p class="message">Loading snapshots...</p>
  {:else if error}
    <p class="error">{error}</p>
    <p class="hint">
      <strong>Option 1: Server Mode</strong><br>
      Terminal 1: <code>go run . --server</code><br>
      Terminal 2: <code>cd js && pnpm dev</code><br><br>

      <strong>Option 2: Batch Mode</strong><br>
      <code>go run .</code><br>
      <code>cp out/snapshots.json js/public/</code><br>
      <code>cd js && pnpm dev</code>
    </p>
  {:else if snapshots.length > 0}
    <LayoutSelector
      {layoutMode}
      onlayoutChange={handleLayoutChange}
    />

    <Controls
      {currentStep}
      totalSteps={snapshots.length}
      {playing}
      onstepBack={stepBackward}
      onstepForward={stepForward}
      onplayPause={togglePlayback}
      onsliderChange={setStep}
      onreset={reset}
      onjumpToFinal={jumpToFinal}
    />

    <MetricsPanel {metrics} />

    <Plots
      {snapshots}
      {currentStep}
      {layoutMode}
      onStepClick={setStep}
    />

    <div class="pedagogical-section">
      <h2>Step-by-Step Breakdown</h2>
      <p class="section-description">
        Explore what happens at each training step: forward pass, loss computation, gradient calculation, and parameter update.
      </p>
      <StepBreakdown snapshot={snapshot} />
    </div>

    <div class="formula-section">
      <h2>Interactive Formulas</h2>
      <p class="section-description">
        See the math behind gradient descent with live values from the current step.
      </p>
      <div class="formula-grid">
        <FormulaOverlay snapshot={snapshot} mode="loss" />
        <FormulaOverlay snapshot={snapshot} mode="gradient" />
        <FormulaOverlay snapshot={snapshot} mode="update" />
      </div>
    </div>
  {/if}
</main>

<style>
  main {
    padding: 2rem;
  }

  h1 {
    margin-bottom: 2rem;
    color: #1a1a1a;
  }

  .message {
    padding: 1rem;
    background: #e3f2fd;
    border-radius: 4px;
    color: #1976d2;
  }

  .error {
    padding: 1rem;
    background: #ffebee;
    border-radius: 4px;
    color: #c62828;
    margin-bottom: 1rem;
  }

  .hint {
    padding: 1rem;
    background: #f5f5f5;
    border-radius: 4px;
    font-family: monospace;
    font-size: 0.9rem;
    line-height: 1.6;
  }

  .hint code {
    background: #e0e0e0;
    padding: 0.2rem 0.4rem;
    border-radius: 2px;
  }

  .pedagogical-section,
  .formula-section {
    margin-top: 3rem;
    padding: 2rem;
    background: #f8fafc;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
  }

  .pedagogical-section h2,
  .formula-section h2 {
    margin: 0 0 0.5rem 0;
    color: #0f172a;
    font-size: 1.5rem;
  }

  .section-description {
    margin: 0 0 1.5rem 0;
    color: #64748b;
    font-size: 1rem;
    line-height: 1.6;
  }

  .formula-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
  }

  @media (max-width: 768px) {
    .formula-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
