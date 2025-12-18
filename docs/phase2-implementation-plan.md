# Phase 2 Implementation Plan: 2-Parameter Gradient Descent

## Overview

Implement Phase 2 of the ML visualization application: **2-Parameter Linear Model (y = w1*x1 + w2*x2)**. This phase teaches learners how gradient descent operates as vector movement in parameter space, with emphasis on learning rate effects on step geometry and convergence patterns.

**Current State:** Phase 1 (scalar GD) is 85% complete with all core pedagogical features implemented.

**Goal:** Build complete Phase 2 with parameter space trajectories, gradient field visualization, and loss surface contour plots.

---

## Architecture Decisions

### File Organization: Integrated with Routing

Use a single application with phase routing rather than separate apps. This maximizes code reuse while maintaining clear separation.

**Directory structure:**
```
go/
  core/           # Existing Phase 1 code
  linear/         # NEW: Phase 2 package
js/src/
  shared/         # NEW: Extract reusable components
  phase1/         # Move existing Phase 1 components here
  phase2/         # NEW: Phase 2 components
  educational/    # Existing (shared across phases)
  App.svelte      # MODIFY: Add phase router
  PhaseSelector.svelte  # NEW: Phase selector tabs
```

### Backend: Separate Go Package

Create `go/linear/` package for Phase 2:
- Keeps Phase 1 code unchanged
- Clean separation of concerns
- Allows parallel development
- Extensible for Phase 3 and 4

### Frontend: Phase-Specific Directories

Organize components by phase with shared components extracted:
- `/phase1/` - Scalar GD components
- `/phase2/` - 2-parameter linear components
- `/shared/` - Reusable across phases (Controls, utils)

---

## Implementation Sequence

### Stage 1: Backend Foundation (2-3 days)

#### 1.1 Create Go Package Structure

**File:** `go/linear/model.go`
- `Forward(w1, w2, x1, x2)` - Compute y_pred = w1*x1 + w2*x2
- `Loss(yPred, yTrue)` - MSE: (y_pred - y_true)²
- `GradW1()`, `GradW2()` - Partial derivatives
- `GradientMagnitude()` - ||∇L|| = sqrt(grad_w1² + grad_w2²)
- `GradientDirection()` - atan2(grad_w2, grad_w1)

**File:** `go/linear/snapshots.go`
- `DataPoint2D` struct (x1, x2, y_true)
- `PointSnapshot2D` struct (per-point breakdown with grad_w1, grad_w2)
- `UpdateDetails2D` struct (parameter update breakdown)
- `LinearSnapshot` struct (complete state at one step)

**File:** `go/linear/training.go`
- `TrainingConfig2D` struct (w1_init, w2_init, lr, steps)
- `RunTraining()` - Main training loop generating snapshots
- Per-point gradient computation and aggregation
- Vector update: (w1, w2) ← (w1, w2) - lr * (grad_w1, grad_w2)

**File:** `go/linear/loss_surface.go`
- `LossGrid` struct (w1_min, w1_max, w2_min, w2_max, resolution, points)
- `LossGridPoint` struct (w1, w2, loss)
- `ComputeLossGrid()` - Generate 50×50 grid for contour plots
- Compute loss at each grid point for visualization

**File:** `go/linear/dataset.go`
- `DataGenConfig2D` struct (num_points, x1/x2 ranges, true_w1, true_w2, noise, seed)
- `GenerateRandomData()` - Generate y = w1*x1 + w2*x2 + noise

#### 1.2 Generate Pre-computed Cases

**File:** `go/linear/generate_cases.go`
- Generate 7 pre-computed cases with different configurations:
  1. `lr-small` - Learning rate too small (0.0001)
  2. `lr-optimal` - Optimal learning rate (0.01)
  3. `lr-large` - Learning rate too large (0.1)
  4. `anisotropic-easy` - Mild scale difference between x1 and x2
  5. `anisotropic-hard` - Extreme scale difference
  6. `saddle-point` - Near-saddle geometry
  7. `zigzag-convergence` - High LR with anisotropy

