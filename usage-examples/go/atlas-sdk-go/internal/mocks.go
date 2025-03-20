package internal

import (
	"context"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"net/http"
	"strings"
)

// NOTE: We're using mocked tests because Monitoring and Logging functionality requires a dedicated cluster (M10+)

// MockAtlasClient is a fake implementation of AtlasClient for testing.
type MockAtlasClient struct {
	FakeHostLogsResponse       string
	FakeHostLogsError          error
	FakeProcessMetricsResponse *admin.ApiMeasurementsGeneralViewAtlas
	FakeProcessMetricsError    error
	FakeDiskMetricsResponse    *admin.ApiMeasurementsGeneralViewAtlas
	FakeDiskMetricsError       error
}

func (m *MockAtlasClient) GetHostLogs(context.Context, *admin.GetHostLogsApiParams) (io.ReadCloser, error) {
	if m.FakeHostLogsError != nil {
		return nil, m.FakeHostLogsError
	}
	return io.NopCloser(strings.NewReader(m.FakeHostLogsResponse)), nil
}

// GetProcessMetrics returns fake process metrics or an error.
func (m *MockAtlasClient) GetProcessMetrics(context.Context, *admin.GetHostMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	if m.FakeProcessMetricsError != nil {
		return nil, nil, m.FakeProcessMetricsError
	}
	return m.FakeProcessMetricsResponse, nil, nil
}

// GetDiskMetrics returns fake disk metrics or an error.
func (m *MockAtlasClient) GetDiskMetrics(context.Context, *admin.GetDiskMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	if m.FakeDiskMetricsError != nil {
		return nil, nil, m.FakeDiskMetricsError
	}
	return m.FakeDiskMetricsResponse, nil, nil
}
