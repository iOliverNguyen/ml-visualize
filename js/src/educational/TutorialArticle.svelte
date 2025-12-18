<script lang="ts">
  import { slide } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import type { TutorialContent } from '../types';

  interface Props {
    tutorial: TutorialContent;
    tocOpen: boolean;
    toggleTOC: () => void;
  }

  let { tutorial, tocOpen, toggleTOC }: Props = $props();

  // Track reading progress
  let activeChapter = $state<number>(1);

  // Smooth scroll to chapter and close TOC
  function scrollToChapter(chapterNumber: number) {
    activeChapter = chapterNumber;
    if (tocOpen) {
      toggleTOC(); // Close TOC to show content
    }
    // Wait for content to render, then scroll
    setTimeout(() => {
      const element = document.querySelector(`[data-chapter="${chapterNumber}"]`);
      if (element) {
        element.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
    }, 100);
  }
</script>

<div class="tutorial-article">
  <!-- Show EITHER TOC or content, not both -->
  {#if tocOpen}
    <nav class="tutorial-toc">
      <div class="toc-header">
        <h4>Contents</h4>
      </div>
      <ol class="chapter-list">
        {#each tutorial.chapters as chapter}
          <li class:active={activeChapter === chapter.number}>
            <button onclick={() => scrollToChapter(chapter.number)} type="button">
              <span class="chapter-number">{chapter.number}.</span>
              <span class="chapter-title">{chapter.title}</span>
            </button>
          </li>
        {/each}
      </ol>
    </nav>
  {:else}
    <article class="tutorial-content">
    {#if tutorial.chapters.length === 0}
      <div class="empty-state">
        <p>No tutorial content available.</p>
      </div>
    {:else}
      {#each tutorial.chapters as chapter}
        <section class="tutorial-chapter" data-chapter={chapter.number}>
          <h3 class="chapter-heading">
            <span class="chapter-number">{chapter.number}.</span>
            {chapter.title}
          </h3>

          {#each chapter.sections as section}
            {#if section.type === 'text'}
              <p class="chapter-text">{section.content}</p>

            {:else if section.type === 'highlight'}
              <div class="chapter-highlight">
                {section.content}
              </div>

            {:else if section.type === 'code'}
              <div class="chapter-code">
                <pre><code>{section.code?.snippet || section.content}</code></pre>
              </div>

            {:else if section.type === 'callout'}
              <div class="chapter-callout" data-type={section.calloutType || 'info'}>
                <span class="callout-icon">
                  {#if section.calloutType === 'info'}ℹ️{/if}
                  {#if section.calloutType === 'tip'}💡{/if}
                  {#if section.calloutType === 'warning'}⚠️{/if}
                </span>
                <div class="callout-content">
                  <p>{section.content}</p>
                </div>
              </div>

            {:else if section.type === 'visual-reference'}
              <div class="visual-reference">
                <span class="reference-icon">👁️</span>
                <div class="reference-content">
                  <p>{section.content}</p>
                </div>
              </div>
            {/if}
          {/each}

          {#if chapter.visualCues?.highlightElements && chapter.visualCues.highlightElements.length > 0}
            <div class="visual-cues">
              <p class="cue-prompt">
                <strong>Try this:</strong> Look at the visualization above and find: {chapter.visualCues.highlightElements.join(', ')}
              </p>
            </div>
          {/if}
        </section>
      {/each}
    {/if}
  </article>
  {/if}
</div>

<style>
  .tutorial-article {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    height: 100%;
  }

  /* Table of Contents - Inline Layout */
  .tutorial-toc {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    padding: 1rem;
  }

  .toc-header {
    margin-bottom: 0.75rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid #e2e8f0;
  }

  .toc-header h4 {
    margin: 0;
    font-size: 0.85rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: #64748b;
    font-weight: 600;
  }

  .chapter-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .chapter-list li button {
    display: flex;
    align-items: baseline;
    gap: 0.5rem;
    width: 100%;
    padding: 0.5rem 0.75rem;
    background: transparent;
    border: none;
    border-radius: 6px;
    text-align: left;
    cursor: pointer;
    transition: all 0.1s ease;
    font-size: 0.875rem;
  }

  .chapter-list li button:hover {
    background: #f1f5f9;
  }

  .chapter-list li.active button {
    background: #dbeafe;
    color: #1e40af;
    font-weight: 500;
  }

  .chapter-number {
    font-weight: 600;
    color: #3b82f6;
    flex-shrink: 0;
    min-width: 1.5rem;
  }

  .chapter-title {
    color: #334155;
    line-height: 1.4;
  }

  /* Tutorial Content */
  .tutorial-content {
    flex: 1;
    overflow-y: auto;
    padding: 2rem 1rem 1rem 1rem;
  }

  .empty-state {
    padding: 2rem;
    text-align: center;
    color: #64748b;
  }

  .tutorial-chapter {
    margin-bottom: 2.5rem;
    scroll-margin-top: 2rem;
  }

  .chapter-heading {
    margin: 0 0 1.25rem 0;
    font-size: 1.25rem;
    font-weight: 700;
    color: #0f172a;
    display: flex;
    align-items: baseline;
    gap: 0.5rem;
    line-height: 1.3;
  }

  .chapter-heading .chapter-number {
    color: #3b82f6;
    font-size: 1.1rem;
  }

  /* Content Types */
  .chapter-text {
    margin: 0 0 1rem 0;
    line-height: 1.7;
    color: #334155;
    font-size: 0.95rem;
  }

  .chapter-highlight {
    margin: 1.25rem 0;
    padding: 1rem 1.25rem;
    background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
    border-left: 4px solid #f59e0b;
    border-radius: 6px;
    font-weight: 500;
    color: #92400e;
    line-height: 1.6;
    font-size: 0.95rem;
  }

  .chapter-code {
    margin: 1.25rem 0;
    padding: 1.25rem;
    background: #1e293b;
    border-radius: 8px;
    overflow-x: auto;
  }

  .chapter-code pre {
    margin: 0;
    font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
  }

  .chapter-code code {
    color: #e2e8f0;
    font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
    font-size: 0.875rem;
    line-height: 1.6;
  }

  .chapter-callout {
    margin: 1.25rem 0;
    padding: 1rem 1.25rem;
    border-radius: 8px;
    display: flex;
    gap: 1rem;
    align-items: flex-start;
  }

  .chapter-callout[data-type='info'] {
    background: #eff6ff;
    border-left: 4px solid #3b82f6;
  }

  .chapter-callout[data-type='tip'] {
    background: #f0fdf4;
    border-left: 4px solid #10b981;
  }

  .chapter-callout[data-type='warning'] {
    background: #fef2f2;
    border-left: 4px solid #ef4444;
  }

  .callout-icon {
    font-size: 1.25rem;
    flex-shrink: 0;
    line-height: 1;
  }

  .callout-content {
    flex: 1;
  }

  .chapter-callout p {
    margin: 0;
    color: #334155;
    line-height: 1.6;
    font-size: 0.9rem;
  }

  .visual-reference {
    margin: 1.25rem 0;
    padding: 1rem 1.25rem;
    background: #f8fafc;
    border: 1px dashed #cbd5e1;
    border-radius: 8px;
    display: flex;
    gap: 1rem;
    align-items: flex-start;
  }

  .reference-icon {
    font-size: 1.25rem;
    flex-shrink: 0;
    line-height: 1;
  }

  .reference-content {
    flex: 1;
  }

  .visual-reference p {
    margin: 0;
    color: #475569;
    font-style: italic;
    line-height: 1.6;
    font-size: 0.9rem;
  }

  .visual-cues {
    margin: 1.5rem 0;
    padding: 1rem 1.25rem;
    background: #f0f9ff;
    border-radius: 8px;
    border: 1px solid #bae6fd;
  }

  .cue-prompt {
    margin: 0;
    color: #0c4a6e;
    font-size: 0.9rem;
    line-height: 1.6;
  }

  .cue-prompt strong {
    color: #0369a1;
    font-weight: 600;
  }

  /* Responsive */
  @media (max-width: 640px) {
    .tutorial-article {
      gap: 1rem;
    }

    .tutorial-toc {
      width: 85%;
      max-width: 300px;
    }

    .chapter-list li button {
      font-size: 0.8125rem;
      padding: 0.4rem 0.5rem;
    }

    .chapter-heading {
      font-size: 1.1rem;
    }

    .chapter-text {
      font-size: 0.9rem;
    }
  }
</style>
