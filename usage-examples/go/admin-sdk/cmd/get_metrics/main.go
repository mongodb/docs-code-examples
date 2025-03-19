package main

import (
	"admin-sdk/internal"
	"admin-sdk/internal/auth"
	"admin-sdk/utils"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"time"
)

const (
	granularity    = "P1D"
	period         = "P1D"
	partitionName  = "data"
	processMetrics = "DB_DATA_SIZE_TOTAL,MAX_SYSTEM_MEMORY_AVAILABLE"
	diskMetrics    = "DISK_PARTITION_SPACE_FREE, DISK_PARTITION_SPACE_USED"
)

type GetProcessMetricParams struct {
	GroupID     string     `json:"groupId"` // GroupID == ProjectID
	ProcessID   string     `json:"processId"`
	Granularity *string    `json:"granularity"`
	M           *[]string  `json:"metrics"`
	Period      *string    `json:"period"`
	Start       *time.Time `json:"start,omitempty"`
	End         *time.Time `json:"end,omitempty"`
}

type GetDiskMetricParams struct {
	GroupID       string     `json:"groupId"`
	ProcessID     string     `json:"processId"`
	PartitionName string     `json:"partitionName"`
	Granularity   *string    `json:"granularity"`
	M             *[]string  `json:"metrics,omitempty"`
	Period        *string    `json:"period,omitempty"`
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

	resp, _, err := client.GetProcessMetrics(ctx, params)
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

// Fetch metrics for a specified disk partition in a project
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

	resp, _, err := client.GetDiskMetrics(ctx, params)
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
	client, _, config, err := auth.CreateAtlasClient()
	if err != nil {
		fmt.Printf("Failed to create Atlas client: %v", err)
	}

	getProcessMetricParams := &GetProcessMetricParams{
		GroupID:     config.AtlasProjectID,
		ProcessID:   config.AtlasProcessID,
		M:           &[]string{processMetrics},
		Granularity: admin.PtrString(granularity),
		Period:      admin.PtrString(period),
	}
	err = getProcessMetrics(ctx, *client, getProcessMetricParams)
	utils.HandleError(err, "Error fetching host process metrics")

	getDiskMetricParams := &GetDiskMetricParams{
		GroupID:       config.AtlasProjectID,
		ProcessID:     config.AtlasProcessID,
		PartitionName: partitionName,
		M:             &[]string{diskMetrics},
		Granularity:   admin.PtrString(granularity),
		Period:        admin.PtrString(period),
	}
	err = getDiskMetrics(ctx, *client, getDiskMetricParams)
	utils.HandleError(err, "Error fetching partition disk metrics")
}
