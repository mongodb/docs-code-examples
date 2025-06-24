package metrics

import (
	"context"
	"fmt"

	"atlas-sdk-go/internal"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// FetchProcessMetrics returns measurements for a specified host process
func FetchProcessMetrics(ctx context.Context, sdk admin.MonitoringAndLogsApi, p *admin.GetHostMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, error) {
	req := sdk.GetHostMeasurements(ctx, p.GroupId, p.ProcessId)
	req = req.Granularity(*p.Granularity).Period(*p.Period).M(*p.M)

	r, _, err := req.Execute()
	if err != nil {
		return nil, internal.FormatAPIError("fetch process metrics", p.GroupId, err)
	}
	if r == nil || !r.HasMeasurements() || len(r.GetMeasurements()) == 0 {
		return nil, fmt.Errorf("no metrics for process %q", p.ProcessId)
	}
	return r, nil
}
