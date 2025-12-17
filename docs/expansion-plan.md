# ML Visualization: Comprehensive Expansion Plan

## Vision

Transform the current Phase 1 scalar gradient descent visualization into a complete pedagogical platform covering all 4 phases of ML learning, with emphasis on **mechanistic clarity**, **inspectable values**, and **visceral understanding**.

---

## Current State (Phase 1 - Scalar GD: y = w*x)

**Implemented:**
- Go engine: explicit gradients, 100-step training, 10 hardcoded datapoints
- Svelte 5 frontend: 3 SVG visualizations (Loss, Param, Gradient arrow)
- Timeline controls: play/pause/scrub with tweened animations
- Snapshot schema: `{step, w, grad_w, loss}`

**Missing:**
- Pedagogical overlays (formulas, annotations, step-by-step breakdown)
- Interactive inspection (click points to see values)
- Case library (multiple scenarios to explore)
- Dataset integration (generate/import data)
- Live training (run calculations on demand)
- Comparative view (side-by-side experiments)

---

## Implementation Strategy

### Priority 1: Deep Phase 1 Enhancements
Make scalar gradient descent pedagogically complete before expanding to multi-parameter.

### Priority 2: Phases 2-4 Roadmap
Design consistent architectural patterns that scale across increasing complexity.

---

## PHASE 1 ENHANCEMENTS

### 1. Enhanced Snapshot Schema

**Extend Go snapshot to capture per-point breakdown:**

```go
// go/snapshots.go (MODIFY)
type Snapshot struct {
    Step  int     `json:"step"`
    W     float64 `json:"w"`
    GradW float64 `json:"grad_w"`
    Loss  float64 `json:"loss"`

    // NEW: Per-point details for inspection
    PointDetails []PointSnapshot `json:"point_details"`

    // NEW: Update breakdown for pedagogy
    UpdateComponents UpdateDetails `json:"update_components"`
}

type PointSnapshot struct {
    X         float64 `json:"x"`
    YTrue     float64 `json:"y_true"`
    YPred     float64 `json:"y_pred"`
    PointLoss float64 `json:"point_loss"`  // (y_pred - y_true)^2
    PointGrad float64 `json:"point_grad"`  // 2(y_pred - y_true)*x
}

type UpdateDetails struct {
    WOld   float64 `json:"w_old"`
    LR     float64 `json:"lr"`
    GradW  float64 `json:"grad_w"`
    DeltaW float64 `json:"delta_w"`  // -lr * grad_w
    WNew   float64 `json:"w_new"`
}
```

**Update TypeScript types:**
```typescript
// js/src/types.ts (MODIFY)
export interface Snapshot {
  step: number;
  w: number;
  grad_w: number;
  loss: number;
  point_details: PointSnapshot[];
  update_components: UpdateDetails;
}
// ... add PointSnapshot, UpdateDetails interfaces
```

### 2. Pedagogical Components

**NEW FILES:**

**`js/src/FormulaOverlay.svelte`**
- Shows formulas with live value substitution
- Three modes: 'loss', 'gradient', 'update'
- Example: `L = (2.1*1.0 - 2.0)² = 0.01`

**`js/src/StepBreakdown.svelte`**
- Tabbed interface: Forward → Loss → Gradient → Update
- Shows per-point calculations in each stage
- Aggregates to final average values

**`js/src/AnnotationTooltip.svelte`**
- Hover tooltips with contextual explanations
- "Gradient is negative → w increases" when hovering arrow

**`js/src/ComparativeView.svelte`**
- Side-by-side training runs with different hyperparams
- Mini versions of Loss/Param plots
- Final metrics comparison

### 3. Case Library

**Structure:**
```
js/public/cases/
  manifest.json          # Index of all cases
  perfect-linear/
    run.json            # Complete TrainingRun
  noisy-data/
    run.json
  lr-too-high/
    run.json
  lr-too-low/
    run.json
  bad-init/
    run.json
  scale-small/
    run.json
  scale-large/
    run.json
  mae-loss/
    run.json
  huber-loss/
    run.json
```

