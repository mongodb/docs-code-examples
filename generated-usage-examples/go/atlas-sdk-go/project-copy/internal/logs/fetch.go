package logs

import (
	"context"
	"fmt"
	"io"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// FetchHostLogs calls the Atlas SDK and returns the raw, compressed log stream.
func FetchHostLogs(
	ctx context.Context,
	sdk admin.MonitoringAndLogsApi,
	p *admin.GetHostLogsApiParams,
) (io.ReadCloser, error) {
	req := sdk.GetHostLogs(ctx, p.GroupId, p.HostName, p.LogName)
	rc, _, err := req.Execute()
	if err != nil {
		if apiErr, ok := admin.AsError(err); ok {
			return nil, fmt.Errorf("failed to fetch logs: %w – %s", err, apiErr.GetDetail())
		}
		return nil, fmt.Errorf("failed to fetch logs: %w", err)
	}
	return rc, nil
}
