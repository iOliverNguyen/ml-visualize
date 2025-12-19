<script lang="ts">
  import Controls from '../shared/Controls.svelte';
  import NeuronMetricsPanel from './NeuronMetricsPanel.svelte';
  import ActivationCurveViz from './ActivationCurveViz.svelte';
  import ChainRuleBreakdown from './ChainRuleBreakdown.svelte';
  import NeuronDiagram from './NeuronDiagram.svelte';
  import Sidebar from '../educational/Sidebar.svelte';
  import IntroPanel from '../educational/IntroPanel.svelte';
  import * as educationalState from '../stores/educationalState.svelte';
  import { loadAllContent } from '../lib/contentLoader';
  import type { NeuronSnapshot, NeuronTrainingCase } from './types';
  import type { Glossary, FAQData, TutorialContent } from '../types';

  // State management
  let snapshots = $state<NeuronSnapshot[]>([]);
  let currentStep = $state(0);
  let playing = $state(false);
  let intervalId = $state<number | null>(null);
  let loading = $state(true);
  let error = $state('');
  let selectedCase = $state('sigmoid-optimal');
  let caseInfo = $state<{description: string; category: string; activation: string} | null>(null);

  // Educational content
  let glossary = $state<Glossary>({});
  let faqs = $state<FAQData>({ categories: [] });
  let tutorial = $state<TutorialContent>({
    meta: { title: '', description: '', estimatedReadTime: 0, version: '1.0' },
    chapters: []
  });
  let contentLoading = $state(true);

  // List of available cases
  const availableCases = [
    { id: 'sigmoid-vanishing', name: 'Sigmoid Vanishing', emoji: '🔻' },
    { id: 'sigmoid-optimal', name: 'Sigmoid Optimal', emoji: '✓' },
    { id: 'relu-dying', name: 'ReLU Dying', emoji: '💀' },
    { id: 'relu-optimal', name: 'ReLU Optimal', emoji: '⚡' },
    { id: 'tanh-saturation', name: 'Tanh Saturation', emoji: '〰️' },
    { id: 'tanh-optimal', name: 'Tanh Optimal', emoji: '✓' },
    { id: 'activation-comparison', name: 'Activation Comparison', emoji: '🔬' },
    { id: 'lr-saturation-interaction', name: 'LR-Saturation', emoji: '⚠️' },
  ];

  // Derived state
  let snapshot = $derived(snapshots[currentStep] || null);

  interface NeuronMetricsData {
    current: NeuronSnapshot;
    initial: NeuronSnapshot;
    final: NeuronSnapshot;
    deltaParams: {
      w1: number;
      w2: number;
      b: number;
    };
    deltaLoss: number;
    progress: number;
  }

  let metrics = $derived<NeuronMetricsData>(() => {
    // Return empty placeholder when no data loaded yet
    if (snapshots.length === 0 || !snapshot) {
      const emptySnapshot: NeuronSnapshot = {
        step: 0,
        params: { w: [0, 0], b: 0 },
        grads: { grad_w: [0, 0], grad_b: 0 },
        z: 0,
        a: 0,
        dL_dz: 0,
        dL_da: 0,
        local_derivative: 0,
        activation: 'sigmoid',
        in_saturation_zone: false,
        loss: 0,
        point_details: [],
        update_components: {
          w1_component: 0,
          w2_component: 0,
          b_component: 0,
          learning_rate: 0,
          current_loss: 0
        },
        chain_rule_breakdown: {
          components: [
            {
              param_name: 'w1',
              dL_da: 0,
              da_dz: 0,
              dz_dparam: 0,
              dL_dparam: 0
            },
            {
              param_name: 'w2',
              dL_da: 0,
              da_dz: 0,
              dz_dparam: 0,
              dL_dparam: 0
            },
            {
              param_name: 'b',
              dL_da: 0,
              da_dz: 0,
              dz_dparam: 0,
              dL_dparam: 0
            }
          ]
        }
      };

      return {
        current: emptySnapshot,
        initial: emptySnapshot,
        final: emptySnapshot,
        deltaParams: { w1: 0, w2: 0, b: 0 },
        deltaLoss: 0,
        progress: 0
      };
    }

    const initial = snapshots[0];
    const current = snapshot;
    const final = snapshots[snapshots.length - 1];

    const deltaW1 = current.params.w[0] - initial.params.w[0];
    const deltaW2 = current.params.w[1] - initial.params.w[1];
    const deltaB = current.params.b - initial.params.b;

    return {
      current,
      initial,
      final,
      deltaParams: {
        w1: deltaW1,
        w2: deltaW2,
        b: deltaB
      },
      deltaLoss: current.loss - initial.loss,
      progress: snapshots.length > 1 ? (currentStep / (snapshots.length - 1)) * 100 : 0
    };
  });

  // Load case data
  $effect(() => {
    async function loadCase() {
      loading = true;
      error = '';
      try {
        const response = await fetch(`${import.meta.env.BASE_URL}cases-phase3/${selectedCase}/snapshots.json`);
        if (!response.ok) {
          throw new Error(`Failed to load case: ${response.statusText}`);
        }
        const data = await response.json() as NeuronTrainingCase;

        snapshots = data.snapshots;
        caseInfo = {
          description: data.description,
          category: data.category,
          activation: data.config.activation
        };
        currentStep = 0;
        loading = false;
      } catch (err) {
        error = err instanceof Error ? err.message : 'Unknown error';
        console.error('Error loading case:', err);
        loading = false;
      }
    }
    loadCase();
  });

  // Load educational content
  $effect(() => {
    async function loadContent() {
      try {
        const content = await loadAllContent('phase3');
        glossary = content.glossary;
        faqs = content.faqs;
        tutorial = content.tutorial;
      } catch (error) {
        console.error('Failed to load Phase 3 educational content:', error);
        // Use empty defaults if content fails to load
      } finally {
        contentLoading = false;
      }
    }
    loadContent();
  });

  // Playback control functions
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
      togglePlayback();
    }
  }

  function jumpToFinal() {
    currentStep = snapshots.length - 1;
    if (playing) {
      togglePlayback();
    }
  }

  function togglePlayback() {
    if (playing) {
      // Stop playback
      if (intervalId !== null) {
        clearInterval(intervalId);
        intervalId = null;
      }
      playing = false;
    } else {
      // Start playback
      playing = true;
      intervalId = window.setInterval(() => {
        if (currentStep < snapshots.length - 1) {
          currentStep++;
        } else {
          // Reached end, stop playback
          if (intervalId !== null) {
            clearInterval(intervalId);
            intervalId = null;
          }
          playing = false;
        }
      }, 100); // 100ms per step (10 steps/second)
    }
  }

  function handleCaseChange(caseId: string) {
    // Stop playback if active
    if (playing) {
      togglePlayback();
    }
    selectedCase = caseId;
  }
