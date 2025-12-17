<script lang="ts">
  import type { DataPoint } from './types';

  interface Props {
    onDatasetChange: (snapshots: any[]) => void;
    loading?: boolean;
  }

  let { onDatasetChange, loading = false }: Props = $props();

  // Tab state
  let activeTab = $state<'random' | 'function'>('random');

  // Random generation config
  let numPoints = $state(20);
  let xMin = $state(0);
  let xMax = $state(20);
  let trueSlope = $state(2.0);
  let noiseLevel = $state(0.5);
  let seed = $state(42);

  // Training config
  let wInit = $state(0.0);
  let learningRate = $state(0.001);  // Reduced from 0.01 to prevent gradient explosion
  let trainingSteps = $state(100);

  // JS function input
  let jsCode = $state(`// Generate data points
// Return array of {x, y} objects
function generateData(n) {
  const points = [];
  for (let i = 0; i < n; i++) {
    const x = i + 1;
    const y = 2 * x + (Math.random() * 2 - 1);
    points.push({ x, y });
  }
  return points;
}

return generateData(${numPoints});`);

  let jsFunctionPoints = $state(20);

  // Error state
  let error = $state('');
  let isGenerating = $state(false);

  function switchTab(tab: 'random' | 'function') {
    activeTab = tab;
    error = '';
  }

  async function generateRandomData() {
    error = '';
    isGenerating = true;

    try {
      const response = await fetch('http://localhost:8080/api/dataset/random', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          data_config: {
            num_points: numPoints,
            x_min: xMin,
            x_max: xMax,
            true_slope: trueSlope,
            noise_level: noiseLevel,
            seed: seed || Date.now(),
          },
          training_config: {
            w_init: wInit,
            lr: learningRate,
            steps: trainingSteps,
          },
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Failed to generate data');
      }

      const snapshots = await response.json();
      onDatasetChange(snapshots);
    } catch (e) {
      if (e instanceof Error && e.message.includes('fetch')) {
        error = 'Server not running. Start with: go run . --server';
      } else {
        error = e instanceof Error ? e.message : 'Failed to generate data';
      }
    } finally {
      isGenerating = false;
    }
  }

  async function generateFromFunction() {
    error = '';
    isGenerating = true;

    try {
      // Execute user's JS function
      const wrappedCode = `
        (function() {
          ${jsCode}
        })()
      `;

      const dataPoints = eval(wrappedCode);

      // Validate structure
      if (!Array.isArray(dataPoints)) {
        throw new Error('Function must return an array');
      }

      if (dataPoints.length === 0) {
        throw new Error('Function must return at least one data point');
      }

      // Convert to backend format
      const formattedData = dataPoints.map((point: any, i: number) => {
        if (typeof point !== 'object' || point === null) {
          throw new Error(`Point ${i} must be an object with x and y properties`);
        }
        if (typeof point.x !== 'number' || typeof point.y !== 'number') {
          throw new Error(`Point ${i} must have numeric x and y values`);
        }
        return {
          X: point.x,
          YTrue: point.y,
        };
      });

      // Send to backend for training
      const response = await fetch('http://localhost:8080/api/dataset/custom', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          data: formattedData,
          config: {
            w_init: wInit,
            lr: learningRate,
            steps: trainingSteps,
          },
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Failed to train with custom data');
      }

      const snapshots = await response.json();
      onDatasetChange(snapshots);
    } catch (e) {
      if (e instanceof Error && e.message.includes('fetch')) {
        error = 'Server not running. Start with: go run . --server';
      } else if (e instanceof SyntaxError) {
        error = 'Syntax error in JavaScript code: ' + e.message;
      } else {
        error = e instanceof Error ? e.message : 'Failed to generate data from function';
      }
    } finally {
      isGenerating = false;
    }
  }
</script>

