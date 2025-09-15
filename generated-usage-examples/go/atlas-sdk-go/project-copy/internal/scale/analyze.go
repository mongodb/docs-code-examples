package scale

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	clusterutils "atlas-sdk-go/internal/clusters"
	"atlas-sdk-go/internal/metrics"
)

// GetAverageProcessCPU fetches host CPU metrics and returns a simple average percentage over the lookback period.
func GetAverageProcessCPU(ctx context.Context, client *admin.APIClient, projectID, clusterName string, periodMinutes int) (float64, error) {
	// Defensive validation so examples/tests can pass nil or bad inputs without panics
	if client == nil {
		return 0, fmt.Errorf("nil atlas client")
	}
	if projectID == "" {
		return 0, fmt.Errorf("empty project id")
	}
	if clusterName == "" {
		return 0, fmt.Errorf("empty cluster name")
	}
	if periodMinutes <= 0 {
		return 0, fmt.Errorf("invalid period minutes: %d", periodMinutes)
	}

	procID, err := clusterutils.GetProcessIdForCluster(ctx, client.MonitoringAndLogsApi, &admin.ListAtlasProcessesApiParams{GroupId: projectID}, clusterName)
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
		return 0, fmt.Errorf("no usable datapoint values")
	}
	avg := total / count
	// Convert fractional to % if needed
	if avg <= 1.0 {
		avg *= 100.0
	}
	return avg, nil
}
