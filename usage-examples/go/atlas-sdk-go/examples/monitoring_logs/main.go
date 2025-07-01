// :snippet-start: get-logs
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal/fileutils"
	"atlas-sdk-go/internal/logs"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}
	secrets, cfg, err := config.LoadAll("configs/config.json")
	if err != nil {
		log.Fatalf("config: failed to load file: %v", err)
	}

	sdk, err := auth.NewClient(cfg, secrets)
	if err != nil {
		log.Fatalf("auth: failed client init: %v", err)
	}

	ctx := context.Background()
	p := &admin.GetHostLogsApiParams{
		GroupId:  cfg.ProjectID,
		HostName: cfg.HostName,
		LogName:  "mongodb",
	}

	outDir := "logs"
	prefix := fmt.Sprintf("%s_%s", p.HostName, p.LogName)
	gzPath, err := fileutils.GenerateOutputPath(outDir, prefix, "gz")
	if err != nil {
		log.Fatalf("common: failed to generate output path: %v", err)
	}
	txtPath, err := fileutils.GenerateOutputPath(outDir, prefix, ".txt")
	if err != nil {
		log.Fatalf("common: failed to generate output path: %v", err)
	}

	rc, err := logs.FetchHostLogs(ctx, sdk.MonitoringAndLogsApi, p)
	if err != nil {
		log.Fatalf("logs: failed to fetch logs: %v", err)
	}
	defer fileutils.SafeClose(rc)

	if err := fileutils.WriteToFile(rc, gzPath); err != nil {
		log.Fatalf("fileutils: failed to save gz: %v", err)
	}
	fmt.Println("Saved compressed log to", gzPath)

	if err := fileutils.DecompressGzip(gzPath, txtPath); err != nil {
		log.Fatalf("fileutils: failed to decompress gz: %v", err)
	}
	fmt.Println("Uncompressed log to", txtPath)
	// :remove-start:
	// NOTE: Internal-only function to clean up any downloaded files
	if err := fileutils.SafeDelete(outDir); err != nil {
		log.Printf("Cleanup error: %v", err)
	}
	fmt.Println("Deleted generated files from", outDir)
	// :remove-end:
}

// :snippet-end: [get-logs]
