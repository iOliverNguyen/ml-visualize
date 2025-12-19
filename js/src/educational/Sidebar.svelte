<script lang="ts">
  import { slide } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import type { Glossary, FAQData, TutorialContent } from '../types';
  import * as educationalState from '../stores/educationalState.svelte';
  import GlossaryTooltip from './GlossaryTooltip.svelte';
  import QAAccordion from './QAAccordion.svelte';
  import TutorialArticle from './TutorialArticle.svelte';

  interface Props {
    glossary: Glossary;
    faqs: FAQData;
    tutorial: TutorialContent;
  }

  let { glossary, faqs, tutorial }: Props = $props();

  let activeTab = $state<'tutorial' | 'glossary' | 'faq'>('tutorial');
  let searchQuery = $state('');
  let tocOpen = $state(false);

  function handleClose() {
    educationalState.toggleSidebar();
  }

  function switchTab(tab: 'tutorial' | 'glossary' | 'faq') {
    activeTab = tab;
    searchQuery = '';
  }

  function toggleTOC() {
    tocOpen = !tocOpen;
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
  const filteredGlossaryTerms = $derived.by(() => {
    const query = searchQuery.toLowerCase();

    return Object.entries(glossary ?? {}).filter(([_, entry]) => {
      // Filter by search only (show all levels)
      const searchMatch = !query ||
        entry?.term?.toLowerCase().includes(query) ||
        entry?.brief?.toLowerCase().includes(query) ||
        entry?.detailed?.toLowerCase().includes(query);

      return searchMatch;
    });
  });
</script>

{#if educationalState.state.showSidebar}
  <aside
    class="sidebar"
    class:left={educationalState.state.sidebarPosition === 'left'}
    class:right={educationalState.state.sidebarPosition === 'right'}
  >
    <div class="sidebar-header">
      <h2 class="sidebar-title">Help & Reference</h2>

      {#if activeTab === 'tutorial'}
        <button
          class="toc-button"
          onclick={toggleTOC}
          aria-label={tocOpen ? 'Show tutorial content' : 'Show table of contents'}
          type="button"
        >
          {tocOpen ? '📄 Tutorial' : '📖 Contents'}
        </button>
      {/if}

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
        class:active={activeTab === 'tutorial'}
        onclick={() => switchTab('tutorial')}
        type="button"
      >
        📚 Tutorial
      </button>
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

    {#if activeTab !== 'tutorial'}
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
    {/if}

    <div class="sidebar-content" class:tutorial={activeTab === 'tutorial'}>
      {#if activeTab === 'tutorial'}
        <TutorialArticle {tutorial} {tocOpen} {toggleTOC} />
      {:else if activeTab === 'glossary'}
        <div class="glossary-list">
          {#if filteredGlossaryTerms.length === 0}
            <div class="empty-state">
              <p>No terms found.</p>
            </div>
          {:else}
            {#each filteredGlossaryTerms as [termId, entry]}
              <div data-term-id={termId}>
                <GlossaryTooltip
                  term={entry.term}
                  entry={entry}
                  inline={false}
                  onrelatedClick={handleGlossaryTermClick}
                />
              </div>
            {/each}
          {/if}
        </div>
      {:else}
        <QAAccordion
          categories={faqs?.categories ?? []}
          filterByLevel={false}
          searchQuery={searchQuery}
        />
      {/if}
    </div>
  </aside>
{/if}

<style>
  .sidebar {
    position: fixed;
    right: 0;
    top: 0;
    width: 550px;
    max-width: 90vw;
    height: 100vh;
    background: white;
    border-left: 1px solid #e2e8f0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    z-index: 100;
  }

  .sidebar.left {
    border-left: none;
    border-right: 1px solid #e2e8f0;
  }

  .sidebar.right {
    /* Right positioning handled by parent layout */
  }

  .sidebar-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem;
    border-bottom: 1px solid #e2e8f0;
    background: #f8fafc;
    gap: 0.75rem;
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

  .toc-button {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.5rem 0.75rem;
    background: #667eea;
    border: none;
    color: white;
    cursor: pointer;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 600;
    transition: all 0.15s ease;
    margin-left: auto;
    margin-right: 0.5rem;
  }

  .toc-button:hover {
    background: #5568d3;
  }

  .toc-button:focus {
    outline: 2px solid #667eea;
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

  .sidebar-content.tutorial {
    padding: 0;
  }

  .glossary-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .empty-state {
    padding: 2rem;
    text-align: center;
    color: #94a3b8;
  }

  .empty-state p {
    margin: 0;
  }

  @media (max-width: 768px) {
    .sidebar {
      width: 100%;
      max-width: 100vw;
    }
  }
</style>
