package main

import (
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/errors"
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// const (
// 	logsDir = "logs"
// )
//
// type Config struct {
// 	PublicKey    string
// 	PrivateKey   string
// 	ProjectID    string
// 	ClusterName  string
// 	BackupID     string
// 	ScenarioType string
// }

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}

	secrets, cfg, err := config.LoadAll("configs/config.json")
	if err != nil {
		errors.ExitWithError("Failed to load configuration", err)
	}

	client, err := auth.NewClient(cfg, secrets)
	if err != nil {
		errors.ExitWithError("Failed to initialize authentication client", err)
	}

	// Parse command line flags for DR configuration
	drCfg := parseFlags()

	ctx := context.Background()
	setupLogging()

	// Execute the requested disaster recovery scenario
	switch drCfg.ScenarioType {
	case "regional-outage":
		handleRegionalOutage(ctx, client, drCfg)
	case "cloud-provider-outage":
		handleCloudProviderOutage(ctx, client, drCfg)
	case "restore-data":
		handleDataRestoration(ctx, client, drCfg)
	default:
		log.Fatalf("Unknown scenario type: %s", drCfg.ScenarioType)
	}
}

func handleRegionalOutage(ctx context.Context, sdk *admin.APIClient, cfg Config) {
	log.Println("Handling regional outage by adding nodes to unaffected regions...")

	// 1. Get current cluster configuration
	cluster, _, err := sdk.ClustersApi.GetCluster(ctx, cfg.ProjectID, cfg.ClusterName).Execute()
	if err != nil {
		log.Fatalf("Failed to get cluster details: %v", err)
	}

	// 2. Identify regions that are currently not in use and add a node
	var newRegions []admin.ReplicationSpec20240805
	foundRegion := false

	if cluster.ReplicationSpecs != nil {
		for _, region := range *cluster.ReplicationSpecs {
			newRegions = append(newRegions, region)
			if *region.ZoneName == "EU_WEST_1" {
				foundRegion = true
			}
		}
	}

	if !foundRegion {
		// Add a new region that's unaffected by the outage
		priority := int64(5)
		electableNodes := int64(1)
		readOnlyNodes := int64(0)
		analyticsNodes := int64(0)
		regionName := "EU_WEST_1"

		newRegion := admin.ReplicationSpec20240805{
			RegionName:     &regionName,
			Priority:       &priority,
			RegionConfigs:  &admin.CloudRegionConfig20240805{} & electableNodes,
			ReadOnlyNodes:  &readOnlyNodes,
			AnalyticsNodes: &analyticsNodes,
		}
		newRegions = append(newRegions, newRegion)
	}

	// 3. Update cluster with new regions
	updateRequest := admin.AdvancedClusterDescriptionV2{
		ReplicationSpecs: &newRegions,
	}

	_, _, err = sdk.ClustersApi.UpdateCluster(ctx, cfg.ProjectID, cfg.ClusterName, &updateRequest).Execute()
	if err != nil {
		log.Fatalf("Failed to update cluster: %v", err)
	}

	log.Println("Successfully added nodes to unaffected regions")
}

func handleCloudProviderOutage(ctx context.Context, sdk *admin.APIClient, cfg Config) {
	log.Println("Handling cloud provider outage...")

	// 1. Get current cluster configuration
	sourceCluster, _, err := sdk.ClustersApi.GetCluster(ctx, cfg.ProjectID, cfg.ClusterName).Execute()
	if err != nil {
		log.Fatalf("Failed to get source cluster details: %v", err)
	}

	// 2. Create new cluster on alternative cloud provider
	newClusterName := cfg.ClusterName + "-recovery"
	clusterType := "REPLICASET"
	providerName := "GCP" // Switch from AWS to GCP or vice versa

	newCluster := admin.AdvancedClusterDescriptionV2{
		Name:                &newClusterName,
		ClusterType:         &clusterType,
		ProviderName:        &providerName,
		DiskSizeGB:          sourceCluster.DiskSizeGB,
		MongoDBMajorVersion: sourceCluster.MongoDBMajorVersion,
	}

	// Configure replica set based on original configuration
	// (simplified for example - would need more configuration in practice)

	_, _, err = sdk.ClustersApi.CreateCluster(ctx, cfg.ProjectID, &newCluster).Execute()
	if err != nil {
		log.Fatalf("Failed to create recovery cluster: %v", err)
	}

	log.Printf("Created recovery cluster: %s", newClusterName)

	// 3. Wait for cluster to be ready
	log.Println("Waiting for cluster to become available...")
	waitForClusterReady(ctx, sdk, cfg.ProjectID, newClusterName)

	// 4. Restore the most recent snapshot to the new cluster
	log.Println("Restoring backup to recovery cluster...")
	restoreRequest := admin.DiskBackupSnapshotRestoreJob{
		TargetClusterName: &newClusterName,
		SnapshotId:        &cfg.BackupID,
	}

	_, _, err = sdk.CloudBackupsApi.CreateBackupRestoreJob(ctx, cfg.ProjectID, cfg.ClusterName, &restoreRequest).Execute()
	if err != nil {
		log.Fatalf("Failed to restore backup: %v", err)
	}

	log.Printf("Successfully initiated restore to cluster %s from backup %s", newClusterName, cfg.BackupID)
	log.Println("Once restore is complete, update your application connection strings to point to the new cluster")
}

