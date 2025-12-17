<script lang="ts">
  import type { Snapshot } from './types'
  import FormulaOverlay from './FormulaOverlay.svelte'

  interface Props {
    snapshot: Snapshot
  }

  let { snapshot }: Props = $props()

  type Stage = 'forward' | 'loss' | 'gradient' | 'update'
  let activeStage = $state<Stage>('forward')

  // Selected point for detailed formula view
  let selectedPointIndex = $state<number | undefined>(undefined)

  function selectPoint(index: number) {
    selectedPointIndex = selectedPointIndex === index ? undefined : index
  }
</script>

<div class="step-breakdown">
  <div class="stage-tabs">
    <button
      class="stage-tab"
      class:active={activeStage === 'forward'}
      onclick={() => (activeStage = 'forward')}
    >
      1. Forward Pass
    </button>
    <button
      class="stage-tab"
      class:active={activeStage === 'loss'}
      onclick={() => (activeStage = 'loss')}
    >
      2. Compute Loss
    </button>
    <button
      class="stage-tab"
      class:active={activeStage === 'gradient'}
      onclick={() => (activeStage = 'gradient')}
    >
      3. Compute Gradient
    </button>
    <button
      class="stage-tab"
      class:active={activeStage === 'update'}
      onclick={() => (activeStage = 'update')}
    >
      4. Update Parameter
    </button>
  </div>

  <div class="stage-content">
    {#if activeStage === 'forward'}
      <div class="stage-panel">
        <h3>Forward Pass: Compute Predictions</h3>
        <p class="stage-description">
          For each data point, compute the prediction using: <code>ŷ = w · x</code>
        </p>

        <div class="points-table">
          <table>
            <thead>
              <tr>
                <th>Point</th>
                <th>x</th>
                <th>w</th>
                <th>Calculation</th>
                <th>ŷ (prediction)</th>
                <th>y (true)</th>
              </tr>
            </thead>
            <tbody>
              {#each snapshot.point_details as point, i}
                <tr
                  class:selected={selectedPointIndex === i}
                  onclick={() => selectPoint(i)}
                >
                  <td>{i + 1}</td>
                  <td>{point.x}</td>
                  <td>{snapshot.w.toFixed(4)}</td>
                  <td class="calculation">
                    {snapshot.w.toFixed(4)} × {point.x}
                  </td>
                  <td class="result">{point.y_pred.toFixed(4)}</td>
                  <td>{point.y_true.toFixed(4)}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>

        {#if selectedPointIndex !== undefined}
          <div class="formula-detail">
            <FormulaOverlay
              {snapshot}
              mode="gradient"
              pointIndex={selectedPointIndex}
            />
          </div>
        {/if}
      </div>
    {:else if activeStage === 'loss'}
      <div class="stage-panel">
        <h3>Compute Loss: Measure Error</h3>
        <p class="stage-description">
          For each point, compute squared error: <code>L = (ŷ - y)²</code>
        </p>

        <div class="points-table">
          <table>
            <thead>
              <tr>
                <th>Point</th>
                <th>ŷ (pred)</th>
                <th>y (true)</th>
                <th>Error (ŷ - y)</th>
                <th>Calculation</th>
                <th>Loss</th>
              </tr>
            </thead>
            <tbody>
              {#each snapshot.point_details as point, i}
                {@const error = point.y_pred - point.y_true}
                <tr
                  class:selected={selectedPointIndex === i}
                  onclick={() => selectPoint(i)}
                >
                  <td>{i + 1}</td>
                  <td>{point.y_pred.toFixed(4)}</td>
                  <td>{point.y_true.toFixed(4)}</td>
                  <td>{error.toFixed(4)}</td>
                  <td class="calculation">({error.toFixed(4)})²</td>
                  <td class="result">{point.point_loss.toFixed(4)}</td>
                </tr>
              {/each}
              <tr class="aggregate-row">
                <td colspan="5">Average Loss</td>
                <td class="result aggregate">{snapshot.loss.toFixed(4)}</td>
              </tr>
            </tbody>
          </table>
        </div>

        {#if selectedPointIndex !== undefined}
          <div class="formula-detail">
            <FormulaOverlay
              {snapshot}
              mode="loss"
              pointIndex={selectedPointIndex}
            />
          </div>
        {:else}
          <div class="formula-detail">
            <FormulaOverlay {snapshot} mode="loss" />
          </div>
        {/if}
      </div>
    {:else if activeStage === 'gradient'}
      <div class="stage-panel">
        <h3>Compute Gradient: Direction to Improve</h3>
        <p class="stage-description">
          For each point, compute gradient: <code>∂L/∂w = 2(ŷ - y) · x</code>
        </p>

        <div class="points-table">
          <table>
            <thead>
              <tr>
                <th>Point</th>
                <th>Error (ŷ - y)</th>
                <th>x</th>
                <th>Calculation</th>
                <th>Gradient</th>
              </tr>
            </thead>
            <tbody>
              {#each snapshot.point_details as point, i}
                {@const error = point.y_pred - point.y_true}
                <tr
                  class:selected={selectedPointIndex === i}
                  onclick={() => selectPoint(i)}
                >
                  <td>{i + 1}</td>
                  <td>{error.toFixed(4)}</td>
                  <td>{point.x}</td>
                  <td class="calculation">
                    2 × {error.toFixed(4)} × {point.x}
                  </td>
                  <td class="result">{point.point_grad.toFixed(4)}</td>
                </tr>
              {/each}
              <tr class="aggregate-row">
                <td colspan="4">Average Gradient</td>
                <td class="result aggregate">{snapshot.grad_w.toFixed(4)}</td>
              </tr>
            </tbody>
          </table>
        </div>

        {#if selectedPointIndex !== undefined}
          <div class="formula-detail">
            <FormulaOverlay
              {snapshot}
              mode="gradient"
              pointIndex={selectedPointIndex}
            />
          </div>
        {:else}
          <div class="formula-detail">
            <FormulaOverlay {snapshot} mode="gradient" />
          </div>
        {/if}
      </div>
    {:else if activeStage === 'update'}
      <div class="stage-panel">
        <h3>Update Parameter: Gradient Descent Step</h3>
        <p class="stage-description">
          Update parameter using gradient descent: <code
            >w_new = w_old - lr · ∂L/∂w</code
          >
        </p>

        <div class="update-breakdown">
          <div class="update-row">
            <span class="update-label">Current parameter:</span>
            <span class="update-value"
              >w = {snapshot.update_components.w_old.toFixed(4)}</span
            >
          </div>
          <div class="update-row">
            <span class="update-label">Learning rate:</span>
            <span class="update-value"
              >lr = {snapshot.update_components.lr.toFixed(4)}</span
            >
          </div>
          <div class="update-row">
            <span class="update-label">Gradient:</span>
            <span class="update-value"
              >∂L/∂w = {snapshot.update_components.grad_w.toFixed(4)}</span
            >
          </div>
          <div class="update-row">
            <span class="update-label">Step size:</span>
            <span class="update-value"
              >Δw = -lr · ∂L/∂w = {snapshot.update_components.delta_w.toFixed(
                4
              )}</span
            >
          </div>
          <div class="update-row highlight">
            <span class="update-label">New parameter:</span>
            <span class="update-value"
              >w_new = {snapshot.update_components.w_old.toFixed(4)} + {snapshot.update_components.delta_w.toFixed(
                4
              )} = {snapshot.update_components.w_new.toFixed(4)}</span
            >
          </div>
        </div>

        <div class="formula-detail">
          <FormulaOverlay {snapshot} mode="update" />
        </div>

        <div class="update-insight">
          {#if snapshot.grad_w < 0}
            <p>
              ✓ Gradient is <strong>negative</strong> → parameter will
              <strong>increase</strong> (moving towards optimum)
            </p>
          {:else if snapshot.grad_w > 0}
            <p>
              ✓ Gradient is <strong>positive</strong> → parameter will
              <strong>decrease</strong> (moving towards optimum)
            </p>
          {:else}
            <p>✓ Gradient is <strong>zero</strong> → at optimal point!</p>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .step-breakdown {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }

  .stage-tabs {
    display: flex;
    background: #f1f5f9;
    border-bottom: 2px solid #e2e8f0;
  }

  .stage-tab {
    flex: 1;
    padding: 1rem;
    border: none;
    background: transparent;
    cursor: pointer;
    font-size: 0.95rem;
    font-weight: 500;
    color: #64748b;
    transition: all 0.2s;
    border-bottom: 3px solid transparent;
  }

  .stage-tab:hover {
    background: #e2e8f0;
    color: #334155;
  }

  .stage-tab.active {
    background: white;
    color: #2563eb;
    border-bottom-color: #2563eb;
  }

  .stage-content {
    padding: 1.5rem;
  }

  .stage-panel h3 {
    margin: 0 0 0.5rem 0;
    color: #0f172a;
    font-size: 1.2rem;
  }

  .stage-description {
    margin: 0 0 1rem 0;
    color: #64748b;
    font-size: 0.95rem;
  }

  .stage-description code {
    background: #f1f5f9;
    padding: 0.2rem 0.4rem;
    border-radius: 3px;
    font-family: 'SF Mono', monospace;
    color: #2563eb;
  }

  .points-table {
    overflow-x: auto;
    margin-bottom: 1rem;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.9rem;
  }

  thead {
    background: #f8fafc;
    border-bottom: 2px solid #e2e8f0;
  }

  th {
    padding: 0.75rem;
    text-align: left;
    font-weight: 600;
    color: #475569;
    font-size: 0.85rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  tbody tr {
    border-bottom: 1px solid #e2e8f0;
    transition: background 0.15s;
    cursor: pointer;
  }

  tbody tr:hover {
    background: #f8fafc;
  }

  tbody tr.selected {
    background: #dbeafe;
    border-color: #3b82f6;
  }

  td {
    padding: 0.75rem;
    color: #334155;
  }

  td.calculation {
    font-family: 'SF Mono', monospace;
    color: #64748b;
    font-size: 0.85rem;
  }

  td.result {
    font-family: 'SF Mono', monospace;
    font-weight: 600;
    color: #059669;
  }

  .aggregate-row {
    background: #f1f5f9;
    font-weight: 600;
    cursor: default;
  }

  .aggregate-row:hover {
    background: #f1f5f9;
  }

  .aggregate-row td {
    padding: 1rem 0.75rem;
  }

  .aggregate-row td.result {
    font-size: 1.1rem;
    color: #2563eb;
  }

  .formula-detail {
    margin-top: 1rem;
  }

  .update-breakdown {
    background: #f8fafc;
    padding: 1.5rem;
    border-radius: 6px;
    margin-bottom: 1rem;
  }

  .update-row {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem 0;
    border-bottom: 1px solid #e2e8f0;
  }

  .update-row:last-child {
    border-bottom: none;
  }

  .update-row.highlight {
    margin-top: 1rem;
    padding: 1rem;
    background: #dbeafe;
    border-radius: 4px;
    border-bottom: none;
  }

  .update-label {
    color: #64748b;
    font-weight: 500;
  }

  .update-value {
    font-family: 'SF Mono', monospace;
    color: #0f172a;
    font-weight: 600;
  }

  .update-row.highlight .update-value {
    color: #2563eb;
    font-size: 1.05rem;
  }

  .update-insight {
    margin-top: 1rem;
    padding: 1rem;
    background: #f0fdf4;
    border-left: 4px solid #22c55e;
    border-radius: 4px;
  }

  .update-insight p {
    margin: 0;
    color: #166534;
    font-size: 0.95rem;
  }

  .update-insight strong {
    color: #15803d;
  }
</style>
