package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/metrics"
)

func main() {
	_ = godotenv.Load()
	secrets, cfg, err := config.LoadAll("configs/config.json")
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	sdk, err := auth.NewClient(cfg, secrets)
	if err != nil {
		log.Fatalf("client init: %v", err)
	}

	ctx := context.Background()
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

	view, err := metrics.FetchProcessMetrics(ctx, sdk.MonitoringAndLogsApi, p)
	if err != nil {
		log.Fatalf("process metrics: %v", err)
	}

	out, _ := json.MarshalIndent(view, "", "  ")
	fmt.Println(string(out))
}

