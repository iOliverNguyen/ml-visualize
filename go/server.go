package core

// DEPRECATED: This server is no longer needed for the frontend.
// The frontend now runs training client-side using JavaScript (see /js/src/shared/training.ts).
// This file is kept for reference and potential future API needs.
//
// To run the frontend without this server:
//   cd js && pnpm dev
//
// The server is only needed if you want to use the API endpoints:
//   - GET  /api/snapshots        - Load pre-computed snapshots from file
//   - POST /api/dataset/random   - Generate random data + train on-the-fly
//   - POST /api/dataset/custom   - Train with user-provided custom data
//
// However, the frontend now has equivalent functionality client-side.

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

// Server state to store current snapshots
var (
	currentSnapshots []Snapshot
	snapshotsMutex   sync.RWMutex
)

// TrainingRequest combines dataset and training configuration
type TrainingRequest struct {
	Data   []DataPoint    `json:"data"`
	Config TrainingConfig `json:"config"`
}

// RandomDataRequest combines data generation config and training config
type RandomDataRequest struct {
	DataConfig     DataGenConfig  `json:"data_config"`
	TrainingConfig TrainingConfig `json:"training_config"`
}

// StartServer starts an HTTP server that serves training snapshots
func StartServer(addr string) error {
	// Initialize with default snapshots from file if available
	if data, err := os.ReadFile("output/snapshots.json"); err == nil {
		json.Unmarshal(data, &currentSnapshots)
	}

	mux := http.NewServeMux()

	// CORS middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}

	// GET /api/snapshots - Get current training snapshots
	mux.HandleFunc("/api/snapshots", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		snapshotsMutex.RLock()
		defer snapshotsMutex.RUnlock()

		if len(currentSnapshots) == 0 {
			// Try to read from file as fallback
			data, err := os.ReadFile("output/snapshots.json")
			if err != nil {
				http.Error(w, "Snapshots not found. Run training first.", 404)
				return
			}
			w.Write(data)
			return
		}

		json.NewEncoder(w).Encode(currentSnapshots)
	}))

	// POST /api/dataset/random - Generate random data and train
	mux.HandleFunc("/api/dataset/random", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RandomDataRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Generate random data
		data, err := GenerateRandomData(req.DataConfig)
		if err != nil {
			http.Error(w, "Failed to generate data: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validate dataset
		if err := ValidateDataset(data); err != nil {
			http.Error(w, "Invalid dataset: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Run training
		snapshots := RunTrainingWithDataset(data, req.TrainingConfig)

		// Update server state
		snapshotsMutex.Lock()
		currentSnapshots = snapshots
		snapshotsMutex.Unlock()

		// Return snapshots
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(snapshots)
	}))

	// POST /api/dataset/custom - Train with custom data
	mux.HandleFunc("/api/dataset/custom", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req TrainingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validate dataset
		if err := ValidateDataset(req.Data); err != nil {
			http.Error(w, "Invalid dataset: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Run training
		snapshots := RunTrainingWithDataset(req.Data, req.Config)

		// Update server state
		snapshotsMutex.Lock()
		currentSnapshots = snapshots
		snapshotsMutex.Unlock()

		// Return snapshots
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(snapshots)
	}))

	log.Printf("Server listening on %s", addr)
	log.Println("Endpoints:")
	log.Println("  GET  /api/snapshots        - Get current training snapshots")
	log.Println("  POST /api/dataset/random   - Generate random data and train")
	log.Println("  POST /api/dataset/custom   - Train with custom data")
	return http.ListenAndServe(addr, mux)
}
