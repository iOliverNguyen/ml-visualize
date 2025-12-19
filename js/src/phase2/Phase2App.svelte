<script lang="ts">
  import Controls from '../shared/Controls.svelte';
  import LinearMetricsPanel from './LinearMetricsPanel.svelte';
  import GradientVectorViz from './GradientVectorViz.svelte';
  import ParameterSpace2D from './ParameterSpace2D.svelte';
  import GradientField from './GradientField.svelte';
  import LossContour from './LossContour.svelte';
  import DataScatterPlot2D from './DataScatterPlot2D.svelte';
  import type { LinearSnapshot, LossGrid, LinearMetricsData } from './types';
  import {
    trainModel2D,
    generateRandomData2D,
    computeLossGrid
  } from '../shared/training-phase2';
  import IntroPanel from '../educational/IntroPanel.svelte';
  import Sidebar from '../educational/Sidebar.svelte';
  import type { Glossary, FAQData, TutorialContent } from '../types';
  import { loadAllContent } from '../lib/contentLoader';
  import * as educationalState from '../stores/educationalState.svelte';

  let snapshots = $state<LinearSnapshot[]>([]);
  let lossGrid = $state<LossGrid | null>(null);
  let currentStep = $state(0);
  let playing = $state(false);
  let intervalId = $state<number | null>(null);
  let loading = $state(true);
  let error = $state('');
  let selectedPointIndex = $state<number | null>(null);

  // Progressive rendering states
  let lossGridReady = $state(false);
  let snapshotsReady = $state(false);

  // Educational content
  let glossary = $state<Glossary>({});
  let faqs = $state<FAQData>({ categories: [] });
  let tutorial = $state<TutorialContent>({
    meta: { title: '', description: '', estimatedReadTime: 0, version: '1.0' },
    chapters: []
  });
  let contentLoading = $state(true);

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

  // Load data on mount with client-side training
  $effect(() => {
    async function loadData() {
      try {
        loading = true;
        lossGridReady = false;
        snapshotsReady = false;

        // Fetch only 5KB config file
        const configResp = await fetch(`${import.meta.env.BASE_URL}cases-phase2/lr-optimal/config.json`);
        if (!configResp.ok) {
          throw new Error('Failed to load config');
        }
        const config = await configResp.json();

        // Generate training data from config
        const data = generateRandomData2D(config.data_config);

        // Compute loss grid first (100-150ms) - enables contour visualization
        lossGrid = computeLossGrid(data, config.loss_grid_config);
        lossGridReady = true; // LossContour can render now!

        // Train model (50-100ms) - generates trajectory
        snapshots = trainModel2D({
          data,
          w1_init: config.training_config.w1_init,
          w2_init: config.training_config.w2_init,
          lr: config.training_config.lr,
          max_steps: config.training_config.max_steps
        });
        snapshotsReady = true; // Trajectory can render now!

        loading = false;
      } catch (e) {
        console.error('Failed to load Phase 2 data:', e);
        error = e instanceof Error ? e.message : 'Unknown error';
        loading = false;
      }
    }
    loadData();
  });

  // Load educational content on mount
  $effect(() => {
    async function loadContent() {
      try {
        const content = await loadAllContent('phase2');
        glossary = content.glossary;
        faqs = content.faqs;
        tutorial = content.tutorial;
        contentLoading = false;
      } catch (e) {
        console.error('Failed to load Phase 2 educational content:', e);
        contentLoading = false;
      }
    }
    loadContent();
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

  function handleSelectPoint(index: number | null) {
    selectedPointIndex = index;
  }
</script>

<div class="app-layout">
  <main class:with-sidebar={educationalState.state.showSidebar && !contentLoading}>
    <header class="app-header">
      <div class="header-content">
        <h1>Phase 2: 2-Parameter Gradient Descent</h1>
        <p class="subtitle">Visualizing y = w₁·x₁ + w₂·x₂ optimization</p>
      </div>
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
    </header>

    <IntroPanel />

    {#if loading}
      <div class="loading">
        <div>Loading Phase 2 data...</div>
        <div class="patience-text">Please wait, this may take up to 15 seconds.</div>
      </div>
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

      <!-- NEW: Data Scatter Plot Section -->
      <div class="viz-section data-view">
        <h2>Training Data Visualization</h2>
        {#if snapshot}
          <DataScatterPlot2D
            {snapshot}
            {selectedPointIndex}
            onSelectPoint={handleSelectPoint}
          />
        {/if}
      </div>

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
  </main>

  {#if !contentLoading && educationalState.state.showSidebar}
    <Sidebar {glossary} {faqs} {tutorial} />
  {/if}
</div>

<style>
  .app-layout {
    display: flex;
    min-height: 100vh;
  }

  main {
    flex: 1;
    min-width: 0;
    padding: 2rem;
    max-width: 1800px;
    margin: 0 auto;
    transition: margin-right 0.3s ease;
  }

  main.with-sidebar {
    margin-right: 550px;
  }

  .app-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    gap: 1rem;
  }

  .header-content h1 {
    margin: 0 0 0.5rem 0;
    font-size: 2rem;
    color: #2563eb;
  }

  .header-content .subtitle {
    margin: 0;
    font-size: 1.1rem;
    color: #64748b;
  }

  .header-buttons {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    flex-shrink: 0;
  }

  .intro-button,
  .help-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    font-size: 1rem;
    transition: all 0.15s ease;
    white-space: nowrap;
  }

  .intro-button {
    background: #10b981;
    color: white;
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
    background: #667eea;
    color: white;
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

  .loading,
  .error {
    text-align: center;
    padding: 3rem;
    font-size: 1.2rem;
  }

  .patience-text {
    margin-top: 0.5rem;
    font-size: 1rem;
    color: #0c326a;
    font-style: italic;
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

  .data-view {
    background: #f8fafc;
    padding: 1.5rem;
    border-radius: 8px;
    border: 1px solid #e2e8f0;
    margin-bottom: 2rem;
  }

  @media (max-width: 1400px) {
    .visualizations {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 768px) {
    main {
      padding: 1rem;
    }

    main.with-sidebar {
      margin-right: 0;
    }

    .app-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .header-buttons {
      width: 100%;
    }

    .intro-button,
    .help-button {
      flex: 1;
      justify-content: center;
    }
  }
</style>
