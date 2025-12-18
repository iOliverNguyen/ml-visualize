<script lang="ts">
  import type { Snapshot } from './types'

  interface Props {
    snapshot: Snapshot
    mode: 'loss' | 'gradient' | 'update'
    pointIndex?: number // Optional: show formula for specific point
  }

  let { snapshot, mode, pointIndex }: Props = $props()

  // Get point details if pointIndex is specified
  let point = $derived(
    pointIndex !== undefined && snapshot.point_details[pointIndex]
      ? snapshot.point_details[pointIndex]
      : null
  )

  // Aggregate formulas (averaged over all points)
  let aggregateFormulas = $derived({
    loss: {
      formula: 'L',
      expansion: `(1/n) Σ (ŷᵢ - yᵢ)²`,
      calculation: `${snapshot.loss.toFixed(4)}`,
      description: 'Mean Squared Error averaged over all points'
    },
    gradient: {
      formula: '∂L/∂w',
      expansion: `(1/n) Σ 2(ŷᵢ - yᵢ) · xᵢ`,
      calculation: `${snapshot.grad_w.toFixed(4)}`,
      description: 'Average gradient across all data points'
    },
    update: {
      formula: 'w_new',
      expansion: `w_old - lr · ∂L/∂w`,
      calculation: `${snapshot.update_components.w_old.toFixed(4)} - ${snapshot.update_components.lr.toFixed(4)} · ${snapshot.update_components.grad_w.toFixed(4)}`,
      result: `${snapshot.update_components.w_new.toFixed(4)}`,
      description: 'Gradient descent parameter update'
    }
  })

  // Per-point formulas (if point is selected)
  let pointFormulas = $derived(
    point
      ? {
          loss: {
            formula: 'Lᵢ',
            expansion: `(ŷᵢ - yᵢ)²`,
            calculation: `(${point.y_pred.toFixed(4)} - ${point.y_true.toFixed(4)})²`,
            result: `${point.point_loss.toFixed(4)}`,
            description: `Loss for point with x=${point.x}`
          },
          gradient: {
            formula: '∂Lᵢ/∂w',
            expansion: `2(ŷᵢ - yᵢ) · xᵢ`,
            calculation: `2(${point.y_pred.toFixed(4)} - ${point.y_true.toFixed(4)}) · ${point.x}`,
            result: `${point.point_grad.toFixed(4)}`,
            description: `Gradient contribution from point with x=${point.x}`
          },
          forward: {
            formula: 'ŷᵢ',
            expansion: `w · xᵢ`,
            calculation: `${snapshot.w.toFixed(4)} · ${point.x}`,
            result: `${point.y_pred.toFixed(4)}`,
            description: `Prediction for x=${point.x}`
          }
        }
      : null
  )

  // Select which formula to display based on mode
  let displayFormula = $derived(() => {
    if (point && pointFormulas) {
      // Show per-point formula if point is selected
      if (mode === 'loss') return pointFormulas.loss
      if (mode === 'gradient') return pointFormulas.gradient
      return null // Update mode doesn't make sense per-point
    } else {
      // Show aggregate formula
      if (mode === 'loss') return aggregateFormulas.loss
      if (mode === 'gradient') return aggregateFormulas.gradient
      if (mode === 'update') return aggregateFormulas.update
      return null
    }
  })

  let formula = $derived(displayFormula())
</script>

<div class="formula-overlay">
  {#if formula}
    <div class="formula-card">
      <div class="formula-header">
        <span class="formula-symbol">{formula.formula}</span>
        <span class="formula-equals">=</span>
        <span class="formula-expansion">{formula.expansion}</span>
      </div>

      <div class="formula-calculation">
        <span class="calculation-label">Substituted:</span>
        <span class="calculation-value">{formula.calculation}</span>
        {#if formula.result}
          <span class="calculation-equals">=</span>
          <span class="calculation-result">{formula.result}</span>
        {/if}
      </div>

      <div class="formula-description">{formula.description}</div>
    </div>
  {/if}
</div>

<style>
  .formula-overlay {
    padding: 1rem;
    background: rgba(255, 255, 255, 0.95);
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    font-family: 'SF Mono', 'Monaco', 'Cascadia Code', 'Courier New', monospace;
  }

  .formula-card {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .formula-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.1rem;
    font-weight: 500;
  }

  .formula-symbol {
    color: #2563eb;
    font-weight: 600;
  }

  .formula-equals {
    color: #64748b;
  }

  .formula-expansion {
    color: #0f172a;
    font-style: italic;
  }

  .formula-calculation {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem;
    background: #f1f5f9;
    border-radius: 4px;
    font-size: 0.95rem;
  }

  .calculation-label {
    color: #64748b;
    font-size: 0.85rem;
  }

  .calculation-value {
    color: #0f172a;
    font-family: 'SF Mono', monospace;
  }

  .calculation-equals {
    color: #64748b;
  }

  .calculation-result {
    color: #059669;
    font-weight: 600;
    font-family: 'SF Mono', monospace;
  }

  .formula-description {
    font-size: 0.85rem;
    color: #64748b;
    font-style: italic;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  }
</style>
