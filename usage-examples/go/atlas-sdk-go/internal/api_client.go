// :snippet-start: api-client-function-full-example
package internal

import (
	"context"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"net/http"
)

type AtlasClient interface {
	GetHostLogs(ctx context.Context, params *admin.GetHostLogsApiParams) (io.ReadCloser, error)
	GetProcessMetrics(ctx context.Context, params *admin.GetHostMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, *http.Response, error)
	GetDiskMetrics(ctx context.Context, params *admin.GetDiskMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, *http.Response, error)
}

type HTTPClient struct {
	sdk *admin.APIClient
}

// NewAtlasClient creates a new Atlas API client using the provided SDK client
func NewAtlasClient(sdk *admin.APIClient) *HTTPClient {
	return &HTTPClient{sdk: sdk}
}

// GetHostLogs fetches MongoDB logs for the specified host in your project
func (c *HTTPClient) GetHostLogs(ctx context.Context, params *admin.GetHostLogsApiParams) (io.ReadCloser, error) {
	resp, _, err := c.sdk.MonitoringAndLogsApi.GetHostLogsWithParams(ctx, params).Execute()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetProcessMetrics fetches metrics for a specified host process in a project
func (c *HTTPClient) GetProcessMetrics(ctx context.Context, params *admin.GetHostMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	resp, r, err := c.sdk.MonitoringAndLogsApi.GetHostMeasurementsWithParams(ctx, params).Execute()
	if err != nil {
		return nil, r, err
	}
	return resp, nil, nil
}

// GetDiskMetrics fetches disk metrics for a specified disk partition in a project
func (c *HTTPClient) GetDiskMetrics(ctx context.Context, params *admin.GetDiskMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	resp, r, err := c.sdk.MonitoringAndLogsApi.GetDiskMeasurementsWithParams(ctx, params).Execute()
	if err != nil {
		return resp, r, err
	}
	return resp, nil, nil
}

// :snippet-end: [api-client-function-full-example]
