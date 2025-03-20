package main

import (
	"atlas-sdk-go/internal"
	"atlas-sdk-go/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"log"
	"time"
)

type GetProcessMetricParams struct {
	GroupID     string     `json:"groupId"` // Note: GroupID == ProjectID
	ProcessID   string     `json:"processId"`
	Granularity *string    `json:"granularity"`
	M           *[]string  `json:"metrics,omitempty"`
	Period      *string    `json:"diskMetricsPeriod,omitempty"`
	Start       *time.Time `json:"start,omitempty"`
	End         *time.Time `json:"end,omitempty"`
}

type GetDiskMetricParams struct {
	GroupID       string     `json:"groupId"` // Note: GroupID == ProjectID
	ProcessID     string     `json:"processId"`
	PartitionName string     `json:"partitionName"`
	Granularity   *string    `json:"granularity"`
	M             *[]string  `json:"metrics,omitempty"`
	Period        *string    `json:"diskMetricsPeriod,omitempty"`
	Start         *time.Time `json:"start,omitempty"`
	End           *time.Time `json:"end,omitempty"`
}

// Fetches metrics for a specified host process in a project
func getProcessMetrics(ctx context.Context, client internal.HTTPClient, hostParams *GetProcessMetricParams) error {
	fmt.Printf("Fetching metrics for process %s in project %s", hostParams.ProcessID, hostParams.GroupID)

	params := &admin.GetHostMeasurementsApiParams{
		GroupId:     hostParams.GroupID,
		ProcessId:   hostParams.ProcessID,
		Granularity: hostParams.Granularity,
		M:           hostParams.M,
		Period:      hostParams.Period,
		Start:       hostParams.Start,
		End:         hostParams.End,
	}

	resp, err := client.GetProcessMetrics(ctx, params)
	if err != nil {
		if apiError, ok := admin.AsError(err); ok {
			return fmt.Errorf("failed to get measurements for process in host: %s (API error: %v)", err, apiError.GetDetail())
		}
		return fmt.Errorf("failed to get measurements: %w", err)
	}

	if resp == nil || resp.HasMeasurements() == false {
		return fmt.Errorf("no measurements found for process %s in project %s", params.ProcessId, params.GroupId)
	}

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	fmt.Println(string(jsonData))

	return nil
}

// Fetch metrics for a specified disk partition in a project and print the results to the console as JSON
func getDiskMetrics(ctx context.Context, client internal.HTTPClient, diskParams *GetDiskMetricParams) error {

	params := &admin.GetDiskMeasurementsApiParams{
		GroupId:       diskParams.GroupID,
		ProcessId:     diskParams.ProcessID,
		PartitionName: diskParams.PartitionName,
		Granularity:   diskParams.Granularity,
		M:             diskParams.M,
		Period:        diskParams.Period,
		Start:         diskParams.Start,
		End:           diskParams.End,
	}

	resp, err := client.GetDiskMetrics(ctx, params)
	if err != nil {
		if apiError, ok := admin.AsError(err); ok {
			return fmt.Errorf("failed to get measurements for partition: %s (API error: %v)", err, apiError.GetDetail())
		}
		return fmt.Errorf("failed to get measurements: %w", err)
	}
	if resp == nil || resp.HasMeasurements() == false {
		return fmt.Errorf("no measurements found for partition %s in project %s", params.PartitionName, params.GroupId)
	}

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

func main() {
	ctx := context.Background()

	// Create an Atlas client authenticated using OAuth2 with service account credentials
	client, _, config, err := auth.CreateAtlasClient()
	if err != nil {
		log.Fatalf("Failed to create Atlas client: %v", err)
	}
	// Fetch process metrics using the following parameters
	processMetricGranularity := admin.PtrString("PT1H")
	processMetricPeriod := admin.PtrString("P7D")
	processMetrics := []string{
		"OPCOUNTER_INSERT", "OPCOUNTER_QUERY", "OPCOUNTER_UPDATE", "TICKETS_AVAILABLE_READS",
		"TICKETS_AVAILABLE_WRITE", "CONNECTIONS", "QUERY_TARGETING_SCANNED_OBJECTS_PER_RETURNED",
		"QUERY_TARGETING_SCANNED_PER_RETURNED", "SYSTEM_CPU_GUEST", "SYSTEM_CPU_IOWAIT",
		"SYSTEM_CPU_IRQ", "SYSTEM_CPU_KERNEL", "SYSTEM_CPU_NICE", "SYSTEM_CPU_SOFTIRQ",
		"SYSTEM_CPU_STEAL", "SYSTEM_CPU_USER",
	}
	getProcessMetricParams := &GetProcessMetricParams{
		GroupID:     config.AtlasProjectID,
		ProcessID:   config.AtlasProcessID,
		M:           &processMetrics,
		Granularity: processMetricGranularity,
		Period:      processMetricPeriod,
	}
	if err := getProcessMetrics(ctx, *client, getProcessMetricParams); err != nil {
		fmt.Printf("Error fetching host process metrics: %v", err)
	}
}