**Manifest schema:**
```json
{
  "cases": [
    {
      "id": "perfect-linear",
      "category": "convergence",
      "name": "Perfect Linear Data",
      "description": "Data exactly on y=2x. Shows clean convergence.",
      "difficulty": "beginner",
      "tags": ["convergence", "ideal"],
      "runFile": "/cases/perfect-linear/run.json"
    }
  ],
  "categories": {
    "convergence": "Basic Convergence Patterns",
    "learning-rate": "Learning Rate Effects",
    "initialization": "Initialization Strategies",
    "scaling": "Data Scaling Issues",
    "loss-functions": "Different Loss Functions"
  }
}
```

**NEW FILE:** `js/src/CaseLibrary.svelte`
- Grid view of cases with filters (category, search)
- Click to load case into visualization
- Dispatches 'case-selected' event

### 4. Dataset Integration

**Go Backend:**

**NEW FILES:**
- `go/generator.go` - Synthetic dataset generator with configurable noise/scale
- `go/import.go` - CSV import functionality
- `go/config.go` - TrainingConfig, Dataset, DatasetMeta types

**Frontend:**

**NEW FILES:**
- `js/src/DatasetGenerator.svelte` - UI for synthetic generation (slope, points, noise, x-range)
- `js/src/CSVImport.svelte` - File upload for CSV datasets
- `js/src/DatasetViz.svelte` - Scatter plot showing data + fitted line

### 5. Live Training

**Recommendation: Server-side with SSE streaming**

**Go Backend:**

**NEW FILE:** `go/streaming.go`
```go
func RunTrainingStreaming(config TrainingConfig, dataset Dataset, w http.ResponseWriter) {
    // Server-sent events
    // Send snapshot after each step
    // Client updates visualization in real-time
}
```

**Frontend:**

**NEW FILE:** `js/src/LiveTraining.svelte`
```svelte
<script>
  // 1. Choose/generate dataset
  // 2. Configure hyperparams (w_init, lr, steps)
  // 3. Start training (EventSource connection)
  // 4. Stream snapshots to live visualization
</script>
```

**Alternative (simpler):** Implement training in TypeScript for Phase 1 only
- `js/src/training.ts` with scalarForward(), scalarLoss(), scalarGradW()
- Fully client-side, no server required
- Trade-off: code duplication vs architectural consistency

### 6. Interactive Enhancements

**MODIFY existing plot components:**

**`js/src/LossPlot.svelte`**
- Add clickable data points
- Click to jump to that step
- Hover to show tooltip with values

**`js/src/ParamPlot.svelte`**
- Same interactivity as LossPlot
- Add "mini" mode prop for comparative view

**`js/src/GradientViz.svelte`**
- Hover on arrow for annotation
- Show formula overlay on click

### 7. App Structure

**MODIFY:** `js/src/App.svelte`

Add mode switching:
```svelte
<nav>
  <button onclick={() => mode = 'case-library'}>Case Library</button>
  <button onclick={() => mode = 'visualization'}>Visualization</button>
  <button onclick={() => mode = 'live-training'}>Live Training</button>
</nav>

{#if mode === 'case-library'}
  <CaseLibrary />
{:else if mode === 'visualization'}
  <Controls />
  <Plots />
  <StepBreakdown />
  <FormulaOverlay />
{:else if mode === 'live-training'}
  <LiveTraining />
{/if}
```

---

## PHASES 2-4 ARCHITECTURAL ROADMAP

### Phase 2: 2-Parameter Geometry (y = w1*x1 + w2*x2)

**Key Insight:** Gradient descent is vector movement in parameter space. Learning rate controls step geometry.

**Snapshot Schema:**
```typescript
interface LinearSnapshot {
  step: number;
  params: { w1: number; w2: number };
  grads: { grad_w1: number; grad_w2: number };
  loss: number;
  gradient_magnitude: number;  // sqrt(grad_w1^2 + grad_w2^2)
  gradient_direction: number;  // atan2(grad_w2, grad_w1)
}
```

