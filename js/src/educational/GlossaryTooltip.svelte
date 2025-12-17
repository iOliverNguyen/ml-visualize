<script lang="ts">
  import { fade, scale } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
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
    inline = true,
    onrelatedClick
  }: Props = $props();

  let isOpen = $state(false);
  let tooltipElement: HTMLDivElement | null = $state(null);
  let triggerElement: HTMLSpanElement | null = $state(null);

  function toggle() {
    isOpen = !isOpen;
    if (isOpen) {
      educationalState.state.activeGlossaryTerm = term;
    } else {
      educationalState.state.activeGlossaryTerm = null;
    }
  }

  function close() {
    isOpen = false;
    educationalState.state.activeGlossaryTerm = null;
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      toggle();
    } else if (e.key === 'Escape' && isOpen) {
      e.preventDefault();
      close();
    }
  }

  function handleRelatedClick(relatedTermId: string) {
    close();
    onrelatedClick?.(relatedTermId);
  }

  // Get content based on user level
  const content = $derived(() => {
    const level = educationalState.state.userLevel;
    if (level === 'advanced' && entry.comprehensive) {
      return entry.comprehensive;
    } else if (level === 'intermediate' || level === 'advanced') {
      return entry.detailed;
    } else {
      return entry.brief;
    }
  });

  // Close when clicking outside
  function handleClickOutside(event: MouseEvent) {
    if (
      isOpen &&
      tooltipElement &&
      triggerElement &&
      !tooltipElement.contains(event.target as Node) &&
      !triggerElement.contains(event.target as Node)
    ) {
      close();
    }
  }

  $effect(() => {
    if (isOpen) {
      document.addEventListener('click', handleClickOutside);
      return () => {
        document.removeEventListener('click', handleClickOutside);
      };
    }
  });
</script>

<span class="glossary-wrapper" class:inline>
  <span
    bind:this={triggerElement}
    class="glossary-term"
    class:active={isOpen}
    onclick={toggle}
    onkeydown={handleKeyDown}
    role="button"
    tabindex="0"
    aria-expanded={isOpen}
    aria-haspopup="true"
    aria-label="Definition for {entry.term}"
  >
    {term}
  </span>

  {#if isOpen}
    <div
      bind:this={tooltipElement}
      class="glossary-tooltip"
      transition:scale={{ duration: 150, easing: cubicOut, start: 0.95 }}
      role="dialog"
      aria-label="Definition of {entry.term}"
    >
      <div class="tooltip-header">
        <h4 class="tooltip-title">{entry.term}</h4>
        <button
          class="close-button"
          onclick={close}
          aria-label="Close definition"
          type="button"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path
              d="M12 4L4 12M4 4L12 12"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>
      </div>

      <div class="tooltip-content">
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

        <div class="level-indicator">
          <span class="level-badge" data-level={educationalState.state.userLevel}>
            {educationalState.state.userLevel} level
          </span>
        </div>
      </div>
    </div>
  {/if}
</span>

<style>
  .glossary-wrapper {
    position: relative;
    display: inline-block;
  }

  .glossary-wrapper.inline {
    display: inline;
  }

  .glossary-term {
    color: #2563eb;
    text-decoration: underline;
    text-decoration-style: dotted;
    text-decoration-thickness: 1px;
    text-underline-offset: 2px;
    cursor: pointer;
    transition: all 0.15s ease;
    font-weight: 500;
  }

  .glossary-term:hover {
    color: #1d4ed8;
    text-decoration-style: solid;
  }

  .glossary-term:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 2px;
    border-radius: 2px;
  }

  .glossary-term.active {
    color: #1d4ed8;
    background: #dbeafe;
    padding: 0 0.2rem;
    border-radius: 3px;
  }

  .glossary-tooltip {
    position: absolute;
    top: calc(100% + 0.5rem);
    left: 0;
    min-width: 300px;
    max-width: 400px;
    background: white;
    border: 1px solid #cbd5e1;
    border-radius: 8px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    z-index: 1000;
    padding: 0;
    overflow: hidden;
  }

  /* Adjust position if near right edge */
  @media (max-width: 640px) {
    .glossary-tooltip {
      left: 50%;
      transform: translateX(-50%);
      min-width: 280px;
      max-width: calc(100vw - 2rem);
    }
  }

  .tooltip-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    background: #f8fafc;
    border-bottom: 1px solid #e2e8f0;
  }

  .tooltip-title {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
    color: #0f172a;
  }

  .close-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    padding: 0;
    background: transparent;
    border: none;
    color: #64748b;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.15s ease;
  }

  .close-button:hover {
    background: #e2e8f0;
    color: #334155;
  }

  .close-button:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 1px;
  }

  .tooltip-content {
    padding: 1rem;
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
    background: #f8fafc;
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
    font-size: 0.8rem;
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
    font-size: 0.8rem;
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

  .level-indicator {
    display: flex;
    justify-content: flex-end;
    margin-top: 0.75rem;
    padding-top: 0.5rem;
    border-top: 1px solid #e2e8f0;
  }

  .level-badge {
    display: inline-flex;
    align-items: center;
    padding: 0.15rem 0.5rem;
    border-radius: 4px;
    font-size: 0.65rem;
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
</style>
