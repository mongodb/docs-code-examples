// :snippet-start: get-metrics-dev
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/metrics"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"log"
)

func main() {
	_ = godotenv.Load()
	secrets, cfg, err := config.LoadAll("configs/.config.json")
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	sdk, err := auth.NewClient(cfg, secrets)
	if err != nil {
		log.Fatalf("client init: %v", err)
	}

	ctx := context.Background()
	p := &admin.GetDiskMeasurementsApiParams{
		GroupId:       cfg.ProjectID,
		ProcessId:     cfg.ProcessID,
		PartitionName: "data",
		M:             &[]string{"DISK_PARTITION_SPACE_FREE", "DISK_PARTITION_SPACE_USED"},
		Granularity:   admin.PtrString("P1D"),
		Period:        admin.PtrString("P1D"),
	}

	view, err := metrics.FetchDiskMetrics(ctx, sdk.MonitoringAndLogsApi, p)
	if err != nil {
		log.Fatalf("disk metrics: %v", err)
	}

	out, _ := json.MarshalIndent(view, "", "  ")
	fmt.Println(string(out))
}

// :snippet-end: [get-metrics-dev]
