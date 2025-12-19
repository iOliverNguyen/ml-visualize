<script lang="ts">
  import Phase1App from './phase1/Phase1App.svelte';
  import Phase2App from './phase2/Phase2App.svelte';
  import Phase3App from './phase3/Phase3App.svelte';
  import PhaseSelector from './PhaseSelector.svelte';

  // Load current phase from localStorage or default to phase1
  let currentPhase = $state<'phase1' | 'phase2' | 'phase3'>(
    (typeof localStorage !== 'undefined' ? localStorage.getItem('currentPhase') : null) as 'phase1' | 'phase2' | 'phase3' || 'phase1'
  );

  function handlePhaseChange(phase: 'phase1' | 'phase2' | 'phase3') {
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
    <Phase2App />
  {:else if currentPhase === 'phase3'}
    <Phase3App />
  {/if}
</div>

<style>
  .app-container {
    min-height: 100vh;
    background: white;
  }
</style>
