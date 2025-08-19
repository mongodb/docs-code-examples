// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/metrics"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func main() {
	_ = godotenv.Load() // or godotenv.Load(".env.development")

	ctx := context.Background()
	envName := config.Environment("development")    // Cast string to config.Environment
	configPath := "configs/config.development.json" // Optional explicit config file path; if empty, uses environment-based path
	secrets, cfg, err := config.LoadAll(envName, configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
	}

	client, err := auth.NewClient(ctx, cfg, secrets) // Pass pointers
	if err != nil {
		log.Fatalf("Failed to initialize authentication client: %v", err)
	}

	// Fetch disk metrics with the provided parameters
	p := &admin.GetDiskMeasurementsApiParams{
		GroupId:       cfg.ProjectID,
		ProcessId:     cfg.ProcessID,
		PartitionName: "data",
		M:             &[]string{"DISK_PARTITION_SPACE_FREE", "DISK_PARTITION_SPACE_USED"},
		Granularity:   admin.PtrString("P1D"),
		Period:        admin.PtrString("P1D"),
	}
	view, err := metrics.FetchDiskMetrics(ctx, client.MonitoringAndLogsApi, p)
	if err != nil {
		log.Fatalf("Failed to fetch disk metrics: %v", err)
	}

	// Output metrics
	out, err := json.MarshalIndent(view, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format metrics data: %v", err)
	}
	fmt.Println(string(out))
}

