<script lang="ts">
  import type { NeuronSnapshot } from './types';

  interface Props {
    snapshot: NeuronSnapshot;
  }

  let { snapshot }: Props = $props();

  let activeTab = $state<'w1' | 'w2' | 'b'>('w1');

  // Get chain rule components for the selected parameter
  let chainComponents = $derived.by(() => {
    const components = snapshot?.chain_rule_breakdown?.components;
    if (!components || components.length === 0) return null;
    return components.find(c => c.param_name === activeTab) ?? components[0];
  });

  // Color coding based on magnitude
  function getMagnitudeColor(value: number): string {
    const absValue = Math.abs(value);
    if (absValue > 0.5) {
      return '#16a34a'; // green - healthy
    } else if (absValue > 0.01) {
      return '#eab308'; // yellow - moderate
    } else {
      return '#dc2626'; // red - vanishing!
    }
  }

  function getMagnitudeLabel(value: number): string {
    const absValue = Math.abs(value);
    if (absValue > 0.5) {
      return 'Healthy gradient flow';
    } else if (absValue > 0.01) {
      return 'Moderate gradient';
    } else {
      return 'Vanishing gradient!';
    }
  }

  function formatNumber(n: number): string {
    return n.toFixed(6);
  }

  function getParamLabel(param: string): string {
    switch (param) {
      case 'w1':
        return 'Weight w₁';
      case 'w2':
        return 'Weight w₂';
      case 'b':
        return 'Bias b';
      default:
        return param;
    }
  }
</script>

