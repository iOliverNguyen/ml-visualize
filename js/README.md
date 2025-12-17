# TypeScript/Svelte Visualization

Interactive timeline for scrubbing through gradient descent training like a video.

## What It Does

Loads training snapshots from the Go engine and provides:
- **Loss plot**: See loss decreasing over time
- **Parameter plot**: Watch w converge to the true value (~2)
- **Gradient visualization**: Arrow showing gradient direction and magnitude
- **Controls**: Step forward/back, play/pause autoplay

## Files

- **src/App.svelte**: Main component, loads snapshots, manages playback state
- **src/Controls.svelte**: Timeline controls (step buttons, play/pause)
- **src/Plots.svelte**: Container for three canvas visualizations
- **src/plotting.ts**: Canvas rendering functions (pure functions, no framework)
- **src/types.ts**: TypeScript interfaces
- **src/main.ts**: Entry point

## Setup

```bash
pnpm install
```

## Run

### 1. Copy snapshots from Go engine

```bash
mkdir -p public
cp ../go/snapshots.json public/
```

### 2. Start dev server

```bash
pnpm dev
```

### 3. Open browser

Navigate to http://localhost:5173

## Controls

- **← Step Back**: Go to previous training step
- **▶ Play / ⏸ Pause**: Autoplay through training (10 steps/second)
- **Step Forward →**: Go to next training step
- **Step counter**: Shows current step / total steps

## Architecture

Built with Svelte for reactive state management:
- When `currentStep` changes, plots automatically re-render
- Playback is controlled via `setInterval` in App.svelte
- Canvas rendering is in plain TypeScript (plotting.ts)

## Why Svelte?

Svelte's reactivity makes the timeline scrubbing natural:

```svelte
$: if (canvas && snapshots.length > 0) {
  drawPlot(canvas, snapshots, currentStep)
}
```

When `currentStep` changes, the plot automatically updates. This is the core of "scrubbing learning like a video."
