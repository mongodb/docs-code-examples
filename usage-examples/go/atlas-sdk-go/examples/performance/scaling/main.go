package main

import (
	"atlas-sdk-go/internal/clusters"
	"atlas-sdk-go/internal/scale"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/errors"
)

// Constants for the scaling thresholds and instance sizes
const (
	currentInstanceSize = "M30"
	targetInstanceSize  = "M40"
	cpuMonitoringPeriod = "P1D" // Look at last 24 hours of CPU data
	scaleUpThreshold    = 70.0  // Scale up if CPU utilization is above 70%
	scaleDownThreshold  = 30.0  // Scale down if CPU utilization is below 30%
)

// CPUMetrics represents CPU utilization metrics for a cluster
type CPUMetrics struct {
	AverageCPUUsage float64
	MaxCPUUsage     float64
	SampleCount     int
}

// CPUThresholds defines the thresholds for scaling decisions
type CPUThresholds struct {
	ScaleUpThreshold   float64
	ScaleDownThreshold float64
}

// ScalingDecision represents a decision on whether to scale a cluster
type ScalingDecision struct {
	ShouldScale bool
	Direction   string
	Reason      string
}

func main() {
	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Load application context with configuration
	appCtx, err := config.LoadAppContextWithContext(ctx, "internal", false)
	if err != nil {
		errors.ExitWithError("Failed to load configuration", err)
	}

	// Initialize the Atlas API client
	client, err := auth.NewClient(appCtx.Config, appCtx.Secrets)
	if err != nil {
		errors.ExitWithError("Failed to initialize Atlas client", err)
	}
	// Get the project ID from configuration
	projectId := appCtx.Config.ProjectID
	if projectId == "" {
		errors.ExitWithError("Project ID not found in configuration", nil)
	}

	// Set up CPU thresholds based on scalability recommendations
	cpuThresholds := CPUThresholds{
		ScaleUpThreshold:   scaleUpThreshold,
		ScaleDownThreshold: scaleDownThreshold,
	}
	log.Printf("Using CPU thresholds - Scale up: %.1f%%, Scale down: %.1f%%",
		cpuThresholds.ScaleUpThreshold, cpuThresholds.ScaleDownThreshold)

	// Get list of all clusters in the project
	clusterNames, err = clusters.ListClusterNames(ctx, client, projectId)
	if err != nil {
		errors.ExitWithError("Failed to list clusters", err)
	}

	if len(clusterNames) == 0 {
		fmt.Printf("No clusters found for the project ID: %s\n", projectId)
		return
	}

	// Evaluate each cluster's details to determine eligibility for scaling
	for _, clusterName := range clusterNames {
		clusterDetails, resp, err := client.ClustersApi.GetCluster(ctx, projectId, clusterName).Execute()
		if err != nil {
			log.Printf("Error getting details for cluster %s: %v", clusterName, err)
			continue
		}

		// Check if the cluster matches the target instance size
		if !isEligibleForScaling(clusterDetails, currentInstanceSize) {
			log.Printf("Cluster %s instance size doesn't match criteria (%s), skipping",
				clusterName, currentInstanceSize)
			continue
		}

		// Get CPU metrics for the cluster
		processId, err := clusters.GetProcessIdForCluster(ctx, client, projectId, clusterName)
		if err != nil {
			log.Printf("Could not get process ID for cluster %s: %v", clusterName, err)
			continue
		}

		cpuMetrics, err := getClusterCPUMetrics(ctx, client, projectId, processId, cpuMonitoringPeriod)
		if err != nil {
			log.Printf("Could not fetch CPU metrics for cluster %s: %v", clusterName, err)
			continue
		}

		// Evaluate scaling decision based on CPU usage
		scalingDecision := evaluateCPUBasedScaling(cpuMetrics, cpuThresholds)

		log.Printf("Cluster %s - CPU: avg=%.2f%%, max=%.2f%%, samples=%d",
			clusterName, cpuMetrics.AverageCPUUsage, cpuMetrics.MaxCPUUsage, cpuMetrics.SampleCount)
		log.Printf("Scaling decision: %s", scalingDecision.Reason)

		// Perform scaling if needed
		if scalingDecision.ShouldScale && scalingDecision.Direction == "up" {
			log.Printf("Scaling cluster %s UP from %s to %s due to high CPU usage",
				clusterName, currentInstanceSize, targetInstanceSize)
			err := scale.UpdateClusterSize(ctx, client, projectId, clusterName, clusterDetails, targetInstanceSize)
			if err != nil {
				log.Printf("Error during scaling: %v", err)
			}
		} else if scalingDecision.ShouldScale && scalingDecision.Direction == "down" {
			// Define a smaller instance size for scale down
			scaleDownSize := getScaleDownSize(currentInstanceSize)
			log.Printf("Scaling cluster %s DOWN from %s to %s due to low CPU usage",
				clusterName, currentInstanceSize, scaleDownSize)
			err := scale.UpdateClusterSize(ctx, client, projectId, clusterName, clusterDetails, scaleDownSize)
			if err != nil {
				log.Printf("Error during scaling: %v", err)
			}
		} else {
			log.Printf("No scaling needed for cluster %s", clusterName)
		}
	}

	log.Println("Cluster scaling process completed successfully.")
}

