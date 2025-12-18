<script lang="ts">
  import { slide } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import * as educationalState from '../stores/educationalState.svelte';

  interface Props {
    onclose?: () => void;
    onstartTour?: () => void;
  }

  let { onclose, onstartTour }: Props = $props();

  function handleClose() {
    educationalState.state.showIntroPanel = false;
    educationalState.saveToLocalStorage();
    onclose?.();
  }

  function handleStartTour() {
    educationalState.state.tourActive = true;
    educationalState.state.currentTourStep = 0;
    educationalState.state.showIntroPanel = false;
    educationalState.saveToLocalStorage();
    onstartTour?.();
  }
</script>

{#if educationalState.state.showIntroPanel}
  <div
    class="intro-panel"
    transition:slide={{ duration: 300, easing: cubicOut }}
  >
    <div class="intro-content">
      <button
        class="close-button"
        onclick={handleClose}
        aria-label="Close introduction panel"
        type="button"
      >
        <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
          <path
            d="M15 5L5 15M5 5L15 15"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          />
        </svg>
      </button>

      <div class="intro-header">
        <h1 class="intro-title">Welcome to Gradient Descent Visualization</h1>
        <p class="intro-subtitle">
          Watch machine learning in action as an algorithm learns to fit a line to data
        </p>
      </div>

      <div class="intro-body">
        <div class="what-section">
          <h2>What You'll Learn</h2>
          <ul>
            <li>How gradient descent finds optimal parameters automatically</li>
            <li>Why loss functions matter and how they guide learning</li>
            <li>How learning rate controls the training process</li>
            <li>The intuition behind modern AI training</li>
          </ul>
        </div>

        <div class="how-section">
          <h2>How to Use This</h2>
          <div class="steps">
            <div class="step">
              <span class="step-number">1</span>
              <div class="step-content">
                <strong>Watch the training</strong>
                <p>Click play and observe how the algorithm learns</p>
              </div>
            </div>
            <div class="step">
              <span class="step-number">2</span>
              <div class="step-content">
                <strong>Explore interactively</strong>
                <p>Click terms for definitions, hover for insights</p>
              </div>
            </div>
          </div>
        </div>

        <div class="actions">
          <button class="primary-button" onclick={handleStartTour} type="button">
            Start Guided Tour
          </button>
          <button class="secondary-button" onclick={handleClose} type="button">
            Explore on My Own
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .intro-panel {
    position: relative;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 2rem;
    margin-bottom: 1.5rem;
    border-radius: 12px;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  }

  .intro-content {
    max-width: 900px;
    margin: 0 auto;
    position: relative;
  }

  .close-button {
    position: absolute;
    top: 0;
    right: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    padding: 0;
    background: rgba(255, 255, 255, 0.2);
    border: none;
    color: white;
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.15s ease;
  }

  .close-button:hover {
    background: rgba(255, 255, 255, 0.3);
  }

  .close-button:focus {
    outline: 2px solid white;
    outline-offset: 2px;
  }

  .intro-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .intro-title {
    margin: 0 0 0.5rem 0;
    font-size: 2rem;
    font-weight: 700;
    color: white;
  }

  .intro-subtitle {
    margin: 0;
    font-size: 1.1rem;
    color: rgba(255, 255, 255, 0.9);
    font-weight: 400;
  }

  .intro-body {
    background: white;
    color: #0f172a;
    padding: 1.5rem;
    border-radius: 8px;
  }

  .what-section,
  .how-section {
    margin-bottom: 1.5rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #e2e8f0;
  }

  .what-section:last-child,
  .how-section:last-child {
    border-bottom: none;
    margin-bottom: 0;
    padding-bottom: 0;
  }

  h2 {
    margin: 0 0 1rem 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: #0f172a;
  }

  .what-section ul {
    margin: 0;
    padding-left: 1.5rem;
    list-style-type: none;
  }

  .what-section li {
    position: relative;
    margin: 0.5rem 0;
    padding-left: 1.5rem;
    line-height: 1.6;
    color: #334155;
  }

  .what-section li::before {
    content: '✓';
    position: absolute;
    left: 0;
    color: #10b981;
    font-weight: 700;
  }

  .steps {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .step {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
  }

  .step-number {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    flex-shrink: 0;
    background: #667eea;
    color: white;
    border-radius: 50%;
    font-weight: 700;
    font-size: 0.9rem;
  }

  .step-content {
    flex: 1;
  }

  .step-content strong {
    display: block;
    margin-bottom: 0.25rem;
    color: #0f172a;
  }

  .step-content p {
    margin: 0;
    color: #64748b;
    font-size: 0.9rem;
    line-height: 1.5;
  }

  .actions {
    display: flex;
    gap: 1rem;
    justify-content: center;
    margin-top: 1.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid #e2e8f0;
  }

  .primary-button,
  .secondary-button {
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    font-weight: 600;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.15s ease;
    border: none;
  }

  .primary-button {
    background: #667eea;
    color: white;
  }

  .primary-button:hover {
    background: #5568d3;
    transform: translateY(-1px);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }

  .primary-button:focus {
    outline: 2px solid #667eea;
    outline-offset: 2px;
  }

  .secondary-button {
    background: white;
    color: #334155;
    border: 2px solid #cbd5e1;
  }

  .secondary-button:hover {
    background: #f8fafc;
    border-color: #94a3b8;
  }

  .secondary-button:focus {
    outline: 2px solid #cbd5e1;
    outline-offset: 2px;
  }

  @media (max-width: 640px) {
    .intro-panel {
      padding: 1.5rem;
    }

    .intro-title {
      font-size: 1.5rem;
    }

    .intro-subtitle {
      font-size: 1rem;
    }

    .actions {
      flex-direction: column;
    }
  }
</style>
