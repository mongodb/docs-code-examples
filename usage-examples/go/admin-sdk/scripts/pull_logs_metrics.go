package main

import (
	"admin-sdk/internal"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	// CLI Flags
	projectID := flag.String("project", "", "Atlas project ID")
	hostName := flag.String("host", "", "Hostname of the process")
	logName := flag.String("log", "mongos", "Log name to fetch")
	interval := flag.Int("interval", 0, "Fetch interval in minutes (0 for one-time run)")
	// configPath := flag.String("config", "config/config.json", "Path to JSON config file")
	flag.Parse()

	// Load Configuration and Authenticate
	sdk, config, err := internal.CreateAtlasClient()
	if err != nil {
		log.Fatalf("Failed to create Atlas client: %v", err)
	}

	// Override config values with CLI flags if provided
	cfgProjectID := firstNonEmpty(*projectID, config.GroupID)
	cfgHostName := *hostName // No default, must be provided
	cfgLogName := firstNonEmpty(*logName, "mongos")
	cfgInterval := firstNonZero(*interval, 0)

	// Ensure required values are set
	if cfgProjectID == "" {
		log.Fatal("Project ID is required (pass --project flag or set GROUP_ID in config)")
	}
	if cfgHostName == "" {
		log.Fatal("Host name is required (pass --host flag)")
	}

	// Loop for repeated execution if interval is set
	for {
		if err := fetchLogs(ctx, sdk, cfgProjectID, cfgHostName, cfgLogName); err != nil {
			log.Printf("Error fetching logs: %v", err)
		}

		if cfgInterval == 0 {
			break
		}
		time.Sleep(time.Duration(cfgInterval) * time.Minute)
	}
	// Fetch metrics
	if err := fetchMetrics(ctx, sdk, cfgProjectID, cfgHostName); err != nil {
		log.Printf("Error fetching metrics: %v", err)
	}
}

func fetchLogs(ctx context.Context, sdk *admin.APIClient, projectID, hostName, logName string) error {
	fmt.Printf("Fetching logs for project %s, host %s, log %s...\n", projectID, hostName, logName)

	resp, _, err := sdk.MonitoringAndLogsApi.GetHostLogsWithParams(ctx, &admin.GetHostLogsApiParams{
		GroupId:  projectID,
		HostName: hostName,
		LogName:  logName,
	}).Execute()
	if err != nil {
		return fmt.Errorf("failed to get logs: %w", err)
	}
	defer resp.Close()

	logFile, err := os.Create(fmt.Sprintf("logs_%s_%s.log", projectID, hostName))
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

func fetchMetrics(ctx context.Context, sdk *admin.APIClient, projectID, hostName string) error {
	fmt.Printf("Fetching metrics for project %s, host %s...\n\n", projectID, hostName)

	resp, _, err := sdk.MonitoringAndLogsApi.GetMeasurements(ctx, &admin.GetMeasurementsApiParams{
		GroupId:  projectID,
		HostName: hostName,
		"SYSTEM"}).Execute()
	if err != nil {
		return fmt.Errorf("failed to fetch metrics: %w", err)
	}

	metricsFile, err := os.Create(fmt.Sprintf("metrics_%s_%s.json", projectID, hostName))
	if err != nil {
		return fmt.Errorf("failed to create metrics file: %w", err)
	}
	defer metricsFile.Close()

	json.NewEncoder(metricsFile).Encode(resp)
	fmt.Println("Metrics saved.")
	return nil
}

func getHostNameFromID(ctx context.Context, sdk *admin.APIClient, projectID string) (string, error) {
	resp, _, err := sdk.MonitoringAndLogsApi.ListAtlasProcesses(ctx, projectID).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to get host: %w", err)
	}
	return resp.HostName, nil
}

// Utility functions
func firstNonEmpty(cli, config string) string {
	if cli != "" {
		return cli
	}
	return config
}

func firstNonZero(cli, config int) int {
	if cli != 0 {
		return cli
	}
	return config
}
