package logs

import (
	"atlas-sdk-go/internal"
	"context"
	"io"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// FetchHostLogs calls the Atlas SDK and returns the raw, compressed log stream.
func FetchHostLogs(ctx context.Context, sdk admin.MonitoringAndLogsApi, p *admin.GetHostLogsApiParams) (io.ReadCloser, error) {
	req := sdk.GetHostLogs(ctx, p.GroupId, p.HostName, p.LogName)
	rc, _, err := req.Execute()
	if err != nil {
		return nil, internal.FormatAPIError("fetch logs", p.HostName, err)
	}
	return rc, nil
}
