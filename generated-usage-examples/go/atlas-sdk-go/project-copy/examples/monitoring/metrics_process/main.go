package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"atlas-sdk-go/internal/errors"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal/metrics"
)

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

	ctx := context.Background()

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
		errors.ExitWithError("Failed to fetch process metrics", err)
	}

	// Output metrics
	out, err := json.MarshalIndent(view, "", "  ")
	if err != nil {
		errors.ExitWithError("Failed to format metrics data", err)
	}
	fmt.Println(string(out))
}

