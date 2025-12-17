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
    const target = e.target as HTMLInputElement
    onsliderChange(parseInt(target.value))
  }
</script>

<div class="controls">
  <div class="button-row">
    <button
      onclick={onreset}
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
    >
      {#if playing}
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M2 2h4v12H2V2zm8 0h4v12h-4V2z"/>
        </svg>
      {:else}
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
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
  </div>

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

  .reset-btn {
    background: #fff5f5;
    border-color: #fed7d7;
    color: #c53030;
  }

  .reset-btn:hover:not(:disabled) {
    background: #fed7d7;
  }

  .step-display {
    margin-left: auto;
    font-weight: bold;
    font-family: monospace;
    font-size: 1.1rem;
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