func handleDataRestoration(ctx context.Context, sdk *admin.APIClient, cfg Config) {
	log.Println("Handling data restoration after accidental deletion...")

	// Restore from point-in-time backup
	restoreRequest := admin.DiskBackupSnapshotRestoreJob{
		TargetClusterName: &cfg.ClusterName,
		SnapshotId:        &cfg.BackupID,
	}

	_, _, err := sdk.CloudBackupsApi.CreateBackupRestoreJob(ctx, cfg.ProjectID, cfg.ClusterName, &restoreRequest).Execute()
	if err != nil {
		log.Fatalf("Failed to restore backup: %v", err)
	}

	log.Printf("Successfully initiated restore to cluster %s from backup %s", cfg.ClusterName, cfg.BackupID)
	log.Println("After restoration, verify data integrity and reimport any data collected since the backup")
}

// Helper function to wait for cluster to be ready
func waitForClusterReady(ctx context.Context, sdk *admin.APIClient, projectID, clusterName string) {
	for {
		cluster, _, err := sdk.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
		if err != nil {
			log.Printf("Error checking cluster status: %v", err)
		} else if cluster.StateName != nil && *cluster.StateName == "IDLE" {
			log.Println("Cluster is ready")
			return
		}

		log.Printf("Cluster status: %s. Waiting 30 seconds...", *cluster.StateName)
		time.Sleep(30 * time.Second)
	}
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.PublicKey, "public-key", os.Getenv("ATLAS_PUBLIC_KEY"), "MongoDB Atlas public API key")
	flag.StringVar(&cfg.PrivateKey, "private-key", os.Getenv("ATLAS_PRIVATE_KEY"), "MongoDB Atlas private API key")
	flag.StringVar(&cfg.ProjectID, "project-id", "", "MongoDB Atlas project ID")
	flag.StringVar(&cfg.ClusterName, "cluster-name", "", "MongoDB Atlas cluster name")
	flag.StringVar(&cfg.BackupID, "backup-id", "", "MongoDB Atlas backup snapshot ID (for restore operations)")
	flag.StringVar(&cfg.ScenarioType, "scenario", "", "Disaster recovery scenario type: regional-outage, cloud-provider-outage, restore-data")

	flag.Parse()

	// Validate required parameters
	if cfg.PublicKey == "" || cfg.PrivateKey == "" || cfg.ProjectID == "" || cfg.ClusterName == "" || cfg.ScenarioType == "" {
		flag.Usage()
		os.Exit(1)
	}

	return cfg
}

func setupLogging() {
	// Ensure logs directory exists
	defaultDir := os.Getenv("ATLAS_DOWNLOADS_DIR")
	logDir := logsDir
	if defaultDir != "" {
		logDir = fmt.Sprintf("%s/%s", defaultDir, logsDir)
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Set up logging to file
	logFile := fmt.Sprintf("%s/disaster_recovery_%s.log", logDir, time.Now().Format("20060102_150405"))
	f, err := os.Create(logFile)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	log.SetOutput(f)
	log.Printf("Starting disaster recovery script at %s", time.Now().Format(time.RFC3339))
}
