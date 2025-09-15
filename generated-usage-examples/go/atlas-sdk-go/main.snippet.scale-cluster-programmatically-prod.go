// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"atlas-sdk-go/internal/auth"
	clusterutils "atlas-sdk-go/internal/clusters"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/scale"

	"github.com/joho/godotenv"
)

func main() {
	envFile := ".env.production"
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: could not load %s file: %v", envFile, err)
	}

	secrets, cfg, err := config.LoadAllFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	client, err := auth.NewClient(ctx, cfg, secrets)
	if err != nil {
		log.Fatalf("Failed to initialize authentication client: %v", err)
	}

	projectID := cfg.ProjectID
	if projectID == "" {
		log.Fatal("Failed to find Project ID in configuration")
	}

	// Based on the env settings, perform the following programmatic scaling:
	//   - Pre-scale ahead of a known traffic spike (e.g. planned bulk inserts)
	//   - Reactive scale when sustained compute utilization exceeds a threshold
	//
	// NOTE: Prefer Atlas built-in auto-scaling for gradual growth. Use programmatic scaling for exceptional events or custom logic.
	scaling := loadScalingConfigFromEnv()
	fmt.Printf("Starting scaling analysis for project: %s\n", projectID)
	fmt.Printf("Configuration - Target tier: %s, Pre-scale: %v, CPU threshold: %.1f%%, Period: %d min, Dry run: %v\n",
		scaling.TargetTier, scaling.PreScale, scaling.CPUThreshold, scaling.PeriodMinutes, scaling.DryRun)

	// Get all clusters in the project
	clusterList, _, err := client.ClustersApi.ListClusters(ctx, projectID).Execute()
	if err != nil {
		log.Fatalf("Failed to list clusters: %v", err)
	}

	clusters := clusterList.GetResults()
	fmt.Printf("\nFound %d clusters to analyze for scaling\n", len(clusters))

	// Track scaling operations across all clusters
	scalingCandidates := 0
	successfulScales := 0
	failedScales := 0
	skippedClusters := 0

	for _, cluster := range clusters {
		clusterName := cluster.GetName()
		fmt.Printf("\n=== Analyzing cluster: %s ===\n", clusterName)

		// Skip clusters that are not in IDLE state
		if cluster.HasStateName() && cluster.GetStateName() != "IDLE" {
			fmt.Printf("- Skipping cluster %s: not in IDLE state (current: %s)\n", clusterName, cluster.GetStateName())
			skippedClusters++
			continue
		}

		// Extract current tier
		currentTier, err := clusterutils.ExtractInstanceSize(&cluster)
		if err != nil {
			fmt.Printf("- Skipping cluster %s: failed to extract current tier: %v\n", clusterName, err)
			skippedClusters++
			continue
		}

		fmt.Printf("- Current tier: %s, Target tier: %s\n", currentTier, scaling.TargetTier)

		// Skip if already at target tier
		if strings.EqualFold(currentTier, scaling.TargetTier) {
			fmt.Printf("- No action needed: cluster already at target tier %s\n", scaling.TargetTier)
			continue
		}

		// Evaluate scaling decision
		shouldScale, reason := scale.EvaluateDecision(ctx, client, projectID, clusterName, scaling)
		if !shouldScale {
			fmt.Printf("- Conditions not met: %s\n", reason)
			continue
		}

		scalingCandidates++
		fmt.Printf("- Scaling decision: proceed -> %s\n", reason)

		if scaling.DryRun {
			fmt.Printf("- DRY_RUN=true: would scale cluster %s from %s to %s\n",
				clusterName, currentTier, scaling.TargetTier)
			successfulScales++
			continue
		}

		// Execute scaling operation
		if err := scale.ExecuteClusterScaling(ctx, client, projectID, clusterName, &cluster, scaling.TargetTier); err != nil {
			fmt.Printf("- ERROR: Failed to scale cluster %s: %v\n", clusterName, err)
			failedScales++
			continue
		}

		fmt.Printf("- Successfully initiated scaling for cluster %s from %s to %s\n",
			clusterName, currentTier, scaling.TargetTier)
		successfulScales++
	}

	// Summary
	fmt.Printf("\n=== Scaling Operation Summary ===\n")
	fmt.Printf("Total clusters analyzed: %d\n", len(clusters))
	fmt.Printf("Scaling candidates identified: %d\n", scalingCandidates)
	fmt.Printf("Successful scaling operations: %d\n", successfulScales)
	fmt.Printf("Failed scaling operations: %d\n", failedScales)
	fmt.Printf("Skipped clusters: %d\n", skippedClusters)

	if failedScales > 0 {
		fmt.Printf("WARNING: %d of %d scaling operations failed\n", failedScales, scalingCandidates)
	}

	if successfulScales > 0 && !scaling.DryRun {
		fmt.Println("\nAtlas will perform rolling resizes with zero-downtime semantics.")
		fmt.Println("Monitor status in the Atlas UI or poll cluster states until STATE_NAME becomes IDLE.")
	}

	fmt.Println("Scaling analysis and operations completed.")
}

// loadScalingConfigFromEnv reads scaling configuration from environment variables with defaults:
//
//	SCALE_TO_TIER        target tier for scaling ops (default: M50)
//	PRE_SCALE_EVENT      "true" triggers immediate scale for all clusters (default: false)
//	CPU_THRESHOLD        avg CPU % threshold to trigger scaling (default: 75, aligned with Atlas auto-scaling)
//	CPU_PERIOD_MINUTES   minutes lookback for CPU avg (default: 60, aligned with Atlas)
//	DRY_RUN              if "true", do not execute scaling operations (default: false)
func loadScalingConfigFromEnv() scale.ScalingConfig {
	cfg := scale.ScalingConfig{
		TargetTier:    defaultIfBlank(strings.TrimSpace(os.Getenv("SCALE_TO_TIER")), "M50"),
		PreScale:      strings.EqualFold(strings.TrimSpace(os.Getenv("PRE_SCALE_EVENT")), "false"),
		CPUThreshold:  75.0,
		PeriodMinutes: 60,
		DryRun:        strings.EqualFold(strings.TrimSpace(os.Getenv("DRY_RUN")), "false"),
	}

	if v := strings.TrimSpace(os.Getenv("CPU_THRESHOLD")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			cfg.CPUThreshold = f
		}
	}
	if v := strings.TrimSpace(os.Getenv("CPU_PERIOD_MINUTES")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.PeriodMinutes = n
		}
	}
	return cfg
}

func defaultIfBlank(v, d string) string {
	if v == "" {
		return d
	}
	return v
}

