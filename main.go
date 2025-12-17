package main

import (
	"flag"
	"fmt"
	"log"

	core "github.com/iOliverNguyen/ml-viz/go"
)

func main() {
	server := flag.Bool("server", false, "Start HTTP server after training")
	flag.BoolVar(server, "s", false, "Start HTTP server after training (shorthand)")
	flag.Parse()

	// Run training and get snapshots
	fmt.Println("Running training...")
	snapshots := core.RunTraining()

	// Write to output/ directory
	filepath := "output/snapshots.json"
	if err := core.WriteSnapshots(snapshots, filepath); err != nil {
		log.Fatalf("Failed to write snapshots: %v", err)
	}

	fmt.Printf("\nSnapshots written to %s (%d steps)\n", filepath, len(snapshots))
	fmt.Println("\nVisualization options:")
	fmt.Println("  1. Server mode: go run . --server")
	fmt.Println("  2. Batch mode:  cp output/snapshots.json js/public/")
	fmt.Println()

	// Start server if requested
	if *server {
		fmt.Println("Starting HTTP server...")
		if err := core.StartServer(":8080"); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