**Output:** Save to `js/public/cases-phase2/{case-id}/`
- `snapshots.json` - Training snapshots
- `loss_grid.json` - Pre-computed loss surface
- `manifest.json` - Case metadata

#### 1.3 Add Server Endpoints

**Modify:** `go/core/server.go`

Add Phase 2 endpoints:
- `GET /api/phase2/cases` - List all Phase 2 cases
- `GET /api/phase2/cases/:id/snapshots` - Get case snapshots
- `GET /api/phase2/cases/:id/loss_grid` - Get loss surface
- `POST /api/phase2/train` - Train with custom dataset and config

**Request/Response formats:**
```json
POST /api/phase2/train
{
  "data": [{"x1": 1.0, "x2": 2.0, "y_true": 5.0}, ...],
  "config": {"w1_init": 0, "w2_init": 0, "lr": 0.01, "steps": 100}
}

Response:
{
  "snapshots": [...],
  "loss_grid": {...}
}
```

---

### Stage 2: Frontend Reorganization (1-2 days)

#### 2.1 Extract Shared Components

**Create:** `js/src/shared/`

Move Phase 1 components that work for both phases:
- `Controls.svelte` - Timeline control (play/pause/scrub)
- `types.ts` - Base interfaces
- `utils.ts` - Formatting helpers (formatNumber, formatPercent)
- `scaling.ts` - SVG coordinate transformations

#### 2.2 Reorganize Phase 1

**Create:** `js/src/phase1/`

Move existing Phase 1 components:
- `Phase1App.svelte` (rename from App.svelte)
- `ScalarMetricsPanel.svelte` (rename from MetricsPanel.svelte)
- `LossPlot.svelte`
- `ParamPlot.svelte`
- `GradientViz.svelte`
- `StepBreakdown.svelte`
- `FormulaOverlay.svelte`
- `DataScatterPlot.svelte`
- `FitComparisonView.svelte`
- `DataInputPanel.svelte`
- `CaseLibrary.svelte`
- `types.ts` - Phase 1 specific types

Update imports in all components to reference `/shared/` for common components.

#### 2.3 Add Phase Routing

**Modify:** `js/src/App.svelte`

```svelte
<script lang="ts">
  import Phase1App from './phase1/Phase1App.svelte';
  import Phase2App from './phase2/Phase2App.svelte';
  import PhaseSelector from './PhaseSelector.svelte';

  let currentPhase = $state<'phase1' | 'phase2'>(
    (localStorage.getItem('currentPhase') as 'phase1' | 'phase2') || 'phase1'
  );

  function handlePhaseChange(phase: 'phase1' | 'phase2') {
    currentPhase = phase;
    localStorage.setItem('currentPhase', phase);
  }
</script>

<PhaseSelector {currentPhase} onchange={handlePhaseChange} />

{#if currentPhase === 'phase1'}
  <Phase1App />
{:else if currentPhase === 'phase2'}
  <Phase2App />
{/if}
```

**Create:** `js/src/PhaseSelector.svelte`

Tab-based selector with two options:
- Phase 1: Scalar GD (y = w*x)
- Phase 2: 2-Parameter Linear (y = w1*x1 + w2*x2)

Persist selection to localStorage.

#### 2.4 Verification

Test that Phase 1 still works after reorganization:
- All visualizations render
- Controls function correctly
- Case library loads
- Dataset generation works

---

### Stage 3: Phase 2 Core Visualizations (4-5 days)

#### 3.1 TypeScript Types

**Create:** `js/src/phase2/types.ts`

