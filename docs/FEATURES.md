# ML Visualization: Feature Documentation

## Overview

An interactive web application for learning gradient descent through visual exploration. Designed for ML learners, educators, and practitioners who want to build intuition about how gradient descent optimizes parameters. Watch training unfold step-by-step with live formula substitutions, per-point calculations, and progressive educational content.

**Target Audience**: Students learning ML fundamentals, educators demonstrating optimization, practitioners reviewing basics.

**Core Value**: Transform abstract gradient descent equations into concrete, explorable visualizations with immediate feedback.

## Core Features

### Visualization Components

**Loss Over Time Plot**
- Track loss decrease throughout training
- Click any point to jump to that step
- Hover for exact loss values
- Visual trend line showing optimization progress

**Parameter Trajectory Plot**
- Watch weight parameter converge to optimal value
- See gradient magnitude shrink over time
- Understand relationship between parameter updates and loss

**Gradient Visualization**
- Real-time gradient values at each step
- Visual indication of gradient direction and magnitude
- Link gradient computation to parameter updates

**Data Scatter Plot with Fitted Line**
- View training data points with true vs predicted values
- See model's fitted line improve over training
- Optional residual lines showing per-point errors
- Interactive hover for point details (x, y_true, y_pred, loss, gradient)

**Training Progress Comparison**
- Side-by-side view of initial, current, and final model fits
- Instantly compare where training started and where it's heading
- Visual progress indicator showing training completion percentage

**Multi-Layout Support**
- Four layout modes: 1-column, 2-column, 3-column, auto-responsive
- Persisted preference saved to browser localStorage
- Optimized viewing for different screen sizes and workflows

### Interactive Controls

**Timeline Playback**
- Play/pause animation through training steps
- Adjustable playback speed (10 steps per second default)
- Automatic pause at final step

**Navigation Controls**
- Step forward: advance one iteration
- Step backward: return to previous iteration
- Jump to start: reset to initial state
- Jump to final: skip to converged solution
- Slider scrubbing: drag to any training step instantly

**Keyboard Shortcuts**
- Accessible navigation without mouse
- Standard media controls (play/pause, step)

### Data & Training

**Pre-loaded Example Dataset**
- Default dataset demonstrates successful convergence
- Ready to explore immediately on application start

**Random Data Generation**
- Configurable slope (true parameter value)
- Adjustable noise level (data variance)
- Custom x-range (min/max values)
- Seed control for reproducibility
- Number of data points selection

**Custom JavaScript Function Input**
- Write arbitrary functions to generate y-values
- Access to Math library and x values
- Test edge cases and non-linear patterns
- Immediate training on custom data

**Real-time Training Visualization**
- Server-side gradient computation with explicit derivatives
- Per-step snapshots with complete calculation breakdown
- Forward pass → loss → gradient → update pipeline captured

**Flexible Hyperparameters**
- Learning rate: control optimization speed and stability
- Initial weight: explore different starting points
- Training steps: balance convergence vs. computation time

### Educational Features

**Multi-Chapter Tutorial**
- 8 progressive chapters covering gradient descent fundamentals:
  1. Introduction to the Learning Problem
  2. Understanding Loss Functions
  3. The Gradient Concept
  4. Parameter Updates
  5. Learning Rate Effects
  6. Convergence Behavior
  7. Overfitting vs. Underfitting
  8. Advanced Optimization Concepts
- Table of contents with chapter navigation
- Estimated reading time per chapter
- Visual cues linking tutorial text to UI elements
- Auto-closing TOC after chapter selection

**Interactive Glossary**
- 8 core ML terms with definitions:
  - Gradient Descent, Learning Rate, Loss Function, Parameter/Weight
  - Convergence, Optimization, Gradient, Epoch/Step
- Three detail levels per term (brief, detailed, comprehensive)
- Mathematical formulas where applicable
- Related term navigation (click to jump between connected concepts)
- Examples demonstrating each term
- Level-based filtering (beginner/intermediate/advanced)

**FAQ System**
- 20+ common questions organized by category
- Categories: Getting Started, Understanding Plots, Training Behavior, Troubleshooting
- Search functionality across all questions
- Level-based content filtering
- Expandable accordion interface

**Step-by-Step Formula Breakdown**
- Four-stage training breakdown for each iteration:
  1. **Forward Pass**: Compute predictions (ŷ = w·x)
  2. **Loss Computation**: Calculate mean squared error
  3. **Gradient Calculation**: Derive ∂L/∂w
  4. **Parameter Update**: Apply w_new = w_old - lr·grad
