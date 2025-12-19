<script lang="ts">
  import Phase1App from './phase1/Phase1App.svelte';
  import Phase2App from './phase2/Phase2App.svelte';
  import Phase3App from './phase3/Phase3App.svelte';
  import PhaseSelector from './PhaseSelector.svelte';

  type PhaseType = 'phase1' | 'phase2' | 'phase3';

  // Parse ?phase=1,2,3 from URL, returns null if invalid
  function parsePhaseFromUrl(): PhaseType | null {
    if (typeof window === 'undefined') return null;
    const params = new URLSearchParams(window.location.search);
    const phase = params.get('phase');
    if (phase === '1') return 'phase1';
    if (phase === '2') return 'phase2';
    if (phase === '3') return 'phase3';
    return null;
  }

  // Convert 'phase1' → 1, 'phase2' → 2, 'phase3' → 3
  function phaseToNumber(phase: PhaseType): number {
    return parseInt(phase.replace('phase', ''), 10);
  }

  // Update URL without page reload
  function updateUrlParam(phase: PhaseType): void {
    if (typeof window === 'undefined') return;
    const url = new URL(window.location.href);
    url.searchParams.set('phase', phaseToNumber(phase).toString());
    window.history.replaceState({}, '', url);
  }

  // Initialize phase with priority: URL param → localStorage → default
  let currentPhase = $state<PhaseType>(
    (() => {
      // Priority 1: URL query parameter
      const urlPhase = parsePhaseFromUrl();
      if (urlPhase !== null) return urlPhase;

      // Priority 2: localStorage
      if (typeof localStorage !== 'undefined') {
        const stored = localStorage.getItem('currentPhase') as PhaseType | null;
        if (stored === 'phase1' || stored === 'phase2' || stored === 'phase3') {
          return stored;
        }
      }

      // Priority 3: Default
      return 'phase1';
    })()
  );

  function handlePhaseChange(phase: PhaseType) {
    currentPhase = phase;

    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('currentPhase', phase);
    }

    updateUrlParam(phase);
  }

  // Keep URL in sync when currentPhase changes
  $effect(() => {
    updateUrlParam(currentPhase);
  });

  // Handle browser back/forward navigation
  $effect(() => {
    if (typeof window === 'undefined') return;

    function handlePopState() {
      const urlPhase = parsePhaseFromUrl();
      if (urlPhase !== null && urlPhase !== currentPhase) {
        currentPhase = urlPhase;
        if (typeof localStorage !== 'undefined') {
          localStorage.setItem('currentPhase', urlPhase);
        }
      }
    }

    window.addEventListener('popstate', handlePopState);
    return () => window.removeEventListener('popstate', handlePopState);
  });
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
