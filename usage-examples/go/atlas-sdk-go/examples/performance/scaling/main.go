// :snippet-start: scale-cluster-programmatically
// :state-remove-start: copy
// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
// :state-remove-end: [copy]
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/clusters"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/metrics"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// This example shows how to programmatically trigger a cluster tier change (e.g., M30 -> M50)
// using the Atlas Go SDK. It demonstrates two example conditions inspired by scalability.md:
//  1. Pre-scale for an expected event/spike: set PRE_SCALE_EVENT=true
//  2. Reactive scale when average CPU > CPU_THRESHOLD (default 70%) over CPU_PERIOD_MINUTES (default 15)
//
// Inputs via environment variables:
//   - SCALE_CLUSTER_NAME:   the name of the cluster to scale (required)
//   - SCALE_TO_TIER:        the target instance size, e.g., M50 (default: M50)
//   - PRE_SCALE_EVENT:      if "true", trigger scaling immediately (default: false)
//   - CPU_THRESHOLD:        average CPU percentage threshold to trigger scaling (default: 70)
//   - CPU_PERIOD_MINUTES:   lookback minutes for CPU average (default: 15)
//   - CONFIG_PATH:          path to JSON config file (defaults to configs/config.json)
func main() {
	envFile := ".env.production"
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: could not load %s file: %v", envFile, err)
	}

	secrets, cfg, err := config.LoadAllFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
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

	// replace with GetClusterName()
	clusterName := strings.TrimSpace(os.Getenv("SCALE_CLUSTER_NAME"))
	if clusterName == "" {
		log.Fatal("SCALE_CLUSTER_NAME is required. Set it to the Atlas cluster you want to scale.")
	}
	//
	targetTier := strings.TrimSpace(os.Getenv("SCALE_TO_TIER"))
	if targetTier == "" {
		targetTier = "M50"
	}

	preScale := strings.EqualFold(strings.TrimSpace(os.Getenv("PRE_SCALE_EVENT")), "true")
	cpuThreshold := 70.0
	if v := strings.TrimSpace(os.Getenv("CPU_THRESHOLD")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cpuThreshold = f
		}
	}
	periodMinutes := 15
	if v := strings.TrimSpace(os.Getenv("CPU_PERIOD_MINUTES")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			periodMinutes = n
		}
	}

	// Fetch current cluster configuration
	cur, _, err := client.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		log.Fatalf("Failed to get cluster '%s': %v", clusterName, err)
	}

	// Determine the current instance size (electable specs of the first region/shard)
	currentTier := ""
	if cur != nil && cur.HasReplicationSpecs() {
		repl := cur.GetReplicationSpecs()
		if len(repl) > 0 {
			rcs := repl[0].GetRegionConfigs()
			if len(rcs) > 0 && rcs[0].HasElectableSpecs() {
				if sz, ok := rcs[0].GetElectableSpecsOk(); ok && sz.HasInstanceSize() {
					currentTier = sz.GetInstanceSize()
				}
			}
		}
	}

	fmt.Printf("Project: %s\nCluster: %s\nCurrent tier: %s\nTarget tier: %s\n", projectID, clusterName, currentTier, targetTier)
	if strings.EqualFold(currentTier, targetTier) {
		fmt.Println("No action: cluster already at target tier.")
		return
	}

	// Example conditions
	shouldScale := false
	reason := ""
	if preScale {
		shouldScale = true
		reason = "pre-scale for upcoming event"
	} else {
		// Try reactive condition: avg CPU over lookback > threshold
		avgCPU, err := averageProcessCPU(ctx, client, projectID, clusterName, periodMinutes)
		if err != nil {
			fmt.Printf("Warning: couldn't evaluate CPU condition: %v\n", err)
		} else {
			fmt.Printf("Average CPU over last %d minutes: %.1f%% (threshold %.1f%%)\n", periodMinutes, avgCPU, cpuThreshold)
			if avgCPU > cpuThreshold {
				shouldScale = true
				reason = fmt.Sprintf("avg CPU %.1f%% > %.1f%%", avgCPU, cpuThreshold)
			}
		}
	}

	if !shouldScale {
		fmt.Println("Conditions not met; no scaling requested.")
		return
	}

	fmt.Printf("Triggering scale to %s (%s) using Atlas Go SDK...\n", targetTier, reason)

	// Build an update payload based on current config and override instance size(s)
	payload := admin.NewClusterDescription20240805()
	if cur.HasReplicationSpecs() {
		repl := cur.GetReplicationSpecs()
		for i := range repl {
			rc := repl[i].GetRegionConfigs()
			for j := range rc {
				// Update electable specs
				if rc[j].HasElectableSpecs() {
					es := rc[j].GetElectableSpecs()
					es.SetInstanceSize(targetTier)
					rc[j].SetElectableSpecs(es)
				}
				// Keep read-only and analytics tiers in sync if present
				if rc[j].HasReadOnlySpecs() {
					ros := rc[j].GetReadOnlySpecs()
					ros.SetInstanceSize(targetTier)
					rc[j].SetReadOnlySpecs(ros)
				}
				if rc[j].HasAnalyticsSpecs() {
					as := rc[j].GetAnalyticsSpecs()
					as.SetInstanceSize(targetTier)
					rc[j].SetAnalyticsSpecs(as)
				}
			}
			repl[i].SetRegionConfigs(rc)
		}
		payload.SetReplicationSpecs(repl)
	}

	_, _, err = client.ClustersApi.UpdateCluster(ctx, projectID, clusterName, payload).Execute()
	if err != nil {
		log.Fatalf("Failed to request scaling via SDK: %v", err)
	}

	fmt.Println("Scaling request submitted. Atlas performs scaling in a rolling fashion to avoid downtime.")
	fmt.Println("Monitor progress in the Atlas UI (Deployments) or poll the API for cluster state.")
}