- Live value substitution in formulas
- Color-coded components matching visualization

**Per-Point Calculation Inspection**
- Drill down to individual data point contributions
- See x, y_true, y_pred, point_loss, point_gradient for each sample
- Understand how batch gradient averages individual contributions

**Level-Based Content Adaptation**
- Three user levels: beginner, intermediate, advanced
- Content automatically filtered by selected level
- More technical depth for advanced users
- Simplified explanations for beginners
- Persistent level selection saved to localStorage

### Pedagogical Tools

**Formula Overlay with Live Values**
- Three interactive formula cards:
  1. Loss formula: L = (1/n)Σ(y_true - y_pred)²
  2. Gradient formula: ∂L/∂w = -(2/n)Σx·(y_true - y_pred)
  3. Update formula: w_new = w_old - lr·∂L/∂w
- Current step values substituted directly into equations
- See abstract math become concrete numbers

**Visual Annotations**
- Highlighted elements when referenced in tutorial
- Suggested step numbers for concept exploration
- Focus metrics that demonstrate specific phenomena

**Progress Indicators**
- Training progress percentage (0-100%)
- Delta from start: Δw, Δloss since initialization
- Delta to target: remaining distance to optimal solution
- Current vs. initial vs. final metrics comparison

## Technical Capabilities

### Backend (Go Server)

**RESTful API**
- `/api/snapshots` - Retrieve pre-computed training run
- `/api/random-dataset` - Generate random data and train (POST)
- `/api/train-custom` - Train on custom dataset (POST)

**Dataset Generation Engine**
- Linear data: y = true_slope·x + noise
- Gaussian noise with configurable standard deviation
- Seeded random generation for reproducibility

**Gradient Computation**
- Explicit analytical derivatives (not numeric approximation)
- Per-point gradient tracking for educational transparency
- Full computation breakdown captured at every step

**Snapshot System**
- Complete state serialization per training iteration
- Includes: step number, w, grad_w, loss, point_details, update_components
- JSON format for easy consumption by frontend

### Frontend Architecture

**Svelte 5 with Runes**
- `$state` for reactive variables
- `$derived` for computed values
- `$effect` for side effects and lifecycle
- Modern reactive paradigm with fine-grained updates

**Component-Based Design**
- 20+ specialized components (plots, controls, educational panels)
- Clear separation of concerns (data, visualization, interaction, education)
- Reusable building blocks (tooltips, accordions, overlays)

**localStorage Persistence**
- Layout preference (1col/2col/3col/auto)
- User level selection (beginner/intermediate/advanced)
- Sidebar open/close state
- Introduction panel dismissed state

**Responsive Layouts**
- Mobile-optimized views (< 768px)
- Tablet adaptations (768-1024px)
- Desktop full layouts (> 1024px)
- CSS Grid for flexible plot arrangements

**Smooth Animations**
- Tweened motion for parameter transitions
- Slide transitions for sidebar and panels
- Easing functions (cubicOut) for natural feel
- 300ms transition duration standard

### Data Loading Modes

**Server Mode** (Development)
- Go server at `localhost:5050` provides live API
- Real-time data generation and training
- Dynamic experimentation with parameters

**Batch Mode** (Production/Static)
- Pre-generated snapshots.json served as static file
- No backend required for deployment
- Automatic fallback if server unavailable

**Automatic Fallback**
- Try server mode first (`/api/snapshots`)
- Fall back to batch mode (`/snapshots.json`)
- Graceful error messages with setup instructions

## User Capabilities

### Learning Path

**Start with Welcome Introduction**
- Overview of application purpose and interface
- Quick orientation to key components
- "Getting Started" checklist
- Dismiss when ready to explore

**Follow 8-Chapter Tutorial**
- Progressive learning from basics to advanced
- Each chapter builds on previous concepts
- Estimated 30-40 minutes total read time
- Can pause and resume at any chapter

**Look Up Terms in Glossary**
- Search by term name or definition
- Jump between related terms
- Three levels of explanation depth
- Mathematical formulas for precision

**Find Answers in FAQ**
- Browse by category or search
- Common misconceptions addressed
- Troubleshooting tips
- Best practices for exploration

**Adjust Reading Level**
- Start at beginner for high-level intuition
- Progress to intermediate for more detail
- Advance to advanced for mathematical rigor
- Content adapts automatically

### Experimentation

**Change Learning Rate**
- Test small values (0.001) for slow, stable convergence
- Try large values (0.5) to observe instability
- Find optimal range for your dataset
- Watch convergence speed vs. stability tradeoff

