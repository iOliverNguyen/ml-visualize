<script lang="ts">
  interface Props {
    currentStep: number;
    totalSteps: number;
    playing: boolean;
    onstepBack: () => void;
    onstepForward: () => void;
    onplayPause: () => void;
    onsliderChange: (step: number) => void;
    onreset: () => void;
    onjumpToFinal: () => void;
  }

  let {
    currentStep,
    totalSteps,
    playing,
    onstepBack,
    onstepForward,
    onplayPause,
    onsliderChange,
    onreset,
    onjumpToFinal
  }: Props = $props();

  function handleSliderInput(e: Event) {
    const target = e.target as HTMLInputElement;
    onsliderChange(parseInt(target.value));
  }

  let showInfo = $state(false);

  // Detect if user is on Mac
  const isMac = typeof navigator !== 'undefined' &&
    (navigator.platform.toUpperCase().indexOf('MAC') >= 0 ||
     navigator.userAgent.toUpperCase().indexOf('MAC') >= 0);

  function handleKeydown(e: KeyboardEvent) {
    // Prevent default only for our shortcuts
    switch(e.key) {
      case 'ArrowLeft':
        e.preventDefault();
        if (e.metaKey) {
          // CMD+Left: Jump to start
          onreset();
        } else if (e.altKey) {
          // ALT+Left: Step back 5
          const targetStep = Math.max(0, currentStep - 5);
          onsliderChange(targetStep);
        } else {
          // Left: Step backward
          if (currentStep > 0) onstepBack();
        }
        break;
      case 'ArrowRight':
        e.preventDefault();
        if (e.metaKey) {
          // CMD+Right: Jump to end
          onjumpToFinal();
        } else if (e.altKey) {
          // ALT+Right: Step forward 5
          const targetStep = Math.min(totalSteps - 1, currentStep + 5);
          onsliderChange(targetStep);
        } else {
          // Right: Step forward
          if (currentStep < totalSteps - 1) onstepForward();
        }
        break;
      case ' ':
        e.preventDefault();
        onplayPause();
        break;
      case 'Home':
        e.preventDefault();
        onreset();
        break;
      case 'End':
        e.preventDefault();
        onjumpToFinal();
        break;
    }
  }

  function toggleInfo() {
    showInfo = !showInfo;
  }

  // Add event listener on mount, cleanup on unmount
  $effect(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  });
</script>

