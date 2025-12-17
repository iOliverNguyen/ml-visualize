<script lang="ts">
  import Controls from './Controls.svelte';
  import Plots from './Plots.svelte';
  import MetricsPanel from './MetricsPanel.svelte';
  import LayoutSelector from './LayoutSelector.svelte';
  import StepBreakdown from './StepBreakdown.svelte';
  import FormulaOverlay from './FormulaOverlay.svelte';
  import IntroPanel from './educational/IntroPanel.svelte';
  import Sidebar from './educational/Sidebar.svelte';
  import QAAccordion from './educational/QAAccordion.svelte';
  import DataScatterPlot from './DataScatterPlot.svelte';
  import DatasetInfoPanel from './DatasetInfoPanel.svelte';
  import FitComparisonView from './FitComparisonView.svelte';
  import DataInputPanel from './DataInputPanel.svelte';
  import type { Snapshot, LayoutMode, MetricsData, Glossary, FAQData } from './types';
  import { loadAllContent } from './lib/contentLoader';
  import * as educationalState from './stores/educationalState.svelte';

  let snapshots = $state<Snapshot[]>([]);
  let currentStep = $state(0);
  let playing = $state(false);
  let intervalId = $state<number | null>(null);
  let loading = $state(true);
  let error = $state('');

  // Educational content
  let glossary = $state<Glossary>({});
  let faqs = $state<FAQData>({ categories: [] });
  let contentLoading = $state(true);

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

  // Load educational content on mount
  $effect(() => {
    async function loadContent() {
      try {
        const content = await loadAllContent();
        glossary = content.glossary;
        faqs = content.faqs;
        contentLoading = false;
      } catch (e) {
        console.error('Failed to load educational content:', e);
        contentLoading = false;
      }
    }
    loadContent();
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

  // Handle dataset change from DataInputPanel
  function handleDatasetChange(newSnapshots: Snapshot[]) {
    snapshots = newSnapshots;
    currentStep = 0;
    if (playing) {
      pause();
    }
  }
</script>

<main>
  <div class="header">
    <h1>Scalar Gradient Descent: y = w*x</h1>
    <div class="header-buttons">
      {#if !educationalState.state.showIntroPanel}
        <button
          class="intro-button"
          onclick={() => {
            educationalState.state.showIntroPanel = true;
            educationalState.saveToLocalStorage();
          }}
          aria-label="Show introduction panel"
          type="button"
        >
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
            <path d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          Welcome
        </button>
      {/if}
      <button
        class="help-button"
        onclick={() => educationalState.toggleSidebar()}
        aria-label="Toggle help sidebar"
        type="button"
      >
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <path d="M12 16V12M12 8H12.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        Help
      </button>
    </div>
  </div>

  <IntroPanel />

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
    <div class="data-section">
      <h2 class="section-title">Training Data & Objective</h2>
      <p class="section-description">
        Understand what we're training: fitting a line <strong>y = w*x</strong> to data points.
        Watch how predictions improve as training progresses.
      </p>

      <div class="data-layout">
        <DatasetInfoPanel
          dataset={snapshot?.point_details || []}
          currentW={snapshot?.w || 0}
          initialW={snapshots[0]?.w || 0}
          targetW={2.0}
          currentLoss={snapshot?.loss || 0}
        />

        <DataScatterPlot
          pointDetails={snapshot?.point_details || []}
          currentW={snapshot?.w || 0}
          showResiduals={true}
        />
      </div>

      <div class="fit-comparison-wrapper">
        <FitComparisonView
          initialSnapshot={snapshots[0]}
          currentSnapshot={snapshot}
          finalSnapshot={snapshots[snapshots.length - 1]}
        />
      </div>

      <details class="dataset-input-details">
        <summary class="dataset-input-summary">
          Change Dataset
        </summary>
        <DataInputPanel onDatasetChange={handleDatasetChange} loading={loading} />
      </details>
    </div>

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

{#if !contentLoading}
  <Sidebar {glossary} {faqs} />
{/if}

<style>
  main {
    padding: 2rem;
    max-width: 1400px;
    margin: 0 auto;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    gap: 1rem;
  }

  h1 {
    margin: 0;
    color: #1a1a1a;
  }

  .header-buttons {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .intro-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    background: #10b981;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    font-size: 1rem;
    transition: all 0.15s ease;
  }

  .intro-button:hover {
    background: #059669;
    transform: translateY(-1px);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }

  .intro-button:focus {
    outline: 2px solid #10b981;
    outline-offset: 2px;
  }

  .help-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    font-size: 1rem;
    transition: all 0.15s ease;
  }

  .help-button:hover {
    background: #5568d3;
    transform: translateY(-1px);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }

  .help-button:focus {
    outline: 2px solid #667eea;
    outline-offset: 2px;
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

  .data-section {
    margin-top: 2rem;
    padding: 2rem;
    background: linear-gradient(135deg, #f8fafc 0%, #e0f2fe 100%);
    border-radius: 12px;
    border: 1px solid #e2e8f0;
  }

  .section-title {
    margin: 0 0 0.5rem 0;
    color: #0f172a;
    font-size: 1.5rem;
    font-weight: 700;
  }

  .data-layout {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    margin-top: 1.5rem;
  }

  .fit-comparison-wrapper {
    margin-top: 1.5rem;
  }

  .dataset-input-details {
    margin-top: 1.5rem;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    overflow: hidden;
  }

  .dataset-input-summary {
    padding: 1rem 1.5rem;
    background: #f8fafc;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.95rem;
    color: #334155;
    list-style: none;
    user-select: none;
    transition: all 0.15s ease;
  }

  .dataset-input-summary:hover {
    background: #f1f5f9;
  }

  .dataset-input-summary::marker,
  .dataset-input-summary::-webkit-details-marker {
    display: none;
  }

  .dataset-input-summary::before {
    content: '▶';
    display: inline-block;
    margin-right: 0.5rem;
    transition: transform 0.15s ease;
    color: #3b82f6;
  }

  details[open] .dataset-input-summary::before {
    transform: rotate(90deg);
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
    main {
      padding: 1rem;
    }

    .header {
      flex-direction: column;
      align-items: flex-start;
    }

    h1 {
      font-size: 1.5rem;
    }

    .header-buttons {
      width: 100%;
      flex-direction: column;
    }

    .intro-button,
    .help-button {
      width: 100%;
      justify-content: center;
    }

    .data-section {
      padding: 1rem;
    }

    .data-layout {
      grid-template-columns: 1fr;
    }

    .formula-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
