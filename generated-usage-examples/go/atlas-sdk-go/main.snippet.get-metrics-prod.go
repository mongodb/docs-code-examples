// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal/metrics"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	ctx := context.Background()
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

	// Fetch process metrics with the provided parameters
	p := &admin.GetHostMeasurementsApiParams{
		GroupId:   cfg.ProjectID,
		ProcessId: cfg.ProcessID,
		M: &[]string{
			"OPCOUNTER_INSERT", "OPCOUNTER_QUERY", "OPCOUNTER_UPDATE", "TICKETS_AVAILABLE_READS",
			"TICKETS_AVAILABLE_WRITE", "CONNECTIONS", "QUERY_TARGETING_SCANNED_OBJECTS_PER_RETURNED",
			"QUERY_TARGETING_SCANNED_PER_RETURNED", "SYSTEM_CPU_GUEST", "SYSTEM_CPU_IOWAIT",
			"SYSTEM_CPU_IRQ", "SYSTEM_CPU_KERNEL", "SYSTEM_CPU_NICE", "SYSTEM_CPU_SOFTIRQ",
			"SYSTEM_CPU_STEAL", "SYSTEM_CPU_USER",
		},
		Granularity: admin.PtrString("PT1H"),
		Period:      admin.PtrString("P7D"),
	}

	view, err := metrics.FetchProcessMetrics(ctx, client.MonitoringAndLogsApi, p)
	if err != nil {
		log.Fatalf("Failed to fetch process metrics: %v", err)
	}

	// Output metrics
	out, err := json.MarshalIndent(view, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format metrics data: %v", err)
	}
	fmt.Println(string(out))
}

