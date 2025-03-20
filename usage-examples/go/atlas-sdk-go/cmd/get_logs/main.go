// :snippet-start: get-logs-full-script
package main

import (
	"atlas-sdk-go/internal"
	"atlas-sdk-go/internal/auth"
	test "atlas-sdk-go/tests" // :remove:
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

	params := &admin.GetHostLogsApiParams{
		GroupId:   hostParams.GroupID,
		HostName:  hostParams.HostName,
		LogName:   hostParams.LogName,
		StartDate: hostParams.StartDate,
		EndDate:   hostParams.EndDate,
	}

	resp, err := client.GetHostLogs(ctx, params)
	if err != nil {
		if apiError, ok := admin.AsError(err); ok {
			return fmt.Errorf("failed to fetch logs for host: %s (API error: %v)", err, apiError.GetDetail())
		}
		return fmt.Errorf("failed to get logs: %w", err)
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
	if err != nil {
		log.Fatalf("Failed to create Atlas client: %v", err)
	}

	params := &GetHostLogsParams{
		GroupID:  config.AtlasProjectID,
		HostName: config.AtlasHostName, // The host to get logs for
		LogName:  LogName,              // The type of log to get ("mongodb" or "mongos")
	}

	// Downloads the specified host's MongoDB logs as a .gz file
	if err := getHostLogs(ctx, *client, params); err != nil {
		fmt.Printf("Error fetching host logs: %v", err)
	}
	// :remove-start:
	// NOTE Internal function to clean up any downloaded files
	if err := test.CleanupGzFiles(); err != nil {
		log.Printf("Cleanup error: %v", err)
	}
	// :remove-end:
}

// :snippet-end: [get-logs-main]
// :snippet-end: [get-logs-full-script]
