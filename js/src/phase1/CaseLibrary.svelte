<script lang="ts">
  import type { CaseManifest, CaseConfig, Snapshot } from './types';
  import DataInputPanel from './DataInputPanel.svelte';

  interface Props {
    onCaseSelect: (caseId: string, snapshots: Snapshot[], caseConfig?: any) => void;
    currentCaseId?: string;
    onDatasetChange: (snapshots: Snapshot[], config?: any) => void;
  }

  let { onCaseSelect, currentCaseId = '', onDatasetChange }: Props = $props();

  let customDataExpanded = $state(false);

  let manifest = $state<CaseManifest>({ version: '1.0', cases: [] });
  let loading = $state(true);
  let error = $state('');
  let selectedCategory = $state<'all' | 'foundational' | 'learning-rate' | 'initialization'>('all');

  // Load manifest on mount
  $effect(() => {
    async function loadManifest() {
      try {
        const response = await fetch(`${import.meta.env.BASE_URL}cases/manifest.json`);
        if (!response.ok) {
          throw new Error('Failed to load case manifest');
        }
        manifest = await response.json();
        loading = false;
      } catch (e) {
        error = 'Failed to load case library';
        loading = false;
        console.error(e);
      }
    }
    loadManifest();
  });

  // Filter cases by category
  const filteredCases = $derived.by(() => {
    if (selectedCategory === 'all') {
      return manifest.cases;
    }
    return manifest.cases.filter(c => c.category === selectedCategory);
  });

  async function handleCaseClick(caseConfig: CaseConfig) {
    try {
      const response = await fetch(`${import.meta.env.BASE_URL}cases/${caseConfig.id}/snapshots.json`);
      if (!response.ok) {
        throw new Error(`Failed to load case ${caseConfig.id}`);
      }
      const snapshots = await response.json();
      onCaseSelect(caseConfig.id, snapshots, caseConfig);
    } catch (e) {
      console.error('Failed to load case:', e);
      alert(`Failed to load case: ${caseConfig.name}`);
    }
  }

  function getCategoryLabel(category: string): string {
    switch (category) {
      case 'foundational': return 'Foundations';
      case 'learning-rate': return 'Learning Rate';
      case 'initialization': return 'Initialization';
      default: return 'All Cases';
    }
  }

  function getCategoryEmoji(category: string): string {
    switch (category) {
      case 'foundational': return '📚';
      case 'learning-rate': return '⚡';
      case 'initialization': return '🎯';
      default: return '📖';
    }
  }
</script>

