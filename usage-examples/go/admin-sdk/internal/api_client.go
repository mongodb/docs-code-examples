package internal

import (
	"context"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"net/http"
)

// HTTPClient is a thin wrapper around admin.APIClient for logs.
type HTTPClient struct {
	sdk *admin.APIClient
}

// NewAtlasClient initializes a new LogsService using admin.APIClient.
func NewAtlasClient(sdk *admin.APIClient) *HTTPClient {
	return &HTTPClient{sdk: sdk}
}

// GetHostLogs fetches logs from MongoDB Atlas.
func (c *HTTPClient) GetHostLogs(ctx context.Context, params *admin.GetHostLogsApiParams) (io.ReadCloser, error) {
	resp, _, err := c.sdk.MonitoringAndLogsApi.GetHostLogsWithParams(ctx, params).Execute()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *HTTPClient) GetProcessMetrics(ctx context.Context, params *admin.GetHostMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	resp, r, err := c.sdk.MonitoringAndLogsApi.GetHostMeasurementsWithParams(ctx, params).Execute()
	if err != nil {
		return nil, r, err
	}
	return resp, nil, nil
}

func (c *HTTPClient) GetDiskMetrics(ctx context.Context, params *admin.GetDiskMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, *http.Response, error) {
	resp, r, err := c.sdk.MonitoringAndLogsApi.GetDiskMeasurementsWithParams(ctx, params).Execute()
	if err != nil {
		return resp, r, err
	}
	return resp, nil, nil
}

//
//// OrganizationService provides organization-related methods
//type OrganizationService interface {
//	ListOrganizations(ctx context.Context) (*admin.PaginatedOrganization, error)
//}
//
//// ProjectService provides project-related methods
//type ProjectService interface {
//	ListProjects(ctx context.Context) (*admin.PaginatedAtlasGroup, error)
//
//	GetProcesses(ctx context.Context, params *admin.ListAtlasProcessesApiParams) (*admin.PaginatedHostViewAtlas, error)
//}
//
//// ClusterService provides cluster-related methods
//type ClusterService interface {
//	ListClusters(ctx context.Context) (*admin.PaginatedOrgGroup, error)
//}
//
//// LogsService provides log retrieval methods
//type LogsService interface {
//	GetHostLogs(ctx context.Context, params *admin.GetHostLogsApiParams) (io.ReadCloser, error)
//
//}
//
//type MetricsService interface {
//	GetProcessMetrics(ctx context.Context, params *admin.GetHostMeasurementsApiParams) (*admin.ApiMeasurementsGeneralViewAtlas, *admin.APIResponse, error)
//
//	GetDiskMetrics(ctx context.Context, params *admin.GetDiskMeasurementsApiParams) (*admin.MeasurementDiskPartition, *admin.APIResponse, error)
//
//
//// AtlasClient combines all services into a single interface
//type AtlasClient interface {
//	OrganizationService
//	ProjectService
//	ClusterService
//	LogsService
//	MetricsService
//}
//
//// HTTPAtlasClient is the concrete implementation of AtlasClient.
//type HTTPAtlasClient struct {
//	sdk *admin.APIClient
//}
//
//func NewAtlasClient(sdk *admin.APIClient) *HTTPAtlasClient {
//	return &HTTPAtlasClient{sdk: sdk}
//}
