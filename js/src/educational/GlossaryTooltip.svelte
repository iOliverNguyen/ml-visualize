<script lang="ts">
  import { slide } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import { onDestroy } from 'svelte';
  import type { GlossaryEntry } from '../types';
  import * as educationalState from '../stores/educationalState.svelte';

  interface Props {
    term: string;
    entry: GlossaryEntry;
    inline?: boolean;
    onrelatedClick?: (termId: string) => void;
  }

  let {
    term,
    entry,
    inline = false,
    onrelatedClick
  }: Props = $props();

  let expanded = $state(false);
  let isHighlighted = $state(false);
  let highlightTimeout: ReturnType<typeof setTimeout> | null = null;

  function toggle() {
    expanded = !expanded;
    if (expanded) {
      educationalState.state.activeGlossaryTerm = term;
    } else {
      educationalState.state.activeGlossaryTerm = null;
    }
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      toggle();
    }
  }

  function handleRelatedClick(relatedTermId: string) {
    onrelatedClick?.(relatedTermId);
  }

  // Watch for external activation via related term clicks
  $effect(() => {
    const activeTermFromState = educationalState.state.activeGlossaryTerm;

    // If this term was activated externally (e.g., from related link)
    if (activeTermFromState === term && !expanded) {
      // Auto-expand
      expanded = true;

      // Trigger highlight
      isHighlighted = true;

      // Clear any existing timeout
      if (highlightTimeout) clearTimeout(highlightTimeout);

      // Remove highlight after 2.5 seconds
      highlightTimeout = setTimeout(() => {
        isHighlighted = false;
      }, 2500);
    }
  });

  // Cleanup timeout on component unmount
  onDestroy(() => {
    if (highlightTimeout) clearTimeout(highlightTimeout);
  });

  // Get content based on user level
  const content = $derived.by(() => {
    const level = educationalState.state.userLevel;
    if (level === 'advanced' && entry.comprehensive) {
      return entry.comprehensive;
    } else if (level === 'intermediate' || level === 'advanced') {
      return entry.detailed;
    } else {
      return entry.brief;
    }
  });
</script>

{#if inline}
  <!-- Inline version for use within text -->
  <span class="glossary-inline">
    <button
      class="glossary-term-inline"
      onclick={toggle}
      onkeydown={handleKeyDown}
      aria-expanded={expanded}
      type="button"
    >
      {term}
    </button>
  </span>
{:else}
  <!-- Block version for glossary list in sidebar -->
  <div class="glossary-accordion" class:highlighted={isHighlighted}>
    <button
      class="glossary-header"
      onclick={toggle}
      onkeydown={handleKeyDown}
      aria-expanded={expanded}
      type="button"
    >
      <div class="header-content">
        <h4 class="term-title">{entry.term}</h4>
      </div>
      <svg
        class="chevron"
        class:rotated={expanded}
        width="20"
        height="20"
        viewBox="0 0 20 20"
        fill="none"
      >
        <path
          d="M5 7.5L10 12.5L15 7.5"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </button>

    {#if expanded}
      <div
        class="glossary-content"
        transition:slide={{ duration: 200, easing: cubicOut }}
      >
        <p class="definition">{content()}</p>

        {#if entry.formula}
          <div class="formula">
            <code>{entry.formula}</code>
          </div>
        {/if}

        {#if entry.example}
          <div class="example">
            <strong>Example:</strong>
            <p>{entry.example}</p>
          </div>
        {/if}

        {#if entry.relatedTerms && entry.relatedTerms.length > 0}
          <div class="related-terms">
            <strong>Related:</strong>
            <div class="related-list">
              {#each entry.relatedTerms as relatedTerm}
                <button
                  class="related-term-link"
                  onclick={() => handleRelatedClick(relatedTerm)}
                  type="button"
                >
                  {relatedTerm}
                </button>
              {/each}
            </div>
          </div>
        {/if}
      </div>
    {/if}
  </div>
{/if}

<style>
  /* Inline version styles */
  .glossary-inline {
    display: inline;
  }

  .glossary-term-inline {
    display: inline;
    padding: 0;
    background: transparent;
    border: none;
    color: #2563eb;
    text-decoration: underline;
    text-decoration-style: dotted;
    text-decoration-thickness: 1px;
    text-underline-offset: 2px;
    cursor: pointer;
    transition: all 0.15s ease;
    font-weight: 500;
    font-size: inherit;
    font-family: inherit;
  }

  .glossary-term-inline:hover {
    color: #1d4ed8;
    text-decoration-style: solid;
  }

  .glossary-term-inline:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 2px;
    border-radius: 2px;
  }

  /* Block accordion version styles */
  .glossary-accordion {
    border: 1px solid #e2e8f0;
    border-radius: 6px;
    overflow: hidden;
    background: white;
    transition: all 0.15s ease;
  }

  .glossary-accordion:hover {
    border-color: #cbd5e1;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  }

  .glossary-accordion.highlighted {
    animation: highlightPulse 2.5s ease-out;
    border-color: #fbbf24;
  }

  @keyframes highlightPulse {
    0% {
      background: #fef3c7;
      box-shadow: 0 0 0 3px rgba(251, 191, 36, 0.3);
    }
    100% {
      background: white;
      box-shadow: none;
    }
  }

  .glossary-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    padding: 0.875rem 1rem;
    background: transparent;
    border: none;
    cursor: pointer;
    transition: background 0.15s ease;
    text-align: left;
  }

  .glossary-header:hover {
    background: #f8fafc;
  }

  .glossary-header:focus {
    outline: 2px solid #3b82f6;
    outline-offset: -2px;
  }

  .header-content {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex: 1;
  }

  .term-title {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
    color: #0f172a;
    flex: 1;
  }

  .chevron {
    flex-shrink: 0;
    color: #64748b;
    transition: transform 0.2s ease;
  }

  .chevron.rotated {
    transform: rotate(180deg);
  }

  .glossary-content {
    padding: 1rem;
    border-top: 1px solid #e2e8f0;
    background: #fafbfc;
  }

  .definition {
    margin: 0 0 0.75rem 0;
    font-size: 0.9rem;
    line-height: 1.6;
    color: #334155;
  }

  .formula {
    margin: 0.75rem 0;
    padding: 0.75rem;
    background: white;
    border-left: 3px solid #3b82f6;
    border-radius: 4px;
  }

  .formula code {
    font-family: 'SF Mono', 'Courier New', monospace;
    font-size: 0.9rem;
    color: #1e40af;
    font-weight: 500;
  }

  .example {
    margin: 0.75rem 0;
    padding: 0.75rem;
    background: #fefce8;
    border-left: 3px solid #eab308;
    border-radius: 4px;
  }

  .example strong {
    display: block;
    margin-bottom: 0.25rem;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: #92400e;
  }

  .example p {
    margin: 0;
    font-size: 0.85rem;
    line-height: 1.5;
    color: #713f12;
  }

  .related-terms {
    margin: 0.75rem 0 0 0;
    padding-top: 0.75rem;
    border-top: 1px solid #e2e8f0;
  }

  .related-terms strong {
    display: block;
    margin-bottom: 0.5rem;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: #64748b;
  }

  .related-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .related-term-link {
    padding: 0.25rem 0.5rem;
    font-size: 0.8rem;
    color: #2563eb;
    background: #eff6ff;
    border: 1px solid #bfdbfe;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.15s ease;
    font-weight: 500;
  }

  .related-term-link:hover {
    background: #dbeafe;
    border-color: #93c5fd;
  }

  .related-term-link:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 1px;
  }
</style>
