package main

import (
	"flag"
	"fmt"
	"log"

	core "github.com/iOliverNguyen/ml-viz/go"
	"github.com/iOliverNguyen/ml-viz/go/linear"
)

func main() {
	server := flag.Bool("server", false, "Start HTTP server after training")
	flag.BoolVar(server, "s", false, "Start HTTP server after training (shorthand)")
	generateCases := flag.Bool("generate-cases", false, "Generate pre-computed Phase 1 training cases")
	generateCases2 := flag.Bool("generate-cases-phase2", false, "Generate pre-computed Phase 2 training cases")
	flag.Parse()

	// Check if generate-cases-phase2 command was requested
	if *generateCases2 {
		fmt.Println("Generating Phase 2 pre-computed training cases...")
		outputDir := "js/public/cases-phase2"
		if err := linear.GenerateCases2D(outputDir); err != nil {
			log.Fatalf("Failed to generate Phase 2 cases: %v", err)
		}
		fmt.Println("✓ Phase 2 cases generated successfully!")
		return
	}

	// Check if generate-cases command was requested
	if *generateCases {
		fmt.Println("Generating Phase 1 pre-computed training cases...")
		if err := core.GenerateCases(); err != nil {
			log.Fatalf("Failed to generate cases: %v", err)
		}
		return
	}

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
		if err := core.StartServer(":5050"); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
