package metrics

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// FetchProcessMetrics returns measurements for a specified host process
func FetchProcessMetrics(ctx context.Context, sdk admin.MonitoringAndLogsApi, p *admin.GetHostMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, error) {
	req := sdk.GetHostMeasurements(ctx, p.GroupId, p.ProcessId)
	req = req.Granularity(*p.Granularity).Period(*p.Period).M(*p.M)

	r, _, err := req.Execute()
	if err != nil {
		if apiErr, ok := admin.AsError(err); ok {
			return nil, fmt.Errorf("failed to fetch process metrics: %w – %s", err, apiErr.GetDetail())
		}
		return nil, fmt.Errorf("failed to fetch process metrics: %w", err)
	}
	if r == nil || !r.HasMeasurements() {
		return nil, fmt.Errorf("no metrics for process %q", p.ProcessId)
	}
	return r, nil
}