func getClusterCPUMetrics(ctx context.Context, client *admin.APIClient, projectID, processID, period string) (CPUMetrics, error) {
	// Configure time window for metrics
	end := time.Now().UTC()
	start := end.Add(-24 * time.Hour) // Default to 1 day
	granularity := "PT1H"             // 1-hour granularity

	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)

	request := client.MonitoringAndLogsApi.GetProcessMeasurements(ctx, projectID, processID)
	request = request.M("CPU_USAGE")
	request = request.Granularity(granularity)
	request = request.Period(period)
	request = request.Start(startStr)
	request = request.End(endStr)

	metrics, httpResp, err := request.Execute()
	if err != nil {
		return CPUMetrics{}, fmt.Errorf("failed to get CPU metrics: %w", err)
	}

	var totalCPU float64
	var maxCPU float64
	var sampleCount int

	// Calculate average and max CPU usage
	if metrics.Measurements != nil {
		for _, measurement := range *metrics.Measurements {
			if measurement.DataPoints != nil {
				for _, dataPoint := range *measurement.DataPoints {
					if dataPoint.Value != nil {
						cpuValue := *dataPoint.Value
						totalCPU += cpuValue
						if cpuValue > maxCPU {
							maxCPU = cpuValue
						}
						sampleCount++
					}
				}
			}
		}
	}

	if sampleCount == 0 {
		return CPUMetrics{}, fmt.Errorf("no CPU metrics available")
	}

	return CPUMetrics{
		AverageCPUUsage: totalCPU / float64(sampleCount),
		MaxCPUUsage:     maxCPU,
		SampleCount:     sampleCount,
	}, nil
}

func evaluateCPUBasedScaling(metrics CPUMetrics, thresholds CPUThresholds) ScalingDecision {
	if metrics.AverageCPUUsage > thresholds.ScaleUpThreshold {
		return ScalingDecision{
			ShouldScale: true,
			Direction:   "up",
			Reason: fmt.Sprintf("Average CPU usage (%.2f%%) exceeds scale-up threshold (%.2f%%)",
				metrics.AverageCPUUsage, thresholds.ScaleUpThreshold),
		}
	}

	if metrics.AverageCPUUsage < thresholds.ScaleDownThreshold {
		return ScalingDecision{
			ShouldScale: true,
			Direction:   "down",
			Reason: fmt.Sprintf("Average CPU usage (%.2f%%) is below scale-down threshold (%.2f%%)",
				metrics.AverageCPUUsage, thresholds.ScaleDownThreshold),
		}
	}

	return ScalingDecision{
		ShouldScale: false,
		Direction:   "",
		Reason: fmt.Sprintf("CPU usage (%.2f%%) is within normal range (%.2f%% - %.2f%%)",
			metrics.AverageCPUUsage, thresholds.ScaleDownThreshold, thresholds.ScaleUpThreshold),
	}
}

func isEligibleForScaling(cluster *admin.ClusterDescription, currentSize string) bool {
	if cluster.ReplicationSpecs == nil || len(*cluster.ReplicationSpecs) == 0 {
		return false
	}

	replicationSpec := (*cluster.ReplicationSpecs)[0]
	if replicationSpec.RegionConfigs == nil || len(*replicationSpec.RegionConfigs) == 0 {
		return false
	}

	regionConfig := (*replicationSpec.RegionConfigs)[0]
	if regionConfig.ElectableSpecs == nil || regionConfig.ElectableSpecs.InstanceSize == nil {
		return false
	}

	return *regionConfig.ElectableSpecs.InstanceSize == currentSize
}

func getScaleDownSize(currentSize string) string {
	scaleDownMap := map[string]string{
		"M40": "M30",
		"M30": "M20",
		"M20": "M10",
		"M50": "M40",
		"M60": "M50",
		"M80": "M60",
	}

	if downSize, exists := scaleDownMap[currentSize]; exists {
		return downSize
	}
	return currentSize
}
