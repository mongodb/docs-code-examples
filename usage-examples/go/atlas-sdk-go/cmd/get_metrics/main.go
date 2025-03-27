// :snippet-start: get-metrics-full-script
package main

import (
	"atlas-sdk-go/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"log"
)

// Fetches metrics for a specified host process in a project
func getProcessMetrics(ctx context.Context, atlasClient admin.APIClient, params *admin.GetHostMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, error) {
	fmt.Printf("Fetching metrics for process %s in project %s", params.ProcessId, params.GroupId)

	hostMeasurementsParams := &admin.GetHostMeasurementsApiParams{
		GroupId:     params.GroupId,
		ProcessId:   params.ProcessId,
		Granularity: params.Granularity,
		M:           params.M,
		Period:      params.Period,
		Start:       params.Start,
		End:         params.End,
	}

	resp, _, err := atlasClient.MonitoringAndLogsApi.GetHostMeasurementsWithParams(ctx, hostMeasurementsParams).Execute()
	if err != nil {
		if apiError, ok := admin.AsError(err); ok {
			return nil, fmt.Errorf("failed to get measurements for process in host: %s (API error: %v)", err, apiError.GetDetail())
		}
		return nil, fmt.Errorf("failed to get measurements: %w", err)
	}

	if resp == nil || resp.HasMeasurements() == false {
		return nil, fmt.Errorf("no measurements found for process %s in project %s", params.ProcessId, params.GroupId)
	}

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	fmt.Println(string(jsonData))

	return resp, nil
}

// Fetch metrics for a specified disk partition in a project and print the results to the console as JSON
func getDiskMetrics(ctx context.Context, atlasClient admin.APIClient, params *admin.GetDiskMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, error) {

	diskMeasurementParams := &admin.GetDiskMeasurementsApiParams{
		GroupId:       params.GroupId,
		ProcessId:     params.ProcessId,
		PartitionName: params.PartitionName,
		Granularity:   params.Granularity,
		M:             params.M,
		Period:        params.Period,
		Start:         params.Start,
		End:           params.End,
	}

	resp, _, err := atlasClient.MonitoringAndLogsApi.GetDiskMeasurementsWithParams(ctx, diskMeasurementParams).Execute()
	if err != nil {
		if apiError, ok := admin.AsError(err); ok {
			return nil, fmt.Errorf("failed to get measurements for partition: %s (API error: %v)", err, apiError.GetDetail())
		}
		return nil, fmt.Errorf("failed to get measurements: %w", err)
	}
	if resp == nil || resp.HasMeasurements() == false {
		return nil, fmt.Errorf("no measurements found for partition %s in project %s", params.PartitionName, params.GroupId)
	}

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	fmt.Println(string(jsonData))
	return resp, nil
}

// :snippet-start: get-metrics-main-dev
// :snippet-start: get-metrics-main-prod
func main() {
	ctx := context.Background()

	// Create an Atlas client authenticated using OAuth2 with service account credentials
	atlasClient, _, config, err := auth.CreateAtlasClient()
	if err != nil {
		log.Fatalf("Failed to create Atlas client: %v", err)
	}
	// :state-start: prod
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
	hostMeasurementsParams := &admin.GetHostMeasurementsApiParams{
		GroupId:     config.AtlasProjectID,
		ProcessId:   config.AtlasProcessID,
		M:           &processMetrics,
		Granularity: processMetricGranularity,
		Period:      processMetricPeriod,
	}
	_, err = getProcessMetrics(ctx, *atlasClient, hostMeasurementsParams)
	if err != nil {
		fmt.Printf("Error fetching host process metrics: %v", err)
	}
	// :state-end: [prod]
	// :state-start: dev
	// Fetch disk metrics using the following parameters
	partitionName := "data"
	diskMetricsGranularity := admin.PtrString("P1D")
	diskMetricsPeriod := admin.PtrString("P1D")
	diskMetrics := []string{
		"DISK_PARTITION_SPACE_FREE", "DISK_PARTITION_SPACE_USED",
	}

	diskMeasurementsParams := &admin.GetDiskMeasurementsApiParams{
		GroupId:       config.AtlasProjectID,
		ProcessId:     config.AtlasProcessID,
		PartitionName: partitionName,
		M:             &diskMetrics,
		Granularity:   diskMetricsGranularity,
		Period:        diskMetricsPeriod,
	}
	_, err = getDiskMetrics(ctx, *atlasClient, diskMeasurementsParams)
	if err != nil {
		fmt.Printf("Error fetching disk metrics: %v", err)
	}
	// :state-end: [dev]
}

// :snippet-end: [get-metrics-main-dev]
// :snippet-end: [get-metrics-main-prod]
// :snippet-end: [get-metrics-full-script]