**New Visualizations:**
- `ParameterSpace2D.svelte` - Trajectory plot showing (w1, w2) path
- `GradientField.svelte` - Vector field showing gradient directions
- `LossContour.svelte` - Heatmap of loss surface

**Case Library:** Different LRs, anisotropic surfaces, saddle points

### Phase 3: Single Neuron (z = w·x + b, a = σ(z))

**Key Insight:** Make activation functions impossible to misunderstand. Show saturation, gradient collapse.

**Snapshot Schema:**
```typescript
interface NeuronSnapshot {
  step: number;
  params: { w: number[]; b: number };
  grads: { grad_w: number[]; grad_b: number };
  z: number;           // pre-activation
  a: number;           // post-activation
  dL_dz: number;       // gradient at z
  dL_da: number;       // gradient at a
  activation: "sigmoid" | "relu" | "tanh";
  in_saturation_zone: boolean;
  local_derivative: number;  // σ'(z)
  loss: number;
}
```

**New Visualizations:**
- `ActivationCurve.svelte` - Plot σ(z) with current point highlighted
- `SaturationIndicator.svelte` - Visual warning for saturation
- `ChainRuleViz.svelte` - Show dL/dw = dL/da * da/dz * dz/dw

**Case Library:** Vanishing gradients with sigmoid, dying ReLU, tanh centered

### Phase 4: Backprop as Signal Flow

**Key Insight:** Backprop is graph traversal, not algebra. Make signal flow visible.

**Snapshot Schema:**
```typescript
interface GraphSnapshot {
  step: number;
  graph: {
    nodes: GraphNode[];  // {id, type, label}
    edges: GraphEdge[];  // {from, to}
  };
  node_values: {
    [node_id: string]: {
      forward_value: number;
      gradient: number;
      local_derivative: number;
    };
  };
  loss: number;
}
```

**New Visualizations:**
- `ComputationGraph.svelte` - Node-edge diagram
- `ForwardPassAnimation.svelte` - Values flow left-to-right
- `BackwardPassAnimation.svelte` - Gradients flow right-to-left
- `NodeInspector.svelte` - Click node to see local derivative

**Case Library:** Single neuron as graph, 2-layer network, skip connections

### Reusable Patterns Across Phases

**Extract shared components:**
1. `TimeSeriesPlot.svelte` - Generic line chart (used for Loss in all phases)
2. `Controls.svelte` - Timeline control (shared across all phases)
3. `MetricsPanel.svelte` - Base for phase-specific metrics
4. `FormulaOverlay.svelte` - Math formulas with value substitution
5. `ComparisonView.svelte` - Side-by-side scenario comparison

**File structure:**
```
js/src/
  shared/           # Reusable components
    TimeSeriesPlot.svelte
    Controls.svelte
    FormulaOverlay.svelte
    types.ts
  phase1/           # Scalar GD
    App.svelte
    LossPlot.svelte
    ...
  phase2/           # 2-param linear
    App.svelte
    ParameterSpace2D.svelte
    ...
  phase3/           # Single neuron
    App.svelte
    ActivationCurve.svelte
    ...
  phase4/           # Backprop graph
    App.svelte
    ComputationGraph.svelte
    ...
```

---

## IMPLEMENTATION SEQUENCE

### Stage 0: Foundation (1-2 days)
**Goal:** Enhanced snapshots + core infrastructure

1. **Go Backend:**
   - Extend Snapshot struct with PointDetails and UpdateComponents
   - Modify training loop to compute per-point values
   - Add TrainingConfig, Dataset, DatasetMeta types

2. **TypeScript:**
   - Update types.ts with new interfaces
   - Verify existing visualizations still work

**Files:**
- `go/snapshots.go` (modify)
- `go/main.go` (modify)
- `go/config.go` (new)
- `js/src/types.ts` (modify)

### Stage 1: Pedagogical Layer (2-3 days)
**Goal:** Make learning process explicit and inspectable