<div class="chain-rule-breakdown">
  <h2>Chain Rule Breakdown: Backpropagation</h2>
  <p class="subtitle">
    How gradients flow backward through the neuron: ∂L/∂{activeTab} = ∂L/∂a × ∂a/∂z × ∂z/∂{activeTab}
  </p>

  <!-- Tabs -->
  <div class="tabs">
    <button
      class="tab"
      class:active={activeTab === 'w1'}
      onclick={() => activeTab = 'w1'}
    >
      Weight w₁
    </button>
    <button
      class="tab"
      class:active={activeTab === 'w2'}
      onclick={() => activeTab = 'w2'}
    >
      Weight w₂
    </button>
    <button
      class="tab"
      class:active={activeTab === 'b'}
      onclick={() => activeTab = 'b'}
    >
      Bias b
    </button>
  </div>

  {#if chainComponents}
    <!-- Flow diagram -->
    <div class="flow-diagram">
      <!-- Step 1: ∂L/∂a -->
      <div class="flow-card" style="--card-color: {getMagnitudeColor(chainComponents?.dL_da ?? 0)}">
        <div class="card-header">
          <div class="card-title">∂L/∂a</div>
          <div class="card-subtitle">Loss gradient</div>
        </div>
        <div class="card-value">{formatNumber(chainComponents?.dL_da ?? 0)}</div>
        <div class="card-description">
          How much does loss change with activation?
          <br />
          Formula: 2(a - y_true)
        </div>
      </div>

      <!-- Multiplication symbol -->
      <div class="operator">×</div>

      <!-- Step 2: ∂a/∂z -->
      <div class="flow-card" style="--card-color: {getMagnitudeColor(chainComponents?.da_dz ?? 0)}">
        <div class="card-header">
          <div class="card-title">∂a/∂z</div>
          <div class="card-subtitle">Local derivative</div>
        </div>
        <div class="card-value">{formatNumber(chainComponents?.da_dz ?? 0)}</div>
        <div class="card-description">
          {#if snapshot?.activation === 'sigmoid'}
            σ'(z) = σ(z) × (1 - σ(z))
          {:else if snapshot?.activation === 'relu'}
            ReLU'(z) = 1 if z > 0, else 0
          {:else if snapshot?.activation === 'tanh'}
            tanh'(z) = 1 - tanh²(z)
          {/if}
        </div>
        {#if Math.abs(chainComponents?.da_dz ?? 0) < 0.01}
          <div class="card-warning">
            ⚠️ Saturated! This is where gradients vanish.
          </div>
        {/if}
      </div>

      <!-- Multiplication symbol -->
      <div class="operator">×</div>

      <!-- Step 3: ∂z/∂param -->
      <div class="flow-card" style="--card-color: {getMagnitudeColor(chainComponents?.dz_dparam ?? 0)}">
        <div class="card-header">
          <div class="card-title">∂z/∂{activeTab}</div>
          <div class="card-subtitle">Input contribution</div>
        </div>
        <div class="card-value">{formatNumber(chainComponents?.dz_dparam ?? 0)}</div>
        <div class="card-description">
          {#if activeTab === 'b'}
            For bias: ∂z/∂b = 1
            <br />
            (bias adds directly to z)
          {:else}
            For weights: ∂z/∂{activeTab} = x{activeTab === 'w1' ? '₁' : '₂'}
            <br />
            (input feature value)
          {/if}
        </div>
      </div>

      <!-- Equals symbol -->
      <div class="operator result">=</div>

      <!-- Final result: ∂L/∂param -->
      <div class="flow-card result" style="--card-color: {getMagnitudeColor(chainComponents?.dL_dparam ?? 0)}">
        <div class="card-header">
          <div class="card-title">∂L/∂{activeTab}</div>
          <div class="card-subtitle">Final gradient</div>
        </div>
        <div class="card-value final">{formatNumber(chainComponents?.dL_dparam ?? 0)}</div>
        <div class="card-description">
          This is what we use to update {getParamLabel(activeTab)}!
          <br />
          Update: {activeTab} ← {activeTab} - lr × ∂L/∂{activeTab}
        </div>
        <div class="magnitude-indicator" style="--indicator-color: {getMagnitudeColor(chainComponents?.dL_dparam ?? 0)}">
          {getMagnitudeLabel(chainComponents?.dL_dparam ?? 0)}
        </div>
      </div>
    </div>

    <!-- Detailed explanation -->
    <div class="explanation-panel">
      <h3>Understanding the Chain Rule</h3>
      <div class="explanation-content">
        <div class="explanation-section">
          <strong>Forward pass:</strong> z = w₁·x₁ + w₂·x₂ + b → a = {snapshot?.activation ?? 'sigmoid'}(z) → loss = (a - y)²
        </div>
        <div class="explanation-section">
          <strong>Backward pass:</strong> We compute gradients by "chaining" derivatives from loss back to parameters.
        </div>
        <div class="explanation-section">
          <strong>Why this matters:</strong>
          {#if Math.abs(chainComponents?.da_dz ?? 0) < 0.01}
            <span class="highlight-danger">
              The middle term (∂a/∂z = {formatNumber(chainComponents?.da_dz ?? 0)}) is very small!
              This is <strong>saturation</strong> - when the activation function's derivative approaches zero,
              gradients "vanish" and learning becomes extremely slow.
            </span>
          {:else}
            <span class="highlight-success">
              All components have reasonable magnitudes. Gradients flow well through the neuron,
              enabling effective learning.
            </span>
          {/if}
        </div>
      </div>
    </div>

    <!-- Current values summary -->
    <div class="summary-panel">
      <h3>Current State (Step {snapshot?.step ?? 0})</h3>
      <div class="summary-grid">
        <div class="summary-item">
          <span class="summary-label">z (pre-activation):</span>
          <span class="summary-value">{formatNumber(snapshot?.z ?? 0)}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">a (post-activation):</span>
          <span class="summary-value">{formatNumber(snapshot?.a ?? 0)}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">σ'(z) magnitude:</span>
          <span class="summary-value">{formatNumber(Math.abs(snapshot?.local_derivative ?? 0))}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">Gradient ∂L/∂{activeTab}:</span>
          <span class="summary-value">{formatNumber(chainComponents?.dL_dparam ?? 0)}</span>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .chain-rule-breakdown {
    background: white;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    padding: 1.5rem;
    margin-top: 2rem;
  }

  h2 {
    margin: 0 0 0.5rem 0;
    font-size: 1.25rem;
    color: #1e293b;
    text-align: center;
  }

  .subtitle {
    text-align: center;
    color: #64748b;
    font-size: 0.95rem;
    margin: 0 0 1.5rem 0;
    font-style: italic;
  }

  .tabs {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
    border-bottom: 2px solid #e2e8f0;
  }

  .tab {
    flex: 1;
    padding: 0.75rem 1rem;
    border: none;
    background: transparent;
    color: #64748b;
    font-weight: 500;
    cursor: pointer;
    border-bottom: 3px solid transparent;
    transition: all 0.2s;
  }

  .tab:hover {
    color: #1e293b;
    background: #f8fafc;
  }

  .tab.active {
    color: #3b82f6;
    border-bottom-color: #3b82f6;
    background: #eff6ff;
  }

  .flow-diagram {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin: 2rem 0;
    padding: 1.5rem;
    background: #f8fafc;
    border-radius: 8px;
    overflow-x: auto;
  }

  .flow-card {
    flex: 1;
    min-width: 180px;
    background: white;
    border: 3px solid var(--card-color, #94a3b8);
    border-radius: 8px;
    padding: 1rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .flow-card.result {
    flex: 1.2;
    background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  }

  .card-header {
    margin-bottom: 0.75rem;
  }

  .card-title {
    font-size: 1.1rem;
    font-weight: bold;
    color: var(--card-color, #1e293b);
    font-family: 'Courier New', monospace;
  }

  .card-subtitle {
    font-size: 0.75rem;
    color: #64748b;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-top: 0.25rem;
  }

  .card-value {
    font-size: 1.5rem;
    font-weight: bold;
    font-family: 'Courier New', monospace;
    color: var(--card-color, #1e293b);
    margin: 0.5rem 0;
    text-align: center;
  }

  .card-value.final {
    font-size: 1.8rem;
  }

  .card-description {
    font-size: 0.8rem;
    color: #475569;
    line-height: 1.4;
    margin-top: 0.5rem;
  }

  .card-warning {
    margin-top: 0.75rem;
    padding: 0.5rem;
    background: #fee2e2;
    border: 1px solid #dc2626;
    border-radius: 4px;
    color: #991b1b;
    font-size: 0.75rem;
    font-weight: 600;
  }

  .magnitude-indicator {
    margin-top: 0.75rem;
    padding: 0.5rem;
    background: var(--indicator-color, #e2e8f0);
    color: white;
    border-radius: 4px;
    text-align: center;
    font-weight: 600;
    font-size: 0.85rem;
  }

  .operator {
    font-size: 2rem;
    font-weight: bold;
    color: #64748b;
    flex-shrink: 0;
  }

  .operator.result {
    color: #3b82f6;
  }

  .explanation-panel {
    margin: 2rem 0;
    padding: 1.5rem;
    background: #eff6ff;
    border: 2px solid #3b82f6;
    border-radius: 8px;
  }

  .explanation-panel h3 {
    margin: 0 0 1rem 0;
    color: #1e40af;
    font-size: 1.1rem;
  }

  .explanation-content {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .explanation-section {
    line-height: 1.6;
    color: #1e293b;
    font-size: 0.95rem;
  }

  .highlight-danger {
    display: block;
    margin-top: 0.5rem;
    padding: 1rem;
    background: #fee2e2;
    border-left: 4px solid #dc2626;
    color: #991b1b;
  }

  .highlight-success {
    display: block;
    margin-top: 0.5rem;
    padding: 1rem;
    background: #dcfce7;
    border-left: 4px solid #16a34a;
    color: #166534;
  }

  .summary-panel {
    margin: 2rem 0 0 0;
    padding: 1.5rem;
    background: #fafaf9;
    border: 2px solid #e7e5e4;
    border-radius: 8px;
  }

  .summary-panel h3 {
    margin: 0 0 1rem 0;
    color: #44403c;
    font-size: 1rem;
  }

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
  }

  .summary-item {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem;
    background: white;
    border-radius: 4px;
  }

  .summary-label {
    font-weight: 500;
    color: #64748b;
    font-size: 0.9rem;
  }

  .summary-value {
    font-family: 'Courier New', monospace;
    font-weight: bold;
    color: #1e293b;
    font-size: 0.9rem;
  }

  @media (max-width: 1024px) {
    .flow-diagram {
      flex-direction: column;
    }

    .operator {
      transform: rotate(90deg);
      margin: 0.5rem 0;
    }

    .flow-card {
      width: 100%;
    }
  }

  @media (max-width: 768px) {
    .tabs {
      flex-direction: column;
    }

    .tab {
      text-align: left;
    }

    .summary-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