```typescript
export interface DataPoint2D {
  x1: number;
  x2: number;
  y_true: number;
}

export interface PointSnapshot2D {
  x1: number;
  x2: number;
  y_true: number;
  y_pred: number;
  point_loss: number;
  grad_w1: number;
  grad_w2: number;
}

export interface UpdateDetails2D {
  w1_old: number;
  w2_old: number;
  lr: number;
  grad_w1: number;
  grad_w2: number;
  delta_w1: number;
  delta_w2: number;
  w1_new: number;
  w2_new: number;
}

export interface LinearSnapshot {
  step: number;
  w1: number;
  w2: number;
  grad_w1: number;
  grad_w2: number;
  loss: number;
  gradient_magnitude: number;
  gradient_direction: number;
  point_details: PointSnapshot2D[];
  update_components: UpdateDetails2D;
}

export interface LossGridPoint {
  w1: number;
  w2: number;
  loss: number;
}

export interface LossGrid {
  w1_min: number;
  w1_max: number;
  w2_min: number;
  w2_max: number;
  resolution: number;
  points: LossGridPoint[];
}
```

#### 3.2 Parameter Space Trajectory Plot

**Create:** `js/src/phase2/ParameterSpace2D.svelte`

**Features:**
- SVG plot showing (w1, w2) trajectory through parameter space
- Line path connecting parameter updates
- Current position marker (animated with tweened)
- Start marker (green circle)
- End marker (red circle)
- Click any point to jump to that step
- Hover tooltips showing (w1, w2, loss)
- Grid with axis labels
- Optional: Show gradient arrows at each point

**Props:**
- `snapshots: LinearSnapshot[]`
- `currentStep: number`
- `onStepClick: (step: number) => void`

**Implementation details:**
- Use `viewBox` for responsive scaling
- Compute axis ranges dynamically from trajectory bounds
- Draw trajectory as SVG `<path>` with `stroke`
- Draw circles for start/end/current positions
- Animate current position marker with Svelte's `tweened`

#### 3.3 Loss Contour Heatmap

**Create:** `js/src/phase2/LossContour.svelte`

**Features:**
- Heatmap showing loss surface landscape
- Color gradient: blue (low loss) → yellow → red (high loss)
- Overlay trajectory path from parameter space
- Current position marker
- Optional contour lines at specific loss values
- Axis labels (w1, w2)

**Props:**
- `lossGrid: LossGrid`
- `snapshots: LinearSnapshot[]`
- `currentStep: number`

