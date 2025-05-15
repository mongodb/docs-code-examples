package metrics

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// FetchDiskMetrics returns measurements for a specified disk partition
func FetchDiskMetrics(ctx context.Context, sdk admin.MonitoringAndLogsApi, p *admin.GetDiskMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, error) {
	req := sdk.GetDiskMeasurements(ctx, p.GroupId, p.PartitionName, p.ProcessId)
	req = req.Granularity(*p.Granularity).Period(*p.Period).M(*p.M)

	r, _, err := req.Execute()
	if err != nil {
		if apiErr, ok := admin.AsError(err); ok {
			return nil, fmt.Errorf("fetch disk metrics: %w – %s", err, apiErr.GetDetail())
		}
		return nil, fmt.Errorf("fetch disk metrics: %w", err)
	}
	if r == nil || !r.HasMeasurements() {
		return nil, fmt.Errorf("no metrics for partition %q on process %q",
			p.PartitionName, p.ProcessId)
	}
	return r, nil
}
