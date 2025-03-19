package main

import (
	"admin-sdk/internal"
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"log"
	"os"
)

// Download a compressed log.gz file that contains the MongoDB logs for the specified host in your project.

func main() {
	ctx := context.Background()

	client, _, config, err := internal.CreateAtlasClient()
	if err != nil {
		log.Fatalf("Failed to create Atlas client: %v", err)
	}

	params := &GetHostLogsParams{
		GroupID:  config.AtlasProjectID,
		HostName: config.AtlasHostName,
		LogName:  "mongodb", // valid values: "mongodb" or "mongos"
	}

	if err := getHostLogs(ctx, *client, params); err != nil {
		log.Fatalf("Error fetching host logs: %v", err)
	}
}

type GetHostLogsParams struct {
	GroupID   string `json:"groupId"` // GroupID == ProjectID
	HostName  string `json:"hostName"`
	LogName   string `json:"logName"`
	EndDate   *int64 `json:"endDate,omitempty"`
	StartDate *int64 `json:"startDate,omitempty"`
}

func getHostLogs(ctx context.Context, client internal.HTTPClient, hostParams *GetHostLogsParams) error {
	fmt.Printf("Fetching logs for project %s, host %s, log %s...\n", hostParams.GroupID, hostParams.HostName, hostParams.LogName)

	// Create request params
	params := &admin.GetHostLogsApiParams{
		GroupId:   hostParams.GroupID,
		HostName:  hostParams.HostName,
		LogName:   hostParams.LogName,
		StartDate: hostParams.StartDate,
		EndDate:   hostParams.EndDate,
	}

	// Call the new GetHostLogs method
	resp, err := client.GetHostLogs(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to fetch logs for host %s in project %s: %w", hostParams.HostName, hostParams.GroupID, err)
	}
	defer func() {
		if resp != nil {
			if closeErr := resp.Close(); closeErr != nil {
				log.Printf("Warning: failed to close response body: %v", closeErr)
			}
		}
	}()

	// Create log file with a .gz extension
	logFileName := fmt.Sprintf("logs_%s_%s.gz", hostParams.GroupID, hostParams.HostName)
	logFile, err := os.Create(logFileName)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer func(logFile *os.File) {
		if logFile != nil {
			if err := logFile.Close(); err != nil {
				log.Printf("Warning: failed to close log file: %v", err)
			}
		}
	}(logFile)

	// Write logs to file
	if _, err = io.Copy(logFile, resp); err != nil {
		return fmt.Errorf("failed to write logs to file %s: %w", logFileName, err)
	}
	fmt.Printf("Logs saved to %s\n", logFileName)
	return nil
}
