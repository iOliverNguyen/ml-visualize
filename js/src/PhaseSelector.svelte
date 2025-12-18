<script lang="ts">
  interface Props {
    currentPhase: 'phase1' | 'phase2';
    onchange: (phase: 'phase1' | 'phase2') => void;
  }

  let { currentPhase, onchange }: Props = $props();
</script>

<div class="phase-selector">
  <button
    class="phase-tab"
    class:active={currentPhase === 'phase1'}
    onclick={() => onchange('phase1')}
    type="button"
  >
    <span class="phase-emoji">📈</span>
    <div class="phase-info">
      <span class="phase-name">Phase 1</span>
      <span class="phase-desc">Scalar GD: y = w·x</span>
    </div>
  </button>

  <button
    class="phase-tab"
    class:active={currentPhase === 'phase2'}
    onclick={() => onchange('phase2')}
    type="button"
  >
    <span class="phase-emoji">📊</span>
    <div class="phase-info">
      <span class="phase-name">Phase 2</span>
      <span class="phase-desc">2-Parameter: y = w₁·x₁ + w₂·x₂</span>
    </div>
  </button>
</div>

<style>
  .phase-selector {
    display: flex;
    gap: 1rem;
    padding: 1rem 2rem;
    background: #f8fafc;
    border-bottom: 2px solid #e2e8f0;
  }

  .phase-tab {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 1rem 1.5rem;
    background: white;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.15s ease;
    flex: 1;
    max-width: 400px;
  }

  .phase-tab:hover {
    border-color: #cbd5e1;
    transform: translateY(-2px);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }

  .phase-tab.active {
    background: #3b82f6;
    border-color: #3b82f6;
    box-shadow: 0 4px 6px -1px rgba(59, 130, 246, 0.3);
  }

  .phase-tab.active:hover {
    transform: translateY(0);
  }

  .phase-emoji {
    font-size: 2rem;
    line-height: 1;
  }

  .phase-tab.active .phase-emoji {
    filter: grayscale(0);
  }

  .phase-tab:not(.active) .phase-emoji {
    filter: grayscale(0.5);
    opacity: 0.7;
  }

  .phase-info {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.25rem;
  }

  .phase-name {
    font-weight: 700;
    font-size: 1.1rem;
    color: #0f172a;
  }

  .phase-tab.active .phase-name {
    color: white;
  }

  .phase-desc {
    font-size: 0.85rem;
    color: #64748b;
    font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
  }

  .phase-tab.active .phase-desc {
    color: #dbeafe;
  }

  @media (max-width: 768px) {
    .phase-selector {
      padding: 0.75rem 1rem;
      flex-direction: column;
    }

    .phase-tab {
      max-width: none;
    }

    .phase-emoji {
      font-size: 1.5rem;
    }

    .phase-name {
      font-size: 1rem;
    }

    .phase-desc {
      font-size: 0.75rem;
    }
  }
</style>
