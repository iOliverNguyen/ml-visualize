<script lang="ts">
  import { slide } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import type { Glossary, FAQData } from '../types';
  import * as educationalState from '../stores/educationalState.svelte';
  import GlossaryTooltip from './GlossaryTooltip.svelte';
  import QAAccordion from './QAAccordion.svelte';

  interface Props {
    glossary: Glossary;
    faqs: FAQData;
  }

  let { glossary, faqs }: Props = $props();

  let activeTab = $state<'glossary' | 'faq'>('glossary');
  let searchQuery = $state('');

  function handleClose() {
    educationalState.toggleSidebar();
  }

  function switchTab(tab: 'glossary' | 'faq') {
    activeTab = tab;
    searchQuery = '';
  }

  function handleGlossaryTermClick(termId: string) {
    educationalState.state.activeGlossaryTerm = termId;
    // Scroll to term if needed
    const element = document.querySelector(`[data-term-id="${termId}"]`);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }
  }

  // Filter glossary terms
  const filteredGlossaryTerms = $derived(() => {
    const query = searchQuery.toLowerCase();
    const userLevel = educationalState.state.userLevel;

    return Object.entries(glossary).filter(([_, entry]) => {
      // Filter by level
      const levelMatch = entry.level.includes(userLevel) || entry.level.includes('all');

      // Filter by search
      const searchMatch = !query ||
        entry.term.toLowerCase().includes(query) ||
        entry.brief.toLowerCase().includes(query) ||
        entry.detailed.toLowerCase().includes(query);

      return levelMatch && searchMatch;
    });
  });
</script>

