package metrics

import (
	"atlas-sdk-go/examples/internal"
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
		return nil, internal.FormatError("fetch disk metrics", p.PartitionName, err)
	}
	if r == nil || !r.HasMeasurements() || len(r.GetMeasurements()) == 0 {
		return nil, fmt.Errorf("no metrics for partition %q on process %q", p.PartitionName, p.ProcessId)
	}
	return r, nil
}