**Generate Random Datasets**
- Vary noise level: clean (0.1) to noisy (2.0)
- Change true slope to see different target values
- Adjust x-range for different scales
- Use seeds to reproduce interesting cases

**Write Custom Data Functions**
- Test non-linear patterns (quadratic, exponential)
- Create challenging edge cases (outliers, clusters)
- Explore where linear models succeed/fail
- Immediate feedback on custom scenarios

**Compare Model Fits Visually**
- Initial fit: see how bad the starting guess is
- Current fit: track improvement in real-time
- Final fit: understand the optimal solution
- Residual visualization shows per-point errors

**Step Through Training**
- Pause at any iteration to inspect state
- Move forward/backward to compare consecutive steps
- Watch individual point contributions change
- Correlate gradient with parameter movement

### Exploration

**Click Plot Points**
- Jump directly to any training step
- Inspect exact state at that iteration
- Useful for investigating specific phenomena (e.g., loss spike)

**Hover for Tooltips**
- Exact numeric values without cluttering display
- Available on all plot elements
- Shows coordinates, loss, gradients, predictions

**Inspect Per-Point Calculations**
- Expand point details table
- See x, y_true, y_pred for each data point
- Understand per-point loss and gradient contributions
- Connect individual points to aggregate metrics

**View Formula Expansions**
- Abstract formulas on left
- Substituted values on right
- Bridge mathematical notation to concrete numbers
- Three cards cover full training pipeline

**Track Progress Metrics**
- Δw from start: how far parameter has moved
- Δloss from start: total loss reduction achieved
- Progress percentage: how close to final step
- Remaining distance: Δw and Δloss to target

## Quick Reference

### Getting Started
1. Start Go server: `go run . --server`
2. Start frontend: `cd js && pnpm dev`
3. Open browser to `localhost:5005`
4. Click "Welcome" button to see introduction
5. Click "Help" to open educational sidebar

### Understanding Training
- **Watch loss plot**: should decrease consistently
- **Monitor parameter plot**: converges to true value (2.0 by default)
- **Check gradient**: magnitude shrinks near optimum
- **View fitted line**: gets closer to true relationship
- **Use playback**: see full training animation

### Generating Custom Data
1. Scroll to "Training Data & Objective" section
2. Click "Change Dataset" to expand options
3. Choose "Random Data" or "JS Function" tab
4. Configure parameters (slope, noise, learning rate, etc.)
5. Click "Generate & Train"
6. Watch new training run with your data

### Using Educational Features
- **Tutorial**: Start at Chapter 1, read sequentially
- **Glossary**: Search for unfamiliar terms, explore related concepts
- **FAQ**: Use search or browse categories
- **Level adjustment**: Footer of sidebar, change as needed
- **Formula cards**: Scroll to "Interactive Formulas" section

### Troubleshooting
- **No data loaded**: Check console for error, verify server is running
- **Training diverges**: Reduce learning rate (try 0.01)
- **Slow performance**: Reduce number of training steps
- **Can't see plots**: Try different layout mode (1col/2col/3col)

## Feature Matrix

| Feature Category      | Capabilities                                                                                                          |
|-----------------------|-----------------------------------------------------------------------------------------------------------------------|
| **Visualizations**    | 6 plot types (loss/time, param trajectory, gradient, scatter, residuals, comparison), multi-layout, responsive design |
| **Controls**          | 5 navigation modes (play/pause, step fwd/back, jump start/end, slider), keyboard shortcuts, auto-pause                |
| **Data Input**        | 2 generation methods (random, JS function), 7 configurable parameters (slope, noise, range, lr, init, steps, seed)    |
| **Education**         | Tutorial (8 chapters, 30-40min), Glossary (8 terms, 3 detail levels), FAQ (20+ Q&A, 4 categories)                     |
| **Pedagogical Tools** | 3 formula overlays with live substitution, 4-stage breakdown, per-point inspection, visual annotations                |
| **Backend API**       | 3 REST endpoints, flexible training config, explicit gradient computation, complete snapshot serialization            |
| **Interactivity**     | Click-to-navigate plots, hover tooltips, smooth 300ms animations, tweened motion, instant feedback                    |
| **Personalization**   | 3 user levels (beginner/intermediate/advanced), 4 layout modes, localStorage persistence for 4 settings               |
| **Content**           | 8 tutorial chapters, 8 glossary terms, 20+ FAQ items, formula cards, progress tracking, related term navigation       |
| **Deployment**        | 2 loading modes (server/batch), automatic fallback, static export capability, no backend required for batch mode      |

---

**Total Word Count**: ~1,950 words
