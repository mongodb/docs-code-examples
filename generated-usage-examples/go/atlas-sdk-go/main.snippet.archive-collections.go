// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"atlas-sdk-go/internal/archive"
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	envFile := ".env.production"
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: could not load %s file: %v", envFile, err)
	}

	configPath := os.Getenv("CONFIG_FILE")
	secrets, cfg, err := config.LoadAll(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	client, err := auth.NewClient(ctx, cfg, secrets)
	if err != nil {
		log.Fatalf("Failed to initialize authentication client: %v", err)
	}

	projectID := cfg.ProjectID
	if projectID == "" {
		log.Fatal("Failed to find Project ID in configuration")
	}

	fmt.Printf("Starting archive analysis for project: %s\n", projectID)

	// Get all clusters in the project
	clusters, _, err := client.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		log.Fatalf("Failed to list clusters: %v", err)
	}

	fmt.Printf("\nFound %d clusters to analyze\n", len(clusters.GetResults()))

	// Connect to each cluster and analyze collections for archiving
	failedArchives := 0
	totalCandidates := 0
	for _, cluster := range clusters.GetResults() {
		clusterName := cluster.GetName()
		fmt.Printf("\n=== Analyzing cluster: %s ===", clusterName)

		// Find collections suitable for archiving based on specific criteria.
		// NOTE: The actual implementation of this function would involve more complex logic
		// to determine which collections are eligible for archiving.
		candidates := archive.CollectionsForArchiving(ctx, client, projectID, clusterName)
		totalCandidates += len(candidates)
		fmt.Printf("\nFound %d collections eligible for archiving in cluster %s\n",
			len(candidates), clusterName)

		// Configure online archive for each candidate collection
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

