package main

import (
	"admin-sdk/internal"
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	client, _, config, err := internal.CreateAtlasClient()
	if err != nil {
		fmt.Printf("Failed to create Atlas client: %v", err)
	}

	processID := admin.PtrString(config.AtlasHostName + config.AtlasPort)
	if config.AtlasProcessID == "" {
		config.AtlasProcessID = *processID
	}

	getProcessMetricParams := &GetProcessMetricParams{
		GroupID:     config.AtlasProjectID,
		ProcessID:   *processID,
		Granularity: admin.PtrString("PT1M"),
		Period:      admin.PtrString("PT10H"),
	}
	err = getProcessMetrics(ctx, client, getProcessMetricParams)
	if err != nil {
		fmt.Printf("Error fetching process metrics: %v", err)
	}
}

type GetProcessMetricParams struct {
	GroupID     string     `json:"groupId"`
	ProcessID   string     `json:"processId"`
	Granularity *string    `json:"granularity"`
	M           *[]string  `json:"metrics"`
	Period      *string    `json:"period"`
	Start       *time.Time `json:"start,omitempty"`
	End         *time.Time `json:"end,omitempty"`
}

// Return Measurements for One MongoDB Process
// ApiMeasurementsGeneralViewAtlas GetHostMeasurements(ctx, groupId, processId).Granularity(granularity).M(m).Period(period).Start(start).End(end).Execute()
func getProcessMetrics(ctx context.Context, client internal.HTTPClient, hostParams *GetProcessMetricParams) error {
	fmt.Printf("Fetching metrics for project %s", hostParams.GroupID)

	params := &admin.GetHostMeasurementsApiParams{
		GroupId:     hostParams.GroupID,
		ProcessId:   hostParams.ProcessID,
		Granularity: hostParams.Granularity,
		M:           hostParams.M,
		Period:      hostParams.Period,
		Start:       hostParams.Start,
		End:         hostParams.End,
	}

	resp, r, err := client.GetProcessMetrics(ctx, params)
	if err != nil {
		if apiError, ok := admin.AsError(err); ok {
			return fmt.Errorf("failed to get measurements for process in host: %s (API error: %v)", err, apiError.GetDetail())
		}
	}
	if resp.HasMeasurements() == false {
		return fmt.Errorf("no measurements found for process %s in group %s", params.ProcessId, params.GroupId)
	}
	fmt.Fprintf(os.Stdout, "Response from `MonitoringAndLogsApi.GetMeasurements`: %v (%v)", resp, r)
	return nil
}

type GetDiskMetricParams struct {
	GroupID       string     `json:"groupId"`
	ProcessID     string     `json:"processId"`
	PartitionName string     `json:"partitionName"`
	M             *[]string  `json:"metrics,omitempty"`
	Period        *string    `json:"period,omitempty"`
	Start         *time.Time `json:"start,omitempty"`
	End           *time.Time `json:"end,omitempty"`
}

// Get /api/atlas/v2/groups/{groupId}/processes/{processId}/disks/{partitionName}/measurements
func getClusterMetrics(ctx context.Context, client internal.HTTPClient, diskParams *GetDiskMetricParams) (*admin.MeasurementDiskPartition, *admin.APIResponse, error) {

	params := &admin.GetDiskMeasurementsApiParams{
		GroupId:       diskParams.GroupID,
		ProcessId:     diskParams.ProcessID,
		PartitionName: diskParams.PartitionName,
		M:             diskParams.M,
		Period:        diskParams.Period,
		Start:         diskParams.Start,
		End:           diskParams.End,
	}

	resp, r, err := client.GetClusterMetrics(ctx, params)
	if err != nil {
		if apiError, ok := admin.AsError(err); ok {
			return nil, nil, fmt.Errorf("failed to get measurements for cluster in group: %s (API error: %v)", err, apiError.GetDetail())
		}
	}

	if resp.HasMeasurements() == false {
		return nil, nil, fmt.Errorf("no measurements found for cluster %s in group %s", diskParams.ProcessID, diskParams.GroupID)
	}
	fmt.Fprintf(os.Stdout, "Response from `MonitoringAndLogsApi.GetDiskMeasurement`: %v (%v)", resp, r)
	return nil, nil, nil

}
