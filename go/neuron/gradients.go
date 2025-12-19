package neuron

// ComputeGradients computes gradients using the chain rule for all parameters
// Chain rule: ∂L/∂w = ∂L/∂a × ∂a/∂z × ∂z/∂w
//             ∂L/∂b = ∂L/∂a × ∂a/∂z × ∂z/∂b
func ComputeGradients(dataset []DataPoint2DNeuron, params NeuronParams, activation string) (NeuronGrads, []PointSnapshotNeuron) {
	numPoints := len(dataset);
	numFeatures := len(params.W);

	// Initialize accumulators for gradients
	gradW := make([]float64, numFeatures);
	gradB := 0.0;

	// Per-point details for visualization
	pointDetails := make([]PointSnapshotNeuron, numPoints);

	// Iterate through each data point
	for i, point := range dataset {
		// Forward pass
		z, a := Forward(point.X, params, activation);

		// Compute loss for this point
		loss := Loss(a, point.Y);

		// Backward pass - Chain rule components:
		// 1. ∂L/∂a = 2(a - y_true)
		dL_da := 2.0 * (a - point.Y);

		// 2. ∂a/∂z = σ'(z)
		da_dz := ApplyActivationDerivative(z, activation);

		// 3. ∂L/∂z = ∂L/∂a × ∂a/∂z
		dL_dz := dL_da * da_dz;

		// 4. ∂z/∂w = x (for each weight)
		//    ∂z/∂b = 1
		dz_dw := make([]float64, numFeatures);
		for j := 0; j < numFeatures; j++ {
			dz_dw[j] = point.X[j];
		}

		// 5. ∂L/∂w = ∂L/∂z × ∂z/∂w = ∂L/∂z × x
		//    ∂L/∂b = ∂L/∂z × ∂z/∂b = ∂L/∂z × 1
		dL_dw := make([]float64, numFeatures);
		for j := 0; j < numFeatures; j++ {
			dL_dw[j] = dL_dz * point.X[j];
			gradW[j] += dL_dw[j];
		}
		dL_db := dL_dz * 1.0;
		gradB += dL_db;

		// Check if this point is in saturation
		inSaturation := IsSaturated(da_dz);

		// Store per-point details
		pointDetails[i] = PointSnapshotNeuron{
			Index:        i,
			X:            point.X,
			YTrue:        point.Y,
			Z:            z,
			A:            a,
			Loss:         loss,
			DLda:         dL_da,
			DaDz:         da_dz,
			DLdz:         dL_dz,
			DzDw:         dz_dw,
			DLdw:         dL_dw,
			DLdb:         dL_db,
			InSaturation: inSaturation,
		};
	}

	// Average gradients over dataset
	for j := 0; j < numFeatures; j++ {
		gradW[j] /= float64(numPoints);
	}
	gradB /= float64(numPoints);

	grads := NeuronGrads{
		GradW: gradW,
		GradB: gradB,
	};

	return grads, pointDetails;
}

// ComputeChainRuleBreakdown computes the chain rule breakdown for visualization
// Shows: dL/da × da/dz × dz/dparam = dL/dparam for each parameter
func ComputeChainRuleBreakdown(dataset []DataPoint2DNeuron, params NeuronParams, activation string) ChainRuleViz {
	numPoints := len(dataset);
	numFeatures := len(params.W);

	// Average chain rule components across dataset
	avgDLda := 0.0;
	avgDaDz := 0.0;
	avgDzDw := make([]float64, numFeatures);

	for _, point := range dataset {
		// Forward pass
		z, a := Forward(point.X, params, activation);

		// Chain rule components
		dL_da := 2.0 * (a - point.Y);
		da_dz := ApplyActivationDerivative(z, activation);

		avgDLda += dL_da;
		avgDaDz += da_dz;

		for j := 0; j < numFeatures; j++ {
			avgDzDw[j] += point.X[j];
		}
	}

	avgDLda /= float64(numPoints);
	avgDaDz /= float64(numPoints);
	for j := 0; j < numFeatures; j++ {
		avgDzDw[j] /= float64(numPoints);
	}

	// Create components for each parameter
	components := make([]ChainRuleComponent, numFeatures+1); // w1, w2, b

	// Components for weights
	for j := 0; j < numFeatures; j++ {
		paramName := "w" + string(rune('1'+j));
		dL_dparam := avgDLda * avgDaDz * avgDzDw[j];

		components[j] = ChainRuleComponent{
			ParamName: paramName,
			DLda:      avgDLda,
			DaDz:      avgDaDz,
			DzDparam:  avgDzDw[j],
			DLdparam:  dL_dparam,
		};
	}

	// Component for bias
	dL_db := avgDLda * avgDaDz * 1.0;
	components[numFeatures] = ChainRuleComponent{
		ParamName: "b",
		DLda:      avgDLda,
		DaDz:      avgDaDz,
		DzDparam:  1.0,
		DLdparam:  dL_db,
	};

	return ChainRuleViz{
		Components: components,
	};
}
