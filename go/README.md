# Go Computation Engine

This is the computation engine for Phase 1: Scalar Gradient Descent.

## What It Does

Trains a single-parameter model `y = w*x` on hardcoded data using explicit gradient descent:
1. Forward pass: compute predictions
2. Loss computation: mean squared error
3. Gradient computation: explicit derivative (no autograd)
4. Parameter update: `w_new = w_old - lr * grad_w`
5. Snapshot creation: captures state at each step

## Files

- **main.go**: Training loop orchestration (explicit, no abstractions)
- **scalar.go**: Pure math functions (Forward, Loss, GradW)
- **dataset.go**: Hardcoded training data
- **snapshots.go**: JSON serialization

## Run

```bash
go run .
```

## Output

Creates `snapshots.json` containing an array of training snapshots:

```json
[
    {
        "step": 0,
        "w": 0.0,
        "grad_w": -123.4,
        "loss": 456.7
    },
    ...
]
```

## Configuration

Edit `main.go` to change:
- Initial parameter value (`w`)
- Learning rate (`lr`)
- Number of training steps (`steps`)

Edit `dataset.go` to change training data.

## Philosophy

No Model or Trainer abstractions. The training loop is explicit imperative code that shows exactly what happens at each step. Every gradient is computed by hand using the chain rule.
