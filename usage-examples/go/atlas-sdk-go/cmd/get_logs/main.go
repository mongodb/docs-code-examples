// :snippet-start: get-logs-full-script
package main

import (
	"admin-sdk/internal"
	"admin-sdk/internal/auth"
	"admin-sdk/utils"
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"log"
	"os"
)

const (
	LogName = "mongodb"
)

type GetHostLogsParams struct {
	GroupID   string `json:"groupId"` // Note: GroupID == ProjectID
	HostName  string `json:"hostName"`
	LogName   string `json:"logName"` // valid values: "mongodb" or "mongos"
	EndDate   *int64 `json:"endDate,omitempty"`
	StartDate *int64 `json:"startDate,omitempty"`
}

// Download a compressed log.gz file that contains the MongoDB logs for the specified host in your project.
func getHostLogs(ctx context.Context, client internal.HTTPClient, hostParams *GetHostLogsParams) error {
	fmt.Printf("Fetching %s log for host %s in project %s \n", hostParams.LogName, hostParams.HostName, hostParams.GroupID)

	// Create request params
	params := &admin.GetHostLogsApiParams{
		GroupId:   hostParams.GroupID,
		HostName:  hostParams.HostName,
		LogName:   hostParams.LogName,
		StartDate: hostParams.StartDate,
		EndDate:   hostParams.EndDate,
	}

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

	// Write compressed logs to file
	if _, err = io.Copy(logFile, resp); err != nil {
		return fmt.Errorf("failed to write logs to file %s: %w", logFileName, err)
	}
	fmt.Printf("Logs saved to %s\n", logFileName)
	return nil
}

// :snippet-start: get-logs-main
func main() {
	ctx := context.Background()

	// Create an Atlas client authenticated using OAuth2 with service account credentials
	client, _, config, err := auth.CreateAtlasClient()
	utils.HandleError(err, "Failed to create Atlas client")

	params := &GetHostLogsParams{
		GroupID:  config.AtlasProjectID,
		HostName: config.AtlasHostName, // The host to get logs for
		LogName:  LogName,              // The type of log to get ("mongodb" or "mongos")
	}

	// Downloads the specified host's MongoDB logs as a .gz file
	utils.HandleError(getHostLogs(ctx, *client, params), "Error fetching host logs")
}

// :snippet-end: [get-logs-main]
// :snippet-end: [get-logs-full-script]
