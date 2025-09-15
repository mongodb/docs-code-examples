package scale

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

type ScalingConfig struct {
	TargetTier    string  // Desired tier for scaling operations (e.g. M50)
	PreScale      bool    // Immediate scale for all clusters (e.g. planned launch/event)
	CPUThreshold  float64 // Average CPU % threshold to trigger reactive scale
	PeriodMinutes int     // Lookback window in minutes for CPU averaging
	DryRun        bool    // If true, only log intended actions without executing
}

// EvaluateDecision returns true if scaling should occur and a human-readable reason.
func EvaluateDecision(ctx context.Context, client *admin.APIClient, projectID, clusterName string, sc ScalingConfig) (bool, string) {
	// Pre-scale always wins (explicit operator intent for predictable events)
	if sc.PreScale {
		return true, "pre-scale event flag set (predictable traffic spike)"
	}

	// Reactive scaling based on sustained CPU utilization
	// Aligned with Atlas auto-scaling guidance: 75% for 1 hour triggers upscaling
	avgCPU, err := GetAverageProcessCPU(ctx, client, projectID, clusterName, sc.PeriodMinutes)
	if err != nil {
		fmt.Printf("  Warning: unable to compute average CPU for reactive scaling: %v\n", err)
		return false, "metrics unavailable for reactive scaling decision"
	}

	fmt.Printf("  Average CPU last %d minutes: %.1f%% (threshold: %.1f%%)\n",
		sc.PeriodMinutes, avgCPU, sc.CPUThreshold)

	if avgCPU > sc.CPUThreshold {
		return true, fmt.Sprintf("sustained CPU utilization %.1f%% > %.1f%% threshold over %d minutes",
			avgCPU, sc.CPUThreshold, sc.PeriodMinutes)
	}

	return false, fmt.Sprintf("CPU utilization %.1f%% below threshold %.1f%%", avgCPU, sc.CPUThreshold)
}
