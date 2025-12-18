# ML Visualization: Learning from First Principles

Make the learning dynamics of ML viscerally understandable.

## Philosophy

This project optimizes for **mechanistic clarity over abstraction**:
- No ML frameworks (TensorFlow, PyTorch, etc.)
- No automatic differentiation
- Every value is inspectable, serializable, and replayable
- Explicit gradients computed by hand

## Architecture

```
Go Engine → JSON Snapshots → TypeScript/Svelte Visualization
```

- **Go**: Computes forward pass, loss, gradients, and parameter updates
- **JSON**: Event-sourced snapshots of every training step
- **TypeScript/Svelte**: Interactive timeline to scrub through learning like a video

## Phase 1: Scalar Gradient Descent

Current implementation: Single parameter learning `y = w*x`

### Model
```
y_pred = w * x
loss = (y_pred - y_true)^2
grad_w = 2(w*x - y_true) * x
```

### What You Can Do
- See loss decreasing over time
- Watch parameter w converge to the true value
- Observe gradient magnitude shrinking
- Step forward/backward through training
- Play/pause autoplay
- Scrub learning like a video

## Quick Start

### 1. Run the Go Server
```bash
go run . --server
```

### 2. Run the Visualization
```bash
cd ts
mkdir -p public
cp ../go/snapshots.json public/
pnpm install
pnpm dev
# Open http://localhost:5005
```