<div class="input-panel">
  <div class="tabs">
    <button
      class="tab"
      class:active={activeTab === 'random'}
      onclick={() => switchTab('random')}
      type="button"
    >
      🎲 Random Data
    </button>
    <button
      class="tab"
      class:active={activeTab === 'function'}
      onclick={() => switchTab('function')}
      type="button"
    >
      ⚡ JS Function
    </button>
  </div>

  {#if error}
    <div class="error-message">
      {error}
    </div>
  {/if}

  {#if activeTab === 'random'}
    <div class="tab-content">
      <h3 class="section-title">Generate Random Dataset</h3>
      <p class="section-description">
        Create synthetic data following y = slope*x + noise
      </p>

      <div class="form-section">
        <h4>Data Configuration</h4>

        <div class="form-group">
          <label>
            <span class="label-text">Number of points</span>
            <span class="label-value">{numPoints}</span>
          </label>
          <input
            type="range"
            min="5"
            max="100"
            step="5"
            bind:value={numPoints}
            class="slider"
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>
              <span class="label-text">X min</span>
            </label>
            <input
              type="number"
              bind:value={xMin}
              step="0.1"
              class="number-input"
            />
          </div>

          <div class="form-group">
            <label>
              <span class="label-text">X max</span>
            </label>
            <input
              type="number"
              bind:value={xMax}
              step="0.1"
              class="number-input"
            />
          </div>
        </div>

        <div class="form-group">
          <label>
            <span class="label-text">True slope</span>
            <span class="label-value">{trueSlope.toFixed(2)}</span>
          </label>
          <input
            type="range"
            min="0"
            max="5"
            step="0.1"
            bind:value={trueSlope}
            class="slider"
          />
        </div>

        <div class="form-group">
          <label>
            <span class="label-text">Noise level</span>
            <span class="label-value">{noiseLevel.toFixed(2)}</span>
          </label>
          <input
            type="range"
            min="0"
            max="3"
            step="0.1"
            bind:value={noiseLevel}
            class="slider"
          />
        </div>

        <div class="form-group">
          <label>
            <span class="label-text">Random seed (0 = random)</span>
          </label>
          <input
            type="number"
            bind:value={seed}
            class="number-input"
          />
        </div>
      </div>

      <div class="form-section">
        <h4>Training Configuration</h4>

        <div class="form-row">
          <div class="form-group">
            <label>
              <span class="label-text">Initial w</span>
            </label>
            <input
              type="number"
              bind:value={wInit}
              step="0.1"
              class="number-input"
            />
          </div>

          <div class="form-group">
            <label>
              <span class="label-text">Learning rate</span>
            </label>
            <input
              type="number"
              bind:value={learningRate}
              step="0.001"
              class="number-input"
            />
          </div>

          <div class="form-group">
            <label>
              <span class="label-text">Steps</span>
            </label>
            <input
              type="number"
              bind:value={trainingSteps}
              step="10"
              class="number-input"
            />
          </div>
        </div>
      </div>

      <div class="preview-section">
        <strong>Preview:</strong> y = {trueSlope.toFixed(2)}*x + noise({noiseLevel.toFixed(2)})
        <br />
        {numPoints} points from x={xMin} to x={xMax}
      </div>

      <button
        class="generate-button"
        onclick={generateRandomData}
        disabled={isGenerating || loading}
        type="button"
      >
        {#if isGenerating}
          Generating...
        {:else}
          Generate & Train
        {/if}
      </button>
    </div>
  {:else}
    <div class="tab-content">
      <h3 class="section-title">Define Data with JavaScript</h3>
      <p class="section-description">
        Write a function that returns an array of {"{x, y}"} objects
      </p>

      <div class="form-section">
        <h4>JavaScript Code</h4>
        <textarea
          bind:value={jsCode}
          class="code-editor"
          rows="15"
          spellcheck="false"
        ></textarea>

        <p class="code-hint">
          <strong>Tip:</strong> Your code must return an array like:
          <code>[{"{x: 1, y: 2.1}"}, {"{x: 2, y: 4.2}"}, ...]</code>
        </p>
      </div>

      <div class="form-section">
        <h4>Training Configuration</h4>

        <div class="form-row">
          <div class="form-group">
            <label>
              <span class="label-text">Initial w</span>
            </label>
            <input
              type="number"
              bind:value={wInit}
              step="0.1"
              class="number-input"
            />
          </div>

          <div class="form-group">
            <label>
              <span class="label-text">Learning rate</span>
            </label>
            <input
              type="number"
              bind:value={learningRate}
              step="0.001"
              class="number-input"
            />
          </div>

          <div class="form-group">
            <label>
              <span class="label-text">Steps</span>
            </label>
            <input
              type="number"
              bind:value={trainingSteps}
              step="10"
              class="number-input"
            />
          </div>
        </div>
      </div>

      <button
        class="generate-button"
        onclick={generateFromFunction}
        disabled={isGenerating || loading}
        type="button"
      >
        {#if isGenerating}
          Generating...
        {:else}
          Generate & Train
        {/if}
      </button>
    </div>
  {/if}
</div>

<style>
  .input-panel {
    background: white;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
    overflow: hidden;
    margin-top: 1rem;
  }

  .tabs {
    display: flex;
    border-bottom: 1px solid #e2e8f0;
    background: #f8fafc;
  }

  .tab {
    flex: 1;
    padding: 1rem;
    background: transparent;
    border: none;
    color: #64748b;
    cursor: pointer;
    transition: all 0.15s ease;
    font-weight: 500;
    font-size: 0.95rem;
    position: relative;
  }

  .tab:hover {
    background: #f1f5f9;
    color: #334155;
  }

  .tab.active {
    color: #3b82f6;
    background: white;
  }

  .tab.active::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: #3b82f6;
  }

  .tab-content {
    padding: 1.5rem;
  }

  .section-title {
    margin: 0 0 0.5rem 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: #0f172a;
  }

  .section-description {
    margin: 0 0 1.5rem 0;
    font-size: 0.9rem;
    color: #64748b;
    line-height: 1.5;
  }

  .form-section {
    margin-bottom: 1.5rem;
  }

  .form-section h4 {
    margin: 0 0 1rem 0;
    font-size: 0.95rem;
    font-weight: 600;
    color: #334155;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #e2e8f0;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  .form-group label {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .label-text {
    font-size: 0.875rem;
    font-weight: 500;
    color: #334155;
  }

  .label-value {
    font-size: 0.875rem;
    font-weight: 600;
    color: #3b82f6;
    font-family: 'SF Mono', monospace;
  }

  .form-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
    gap: 1rem;
  }

  .slider {
    width: 100%;
    height: 6px;
    border-radius: 3px;
    background: #e2e8f0;
    outline: none;
    cursor: pointer;
  }

  .slider::-webkit-slider-thumb {
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: #3b82f6;
    cursor: pointer;
  }

  .slider::-moz-range-thumb {
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: #3b82f6;
    cursor: pointer;
    border: none;
  }

  .number-input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    font-size: 0.875rem;
    font-family: 'SF Mono', monospace;
    transition: all 0.15s ease;
  }

  .number-input:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 0;
    border-color: #3b82f6;
  }

  .code-editor {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
    font-size: 0.85rem;
    line-height: 1.5;
    background: #f8fafc;
    color: #0f172a;
    resize: vertical;
    transition: all 0.15s ease;
  }

  .code-editor:focus {
    outline: 2px solid #3b82f6;
    outline-offset: 0;
    border-color: #3b82f6;
    background: white;
  }

  .code-hint {
    margin: 0.75rem 0 0 0;
    font-size: 0.8rem;
    color: #64748b;
    line-height: 1.5;
  }

  .code-hint code {
    background: #f1f5f9;
    padding: 0.2rem 0.4rem;
    border-radius: 3px;
    font-family: 'SF Mono', monospace;
    font-size: 0.75rem;
  }

  .preview-section {
    padding: 1rem;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 6px;
    font-size: 0.9rem;
    color: #334155;
    line-height: 1.6;
    margin-bottom: 1.5rem;
  }

  .generate-button {
    width: 100%;
    padding: 0.875rem;
    background: #10b981;
    color: white;
    border: none;
    border-radius: 8px;
    font-weight: 600;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .generate-button:hover:not(:disabled) {
    background: #059669;
    transform: translateY(-1px);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  }

  .generate-button:disabled {
    background: #cbd5e1;
    cursor: not-allowed;
  }

  .generate-button:focus {
    outline: 2px solid #10b981;
    outline-offset: 2px;
  }

  .error-message {
    padding: 0.875rem;
    background: #fee2e2;
    border: 1px solid #fecaca;
    border-radius: 6px;
    color: #991b1b;
    font-size: 0.875rem;
    margin: 1rem 1.5rem 0 1.5rem;
    line-height: 1.5;
  }

  @media (max-width: 640px) {
    .tab-content {
      padding: 1rem;
    }

    .form-row {
      grid-template-columns: 1fr;
    }
  }
</style>
