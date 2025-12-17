<script lang="ts">
  import { slide } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';

  interface Props {
    title: string;
    subtitle?: string;
    id: string;
    defaultExpanded?: boolean;
    level?: 'beginner' | 'intermediate' | 'advanced';
    icon?: string;
    onexpandChange?: (expanded: boolean) => void;
  }

  let {
    title,
    subtitle,
    id,
    defaultExpanded = false,
    level,
    icon,
    onexpandChange,
    children
  }: Props = $props();

  let expanded = $state(defaultExpanded);

  function toggle() {
    expanded = !expanded;
    onexpandChange?.(expanded);
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      toggle();
    }
  }

  // Level indicator colors
  const levelColors = {
    beginner: '#dcfce7',
    intermediate: '#fef3c7',
    advanced: '#fee2e2'
  };

  const levelTextColors = {
    beginner: '#166534',
    intermediate: '#92400e',
    advanced: '#991b1b'
  };
</script>

<div class="expandable-section" data-section-id={id}>
  <button
    class="expandable-header"
    onclick={toggle}
    onkeydown={handleKeyDown}
    aria-expanded={expanded}
    aria-controls="{id}-content"
  >
    <div class="header-content">
      {#if icon}
        <span class="icon">{icon}</span>
      {/if}
      <div class="header-text">
        <div class="title-row">
          <h3 class="title">{title}</h3>
          {#if level}
            <span
              class="level-badge"
              style="background: {levelColors[level]}; color: {levelTextColors[level]}"
            >
              {level}
            </span>
          {/if}
        </div>
        {#if subtitle}
          <p class="subtitle">{subtitle}</p>
        {/if}
      </div>
    </div>
    <svg
      class="chevron"
      class:rotated={expanded}
      width="20"
      height="20"
      viewBox="0 0 20 20"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
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
      class="expandable-content"
      id="{id}-content"
      transition:slide={{ duration: 200, easing: cubicOut }}
    >
      {@render children?.()}
    </div>
  {/if}
</div>

<style>
  .expandable-section {
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    overflow: hidden;
    margin-bottom: 0.75rem;
    background: white;
  }

  .expandable-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    padding: 1rem 1.25rem;
    background: #f8fafc;
    border: none;
    cursor: pointer;
    transition: background 0.15s ease;
    text-align: left;
  }

  .expandable-header:hover {
    background: #f1f5f9;
  }

  .expandable-header:focus {
    outline: 2px solid #3b82f6;
    outline-offset: -2px;
  }

  .header-content {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    flex: 1;
  }

  .icon {
    font-size: 1.5rem;
    line-height: 1;
    flex-shrink: 0;
  }

  .header-text {
    flex: 1;
  }

  .title-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .title {
    margin: 0;
    font-size: 1.05rem;
    font-weight: 600;
    color: #0f172a;
  }

  .subtitle {
    margin: 0.25rem 0 0 0;
    font-size: 0.9rem;
    color: #64748b;
    line-height: 1.4;
  }

  .level-badge {
    display: inline-flex;
    align-items: center;
    padding: 0.15rem 0.5rem;
    border-radius: 4px;
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .chevron {
    flex-shrink: 0;
    color: #64748b;
    transition: transform 0.2s ease;
  }

  .chevron.rotated {
    transform: rotate(180deg);
  }

  .expandable-content {
    padding: 1.25rem;
    border-top: 1px solid #e2e8f0;
    background: white;
  }

  /* Allow nested content to have spacing */
  .expandable-content :global(> *:first-child) {
    margin-top: 0;
  }

  .expandable-content :global(> *:last-child) {
    margin-bottom: 0;
  }

  .expandable-content :global(p) {
    line-height: 1.6;
    color: #334155;
    margin: 0.5rem 0;
  }

  .expandable-content :global(code) {
    background: #f1f5f9;
    padding: 0.2rem 0.4rem;
    border-radius: 3px;
    font-family: 'SF Mono', monospace;
    font-size: 0.9em;
    color: #2563eb;
  }

  .expandable-content :global(ul) {
    margin: 0.75rem 0;
    padding-left: 1.5rem;
  }

  .expandable-content :global(li) {
    margin: 0.5rem 0;
    line-height: 1.6;
    color: #334155;
  }
</style>
