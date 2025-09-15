// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/data/recovery"
	"atlas-sdk-go/internal/typeutils"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

const (
	scenarioRegionalOutage = "regional-outage"
	scenarioDataDeletion   = "data-deletion"
)

func main() {
	envFile := ".env.production"
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: could not load %s file: %v", envFile, err)
	}

	secrets, cfg, err := config.LoadAllFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	client, err := auth.NewClient(ctx, cfg, secrets)
	if err != nil {
		log.Fatalf("Failed to initialize authentication client: %v", err)
	}

	opts, err := recovery.LoadDROptionsFromEnv(cfg.ProjectID)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	fmt.Printf("Starting disaster recovery scenario: %s\nProject: %s\nCluster: %s\n", opts.Scenario, opts.ProjectID, opts.ClusterName)

	if opts.DryRun {
		fmt.Println("DRY RUN: no write operations will be performed")
	}

	var summary string
	var opErr error

	switch opts.Scenario {
	case scenarioRegionalOutage:
		summary, opErr = simulateRegionalOutage(ctx, client, opts)
	case scenarioDataDeletion:
		summary, opErr = executeDataDeletionRestore(ctx, client, opts)
	default:
		opErr = fmt.Errorf("unsupported DR_SCENARIO '%s'", opts.Scenario)
	}

	if opErr != nil {
		log.Fatalf("Scenario failed: %v", opErr)
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println(summary)
	fmt.Println("Disaster recovery procedure completed.")
}

// executeDataDeletionRestore initiates a restore job for a specified snapshot in a MongoDB Atlas cluster.
func executeDataDeletionRestore(ctx context.Context, client *admin.APIClient, o recovery.DrOptions) (string, error) {
	job := admin.DiskBackupSnapshotRestoreJob{SnapshotId: &o.SnapshotID, TargetClusterName: &o.ClusterName}
	if o.DryRun {
		return fmt.Sprintf("(dry-run) Would submit restore job for snapshot %s", o.SnapshotID), nil
	}
	_, _, err := client.CloudBackupsApi.CreateBackupRestoreJob(ctx, o.ProjectID, o.ClusterName, &job).Execute()
	if err != nil {
		return "", fmt.Errorf("create restore job: %w", err)
	}
	return fmt.Sprintf("Restore job submitted for snapshot %s", o.SnapshotID), nil
}

// simulateRegionalOutage modifies the electable node count in a target region for a MongoDB Atlas cluster.
func simulateRegionalOutage(ctx context.Context, client *admin.APIClient, o recovery.DrOptions) (string, error) {
	cluster, _, err := client.ClustersApi.GetCluster(ctx, o.ProjectID, o.ClusterName).Execute()
	if err != nil {
		return "", fmt.Errorf("get cluster: %w", err)
	}
	if !cluster.HasReplicationSpecs() {
		return "", fmt.Errorf("cluster has no replication specs")
	}
	repl := cluster.GetReplicationSpecs()
	addedNodes, foundTarget := recovery.AddElectableNodesToRegion(repl, o.TargetRegion, o.AddNodes)
	if !foundTarget {
		return "", fmt.Errorf("target region '%s' not found in replication specs", o.TargetRegion)
	}
	zeroedRegions := 0
	if o.OutageRegion != "" {
		zeroedRegions = recovery.ZeroElectableNodesInRegion(repl, o.OutageRegion)
	}
	payload := admin.NewClusterDescription20240805()
	payload.SetReplicationSpecs(repl)
	if o.DryRun {
		return fmt.Sprintf("(dry-run) Would add %d electable nodes to %s%s", addedNodes, o.TargetRegion, typeutils.SuffixZeroed(zeroedRegions, o.OutageRegion)), nil
	}
	_, _, err = client.ClustersApi.UpdateCluster(ctx, o.ProjectID, o.ClusterName, payload).Execute()
	if err != nil {
		return "", fmt.Errorf("update cluster: %w", err)
	}
	return fmt.Sprintf("Added %d electable nodes to %s%s", addedNodes, o.TargetRegion, typeutils.SuffixZeroed(zeroedRegions, o.OutageRegion)), nil
}

