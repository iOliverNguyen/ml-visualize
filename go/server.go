package core

import (
	"log"
	"net/http"
	"os"
)

// StartServer starts an HTTP server that serves training snapshots
func StartServer(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/snapshots", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers for Vite dev server
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		// Read snapshots file
		data, err := os.ReadFile("output/snapshots.json")
		if err != nil {
			http.Error(w, "Snapshots not found. Run training first.", 404)
			return
		}

		w.Write(data)
	})

	log.Printf("Server listening on %s", addr)
	log.Println("Serving snapshots at /api/snapshots")
	return http.ListenAndServe(addr, mux)
}
