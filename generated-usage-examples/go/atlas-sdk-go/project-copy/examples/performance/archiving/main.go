package main

import (
	"context"
	"fmt"
	"time"

	"atlas-sdk-go/internal/archive"
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/errors"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Load application context with configuration and secrets for the environment
	explicitEnv := "production"
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

	fmt.Println("Starting archive analysis for project:", projectID)

	// Step 1: List all clusters in the project
	clusters, _, err := client.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		errors.ExitWithError("Failed to list clusters", err)
	}

	fmt.Printf("Found %d clusters to analyze", len(clusters.GetResults()))

	// Step 2: Process each cluster
	failedArchives := 0
	for _, cluster := range clusters.GetResults() {
		clusterName := cluster.GetName()
		fmt.Printf("Analyzing cluster: %s", clusterName)

		// Step 3: Find collections suitable for archiving
		// NOTE: In a real production scenario, you would customize the collection analysis logic to match your specific data patterns.
		candidates := archive.CollectionsForArchiving(ctx, client, projectID, clusterName)
		fmt.Printf("Found %d collections eligible for archiving in cluster %s",
			len(candidates), clusterName)

		// Step 4: Configure online archive for each candidate collection
		for _, candidate := range candidates {
			fmt.Printf("Configuring archive for %s.%s",
				candidate.DatabaseName, candidate.CollectionName)

			configureErr := archive.ConfigureOnlineArchive(ctx, client, projectID, clusterName, candidate)
			if configureErr != nil {
				fmt.Printf("Failed to configure archive: %v", configureErr)
				failedArchives++
				continue
			}

			fmt.Printf("Successfully configured online archive for %s.%s",
				candidate.DatabaseName, candidate.CollectionName)
		}
	}

	if failedArchives > 0 {
		fmt.Printf("Warning: %d archive configurations failed", failedArchives)
	}

	fmt.Println("Archive analysis and configuration completed")
}

