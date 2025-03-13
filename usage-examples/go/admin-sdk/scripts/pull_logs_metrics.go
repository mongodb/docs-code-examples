package main

import (
	"admin-sdk/internal"
	"admin-sdk/utils"
	"context"
	"flag"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"log"
	"os"
	"time"
)

//type Metric struct {
//	processID   string
//	groupID     string
//	granularity string
//	metrics     []string
//	period      *string
//	start       *time.Time
//	end         *time.Time
//}

func main() {
	ctx := context.Background()

	// CLI Flags
	projectID := flag.String("project", "", "Atlas project ID")
	logName := flag.String("log", "mongodb", "Log name to fetch")
	interval := flag.Int("interval", 0, "Fetch interval in minutes (0 for one-time run)")
	//	granularity := flag.String("granularity", "PT1M", "Granularity of the metrics")
	//	metrics := flag.String("metrics", "", "Comma-separated list of metrics to fetch")
	flag.Parse()

	// Load Configuration and Authenticate
	sdk, secrets, config, err := internal.CreateAtlasClient()
	if err != nil {
		log.Fatalf("Failed to create Atlas client: %v", err)
	}

	// Override config values with CLI flags if provided
	cfgProjectID := utils.FirstNonEmpty(*projectID, config.GroupID)
	cfgLogName := utils.FirstNonEmpty(*logName, "mongodb")
	cfgInterval := utils.FirstNonZero(*interval, 0)
	//	cfgMetrics := utils.FirstNonEmpty(*metrics, "")
	//	cfgGranularity := utils.FirstNonEmpty(*granularity, "")

	// Ensure required values are set
	if cfgProjectID == "" {
		log.Fatal("Project ID is required (pass --project flag or set GROUP_ID in config)")
	}
	//metrics := []string{"SYSTEM"}
	//measurement := []string{""}
	granularity := "PT1M"
	period := "PT10H"
	//metrics := []String
	//if metrics == "" {
	//	metrics = strings.Split(*metrics, ",")
	//}
	//if cfgGranularity == "" {
	//	cfgGranularity = *granularity
	//}
	fmt.Println("DEBUG: Calling ListAtlasProcesses with GroupID:", cfgProjectID)
	fmt.Printf("DEBUG: SDK Config: %+v\n", sdk)
	fmt.Println("DEBUG: Using ClientID:", secrets.ClientID)
	fmt.Println("DEBUG: Using ClientSecret Length:", len(secrets.ClientSecret))

	// Fetch all processes in the project to get hostnames
	processes, err := fetchProcesses(ctx, sdk, cfgProjectID, false, 100, 1)
	if err != nil {
		log.Fatalf("Failed to fetch processes: %v", err)
	}

	// Extract host names and process IDs from the response
	var processIDs []string
	var hostNames []string
	if processes.Results != nil {
		for _, process := range *processes.Results {
			if process.Id != nil {
				processIDs = append(processIDs, *process.Id)
			}
			if process.Hostname != nil {
				hostNames = append(hostNames, *process.Hostname)
			}
		}
	} else {
		log.Fatal("No processes found in the project")
	}

	// Loop through each host and fetch logs & metrics
	for {
		for i := range processIDs {
			if err := fetchHostLogs(ctx, sdk, cfgProjectID, hostNames[i], cfgLogName); err != nil {
				log.Printf("Error fetching logs for %s: %v", processIDs[i], err)
			}

			if err := fetchHostMetrics(ctx, sdk, cfgProjectID, processIDs[i], granularity, period); err != nil {
				log.Printf("Error fetching metrics for %s: %v", processIDs[i], err)
			}
		}

		if cfgInterval == 0 {
			break
		}
		time.Sleep(time.Duration(cfgInterval) * time.Minute)
	}
}

// https://cloud.mongodb.com/api/atlas/v2/groups/{groupId}/clusters/{hostName}/logs/{logName}.gz
func fetchHostLogs(ctx context.Context, sdk *admin.APIClient, groupID, hostName, logName string) error {
	fmt.Printf("Fetching logs for project %s, host %s, log %s...\n", groupID, hostName, logName)

	req := sdk.MonitoringAndLogsApi.GetHostLogs(ctx, groupID, hostName, logName)
	fmt.Println("DEBUG: Generated Fetch Logs API request:", req)

	resp, _, err := req.Execute()
	//resp, _, err := sdk.MonitoringAndLogsApi.GetHostLogs(ctx, groupID, hostName, logName).Execute()
	//dk.MonitoringAndLogsAPI.GetHostLogs(ctx, groupID, hostName, logName).Execute()
	if err != nil {
		return fmt.Errorf("failed to get logs: %w", err)
	}
	defer resp.Close()

	logFile, err := os.Create(fmt.Sprintf("logs_%s_%s.log", groupID, hostName))
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	_, err = io.Copy(logFile, resp)
	if err != nil {
		return fmt.Errorf("failed to save logs: %w", err)
	}
	fmt.Println("Logs saved.")

	return nil
}

//https://cloud.mongodb.com/api/atlas/v2/groups/{groupId}/processes/{processId}/measurements

func fetchHostMetrics(ctx context.Context, sdk *admin.APIClient, groupID, processID, granularity, period string) error {
	//fmt.Printf("Fetching metrics for project %s ...\n\n", groupID)
	fmt.Printf("DEBUG: Fetching metrics for GroupID: %s, ProcessID: %s, Granularity: %s\n", groupID, processID, granularity, period)
	// Simulate request construction
	req := sdk.MonitoringAndLogsApi.GetHostMeasurements(ctx, groupID, processID).
		Granularity(granularity).Period(period)

	fmt.Println("DEBUG: Generated API request:", req)

	resp, r, err := req.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GetMeasurements`: %v (%v)\n", err, r)
		return err
	}
	//resp, r, err := sdk.MonitoringAndLogsApi.GetMeasurements(ctx, groupID, processID).Granularity(granularity).Metrics(metrics).Execute()
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Error when calling `MonitoringAndLogsApi.GetMeasurements`: %v (%v)\n", err, r)
	//	apiError, ok := admin.AsError(err)
	//	if ok {
	//		fmt.Fprintf(os.Stderr, "API error obj: %v\n", apiError)
	//	}
	//}
	// response from `GetMeasurements`: MeasurementsNonIndex
	fmt.Fprintf(os.Stdout, "Response from `MonitoringAndLogsApi.GetMeasurements`: %v (%v)\n", resp, r)
	return nil
}

// GetHostNameFromID retrieves the hostname of a process in an Atlas project.

func fetchProcesses(ctx context.Context, sdk *admin.APIClient, groupID string, includeCount bool, itemsPerPage, pageNum int) (*admin.PaginatedHostViewAtlas, error) {
	resp, r, err := sdk.MonitoringAndLogsApi.ListAtlasProcesses(ctx, groupID).IncludeCount(includeCount).ItemsPerPage(itemsPerPage).PageNum(pageNum).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MonitoringAndLogsApi.ListAtlasProcesses`: %v (%v)\n", err, r)
		apiError, ok := admin.AsError(err)
		if ok {
			fmt.Fprintf(os.Stderr, "API error obj: %v\n", apiError)
		}
		// response from `ListAtlasProcesses`: PaginatedHostViewAtlas
		fmt.Fprintf(os.Stdout, "Response from `MonitoringAndLogsApi.ListAtlasProcesses`: %v (%v)\n", resp, r)
	}
	return resp, nil
}
