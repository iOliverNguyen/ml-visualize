<script lang="ts">
  import { onMount } from 'svelte';
  import type { NeuronTrainingCase } from './types';

  let caseData = $state<NeuronTrainingCase | null>(null);
  let loading = $state(true);
  let error = $state<string | null>(null);

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

  let selectedCase = $state('sigmoid-vanishing');

  async function loadCase(caseId: string) {
    loading = true;
    error = null;
    try {
      const response = await fetch(`/cases-phase3/${caseId}/snapshots.json`);
      if (!response.ok) {
        throw new Error(`Failed to load case: ${response.statusText}`);
      }
      const data = await response.json();
      caseData = data as NeuronTrainingCase;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Unknown error';
      console.error('Error loading case:', err);
    } finally {
      loading = false;
    }
  }

  function handleCaseChange(caseId: string) {
    selectedCase = caseId;
    loadCase(caseId);
  }

  onMount(() => {
    loadCase(selectedCase);
  });
</script>

<div class="phase3-app">
  <div class="header">
    <h1>🧠 Phase 3: Single Neuron with Activation Functions</h1>
    <p class="subtitle">
      Explore how activation functions (σ, ReLU, tanh) affect training through saturation, vanishing gradients, and the chain rule
    </p>
  </div>

  <div class="case-selector">
    <label for="case-select">Select Case:</label>
    <select id="case-select" value={selectedCase} onchange={(e) => handleCaseChange((e.target as HTMLSelectElement).value)}>
      {#each availableCases as caseOption}
        <option value={caseOption.id}>
          {caseOption.emoji} {caseOption.name}
        </option>
      {/each}
    </select>
  </div>

  {#if loading}
    <div class="loading">Loading case data...</div>
  {:else if error}
    <div class="error">Error: {error}</div>
  {:else if caseData}
    <div class="case-info">
      <h2>{caseData.case_id}</h2>
      <p><strong>Description:</strong> {caseData.description}</p>
      <p><strong>Category:</strong> {caseData.category}</p>
      <p><strong>Activation:</strong> {caseData.config.activation}</p>
      <p><strong>Learning Rate:</strong> {caseData.config.learning_rate}</p>
      <p><strong>Steps:</strong> {caseData.snapshots.length}</p>
      <p><strong>Dataset Size:</strong> {caseData.dataset.length} points</p>
    </div>

    <div class="placeholder">
      <h3>🚧 Phase 3 Visualizations Coming Soon</h3>
      <ul>
        <li>✓ Backend complete (8 cases generated)</li>
        <li>✓ TypeScript types defined</li>
        <li>✓ Case loading functional</li>
        <li>⏸️ Activation curve visualization</li>
        <li>⏸️ Chain rule breakdown</li>
        <li>⏸️ Saturation indicator</li>
        <li>⏸️ Neuron diagram</li>
        <li>⏸️ Z-value tracker</li>
        <li>⏸️ Metrics panel</li>
      </ul>
    </div>
  {/if}
</div>

<style>
  .phase3-app {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }

  .header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .header h1 {
    font-size: 2rem;
    margin: 0 0 0.5rem 0;
    color: #1e293b;
  }

  .subtitle {
    color: #64748b;
    font-size: 1rem;
    margin: 0;
  }

  .case-selector {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 2rem;
    padding: 1rem;
    background: #f8fafc;
    border-radius: 8px;
  }

  .case-selector label {
    font-weight: 600;
    color: #334155;
  }

  .case-selector select {
    flex: 1;
    max-width: 400px;
    padding: 0.5rem;
    font-size: 1rem;
    border: 2px solid #e2e8f0;
    border-radius: 4px;
    background: white;
    cursor: pointer;
  }

  .loading, .error {
    text-align: center;
    padding: 2rem;
    font-size: 1.2rem;
  }

  .error {
    color: #dc2626;
    background: #fee2e2;
    border-radius: 8px;
  }

  .case-info {
    background: white;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 2rem;
  }

  .case-info h2 {
    margin: 0 0 1rem 0;
    color: #1e293b;
  }

  .case-info p {
    margin: 0.5rem 0;
    color: #334155;
  }

  .placeholder {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 2rem;
    border-radius: 12px;
    text-align: center;
  }

  .placeholder h3 {
    margin: 0 0 1rem 0;
    font-size: 1.5rem;
  }

  .placeholder ul {
    list-style: none;
    padding: 0;
    margin: 1rem 0 0 0;
    text-align: left;
    display: inline-block;
  }

  .placeholder li {
    margin: 0.5rem 0;
    font-size: 1rem;
    padding-left: 1rem;
  }
</style>
