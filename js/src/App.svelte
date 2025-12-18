<script lang="ts">
  import Phase1App from './phase1/Phase1App.svelte';
  import PhaseSelector from './PhaseSelector.svelte';

  // Load current phase from localStorage or default to phase1
  let currentPhase = $state<'phase1' | 'phase2'>(
    (typeof localStorage !== 'undefined' ? localStorage.getItem('currentPhase') : null) as 'phase1' | 'phase2' || 'phase1'
  );

  function handlePhaseChange(phase: 'phase1' | 'phase2') {
    currentPhase = phase;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('currentPhase', phase);
    }
  }
</script>

<div class="app-container">
  <PhaseSelector {currentPhase} onchange={handlePhaseChange} />

  {#if currentPhase === 'phase1'}
    <Phase1App />
  {:else if currentPhase === 'phase2'}
    <div class="phase2-placeholder">
      <div class="placeholder-content">
        <h1>🚧 Phase 2: Coming Soon</h1>
        <p>2-Parameter Gradient Descent visualization with:</p>
        <ul>
          <li>Parameter space trajectory plots</li>
          <li>Loss surface contour visualization</li>
          <li>Gradient vector field</li>
          <li>7 pre-computed cases exploring learning rate effects</li>
        </ul>
        <p class="status">Backend complete! Frontend components in development.</p>
      </div>
    </div>
  {/if}
</div>

<style>
  .app-container {
    min-height: 100vh;
    background: white;
  }

  .phase2-placeholder {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: calc(100vh - 120px);
    padding: 2rem;
  }

  .placeholder-content {
    max-width: 600px;
    text-align: center;
    padding: 3rem;
    background: #f8fafc;
    border-radius: 12px;
    border: 2px dashed #cbd5e1;
  }

  .placeholder-content h1 {
    margin: 0 0 1rem 0;
    color: #0f172a;
    font-size: 2rem;
  }

  .placeholder-content p {
    margin: 1rem 0;
    color: #475569;
    font-size: 1.1rem;
    line-height: 1.6;
  }

  .placeholder-content ul {
    text-align: left;
    margin: 1.5rem auto;
    max-width: 400px;
    list-style-position: inside;
  }

  .placeholder-content li {
    margin: 0.5rem 0;
    color: #334155;
  }

  .status {
    margin-top: 2rem;
    padding: 1rem;
    background: #dbeafe;
    border-radius: 8px;
    color: #1e40af;
    font-weight: 600;
  }

  @media (max-width: 768px) {
    .phase2-placeholder {
      min-height: calc(100vh - 180px);
    }

    .placeholder-content {
      padding: 2rem 1.5rem;
    }

    .placeholder-content h1 {
      font-size: 1.5rem;
    }

    .placeholder-content p {
      font-size: 1rem;
    }
  }
</style>