3. Build `FormulaOverlay.svelte` with 3 modes
4. Build `StepBreakdown.svelte` with 4-stage tabs
5. Build `AnnotationTooltip.svelte`
6. Add interactivity to LossPlot/ParamPlot (click/hover)
7. Integrate pedagogical components into App.svelte

**Files:**
- `js/src/FormulaOverlay.svelte` (new)
- `js/src/StepBreakdown.svelte` (new)
- `js/src/AnnotationTooltip.svelte` (new)
- `js/src/LossPlot.svelte` (modify)
- `js/src/ParamPlot.svelte` (modify)
- `js/src/App.svelte` (modify)

### Stage 2: Case Library (2-3 days)
**Goal:** Curated examples for exploration

8. Generate 9 pre-computed training runs using Go
9. Create manifest.json with case metadata
10. Build `CaseLibrary.svelte` selector UI
11. Build `ComparativeView.svelte` for side-by-side
12. Add mode switching to App.svelte

**Files:**
- `js/public/cases/manifest.json` (new)
- `js/public/cases/{case-id}/run.json` (new × 9)
- `js/src/CaseLibrary.svelte` (new)
- `js/src/ComparativeView.svelte` (new)
- `js/src/App.svelte` (modify)

### Stage 3: Dataset Integration (2-3 days)
**Goal:** Generate and import custom datasets

13. Implement synthetic generator (Go + UI)
14. Implement CSV import (Go + UI)
15. Build DatasetViz.svelte (scatter plot)
16. Add dataset selection to app flow

**Files:**
- `go/generator.go` (new)
- `go/import.go` (new)
- `go/server.go` (modify - add endpoints)
- `js/src/DatasetGenerator.svelte` (new)
- `js/src/CSVImport.svelte` (new)
- `js/src/DatasetViz.svelte` (new)

### Stage 4: Live Training (2-3 days)
**Goal:** Interactive training on demand

17. Implement SSE streaming endpoint (Go)
18. Build `LiveTraining.svelte` with config UI
19. Wire EventSource to real-time visualization
20. Add training controls (pause/stop/restart)

**Files:**
- `go/streaming.go` (new)
- `go/server.go` (modify)
- `js/src/LiveTraining.svelte` (new)

**Alternative:** Implement `js/src/training.ts` for client-side training (simpler, but code duplication)

### Stage 5: Phases 2-4 (8-12 weeks)

**Phase 2 (3-4 weeks):**
21. Extract shared components (TimeSeriesPlot, Controls)
22. Implement go/phase2_linear (vector math)
23. Build phase2 visualizations (ParameterSpace2D, GradientField, LossContour)
24. Generate case library (5-7 cases)

**Phase 3 (3-4 weeks):**
25. Implement go/phase3_neuron (activations with explicit derivatives)
26. Build phase3 visualizations (ActivationCurve, ChainRuleViz)
27. Generate case library (sigmoid/ReLU/tanh comparisons)

**Phase 4 (2-4 weeks):**
28. Implement go/phase4_graph (computation graph with backprop)
29. Build phase4 visualizations (ComputationGraph, animation)
30. Generate case library (single neuron, 2-layer, skip connections)

---

## KEY DESIGN DECISIONS

### 1. Live Training: Server-side SSE (Recommended)
- **Pro:** Go remains single source of truth, scales to Phase 2+
- **Con:** Requires running server
- **Alternative:** TypeScript implementation for Phase 1 only (simpler, but duplicates logic)

### 2. Formula Rendering: Unicode → LaTeX
- **MVP:** Unicode math symbols (ŷ, ∂, ∇)
- **Enhancement:** KaTeX for Phase 3+ (chain rule, partial derivatives)

### 3. Loss Surface (Phase 2): Pre-computed Grid
- Go computes loss at grid points, serializes to JSON
- Keeps computation in Go, consistent with philosophy

### 4. Snapshot Size Management
- Phase 1-3: Inline everything in single JSON
- Phase 4: Separate graph.json + snapshots.json (graph structure is static)

