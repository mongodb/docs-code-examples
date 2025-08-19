// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
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

	_ = godotenv.Load()

	envName := config.Environment("")        // Use empty string to load from environment variables
	configPath := "configs/config.test.json" // Optional explicit config file path; if empty, uses environment-based path
	secrets, cfg, err := config.LoadAll(envName, configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
	}

	client, err := auth.NewClient(ctx, &cfg, &secrets) // Pass pointers
	if err != nil {
		log.Fatalf("Failed to initialize authentication client: %v", err)
	}

	projectID := cfg.ProjectID
	if projectID == "" {
		log.Fatal("Failed to find Project ID in configuration")
	}

	fmt.Printf("Starting archive analysis for project: %s\n", projectID)

	// Step 1: List all clusters in the project
	clusters, _, err := client.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		log.Fatalf("Failed to list clusters: %v", err)
	}

	fmt.Printf("Found %d clusters to analyze\n", len(clusters.GetResults()))

	// Step 2: Process each cluster
	failedArchives := 0
	for _, cluster := range clusters.GetResults() {
		clusterName := cluster.GetName()
		fmt.Printf("Analyzing cluster: %s\n", clusterName)

		// Step 3: Find collections suitable for archiving
		// NOTE: This example passes example database/collections.
		// In a real production scenario, you would customize the collection analysis logic to match your specific data patterns.
		candidates := archive.CollectionsForArchiving(ctx, client, projectID, clusterName)
		fmt.Printf("\nFound %d collections eligible for archiving in cluster %s",
			len(candidates), clusterName)

		// Step 4: Configure online archive for each candidate collection
		for _, candidate := range candidates {
			fmt.Printf("\nConfiguring archive for %s.%s ",
				candidate.DatabaseName, candidate.CollectionName)

			configureErr := archive.ConfigureOnlineArchive(ctx, client, projectID, clusterName, candidate)
			if configureErr != nil {
				fmt.Printf("\nFailed to configure archive: %v", configureErr)
				failedArchives++
				continue
			}

			fmt.Printf("\nSuccessfully configured online archive for %s.%s ",
				candidate.DatabaseName, candidate.CollectionName)
		}
	}

	if failedArchives > 0 {
		fmt.Printf("Warning: %d archive configurations failed\n", failedArchives)
	}

	fmt.Println("Archive analysis and configuration completed")
}