<div class="controls">
  <div class="button-row">
    <button
      onclick={onreset}
      disabled={currentStep === 0}
      title="Reset to first step"
      class="reset-btn"
    >
      <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <path d="M1 2v12h2V2H1zm3 6l8 6V2l-8 6z"/>
      </svg>
    </button>

    <button
      onclick={onstepBack}
      disabled={currentStep === 0}
      title="Step backward"
    >
      <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <path d="M10 2L4 8l6 6V2z"/>
      </svg>
    </button>

    <button
      onclick={onplayPause}
      title={playing ? 'Pause autoplay' : 'Play autoplay'}
      class="play-button"
    >
      {#if playing}
        <svg width="20" height="20" viewBox="0 0 16 16" fill="currentColor">
          <path d="M2 2h4v12H2V2zm8 0h4v12h-4V2z"/>
        </svg>
      {:else}
        <svg width="20" height="20" viewBox="0 0 16 16" fill="currentColor">
          <path d="M3 2v12l10-6L3 2z"/>
        </svg>
      {/if}
    </button>

    <button
      onclick={onstepForward}
      disabled={currentStep === totalSteps - 1}
      title="Step forward"
    >
      <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <path d="M6 2l6 6-6 6V2z"/>
      </svg>
    </button>

    <button
      onclick={onjumpToFinal}
      disabled={currentStep === totalSteps - 1}
      title="Jump to final step"
    >
      <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <path d="M4 2v12l8-6-8-6zm9 0v12h2V2h-2z"/>
      </svg>
    </button>

    <span class="step-display">
      Step: {currentStep} / {totalSteps - 1}
    </span>

    <button
      onclick={toggleInfo}
      title="Show keyboard shortcuts"
      class="info-btn"
    >
      <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <circle cx="8" cy="8" r="7" fill="none" stroke="currentColor" stroke-width="1.5"/>
        <text x="8" y="12" text-anchor="middle" font-size="11" font-weight="bold">i</text>
      </svg>
    </button>
  </div>

  {#if showInfo}
    <div class="info-panel">
      <h3>Keyboard Shortcuts</h3>
      <div class="shortcuts">
        <div class="shortcut">
          <kbd>← / →</kbd>
          <span>Step backward / forward</span>
        </div>
        <div class="shortcut">
          <kbd>Alt + ← / →</kbd>
          <span>Step 5 steps backward / forward</span>
        </div>
        <div class="shortcut">
          <kbd>Space</kbd>
          <span>Play / Pause</span>
        </div>
        <div class="shortcut">
          <kbd>{isMac ? 'Cmd + ← / →' : 'Home / End'}</kbd>
          <span>Reset to start / Jump to final</span>
        </div>
      </div>
      <p class="hint">Click the (i) button again to hide this panel</p>
    </div>
  {/if}

  <div class="slider-row">
    <input
      type="range"
      min="0"
      max={totalSteps - 1}
      value={currentStep}
      oninput={handleSliderInput}
      class="timeline-slider"
      title="Scrub through timeline"
    />
  </div>
</div>

<style>
  .controls {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin: 2rem 0;
    padding: 1rem;
    background: #f5f5f5;
    border-radius: 8px;
    border: 1px solid #e0e0e0;
  }

  .button-row {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .slider-row {
    display: flex;
    width: 100%;
    height: 32px;
    align-items: center;
  }

  button {
    padding: 0.5rem 1rem;
    font-size: 1rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    background: white;
    cursor: pointer;
    transition: background 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  button svg {
    display: block;
    width: 16px;
    height: 16px;
  }

  button:hover:not(:disabled) {
    background: #e0e0e0;
  }

  button:active:not(:disabled) {
    background: #d0d0d0;
  }

  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .play-button {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: 2px solid #667eea;
    padding: 0.75rem 1.25rem;
    box-shadow: 0 4px 6px -1px rgba(102, 126, 234, 0.3);
    transform: scale(1.1);
  }

  .play-button:hover:not(:disabled) {
    background: linear-gradient(135deg, #5568d3 0%, #653a8a 100%);
    box-shadow: 0 6px 8px -1px rgba(102, 126, 234, 0.4);
    transform: scale(1.15);
  }

  .play-button:active:not(:disabled) {
    transform: scale(1.08);
  }

  .play-button svg {
    width: 20px;
    height: 20px;
  }

  .step-display {
    margin-left: auto;
    font-weight: bold;
    font-family: monospace;
    font-size: 1.1rem;
  }

  .info-btn {
    padding: 0.4rem 0.6rem;
    margin-left: 0.5rem;
    background: #e3f2fd;
    border-color: #90caf9;
    color: #1976d2;
  }

  .info-btn:hover:not(:disabled) {
    background: #bbdefb;
  }

  .info-panel {
    background: #e8f5e9;
    border: 1px solid #81c784;
    border-radius: 6px;
    padding: 1rem;
    margin-top: 0.5rem;
  }

  .info-panel h3 {
    margin: 0 0 0.75rem 0;
    font-size: 1rem;
    color: #2e7d32;
  }

  .shortcuts {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 0.75rem 1.5rem;
  }

  .shortcut {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .shortcut kbd {
    min-width: 100px;
    padding: 0.25rem 0.5rem;
    background: white;
    border: 1px solid #c5e1a5;
    border-radius: 4px;
    font-family: monospace;
    font-size: 0.85rem;
    text-align: center;
    box-shadow: 0 1px 2px rgba(0,0,0,0.1);
    white-space: nowrap;
    flex-shrink: 0;
  }

  .shortcut span {
    color: #1b5e20;
    font-size: 0.9rem;
  }

  @media (max-width: 768px) {
    .shortcuts {
      grid-template-columns: 1fr;
    }
  }

  .info-panel .hint {
    margin: 0.75rem 0 0 0;
    font-size: 0.85rem;
    color: #558b2f;
    font-style: italic;
  }

  .timeline-slider {
    width: 100%;
    height: 8px;
    border-radius: 4px;
    background: #ddd;
    outline: none;
    cursor: pointer;
    appearance: none;
  }

  .timeline-slider::-webkit-slider-thumb {
    appearance: none;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #2563eb;
    cursor: pointer;
    transition: background 0.2s;
    transform: translateY(-6px);
  }

  .timeline-slider::-webkit-slider-thumb:hover {
    background: #1d4ed8;
  }

  .timeline-slider::-moz-range-thumb {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #2563eb;
    cursor: pointer;
    border: none;
    transition: background 0.2s;
    transform: translateY(-6px);
  }

  .timeline-slider::-moz-range-thumb:hover {
    background: #1d4ed8;
  }

  /* Track styling */
  .timeline-slider::-webkit-slider-runnable-track {
    width: 100%;
    height: 8px;
    border-radius: 4px;
    background: #ddd;
  }

  .timeline-slider::-moz-range-track {
    width: 100%;
    height: 8px;
    border-radius: 4px;
    background: #ddd;
  }
</style>
