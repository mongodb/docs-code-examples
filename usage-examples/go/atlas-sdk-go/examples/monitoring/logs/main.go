// :snippet-start: get-logs
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"context"
	"fmt"
	"log"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/errors"
	"atlas-sdk-go/internal/fileutils"
	"atlas-sdk-go/internal/logs"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}

	secrets, cfg, err := config.LoadAll("configs/config.json")
	if err != nil {
		errors.ExitWithError("Failed to load configuration", err)
	}

	client, err := auth.NewClient(cfg, secrets)
	if err != nil {
		errors.ExitWithError("Failed to initialize authentication client", err)
	}

	ctx := context.Background()

	// Fetch logs with the provided parameters
	p := &admin.GetHostLogsApiParams{
		GroupId:  cfg.ProjectID,
		HostName: cfg.HostName,
		LogName:  "mongodb",
	}
	fmt.Printf("Request parameters: GroupID=%s, HostName=%s, LogName=%s\n",
		cfg.ProjectID, cfg.HostName, p.LogName)
	rc, err := logs.FetchHostLogs(ctx, client.MonitoringAndLogsApi, p)
	if err != nil {
		errors.ExitWithError("Failed to fetch logs", err)
	}
	defer fileutils.SafeClose(rc)

	// Prepare output paths
	// If the ATLAS_DOWNLOADS_DIR env variable is set, it will be used as the base directory for output files
	outDir := "logs"
	prefix := fmt.Sprintf("%s_%s", p.HostName, p.LogName)
	gzPath, err := fileutils.GenerateOutputPath(outDir, prefix, "gz")
	if err != nil {
		errors.ExitWithError("Failed to generate GZ output path", err)
	}
	txtPath, err := fileutils.GenerateOutputPath(outDir, prefix, "txt")
	if err != nil {
		errors.ExitWithError("Failed to generate TXT output path", err)
	}

	// Save compressed logs
	if err := fileutils.WriteToFile(rc, gzPath); err != nil {
		errors.ExitWithError("Failed to save compressed logs", err)
	}
	fmt.Println("Saved compressed log to", gzPath)

	// Decompress logs
	if err := fileutils.DecompressGzip(gzPath, txtPath); err != nil {
		errors.ExitWithError("Failed to decompress logs", err)
	}
	fmt.Println("Uncompressed log to", txtPath)
	// :remove-start:
	// Clean up (internal-only function)
	if err := fileutils.SafeDelete(outDir); err != nil {
		log.Printf("Cleanup error: %v", err)
	}
	fmt.Println("Deleted generated files from", outDir)
	// :remove-end:
}

// :snippet-end: [get-logs]
