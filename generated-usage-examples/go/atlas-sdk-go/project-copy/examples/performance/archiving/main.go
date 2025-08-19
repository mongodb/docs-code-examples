package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"atlas-sdk-go/internal/archive"
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	envName := config.Environment("production")
	configPath := "configs/config.production.json"
	secrets, cfg, err := config.LoadAll(envName, configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
	}

	client, err := auth.NewClient(ctx, cfg, secrets)
	if err != nil {
		log.Fatalf("Failed to initialize authentication client: %v", err)
	}

	projectID := cfg.ProjectID
	if projectID == "" {
		log.Fatal("Failed to find Project ID in configuration")
	}

	fmt.Printf("Starting archive analysis for project: %s\n", projectID)

	// List all clusters in the project
	clusters, _, err := client.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		log.Fatalf("Failed to list clusters: %v", err)
	}

	fmt.Printf("\nFound %d clusters to analyze\n", len(clusters.GetResults()))

	// Process each cluster
	failedArchives := 0
	totalCandidates := 0
	for _, cluster := range clusters.GetResults() {
		clusterName := cluster.GetName()
		fmt.Printf("\n=== Analyzing cluster: %s ===", clusterName)

		// Find collections suitable for archiving
		// NOTE: This function passes example database/collection names.
		// In a real production scenario, you would analyze data patterns and customize the selection logic.
		candidates := archive.CollectionsForArchiving(ctx, client, projectID, clusterName)
		totalCandidates += len(candidates)
		fmt.Printf("\nFound %d collections eligible for archiving in cluster %s\n",
			totalCandidates, clusterName)

		// Step 4: Configure online archive for each candidate collection
		for _, candidate := range candidates {
			fmt.Printf("- Configuring archive for %s.%s\n",
				candidate.DatabaseName, candidate.CollectionName)

			configureErr := archive.ConfigureOnlineArchive(ctx, client, projectID, clusterName, candidate)
			if configureErr != nil {
				fmt.Printf("  Failed to configure archive: %v\n", configureErr)
				failedArchives++
				continue
			}

			fmt.Printf("  Successfully configured online archive for %s.%s\n",
				candidate.DatabaseName, candidate.CollectionName)
		}
	}

	if failedArchives > 0 {
		fmt.Printf("\nWARNING: %d of %d archive configurations failed\n", failedArchives, totalCandidates)
	}

	fmt.Println("Archive analysis and configuration completed.")
}

