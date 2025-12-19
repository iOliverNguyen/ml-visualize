package neuron

import "math"

// Train performs gradient descent training and captures snapshots at each step
func Train(dataset []DataPoint2DNeuron, initParams NeuronParams, config TrainingConfig) []NeuronSnapshot {
	params := NeuronParams{
		W: make([]float64, len(initParams.W)),
		B: initParams.B,
	};
	copy(params.W, initParams.W);

	snapshots := make([]NeuronSnapshot, config.NumSteps);

	for step := 0; step < config.NumSteps; step++ {
		// Compute gradients and per-point details
		grads, pointDetails := ComputeGradients(dataset, params, config.Activation);

		// Compute average metrics across dataset
		avgZ := ComputeAvgZ(dataset, params);
		avgA := ComputeAvgA(dataset, params, config.Activation);
		avgLoss := ComputeAvgLoss(dataset, params, config.Activation);
		avgDerivative := ComputeAvgDerivative(dataset, params, config.Activation);

		// Compute chain rule breakdown
		chainRuleBreakdown := ComputeChainRuleBreakdown(dataset, params, config.Activation);

		// Compute average dL/da and dL/dz
		avgDLda := 0.0;
		avgDLdz := 0.0;
		for _, point := range pointDetails {
			avgDLda += point.DLda;
			avgDLdz += point.DLdz;
		}
		avgDLda /= float64(len(pointDetails));
		avgDLdz /= float64(len(pointDetails));

		// Check if in saturation zone
		inSaturationZone := IsSaturated(avgDerivative);

		// Compute gradient magnitude
		gradMag := GradientMagnitude(grads);

		// Compute updates
		updateW := make([]float64, len(params.W));
		for i := 0; i < len(params.W); i++ {
			updateW[i] = -config.LearningRate * grads.GradW[i];
		}
		updateB := -config.LearningRate * grads.GradB;

		// Compute step size (magnitude of update vector)
		stepSize := 0.0;
		for _, u := range updateW {
			stepSize += u * u;
		}
		stepSize += updateB * updateB;
		stepSize = math.Sqrt(stepSize);

		// Create update details
		updateComponents := UpdateDetailsNeuron{
			LearningRate:  config.LearningRate,
			GradMagnitude: gradMag,
			UpdateW:       updateW,
			UpdateB:       updateB,
			StepSize:      stepSize,
		};

		// Create snapshot before update
		snapshots[step] = NeuronSnapshot{
			Step:               step,
			Params:             NeuronParams{W: append([]float64(nil), params.W...), B: params.B},
			Grads:              grads,
			Z:                  avgZ,
			A:                  avgA,
			DLdz:               avgDLdz,
			DLda:               avgDLda,
			LocalDerivative:    avgDerivative,
			Activation:         config.Activation,
			InSaturationZone:   inSaturationZone,
			Loss:               avgLoss,
			PointDetails:       pointDetails,
			UpdateComponents:   updateComponents,
			ChainRuleBreakdown: chainRuleBreakdown,
		};

		// Update parameters
		for i := 0; i < len(params.W); i++ {
			params.W[i] += updateW[i];
		}
		params.B += updateB;
	}

	return snapshots;
}
