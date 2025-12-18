<script lang="ts">
  import type { LayoutMode } from './types';

  interface Props {
    layoutMode: LayoutMode;
    onlayoutChange: (mode: LayoutMode) => void;
  }

  let { layoutMode, onlayoutChange }: Props = $props();
</script>

<div class="layout-selector">
  <span class="label">Layout:</span>
  <div class="button-group">
    <button
      class:active={layoutMode === '1col'}
      onclick={() => onlayoutChange('1col')}
      title="Single column layout"
    >
      1 Col
    </button>
    <button
      class:active={layoutMode === '2col'}
      onclick={() => onlayoutChange('2col')}
      title="Two column layout"
    >
      2 Col
    </button>
    <button
      class:active={layoutMode === '3col'}
      onclick={() => onlayoutChange('3col')}
      title="Three column layout"
    >
      3 Col
    </button>
    <button
      class:active={layoutMode === 'auto'}
      onclick={() => onlayoutChange('auto')}
      title="Auto-responsive layout"
    >
      Auto
    </button>
  </div>
</div>

<style>
  .layout-selector {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin: 1rem 0;
    padding: 0.75rem 1rem;
    background: #f9fafb;
    border-radius: 8px;
    border: 1px solid #e5e7eb;
  }

  .label {
    font-weight: 600;
    font-size: 0.95rem;
    color: #374151;
  }

  .button-group {
    display: flex;
    gap: 0;
    border-radius: 6px;
    overflow: hidden;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  button {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    border: 1px solid #d1d5db;
    border-right: none;
    background: white;
    cursor: pointer;
    transition: all 0.2s;
    font-weight: 500;
    color: #6b7280;
  }

  button:last-child {
    border-right: 1px solid #d1d5db;
  }

  button:hover:not(.active) {
    background: #f3f4f6;
    color: #374151;
  }

  button.active {
    background: #2563eb;
    color: white;
    border-color: #2563eb;
  }

  button.active + button {
    border-left-color: #2563eb;
  }

  @media (max-width: 640px) {
    .layout-selector {
      flex-direction: column;
      align-items: stretch;
      gap: 0.5rem;
    }

    .button-group {
      width: 100%;
    }

    button {
      flex: 1;
    }
  }
</style>
