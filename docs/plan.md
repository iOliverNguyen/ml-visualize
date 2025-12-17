Project Plan: ML-from-First-Principles Visualization

Goal: Make the learning dynamics of ML viscerally understandable, not performant or production-grade.
Tech: Go (core engine) + TypeScript (visualization)
Optional: Svelte

⸻

0. Prime Directive (tell Claude this explicitly)

This project optimizes for mechanistic clarity over abstraction.
No ML frameworks, no autograd, no hidden math.
Every value must be inspectable, serializable, and replayable.

If Claude violates this, it’s a failure—even if the code “works.”

⸻

1. System boundaries (non-negotiable)

Forbidden
•	TensorFlow, PyTorch, ONNX, sklearn
•	Automatic differentiation
•	“Trainer” or “Model” abstractions that hide math
•	Premature optimization

Required
•	Explicit gradients
•	Step-by-step state snapshots
•	Deterministic execution
•	Replayable training timeline

⸻

2. High-level architecture

Go = Computation Engine

Responsibilities:
•	Forward pass
•	Loss computation
•	Manual gradient computation
•	Parameter update
•	Emit immutable step snapshots

Output format: JSON snapshots (event-sourced)

TypeScript = Cognitive Interface

Responsibilities:
•	Render parameter evolution
•	Animate gradient descent
•	Allow step / play / rewind
•	Make invisible forces (gradients) visible

No ML logic in TS.

⸻

3. Phase 1 — Scalar Gradient Descent (foundation)

Objective

Make gradient descent felt, not explained.

Model

y_pred = w * x
loss = (y_pred - y_true)^2

Go tasks

Claude should:
1.	Implement:
•	Forward(w, x) -> y
•	Loss(y, y_true)
•	GradW(w, x, y_true)
2.	Implement training loop:
•	single parameter w
•	fixed dataset (hardcoded initially)
3.	Emit snapshot per step:

{
"step": 0,
"w": 0.12,
"grad_w": -1.34,
"loss": 2.91
}

	4.	Store snapshots in-memory + write to disk

Hard rule: No structs named Model or Trainer.

⸻

TypeScript tasks

Claude should:
1.	Load snapshot JSON
2.	Visualize:
•	loss vs step (line)
•	parameter value over time
•	gradient arrow (direction + magnitude)
3.	Controls:
•	step forward
•	step backward
•	autoplay / pause

Target: You can scrub learning like a video.

⸻

4. Phase 2 — Geometry of Learning (2 parameters)

Objective

Show gradient descent as vector movement in parameter space.

Model

y = w1*x1 + w2*x2 + b

Go tasks
•	Extend engine to:
•	vector parameters
•	vector gradients
•	Emit:

{
"step": 12,
"params": [w1, w2],
"grads": [g1, g2],
"loss": 0.42
}

TypeScript tasks
•	2D parameter space plot
•	Trajectory trace of (w1, w2)
•	Gradient vector overlay
•	Optional loss surface heatmap (static grid)

Insight to enforce:

Learning rate controls step geometry, not just speed.

⸻

5. Phase 3 — Single Neuron (activation exposed)

Objective

Make activation functions impossible to misunderstand.

Model

z = w·x + b
a = activation(z)

Activation functions to implement:
•	Sigmoid
•	ReLU
•	Tanh

Go tasks
•	Explicit derivative per activation
•	Emit:
•	z
•	a
•	dL/dz
•	dL/dw

TypeScript tasks
•	Plot activation curve
•	Highlight saturation regions
•	Animate gradient collapse when applicable

No batching yet. Single sample only.

⸻

6. Phase 4 — Backprop as Signal Flow

Objective

Backprop is graph traversal, not algebra.

Go tasks
•	Represent computation as explicit nodes:
•	Add
•	Multiply
•	Activation
•	Each node stores:
•	forward value
•	local derivative
•	Backward pass:
•	manual gradient propagation
•	Emit per-node gradient snapshots

TypeScript tasks
•	Render computation graph
•	Animate:
•	forward pass (values flowing)
•	backward pass (gradients flowing)
•	Toggle visibility per node

⸻

7. Engineering discipline for Claude

Tell Claude to:
•	Keep files small and named by concept
•	Write comments explaining why, not what
•	Prefer duplication over abstraction early
•	Commit after each phase

Repo structure:

/engine-go
scalar.go
linear.go
neuron.go
backprop.go

/viz-ts
timeline.ts
plots.ts
controls.ts


⸻

8. Definition of success (very important)

This project is successful only if:
•	A smart engineer can explain why vanishing gradients happen
•	Without formulas
•	After playing with the visualization for 10 minutes

If not, add clarity—not features.

⸻

9. Next action (immediate)

Give Claude only Phase 1 first.

Explicit instruction:

“Do not anticipate later phases.
Make scalar gradient descent so clear it feels embarrassing to misunderstand.”

If you want, next we can:
•	Rewrite this as a single Claude system prompt
•	Tighten snapshot schemas
•	Or define the one insight you want the whole project to orbit around

Your move.