// averageProcessCPU fetches host CPU metrics and returns a simple average percentage over the lookback period.
func averageProcessCPU(ctx context.Context, client *admin.APIClient, projectID, clusterName string, periodMinutes int) (float64, error) {
	// Resolve process ID for this cluster
	procID, err := clusters.GetProcessIdForCluster(ctx, client.MonitoringAndLogsApi, &admin.ListAtlasProcessesApiParams{GroupId: projectID}, clusterName)
	if err != nil {
		return 0, err
	}
	if procID == "" {
		return 0, fmt.Errorf("no process found for cluster %s", clusterName)
	}

	granularity := "PT1M"
	period := fmt.Sprintf("PT%vM", periodMinutes)
	metricsList := []string{"PROCESS_CPU_USER"}
	m, err := metrics.FetchProcessMetrics(ctx, client.MonitoringAndLogsApi, &admin.GetHostMeasurementsApiParams{
		GroupId:     projectID,
		ProcessId:   procID,
		Granularity: &granularity,
		Period:      &period,
		M:           &metricsList,
	})
	if err != nil {
		return 0, err
	}

	// Compute average over all datapoints of the first measurement
	if m == nil || !m.HasMeasurements() {
		return 0, fmt.Errorf("no measurements returned")
	}
	meas := m.GetMeasurements()
	if len(meas) == 0 || !meas[0].HasDataPoints() {
		return 0, fmt.Errorf("no datapoints returned")
	}
	total := 0.0
	count := 0.0
	for _, dp := range meas[0].GetDataPoints() {
		if dp.HasValue() {
			v := float64(dp.GetValue())
			total += v
			count++
		}
	}
	if count == 0 {
		return 0, fmt.Errorf("no datapoint values")
	}
	// Atlas returns CPU fraction (0..1) for some metrics; if values look like 0..1, convert to %.
	avg := total / count
	if avg <= 1.0 {
		avg = avg * 100.0
	}
	return avg, nil
}

// :snippet-end: [scale-cluster-programmatically]