{#if educationalState.state.showSidebar}
  <aside
    class="sidebar"
    class:left={educationalState.state.sidebarPosition === 'left'}
    class:right={educationalState.state.sidebarPosition === 'right'}
    transition:slide={{ duration: 300, easing: cubicOut, axis: 'x' }}
  >
    <div class="sidebar-header">
      <h2 class="sidebar-title">Help & Reference</h2>
      <button
        class="close-button"
        onclick={handleClose}
        aria-label="Close sidebar"
        type="button"
      >
        <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
          <path
            d="M15 5L5 15M5 5L15 15"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          />
        </svg>
      </button>
    </div>

    <div class="tabs">
      <button
        class="tab"
        class:active={activeTab === 'glossary'}
        onclick={() => switchTab('glossary')}
        type="button"
      >
        📖 Glossary
      </button>
      <button
        class="tab"
        class:active={activeTab === 'faq'}
        onclick={() => switchTab('faq')}
        type="button"
      >
        ❓ FAQ
      </button>
    </div>

    <div class="search-box">
      <svg class="search-icon" width="16" height="16" viewBox="0 0 16 16" fill="none">
        <circle cx="7" cy="7" r="5" stroke="currentColor" stroke-width="2"/>
        <path d="M11 11L14 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
      <input
        type="text"
        class="search-input"
        placeholder={activeTab === 'glossary' ? 'Search terms...' : 'Search questions...'}
        bind:value={searchQuery}
      />
      {#if searchQuery}
        <button
          class="clear-search"
          onclick={() => searchQuery = ''}
          aria-label="Clear search"
          type="button"
        >
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
            <path
              d="M10 4L4 10M4 4L10 10"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>
      {/if}
    </div>

    <div class="sidebar-content">
      {#if activeTab === 'glossary'}
        <div class="glossary-list">
          {#if filteredGlossaryTerms().length === 0}
            <div class="empty-state">
              <p>No terms found.</p>
            </div>
          {:else}
            {#each filteredGlossaryTerms() as [termId, entry]}
              <div class="glossary-item" data-term-id={termId}>
                <GlossaryTooltip
                  term={entry.term}
                  entry={entry}
                  inline={false}
                  onrelatedClick={handleGlossaryTermClick}
                />
                <p class="term-brief">{entry.brief}</p>
              </div>
            {/each}
          {/if}
        </div>
      {:else}
        <QAAccordion
          categories={faqs.categories}
          filterByLevel={true}
          searchQuery={searchQuery}
        />
      {/if}
    </div>

    <div class="sidebar-footer">
      <div class="level-indicator">
        <span class="level-label">Current Level:</span>
        <span class="level-badge" data-level={educationalState.state.userLevel}>
          {educationalState.state.userLevel}
        </span>
      </div>
    </div>
  </aside>
{/if}

<style>
  .sidebar {
    position: fixed;
    top: 0;
    bottom: 0;
    width: 550px;
    max-width: 90vw;
    background: white;
    border-left: 1px solid #e2e8f0;
    box-shadow: -4px 0 6px -1px rgba(0, 0, 0, 0.1);
    z-index: 1000;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .sidebar.left {
    left: 0;
    border-left: none;
    border-right: 1px solid #e2e8f0;
    box-shadow: 4px 0 6px -1px rgba(0, 0, 0, 0.1);
  }

  .sidebar.right {
    right: 0;
  }

  .sidebar-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem;
    border-bottom: 1px solid #e2e8f0;
    background: #f8fafc;
  }

  .sidebar-title {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: #0f172a;
  }

  .close-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    padding: 0;
    background: transparent;
    border: none;
    color: #64748b;
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.15s ease;
  }

  .close-button:hover {
    background: #e2e8f0;
    color: #334155;
  }

  .close-button:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 2px;
  }

  .tabs {
    display: flex;
    border-bottom: 1px solid #e2e8f0;
    background: #f8fafc;
  }

  .tab {
    flex: 1;
    padding: 0.75rem 1rem;
    background: transparent;
    border: none;
    color: #64748b;
    cursor: pointer;
    transition: all 0.15s ease;
    font-weight: 500;
    font-size: 0.95rem;
    position: relative;
  }

  .tab:hover {
    background: #f1f5f9;
    color: #334155;
  }

  .tab.active {
    color: #3b82f6;
    background: white;
  }

  .tab.active::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: #3b82f6;
  }

  .tab:focus {
    outline: 2px solid #3b82f6;
    outline-offset: -2px;
  }

  .search-box {
    position: relative;
    padding: 1rem;
    border-bottom: 1px solid #e2e8f0;
  }

  .search-icon {
    position: absolute;
    left: 1.75rem;
    top: 50%;
    transform: translateY(-50%);
    color: #94a3b8;
    pointer-events: none;
  }

  .search-input {
    width: 100%;
    padding: 0.5rem 2.5rem 0.5rem 2.5rem;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    font-size: 0.9rem;
    transition: all 0.15s ease;
  }

  .search-input:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 0;
    border-color: #3b82f6;
  }

  .clear-search {
    position: absolute;
    right: 1.75rem;
    top: 50%;
    transform: translateY(-50%);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    padding: 0;
    background: transparent;
    border: none;
    color: #94a3b8;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.15s ease;
  }

  .clear-search:hover {
    background: #f1f5f9;
    color: #64748b;
  }

  .sidebar-content {
    flex: 1;
    overflow-y: auto;
    padding: 1rem;
  }

  .glossary-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .glossary-item {
    padding: 1rem;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 6px;
    transition: all 0.15s ease;
  }

  .glossary-item:hover {
    background: #f1f5f9;
    border-color: #cbd5e1;
  }

  .term-brief {
    margin: 0.5rem 0 0 0;
    font-size: 0.85rem;
    color: #64748b;
    line-height: 1.5;
  }

  .empty-state {
    padding: 2rem;
    text-align: center;
    color: #94a3b8;
  }

  .empty-state p {
    margin: 0;
  }

  .sidebar-footer {
    padding: 1rem;
    border-top: 1px solid #e2e8f0;
    background: #f8fafc;
  }

  .level-indicator {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .level-label {
    font-size: 0.85rem;
    color: #64748b;
    font-weight: 500;
  }

  .level-badge {
    display: inline-flex;
    align-items: center;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .level-badge[data-level='beginner'] {
    background: #dcfce7;
    color: #166534;
  }

  .level-badge[data-level='intermediate'] {
    background: #fef3c7;
    color: #92400e;
  }

  .level-badge[data-level='advanced'] {
    background: #fee2e2;
    color: #991b1b;
  }

  @media (max-width: 768px) {
    .sidebar {
      width: 100%;
      max-width: 100vw;
    }
  }
</style>
