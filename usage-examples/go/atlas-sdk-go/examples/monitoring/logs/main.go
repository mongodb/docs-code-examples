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
	"atlas-sdk-go/internal/fileutils"
	"atlas-sdk-go/internal/logs"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func main() {
	_ = godotenv.Load() // or godotenv.Load(".env.development")

	ctx := context.Background()
	envName := config.Environment("test")    // Cast string to config.Environment
	configPath := "configs/config.test.json" // Optional explicit config file path; if empty, uses environment-based path
	secrets, cfg, err := config.LoadAll(envName, configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
	}

	client, err := auth.NewClient(ctx, &cfg, &secrets) // Pass pointers
	if err != nil {
		log.Fatalf("Failed to initialize authentication client: %v", err)
	}

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
		log.Fatalf("Failed to fetch logs: %v", err)
	}
	defer fileutils.SafeClose(rc)

	// Prepare output paths
	// If the ATLAS_DOWNLOADS_DIR env variable is set, it will be used as the base directory for output files
	outDir := "logs"
	prefix := fmt.Sprintf("%s_%s", p.HostName, p.LogName)
	gzPath, err := fileutils.GenerateOutputPath(outDir, prefix, "gz")
	if err != nil {
		log.Fatalf("Failed to generate GZ output path: %v", err)
	}
	txtPath, err := fileutils.GenerateOutputPath(outDir, prefix, "txt")
	if err != nil {
		log.Fatalf("Failed to generate TXT output path: %v", err)
	}

	// Save compressed logs
	if err := fileutils.WriteToFile(rc, gzPath); err != nil {
		log.Fatalf("Failed to save compressed logs: %v", err)
	}
	fmt.Println("Saved compressed log to", gzPath)

	// Decompress logs
	if err := fileutils.DecompressGzip(gzPath, txtPath); err != nil {
		log.Fatalf("Failed to decompress logs: %v", err)
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