### 5. Multi-Phase UI: Separate Apps Initially
- Each phase has distinct pedagogical goals
- Easier to develop/test independently
- Can unify later if patterns emerge

---

## CRITICAL FILES

### Phase 1 Enhancements

**Go Backend:**
- `/Users/i/ws/olv/ml-viz/go/snapshots.go` - Extend Snapshot schema
- `/Users/i/ws/olv/ml-viz/go/main.go` - Modify training loop
- `/Users/i/ws/olv/ml-viz/go/config.go` - NEW: Training config types
- `/Users/i/ws/olv/ml-viz/go/generator.go` - NEW: Synthetic data
- `/Users/i/ws/olv/ml-viz/go/import.go` - NEW: CSV import
- `/Users/i/ws/olv/ml-viz/go/streaming.go` - NEW: SSE for live training
- `/Users/i/ws/olv/ml-viz/go/server.go` - Add endpoints

**Frontend:**
- `/Users/i/ws/olv/ml-viz/js/src/types.ts` - Extend type definitions
- `/Users/i/ws/olv/ml-viz/js/src/App.svelte` - Mode switching
- `/Users/i/ws/olv/ml-viz/js/src/FormulaOverlay.svelte` - NEW
- `/Users/i/ws/olv/ml-viz/js/src/StepBreakdown.svelte` - NEW
- `/Users/i/ws/olv/ml-viz/js/src/CaseLibrary.svelte` - NEW
- `/Users/i/ws/olv/ml-viz/js/src/LiveTraining.svelte` - NEW
- `/Users/i/ws/olv/ml-viz/js/src/DatasetGenerator.svelte` - NEW
- `/Users/i/ws/olv/ml-viz/js/src/ComparativeView.svelte` - NEW
- `/Users/i/ws/olv/ml-viz/js/src/LossPlot.svelte` - Add interactivity

### Phases 2-4 (Future)

**Shared:**
- `/Users/i/ws/olv/ml-viz/js/src/shared/TimeSeriesPlot.svelte` - Extract from Phase 1
- `/Users/i/ws/olv/ml-viz/js/src/shared/types.ts` - Base interfaces

**Phase 2:**
- `/Users/i/ws/olv/ml-viz/go/phase2_linear/linear.go` - Vector gradient logic
- `/Users/i/ws/olv/ml-viz/js/src/phase2/ParameterSpace2D.svelte` - Trajectory plot

**Phase 3:**
- `/Users/i/ws/olv/ml-viz/go/phase3_neuron/activations.go` - Sigmoid/ReLU/Tanh + derivatives
- `/Users/i/ws/olv/ml-viz/js/src/phase3/ActivationCurve.svelte` - Saturation viz

**Phase 4:**
- `/Users/i/ws/olv/ml-viz/go/phase4_graph/backward.go` - Explicit backprop
- `/Users/i/ws/olv/ml-viz/js/src/phase4/ComputationGraph.svelte` - Graph animation

---

## SUCCESS METRICS

**Phase 1 Complete When:**
- User can select from 9 pre-computed cases
- Clicking any point in Loss/Param plots shows full state
- Step-by-step breakdown shows Forward → Loss → Gradient → Update
- Formula overlay shows live value substitution
- Live training works (generate dataset, configure, train, visualize)
- Comparative view shows 2-3 runs side-by-side

**Pedagogical Success (per original plan):**
*"A smart engineer can explain why vanishing gradients happen, without formulas, after playing with the visualization for 10 minutes."*

---

## NEXT IMMEDIATE ACTION

Start with **Stage 0: Foundation**

1. Extend `go/snapshots.go` with PointDetails and UpdateComponents
2. Modify `go/main.go` training loop to compute per-point values
3. Update `js/src/types.ts` with new interfaces
4. Verify existing visualizations still render correctly

Then proceed sequentially through Stages 1-4 for deep Phase 1 enhancements before tackling Phases 2-4.