<div class="case-library">
  <div class="library-header">
    <h2>Explore Training Scenarios</h2>
    <p class="library-description">
      Learn by comparing different training conditions. Each case shows what happens when you change data quality, learning rate, or starting point.
    </p>
  </div>

  {#if loading}
    <div class="message">
      <p>Loading case library...</p>
      <p class="patience-text">Please wait, this may take up to 15 seconds.</p>
    </div>
  {:else if error}
    <p class="error">{error}</p>
  {:else}
    <div class="category-filter">
      <button
        class="filter-btn"
        class:active={selectedCategory === 'all'}
        onclick={() => selectedCategory = 'all'}
        type="button"
      >
        📖 All Cases
      </button>
      <button
        class="filter-btn"
        class:active={selectedCategory === 'foundational'}
        onclick={() => selectedCategory = 'foundational'}
        type="button"
      >
        📚 Foundations
      </button>
      <button
        class="filter-btn"
        class:active={selectedCategory === 'learning-rate'}
        onclick={() => selectedCategory = 'learning-rate'}
        type="button"
      >
        ⚡ Learning Rate
      </button>
      <button
        class="filter-btn"
        class:active={selectedCategory === 'initialization'}
        onclick={() => selectedCategory = 'initialization'}
        type="button"
      >
        🎯 Initialization
      </button>
    </div>

    <div class="cases-grid">
      {#each filteredCases as caseConfig (caseConfig.id)}
        <button
          class="case-card"
          class:active={currentCaseId === caseConfig.id}
          onclick={() => handleCaseClick(caseConfig)}
          type="button"
        >
          <div class="case-header">
            <span class="case-emoji">{caseConfig.emoji}</span>
            <h3 class="case-name">{caseConfig.name}</h3>
          </div>

          <p class="case-description">{caseConfig.description}</p>

          <div class="case-category">
            <span class="category-badge">
              {getCategoryEmoji(caseConfig.category)} {getCategoryLabel(caseConfig.category)}
            </span>
          </div>

          <div class="case-insights">
            <p class="insights-title">What you'll see:</p>
            <ul>
              {#each caseConfig.insights as insight}
                <li>{insight}</li>
              {/each}
            </ul>
          </div>

          <div class="case-config">
            <span class="config-item">LR: {caseConfig.training_config.lr}</span>
            <span class="config-item">Steps: {caseConfig.training_config.steps}</span>
            <span class="config-item">w₀: {caseConfig.training_config.w_init}</span>
          </div>
        </button>
      {/each}

      <!-- Custom Data Card -->
      {#if selectedCategory === 'all' || selectedCategory === 'initialization'}
        <button
          class="case-card custom-data-card"
          class:active={currentCaseId === 'custom'}
          class:expanded={customDataExpanded}
          onclick={() => customDataExpanded = !customDataExpanded}
          type="button"
        >
          <div class="case-header">
            <span class="case-emoji">🎨</span>
            <h3 class="case-name">Custom Data</h3>
          </div>

          <p class="case-description">
            Create your own dataset with custom parameters. Generate random data or define it with JavaScript.
          </p>

          <div class="case-category">
            <span class="category-badge custom-badge">
              ✨ Your Own
            </span>
          </div>

          <div class="case-insights">
            <p class="insights-title">You can:</p>
            <ul>
              <li>Generate random data with custom noise</li>
              <li>Define data using JavaScript functions</li>
              <li>Experiment with any configuration</li>
            </ul>
          </div>

          <div class="custom-action">
            <span class="action-text">Click to create →</span>
          </div>
        </button>

        <!-- Expandable Custom Data Section -->
        {#if customDataExpanded}
          <div class="custom-data-expanded">
            <div class="expanded-header">
              <h3>Create Custom Dataset</h3>
              <button
                class="collapse-button"
                onclick={() => customDataExpanded = false}
                type="button"
                aria-label="Collapse"
              >
                ✕
              </button>
            </div>
            <DataInputPanel
              onDatasetChange={onDatasetChange}
              loading={false}
            />
          </div>
        {/if}
      {/if}
    </div>
  {/if}
</div>

<style>
  .case-library {
    padding: 2rem;
    max-width: 1400px;
    margin: 0 auto;
  }

  .library-header {
    margin-bottom: 2rem;
  }

  .library-header h2 {
    margin: 0 0 0.5rem 0;
    color: #0f172a;
    font-size: 1.75rem;
  }

  .library-description {
    margin: 0;
    color: #64748b;
    font-size: 1.05rem;
    line-height: 1.6;
  }

  .message {
    padding: 1.5rem;
    background: #e0f2fe;
    border-radius: 8px;
    color: #0369a1;
    text-align: center;
  }

  .patience-text {
    margin-top: 0.5rem;
    font-size: 1rem;
    color: #0c326a;
    font-style: italic;
  }

  .error {
    padding: 1.5rem;
    background: #fee2e2;
    border-radius: 8px;
    color: #991b1b;
    text-align: center;
  }

  .category-filter {
    display: flex;
    gap: 0.75rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
  }

  .filter-btn {
    padding: 0.75rem 1.5rem;
    background: white;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    font-size: 0.95rem;
    font-weight: 600;
    color: #64748b;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .filter-btn:hover {
    border-color: #cbd5e1;
    background: #f8fafc;
  }

  .filter-btn.active {
    background: #3b82f6;
    border-color: #3b82f6;
    color: white;
  }

  .cases-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 1.5rem;
  }

  .case-card {
    padding: 1.5rem;
    background: white;
    border: 2px solid #e2e8f0;
    border-radius: 12px;
    text-align: left;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .case-card:hover {
    border-color: #3b82f6;
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  }

  .case-card.active {
    border-color: #10b981;
    background: #f0fdf4;
    box-shadow: 0 4px 12px rgba(16, 185, 129, 0.2);
  }

  .case-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .case-emoji {
    font-size: 2rem;
  }

  .case-name {
    margin: 0;
    font-size: 1.25rem;
    color: #0f172a;
  }

  .case-description {
    margin: 0;
    color: #475569;
    line-height: 1.5;
    font-size: 0.95rem;
  }

  .case-category {
    display: flex;
    align-items: center;
  }

  .category-badge {
    padding: 0.375rem 0.75rem;
    background: #f1f5f9;
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 600;
    color: #475569;
  }

  .case-insights {
    background: #f8fafc;
    padding: 1rem;
    border-radius: 8px;
    border-left: 3px solid #3b82f6;
  }

  .insights-title {
    margin: 0 0 0.5rem 0;
    font-weight: 600;
    color: #334155;
    font-size: 0.9rem;
  }

  .case-insights ul {
    margin: 0;
    padding-left: 1.25rem;
    list-style-type: '→';
  }

  .case-insights li {
    padding-left: 0.5rem;
    color: #64748b;
    font-size: 0.875rem;
    line-height: 1.5;
    margin-bottom: 0.25rem;
  }

  .case-insights li:last-child {
    margin-bottom: 0;
  }

  .case-config {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-top: auto;
  }

  .config-item {
    padding: 0.375rem 0.75rem;
    background: #e0f2fe;
    border-radius: 6px;
    font-size: 0.825rem;
    font-weight: 600;
    color: #0369a1;
    font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
  }

  .custom-data-card {
    background: linear-gradient(135deg, #f8fafc 0%, #e0f2fe 100%);
    border: 2px dashed #3b82f6;
  }

  .custom-data-card:hover {
    background: linear-gradient(135deg, #e0f2fe 0%, #bfdbfe 100%);
    border-color: #2563eb;
    border-style: solid;
  }

  .custom-data-card.active {
    background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
    border-color: #3b82f6;
    border-style: solid;
  }

  .custom-badge {
    background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
    color: white;
  }

  .custom-action {
    margin-top: auto;
    padding: 0.75rem;
    background: #3b82f6;
    border-radius: 8px;
    text-align: center;
  }

  .action-text {
    color: white;
    font-weight: 600;
    font-size: 0.95rem;
  }

  .custom-data-expanded {
    grid-column: 1 / -1;
    margin-top: 1rem;
    padding: 2rem;
    background: white;
    border: 2px solid #3b82f6;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
  }

  .expanded-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #e2e8f0;
  }

  .expanded-header h3 {
    margin: 0;
    font-size: 1.25rem;
    color: #0f172a;
  }

  .collapse-button {
    background: none;
    border: none;
    font-size: 1.5rem;
    color: #64748b;
    cursor: pointer;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
    transition: all 0.2s ease;
  }

  .collapse-button:hover {
    background: #f1f5f9;
    color: #0f172a;
  }

  .custom-data-card.expanded {
    border-color: #3b82f6;
    border-style: solid;
    background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
  }

  @media (max-width: 768px) {
    .case-library {
      padding: 1rem;
    }

    .library-header h2 {
      font-size: 1.5rem;
    }

    .cases-grid {
      grid-template-columns: 1fr;
    }

    .category-filter {
      flex-direction: column;
    }

    .filter-btn {
      width: 100%;
    }
  }
</style>
