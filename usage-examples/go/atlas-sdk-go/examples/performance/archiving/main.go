// :snippet-start: archive-collections
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"atlas-sdk-go/internal/archive"
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/errors"
	"context"
	"fmt"
	"time"
)

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
		// NOTE: In a production scenario, you would customize the collection analysis logic to match your specific data patterns.
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

// :snippet-end: [archive-collections]
// :state-remove-start: copy
// NOTE: INTERNAL
// ** OUTPUT EXAMPLE **
//
// Configuration loaded successfully: env=production, baseURL=https://cloud.mongodb.com, orgID=5bfda007553855125605a5cf
// Starting archive analysis for project: 5f60207f14dfb25d24511201
// Found 2 clusters to analyze
// Analyzing cluster: Cluster0
// Found 2 collections eligible for archiving in cluster Cluster0
// Configuring archive for sample_analytics.transactions
// Successfully configured online archive for sample_analytics.transactions
// Configuring archive for sample_analytics.users
// Successfully configured online archive for sample_analytics.users
// Analyzing cluster: Cluster1
// Found 1 collections eligible for archiving in cluster Cluster1
// Configuring archive for sample_analytics.orders
//  Failed to configure archive: validate archive candidate for sample_analytics.transactions: date field transaction_date must be included in partition fields
//  Configuring archive for sample_logs.application_logs
//  Failed to configure archive: validate archive candidate for sample_logs.application_logs: date field timestamp must be included in partition fields
//  Warning: 2 archive configurations failed
//  Archive analysis and configuration completed
// :state-remove-end: [copy]
