package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal"
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/logs"
)

func main() {
	_ = godotenv.Load()
	secrets, cfg, err := config.LoadAll("configs/config.json")
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	sdk, err := auth.NewClient(cfg, secrets)
	if err != nil {
		log.Fatalf("client init: %v", err)
	}

	ctx := context.Background()
	p := &admin.GetHostLogsApiParams{
		GroupId:  cfg.ProjectID,
		HostName: cfg.HostName,
		LogName:  "mongodb",
	}
	ts := time.Now().Format("20060102_150405")
	base := fmt.Sprintf("%s_%s_%s", p.HostName, p.LogName, ts)
	outDir := "logs"
	os.MkdirAll(outDir, 0o755)
	gzPath := filepath.Join(outDir, base+".gz")
	txtPath := filepath.Join(outDir, base+".txt")

	rc, err := logs.FetchHostLogs(ctx, sdk.MonitoringAndLogsApi, p)
	if err != nil {
		log.Fatalf("download logs: %v", err)
	}
	defer internal.SafeClose(rc)

	if err := logs.WriteToFile(rc, gzPath); err != nil {
		log.Fatalf("save gz: %v", err)
	}
	fmt.Println("Saved compressed log to", gzPath)

	if err := logs.DecompressGzip(gzPath, txtPath); err != nil {
		log.Fatalf("decompress: %v", err)
	}
	fmt.Println("Uncompressed log to", txtPath)
}

