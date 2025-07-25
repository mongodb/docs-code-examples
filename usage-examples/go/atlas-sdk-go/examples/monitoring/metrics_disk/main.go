// :snippet-start: get-metrics-dev
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/errors"
	"atlas-sdk-go/internal/metrics"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
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
		errors.ExitWithError("Failed to fetch disk metrics", err)
	}

	// Output metrics
	out, err := json.MarshalIndent(view, "", "  ")
	if err != nil {
		errors.ExitWithError("Failed to format metrics data", err)
	}
	fmt.Println(string(out))
}

// :snippet-end: [get-metrics-dev]
// :state-remove-start: copy
// NOTE: INTERNAL
// ** OUTPUT EXAMPLE **
// {
//   "measurements": [
//     {
//       "name": "DISK_PARTITION_SPACE_FREE",
//       "granularity": "P1D",
//       "period": "P1D",
//       "values": [
//         {
//           "timestamp": "2023-10-01T00:00:00Z",
//           "value": 1234567890
//         },
//         {
//           "timestamp": "2023-10-02T00:00:00Z",
//           "value": 1234567890
//         }
//       ]
//     },
//	 	...
//   ]
// }
// :state-remove-end: [copy]
