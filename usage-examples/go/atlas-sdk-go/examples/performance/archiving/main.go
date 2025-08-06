package main

import (
	"atlas-sdk-go/internal/archive"
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/errors"
	"context"
	"log"
	"time"
)

// This program demonstrates an automated approach to:
// 1. Discover all clusters in an Atlas project
// 2. Analyze collections within each cluster for archiving candidates
// 3. Configure Online Archive for eligible collections
//
// In a production scenario, you would customize the collection analysis
// logic in CollectionsForArchiving() to match your specific data patterns.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Load application context with configuration and secrets for the specified environment
	explicitEnv := "internal"
	appCtx, err := config.LoadAppContextWithContext(ctx, explicitEnv, false)
	if err != nil {
		errors.ExitWithError("Failed to load configuration", err)
	}

	client, err := auth.NewClient(appCtx.Config, appCtx.Secrets)
	if err != nil {
		errors.ExitWithError("Failed to initialize Atlas client", err)
	}

	projectID := appCtx.Config.ProjectID
	if projectID == "" {
		errors.ExitWithError("Project ID not found in configuration", nil)
	}

	log.Println("Starting archive analysis for project:", projectID)

	// Step 1: List all clusters in the project
	clusters, _, err := client.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		errors.ExitWithError("Failed to list clusters", err)
	}

	log.Printf("Found %d clusters to analyze", len(clusters.GetResults()))

	// Step 2: Process each cluster
	failedArchives := 0
	for _, cluster := range clusters.GetResults() {
		clusterName := cluster.GetName()
		log.Printf("Analyzing cluster: %s", clusterName)

		// Step 3: Find collections suitable for archiving
		// Note: Partition fields are ordered by query frequency - most frequently
		// queried field should be first for optimal query performance against
		// archived data. This significantly impacts cost and performance.
		candidates := archive.CollectionsForArchiving(ctx, client, projectID, clusterName)
		log.Printf("Found %d collections eligible for archiving in cluster %s",
			len(candidates), clusterName)

		// Step 4: Configure online archive for each candidate collection
		for _, candidate := range candidates {
			log.Printf("Configuring archive for %s.%s",
				candidate.DatabaseName, candidate.CollectionName)

			configureErr := archive.ConfigureOnlineArchive(ctx, client, projectID, clusterName, candidate)
			if configureErr != nil {
				log.Printf("Failed to configure archive: %v", configureErr)
				failedArchives++
				continue
			}

			log.Printf("Successfully configured online archive for %s.%s",
				candidate.DatabaseName, candidate.CollectionName)
		}
	}

	if failedArchives > 0 {
		log.Printf("Warning: %d archive configurations failed", failedArchives)
	}

	log.Println("Archive analysis and configuration completed")
}
