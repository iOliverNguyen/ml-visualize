# ML Visualization - Interactive Gradient Descent Explorer

Interactive visualization for understanding gradient descent training, step by step. Watch machine learning optimization like a video — scrub through each training step to see exactly how gradients guide parameter updates.

## What It Does

Provides interactive visualizations for gradient descent training:

### Phase 1: Single-Parameter Linear Regression (y = w·x)
- **Case Library**: Pre-computed training scenarios (perfect start, noisy data, different learning rates)
- **Interactive Training**: Generate random data or define custom datasets with JavaScript
- **Client-Side Training**: Run gradient descent entirely in your browser (no server needed)
- **Visualizations**: Loss plots, parameter evolution, gradient arrows, data scatter plots
- **Timeline Controls**: Scrub through training steps, play/pause, jump to key moments

### Phase 2: 2-Parameter Linear Regression (y = w₁·x₁ + w₂·x₂)
- **Parameter Space Trajectory**: Watch optimization path through (w₁, w₂) space
- **Loss Contour Surface**: Heatmap showing loss landscape with trajectory overlay
- **Gradient Vector Field**: Arrows showing descent directions across parameter space
- **7 Pre-computed Cases**: Optimal LR, too slow, too fast, anisotropic, saddle points, zigzag

## Architecture

### Pure Static Site (No Server Required!)
- **Client-Side Training**: Gradient descent implemented in TypeScript (`src/shared/training.ts`)
- **Pre-Computed Cases**: JSON files in `public/cases/` and `public/cases-phase2/`
- **Framework**: Svelte 5 with runes ($state, $derived, $effect)
- **Deployment**: Can be hosted on GitHub Pages, Netlify, Vercel, or any static host

### Key Files
- **src/App.svelte**: Main app with phase selector
- **src/phase1/Phase1App.svelte**: Single-parameter UI container
- **src/phase2/Phase2App.svelte**: 2-parameter UI container
- **src/shared/training.ts**: Client-side gradient descent engine
- **src/shared/Controls.svelte**: Timeline controls (used by both phases)
- **public/cases/**: Pre-computed Phase 1 training scenarios
- **public/cases-phase2/**: Pre-computed Phase 2 training scenarios

## Setup

```bash
pnpm install
```

## Run

```bash
pnpm dev
```

Open http://localhost:5005 in your browser.

**No server needed!** The frontend runs entirely client-side.

## Build for Production

```bash
pnpm build
```

This generates static files in `dist/` ready for deployment to any static host.

### Deploy to Static Hosts

**GitHub Pages:**
```bash
pnpm build
# Push dist/ to gh-pages branch
```

**Netlify:**
```bash
# Build command: cd js && pnpm build
# Publish directory: js/dist
```

**Vercel/Cloudflare Pages:**
Similar configuration, build in `js/`, publish `js/dist/`.

## Interactive Features

### Pre-Computed Cases
Browse curated training scenarios:
- Perfect initialization vs. random start
- Different noise levels
- Learning rate comparisons (too slow, optimal, too fast)
- Pathological cases (saddle points, zigzag convergence)

### Custom Training
Train your own models in the browser:
- **Random Data**: Configure data points, x range, slope, noise level
- **JavaScript Function**: Define data programmatically with custom code
- **Instant Results**: Training runs client-side in <10ms
- **Full Visualization**: See loss, gradients, and parameter evolution

## Development

### Generate Pre-Computed Cases (Optional)
If you want to update or add pre-computed training scenarios:

```bash
# From project root
go run . --generate-cases          # Phase 1 cases → js/public/cases/
go run . --generate-cases-phase2   # Phase 2 cases → js/public/cases-phase2/
```

This is **optional** — the repo already includes pre-computed cases.

## Why Svelte 5?

Svelte's reactive runes make timeline scrubbing natural:

```typescript
let currentStep = $state(0);
let snapshot = $derived(snapshots[currentStep]);

// When currentStep changes, snapshot automatically updates
// and all components using snapshot re-render
```

This is the core of "scrubbing learning like a video" — pure reactive state without manual subscriptions.