**Implementation approach:**
- Use `<canvas>` for heatmap rendering (better performance)
- SVG overlay for trajectory and annotations
- Color scale: Interpolate between blue (#3b82f6) → yellow (#eab308) → red (#ef4444)
- Normalize loss values to [0, 1] for color mapping

**Rendering process:**
1. Create canvas matching grid resolution
2. Draw each grid point as colored pixel
3. Apply smoothing (optional)
4. Overlay trajectory as SVG path
5. Add current position marker

#### 3.4 Gradient Vector Field

**Create:** `js/src/phase2/GradientField.svelte`

**Features:**
- Grid of arrows showing -∇L direction (descent direction)
- Arrow color indicates gradient magnitude
- Arrow length proportional to magnitude (normalized)
- Overlay trajectory path
- Current parameter position highlighted
- Toggle field density (coarse: 10×10, fine: 20×20)

**Props:**
- `lossGrid: LossGrid` (compute gradients from loss surface)
- `snapshots: LinearSnapshot[]`
- `currentStep: number`
- `density: 'coarse' | 'fine'`

**Implementation:**
- Compute gradient vectors at grid points using finite differences
- Normalize arrow lengths (scale factor configurable)
- Use HSL color scale: hsl(240 - magnitude*120, 80%, 50%)
  - Low magnitude: blue
  - High magnitude: red
- Draw arrows as SVG `<line>` with arrowhead markers

#### 3.5 Metrics Panel

**Create:** `js/src/phase2/LinearMetricsPanel.svelte`

**Metrics to display:**

**Current State:**
- w1, w2 (parameter values)
- ||∇L|| (gradient magnitude)
- ∠∇L (gradient direction in degrees)
- Loss

**Comparison:**
- Initial vs Current vs Final
- Δw1, Δw2 from start
- Distance traveled: sqrt((w1-w1_init)² + (w2-w2_init)²)
- Progress percentage

**Step Information:**
- Learning rate
- Step size: lr × ||∇L||
- Update vector: (-lr*grad_w1, -lr*grad_w2)

**Layout:** Similar to Phase 1 ScalarMetricsPanel but adapted for vector quantities.

#### 3.6 Gradient Vector Visualization

**Create:** `js/src/phase2/GradientVectorViz.svelte`

**Features:**
- Large arrow showing current gradient vector ∇L
- Arrow pointing in -∇L direction (descent direction)
- Length proportional to magnitude
- Display vector components: (grad_w1, grad_w2)
- Show step vector: -lr × ∇L
- Comparison overlay: gradient vs update vector

**Props:**
- `snapshot: LinearSnapshot`
- `lr: number`

**Visual design:**
- Gradient arrow: Blue, dashed
- Update arrow: Green, solid
- Origin at center
- Annotations showing component values

---

### Stage 4: Pedagogical Components (2-3 days)

#### 4.1 Step Breakdown

**Create:** `js/src/phase2/StepBreakdown2D.svelte`

**Tabs (same structure as Phase 1):**

1. **Forward Pass**
   - Table showing per-point: x1, x2, y_pred = w1*x1 + w2*x2
   - Aggregate: No aggregation needed (per-point predictions)

2. **Loss Computation**
   - Table showing per-point: (y_pred - y_true)²
   - Aggregate: Average loss across points

3. **Gradient Computation**
   - Table showing per-point: grad_w1, grad_w2
   - Aggregate: Average gradients
   - Show formula: ∂L/∂w1 = 2(y_pred - y_true)*x1

4. **Parameter Update**
   - Show vector update equation
   - w1_new = w1_old - lr * grad_w1
   - w2_new = w2_old - lr * grad_w2
   - Display delta values

**Props:**
- `snapshot: LinearSnapshot`

#### 4.2 Formula Overlay

**Create:** `js/src/phase2/FormulaOverlay2D.svelte`

**Modes (similar to Phase 1):**

1. **Loss Mode**
   - Formula: L = (1/n) Σ (w1*x1i + w2*x2i - yi)²
   - Substitute current values
   - Show expanded calculation

2. **Gradient Mode**
   - Formula: ∇L = [∂L/∂w1, ∂L/∂w2]
   - Show partial derivatives
   - ∂L/∂w1 = (2/n) Σ (y_pred - y_true)*x1
   - ∂L/∂w2 = (2/n) Σ (y_pred - y_true)*x2

3. **Update Mode**
   - Formula: w_new = w_old - lr * ∇L
   - Vector form: [w1_new, w2_new] = [w1_old, w2_old] - lr * [grad_w1, grad_w2]
   - Show component-wise calculation

4. **Magnitude Mode** (NEW)
   - Formula: ||∇L|| = sqrt(grad_w1² + grad_w2²)
   - Direction: θ = atan2(grad_w2, grad_w1)

**Props:**
- `snapshot: LinearSnapshot`
- `mode: 'loss' | 'gradient' | 'update' | 'magnitude'`

#### 4.3 Case Library

**Create:** `js/src/phase2/CaseLibrary2D.svelte`

**Features:**
- Grid view of 7 pre-computed cases
- Category filters: Learning Rate, Anisotropy, Geometry
- Click card to load case
- Display: emoji, name, description, training config, insights

**Props:**
- `onCaseSelect: (caseId: string) => void`

**Data source:** Load from `js/public/cases-phase2/manifest.json`

---

### Stage 5: App Integration (2-3 days)

#### 5.1 Main Phase 2 App

**Create:** `js/src/phase2/Phase2App.svelte`

**Layout:**

```
+--------------------------------------------------+
| Header: "2-Parameter Gradient Descent"          |
|   [Case Library] [Help]                         |
+--------------------------------------------------+
| Controls: [<< | < | Play | > | >>] [Slider]    |
+--------------------------------------------------+
| +----------------------+ +----------------------+ |
| | Parameter Space 2D   | | Loss Contour        | |
| | (trajectory plot)    | | (heatmap + path)    | |
| +----------------------+ +----------------------+ |
| +----------------------+ +----------------------+ |
| | Gradient Field       | | Metrics Panel       | |
| | (vector field)       | | (current values)    | |
| +----------------------+ +----------------------+ |
+--------------------------------------------------+
| Gradient Vector Viz (large arrow)               |
+--------------------------------------------------+
| Step Breakdown (tabs: Fwd | Loss | Grad | Upd) |
+--------------------------------------------------+
| Formula Overlay (modes: loss/gradient/update)   |
+--------------------------------------------------+
```

**State management:**
- `snapshots: LinearSnapshot[]`
- `lossGrid: LossGrid`
- `currentStep: number`
- `playing: boolean`
- Load case from API on mount or case selection

**Props:** None (top-level app)

#### 5.2 Data Loading

Load snapshots and loss grid:
- Default: Load first case from manifest
- User selects case: Fetch `/api/phase2/cases/:id/snapshots` and `/api/phase2/cases/:id/loss_grid`
- Custom training: POST to `/api/phase2/train` and render results

#### 5.3 Synchronization

Wire up Controls to update all Phase 2 views:
- `currentStep` changes → all components re-render
- ParameterSpace2D: Update current position marker
- LossContour: Update current position marker
- GradientField: Update trajectory overlay
- LinearMetricsPanel: Update displayed metrics
- GradientVectorViz: Update arrow direction/length

Use Svelte's reactive `$derived` for computed values.

#### 5.4 Interactivity

Add click handlers:
- Click point in ParameterSpace2D → Jump to that step
- Click in LossContour → (Optional) Show loss value at that point
- Hover on trajectory → Show tooltip with step info

---

### Stage 6: Polish & Educational Content (2-3 days)

#### 6.1 Educational Enhancements

**Extend glossary** (`js/public/content/glossary.json`):
- Add entries for: Gradient Vector, Parameter Space, Learning Rate Geometry, Anisotropic Loss Surface, Contour Plot

**Extend FAQs** (`js/public/content/faqs.json`):
- Why does learning rate cause zigzag patterns?
- What is an anisotropic loss surface?
- How do I interpret the gradient field?
- Why do contour lines show loss landscape?

**Extend tutorial** (`js/public/content/tutorial.json`):
- Chapter 9: From Scalar to Vector Gradients
- Chapter 10: Understanding Parameter Trajectories
- Chapter 11: Learning Rate and Step Geometry

#### 6.2 Performance Optimization

- Canvas rendering for LossContour (50×50 grid)
- Lazy load loss grids (only fetch when case selected)
- Debounce hover interactions
- Use SVG `<g>` transforms for efficient trajectory rendering

#### 6.3 Responsive Design

Add mobile/tablet layouts:
- Stack visualizations vertically on narrow screens
- Reduce gradient field density on mobile
- Collapsible sections for step breakdown

#### 6.4 Testing

- Verify all 7 cases load and render correctly
- Test timeline controls (play/pause/scrub)
- Test phase switching (Phase 1 ↔ Phase 2)
- Check gradient field arrow directions (should point downhill)
- Validate metrics calculations match backend

---

## Critical Files

### Backend (Go)

| File | Purpose |
|------|---------|
| `go/linear/model.go` | Core math: Forward, Loss, GradW1, GradW2, magnitude, direction |
| `go/linear/snapshots.go` | Type definitions: LinearSnapshot, PointSnapshot2D, UpdateDetails2D |
| `go/linear/training.go` | Training loop generating snapshots with per-point breakdown |
| `go/linear/loss_surface.go` | Pre-compute loss grid for contour visualization |
| `go/linear/dataset.go` | Generate random 2D datasets with configurable noise |
| `go/linear/generate_cases.go` | Generate 7 pre-computed cases for case library |
| `go/core/server.go` | Add Phase 2 API endpoints (MODIFY) |

### Frontend (Svelte)

| File | Purpose |
|------|---------|
| `js/src/App.svelte` | Phase router (MODIFY to add PhaseSelector) |
| `js/src/PhaseSelector.svelte` | Tab selector for Phase 1 vs Phase 2 (NEW) |
| `js/src/phase2/Phase2App.svelte` | Main Phase 2 container component |
| `js/src/phase2/types.ts` | TypeScript interfaces for Phase 2 data |
| `js/src/phase2/ParameterSpace2D.svelte` | Trajectory plot showing (w1, w2) path |
| `js/src/phase2/LossContour.svelte` | Heatmap of loss surface with trajectory overlay |
| `js/src/phase2/GradientField.svelte` | Vector field showing gradient directions |
| `js/src/phase2/LinearMetricsPanel.svelte` | Display vector quantities and progress |
| `js/src/phase2/GradientVectorViz.svelte` | Large arrow showing gradient direction |
| `js/src/phase2/StepBreakdown2D.svelte` | 4-stage pedagogical breakdown |
| `js/src/phase2/FormulaOverlay2D.svelte` | Vector math formulas with live values |
| `js/src/phase2/CaseLibrary2D.svelte` | Case selector for Phase 2 scenarios |

### Data

| File | Purpose |
|------|---------|
| `js/public/cases-phase2/manifest.json` | Index of 7 Phase 2 cases |
| `js/public/cases-phase2/{case-id}/snapshots.json` | Training snapshots for each case |
| `js/public/cases-phase2/{case-id}/loss_grid.json` | Pre-computed loss surface (50×50 grid) |

---

## Technical Decisions

### 1. Loss Surface Computation
**Decision:** Pre-compute 50×50 grids in Go, store as JSON
**Rationale:** Keeps computation in Go (single source of truth), acceptable file size (~50KB), consistent with Phase 1 philosophy

### 2. Rendering Strategy
**Decision:** Canvas for LossContour, SVG for ParameterSpace2D and GradientField
**Rationale:** Canvas handles dense heatmaps efficiently, SVG provides clean interactivity for trajectories and vectors

### 3. Phase Navigation
**Decision:** Tab-based selector at top of app
**Rationale:** Clear separation, allows independent URLs, easy to extend to Phase 3/4

### 4. Code Organization
**Decision:** Separate phase directories (`phase1/`, `phase2/`) with shared components in `shared/`
**Rationale:** Clean ownership, avoids type collisions, extensible to future phases

---

## Success Criteria

**Phase 2 is complete when:**

1. User can select from 7 pre-computed cases showing different scenarios
2. ParameterSpace2D displays smooth trajectory with animated current position
3. LossContour renders heatmap with overlaid trajectory path
4. GradientField shows descent directions across parameter space
5. All visualizations update synchronously when scrubbing timeline
6. Clicking trajectory jumps to that step
7. LinearMetricsPanel displays vector quantities (magnitude, direction, etc.)
8. StepBreakdown2D shows per-point 2D gradient computation
9. FormulaOverlay2D displays vector math with live value substitution
10. PhaseSelector allows switching between Phase 1 and Phase 2

**Pedagogical goal achieved when:**

*"A learner can explain why high learning rates cause zigzag patterns in anisotropic loss surfaces, and demonstrate this by comparing lr-optimal vs lr-large vs anisotropic-hard cases."*

---

## Estimated Timeline

- **Stage 1 (Backend):** 2-3 days
- **Stage 2 (Reorganization):** 1-2 days
- **Stage 3 (Visualizations):** 4-5 days
- **Stage 4 (Pedagogy):** 2-3 days
- **Stage 5 (Integration):** 2-3 days
- **Stage 6 (Polish):** 2-3 days

**Total:** 13-19 days (approximately 3-4 weeks)

---

## Next Immediate Action

Start with **Stage 1: Backend Foundation**

1. Create `go/linear/` directory structure
2. Implement `model.go` with core mathematical functions
3. Define snapshot types in `snapshots.go`
4. Build training loop in `training.go`
5. Implement loss surface computation in `loss_surface.go`