</script>

<div class="app-layout">
  <main class:with-sidebar={educationalState.state.showSidebar && !contentLoading}>
    <header class="app-header">
      <div class="header-content">
        <h1>Phase 3: Single Neuron with Activation Functions</h1>
        <p class="subtitle">
          Explore how activation functions (σ, ReLU, tanh) affect training through saturation, vanishing gradients, and the chain rule
        </p>
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

  <div class="case-selector">
    <label for="case-select">Select Case:</label>
    <select id="case-select" value={selectedCase} onchange={(e) => handleCaseChange((e.target as HTMLSelectElement).value)}>
      {#each availableCases as caseOption}
        <option value={caseOption.id}>
          {caseOption.emoji} {caseOption.name}
        </option>
      {/each}
    </select>
    {#if caseInfo}
      <div class="case-description">
        {caseInfo.description}
      </div>
    {/if}
  </div>

  {#if loading}
    <div class="loading">
      <div>Loading case data...</div>
      <div class="patience-text">Please wait, this may take up to 15 seconds.</div>
    </div>
  {:else if error}
    <div class="error">Error: {error}</div>
  {:else if snapshots.length === 0}
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

    {#if metrics}
      <NeuronMetricsPanel metrics={metrics()} />
    {/if}

    {#if snapshot}
      <ActivationCurveViz {snapshot} />
      <ChainRuleBreakdown {snapshot} />
      <NeuronDiagram {snapshot} />
    {/if}

    <div class="placeholder">
      <h3>🎉 Phase 3 Visualizations Complete!</h3>
      <ul>
        <li>✓ Metrics panel</li>
        <li>✓ Activation curve visualization</li>
        <li>✓ Chain rule breakdown</li>
        <li>✓ Neuron diagram</li>
      </ul>
      <p style="margin-top: 1rem; font-size: 0.95rem;">
        All core educational components are now functional. Try different cases to see how activation functions affect training!
      </p>
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

  .case-selector {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 2rem;
    padding: 1rem;
    background: #f8fafc;
    border-radius: 8px;
    flex-wrap: wrap;
  }

  .case-selector label {
    font-weight: 600;
    color: #334155;
  }

  .case-selector select {
    flex: 0 1 auto;
    min-width: 250px;
    padding: 0.5rem;
    font-size: 1rem;
    border: 2px solid #e2e8f0;
    border-radius: 4px;
    background: white;
    cursor: pointer;
  }

  .case-description {
    flex: 1 1 100%;
    margin-top: 0.5rem;
    color: #475569;
    font-size: 0.95rem;
    font-style: italic;
  }

  .loading, .error {
    text-align: center;
    padding: 2rem;
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
    background: #fee2e2;
    border-radius: 8px;
  }

  .placeholder {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 2rem;
    border-radius: 12px;
    margin-top: 2rem;
  }

  .placeholder h3 {
    margin: 0 0 1rem 0;
    font-size: 1.5rem;
  }

  .placeholder p {
    margin: 0.5rem 0;
  }

  .placeholder ul {
    list-style: none;
    padding: 0;
    margin: 1rem 0;
  }

  .placeholder li {
    margin: 0.5rem 0;
    font-size: 1rem;
    padding-left: 1rem;
  }

  @media (max-width: 1024px) {
    main.with-sidebar {
      margin-right: 0;
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

    .case-selector {
      flex-direction: column;
      align-items: stretch;
    }

    .case-selector select {
      width: 100%;
    }
  }
</style>
