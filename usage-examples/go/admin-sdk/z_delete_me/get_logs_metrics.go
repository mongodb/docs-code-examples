package z_delete_me

//
//import (
//	"admin-sdk/internal"
//	"admin-sdk/scripts"
//	"context"
//	"fmt"
//	"go.mongodb.org/atlas-sdk/v20250219001/admin"
//	"io"
//	"log"
//	"os"
//	"testing" // :remove:
//)
//
//func GetHostLogsAndMetrics(t *testing.T) {
//	ctx := context.Background()
//
//	sdk, _, config, err := internal.CreateAtlasClient()
//	if err != nil {
//		log.Fatalf("Error creating Atlas client: %v", err)
//	}
//
//	// Get the list of projects to retrieve the groupID
//	projectParams := &scripts.ListProjectsParams{
//		GroupID:      config.AtlasProjectID,
//		ItemsPerPage: admin.PtrInt(1),
//		IncludeCount: admin.PtrBool(true),
//		PageNum:      admin.PtrInt(1),
//	}
//	projects, err := getProjectList(ctx, sdk, projectParams)
//	if err != nil {
//		log.Fatalf("Error getting projects: %v", err)
//	}
//	projectId := projects.GetResults()[0].Id
//
//	// If groupID is not set in the config, use the retrieved ID
//	if config.AtlasProjectID == "" {
//		config.AtlasProjectID = *projectId
//	}
//
//	// Get list of processes
//	processParams := &scripts.ListAtlasProcessesParams{
//		GroupID:      config.AtlasProjectID,
//		IncludeCount: admin.PtrBool(true),
//		ItemsPerPage: admin.PtrInt(1),
//		PageNum:      admin.PtrInt(1),
//	}
//	hosts, err := getProcessList(ctx, sdk, processParams)
//	if err != nil {
//		log.Fatalf("Error getting processes: %v", err)
//	}
//	hostName := hosts.GetResults()[0].Hostname
//	processID := hosts.GetResults()[0].Id
//
//	hostParams := &scripts.GetHostLogsParams{
//		GroupID:  config.AtlasProjectID,
//		HostName: *hostName,
//		LogName:  "mongodb",
//	}
//	// Get logs for the first host in the list
//	err = scripts.getHostLogs(ctx, sdk, hostParams)
//	if err != nil {
//		log.Fatalf("Error getting host logs: %v", err)
//	}
//	// Get metrics for one disk on the first host in the list
//	hostMetricParams := scripts.HostMetricParams{
//		GroupID:     config.AtlasProjectID,
//		ProcessID:   *processID,
//		Granularity: admin.PtrString("PT1M"),
//		Period:      admin.PtrString("PT10H"),
//	}
//	// Get metrics for the first host in the list
//	err = getHostMetrics(ctx, sdk, hostMetricParams)
//	if err != nil {
//		log.Fatalf("Error getting host metrics: %v", err)
//	}
//
//	clusterMetricParams := scripts.ClusterMetricParams{
//		GroupID:       config.AtlasProjectID,
//		ProcessID:     *processID,
//		PartitionName: "data",
//		Period:        admin.PtrString("P1D"),
//		M:             &[]string{"DISK_PARTITION_SPACE_FREE", "DISK_PARTITION_SPACE_USED"},
//	}
//	err = getClusterMetrics(ctx, sdk, clusterMetricParams)
//	if err != nil {
//		log.Fatalf("Error getting host metrics: %v", err)
//	}
//}
//
//// GetHostLogs
//// Download the logs for a specific host in an Atlas project.
//// Equivalent to atlas logs download CLI command
//// get hostname from the process list for the cluster
//// Get /api/atlas/v2/groups/{groupId}/clusters/{hostName}/logs/{logName}.gz
//
//func getHostLogs(ctx context.Context, sdk *admin.APIClient, hostParams *GetHostLogsParams) error {
//	fmt.Printf("Fetching logs for project %s, host %s, log %s...\n", hostParams.GroupID, hostParams.HostName, hostParams.LogName)
//	resp, _, err := sdk.MonitoringAndLogsApi.GetHostLogsWithParams(ctx, &admin.GetHostLogsApiParams{
//		GroupId:   hostParams.GroupID,
//		HostName:  hostParams.HostName,
//		LogName:   hostParams.LogName,
//		EndDate:   hostParams.EndDate,
//		StartDate: hostParams.StartDate,
//	}).Execute()
//	if err != nil {
//		if apiError, ok := admin.AsError(err); ok {
//			return fmt.Errorf("failed to fetch logs for host %s in group %s: %w (API error: %v)", hostParams.HostName, hostParams.GroupID, err, apiError)
//		}
//		return fmt.Errorf("failed to fetch logs for host %s in group %s: %w", hostParams.HostName, hostParams.GroupID, err)
//	}
//	defer func() {
//		if resp != nil {
//			if closeErr := resp.Close(); closeErr != nil {
//				log.Printf("Warning: failed to close response body: %v", closeErr)
//			}
//		}
//	}()
//	// Create log file
//	logFileName := fmt.Sprintf("logs_%s_%s.log.gz", hostParams.GroupID, hostParams.HostName)
//	logFile, err := os.Create(logFileName)
//	if err != nil {
//		return fmt.Errorf("failed to create log file %s: %w", logFileName, err)
//	}
//	defer func(logFile *os.File) {
//		if logFile != nil {
//			if err := logFile.Close(); err != nil {
//				log.Printf("Warning: failed to close log file: %v", err)
//			}
//		}
//	}(logFile)
//
//	writer := logFile
//	if _, err = io.Copy(writer, resp); err != nil {
//		return fmt.Errorf("failed to write logs to file %s: %w", logFileName, err)
//	}
//	fmt.Printf("Logs saved to %s\n", logFileName)
//	return nil
//}
//
//func getProjectList(ctx context.Context, sdk *admin.APIClient, projectParams *scripts.ListProjectsParams) (*admin.PaginatedAtlasGroup, error) {
//	resp, _, err := sdk.ProjectsApi.ListProjectsWithParams(ctx,
//		&admin.ListProjectsApiParams{
//			ItemsPerPage: projectParams.ItemsPerPage,
//			IncludeCount: projectParams.IncludeCount,
//			PageNum:      projectParams.PageNum,
//		}).Execute()
//	if err != nil {
//		if apiError, ok := admin.AsError(err); ok {
//			return nil, fmt.Errorf("error getting projects: %w (API error: %v)+", err,
//				apiError.GetDetail())
//		}
//	}
//	if resp.GetTotalCount() == 0 {
//		log.Fatal("account should have at least single project")
//	}
//	return resp, nil
//}
//
//// GetProcessList
//// Get the list of processes for a specific Atlas project.
//// Equivalent to atlas processes list CLI command
//// PaginatedHostViewAtlas ListAtlasProcesses(ctx, groupId).IncludeCount(includeCount).ItemsPerPage(itemsPerPage).PageNum(pageNum).Execute()
//
//func getProcessList(ctx context.Context, sdk *admin.APIClient, processParams *scripts.ListAtlasProcessesParams) (*admin.PaginatedHostViewAtlas, error) {
//	resp, _, err := sdk.MonitoringAndLogsApi.ListAtlasProcessesWithParams(ctx,
//		&admin.ListAtlasProcessesApiParams{
//			GroupId:      processParams.GroupID,
//			IncludeCount: processParams.IncludeCount,
//			ItemsPerPage: processParams.ItemsPerPage,
//			PageNum:      processParams.PageNum,
//		}).Execute()
//	if err != nil {
//		if apiError, ok := admin.AsError(err); ok {
//			return nil, fmt.Errorf("failed to list processes in group: %s (API error: %v)", err, apiError.GetDetail())
//		}
//		// Debugging only: Remove or log only at the caller level
//		// fmt.Fprintf(os.Stdout, "Response: %v\n", resp)
//	}
//	if resp.GetTotalCount() == 0 {
//		return nil, fmt.Errorf("no processes found in group %s", processParams.GroupID)
//	}
//	return resp, nil
//}
//
//// Return Measurements for One MongoDB Process
//// ApiMeasurementsGeneralViewAtlas GetHostMeasurements(ctx, groupId, processId).Granularity(granularity).M(m).Period(period).Start(start).End(end).Execute()
//func getHostMetrics(ctx context.Context, sdk *admin.APIClient, metricParams scripts.HostMetricParams) error {
//	resp, r, err := sdk.MonitoringAndLogsApi.GetHostMeasurementsWithParams(ctx, &admin.GetHostMeasurementsApiParams{
//		GroupId:     metricParams.GroupID,
//		ProcessId:   metricParams.ProcessID,
//		Granularity: metricParams.Granularity,
//		M:           metricParams.M,
//		Period:      metricParams.Period,
//		Start:       metricParams.Start,
//		End:         metricParams.End,
//	}).Execute()
//	if err != nil {
//		if apiError, ok := admin.AsError(err); ok {
//			return fmt.Errorf("failed to get measurements for process in group: %s (API error: %v)", err, apiError.GetDetail())
//		}
//	}
//	if resp.HasMeasurements() == false {
//		return fmt.Errorf("no measurements found for process %s in group %s", metricParams.ProcessID, metricParams.GroupID)
//	}
//	fmt.Fprintf(os.Stdout, "Response from `MonitoringAndLogsApi.GetMeasurements`: %v (%v)", resp, r)
//	return nil
//}
//
//// Get /api/atlas/v2/groups/{groupId}/processes/{processId}/disks/{partitionName}/measurements
//func getClusterMetrics(ctx context.Context, sdk *admin.APIClient, clusterMetricParams scripts.ClusterMetricParams) error {
//	resp, r, err := sdk.MonitoringAndLogsApi.GetDiskMeasurementsWithParams(ctx, &admin.GetDiskMeasurementsApiParams{
//		GroupId:       clusterMetricParams.GroupID,
//		ProcessId:     clusterMetricParams.ProcessID,
//		PartitionName: clusterMetricParams.PartitionName,
//		Period:        clusterMetricParams.Period,
//		M:             clusterMetricParams.M,
//		Start:         clusterMetricParams.Start,
//		End:           clusterMetricParams.End,
//	}).Execute()
//	if err != nil {
//		if apiError, ok := admin.AsError(err); ok {
//			return fmt.Errorf("failed to get measurements for cluster in group: %s (API error: %v)", err, apiError.GetDetail())
//		}
//	}
//	if resp.HasMeasurements() == false {
//		return fmt.Errorf("no measurements found for cluster %s in group %s", clusterMetricParams.ProcessID, clusterMetricParams.GroupID)
//	}
//	_, err = fmt.Fprintf(os.Stdout, "Response from `MonitoringAndLogsApi.GetDiskMeasurement`: %v (%v)", resp, r)
//	if err != nil {
//		return err
//	}
//	return nil
//}
