package scale_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// MockClustersApi mocks the ClustersApi interface
type MockClustersApi struct {
	mock.Mock
}

// UpdateCluster mocks the UpdateCluster method
func (m *MockClustersApi) UpdateCluster(ctx context.Context, groupId string, clusterName string, clusterRequest *admin.ClusterDescription20240805) admin.ApiUpdateClusterRequest {
	args := m.Called(ctx, groupId, clusterName, clusterRequest)
	return args.Get(0).(admin.ApiUpdateClusterRequest)
}

// MockApiUpdateClusterRequest mocks the ApiUpdateClusterRequest
type MockApiUpdateClusterRequest struct {
	mock.Mock
}

// Execute mocks the Execute method
func (m *MockApiUpdateClusterRequest) Execute() (*admin.ClusterDescription20240805, *http.Response, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Get(1).(*http.Response), args.Error(2)
	}
	return args.Get(0).(*admin.ClusterDescription20240805), args.Get(1).(*http.Response), args.Error(2)
}

func SuccessfulClusterSizeUpdate(t *testing.T) {
	// Setup
	mockApi := new(MockClustersApi)
	mockRequest := new(MockApiUpdateClusterRequest)

	replicationSpecs := []admin.ReplicationSpec20240805{
		{
			Id: admin.PtrString("rs1"),
			RegionConfigs: &[]admin.CloudRegionConfig20240805{
				{
					ProviderName: admin.PtrString("AWS"),
					RegionName:   admin.PtrString("US_EAST_1"),
					ElectableSpecs: &admin.HardwareSpec20240805{
						InstanceSize: admin.PtrString("M10"),
					},
				},
			},
		},
	}

	cluster := &admin.ClusterDescription20240805{
		ReplicationSpecs: &replicationSpecs,
	}

	mockRequest.On("Execute").Return(&admin.ClusterDescription20240805{}, nil, nil)
	mockApi.On("UpdateCluster", mock.Anything, "groupId", "clusterName", mock.Anything).Return(mockRequest)

	// Act
	err := UpdateClusterSize(context.Background(), mockApi, "groupId", "clusterName", cluster, "M20")

	// Assert
	assert.NoError(t, err)
	mockApi.AssertExpectations(t)
	mockRequest.AssertExpectations(t)
}

func FailsWhenNoReplicationSpecs(t *testing.T) {
	// Setup
	mockApi := new(MockClustersApi)
	cluster := &admin.ClusterDescription20240805{
		ReplicationSpecs: nil,
	}

	// Act
	err := UpdateClusterSize(context.Background(), mockApi, "groupId", "clusterName", cluster, "M20")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no replication specs found")
}

func FailsWhenEmptyReplicationSpecs(t *testing.T) {
	// Setup
	mockApi := new(MockClustersApi)
	emptySpecs := []admin.ReplicationSpec20240805{}
	cluster := &admin.ClusterDescription20240805{
		ReplicationSpecs: &emptySpecs,
	}

	// Act
	err := UpdateClusterSize(context.Background(), mockApi, "groupId", "clusterName", cluster, "M20")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no replication specs found")
}

func FailsWhenNoRegionConfigs(t *testing.T) {
	// Setup
	mockApi := new(MockClustersApi)
	replicationSpecs := []admin.ReplicationSpec20240805{
		{
			Id:            admin.PtrString("rs1"),
			RegionConfigs: nil,
		},
	}

	cluster := &admin.ClusterDescription20240805{
		ReplicationSpecs: &replicationSpecs,
	}

	// Act
	err := UpdateClusterSize(context.Background(), mockApi, "groupId", "clusterName", cluster, "M20")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no region configs found")
}

func FailsWhenEmptyRegionConfigs(t *testing.T) {
	// Setup
	mockApi := new(MockClustersApi)
	emptyRegions := []admin.CloudRegionConfig20240805{}
	replicationSpecs := []admin.ReplicationSpec20240805{
		{
			Id:            admin.PtrString("rs1"),
			RegionConfigs: &emptyRegions,
		},
	}

	cluster := &admin.ClusterDescription20240805{
		ReplicationSpecs: &replicationSpecs,
	}

	// Act
	err := UpdateClusterSize(context.Background(), mockApi, "groupId", "clusterName", cluster, "M20")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no region configs found")
}

func FailsWhenApiCallErrors(t *testing.T) {
	// Setup
	mockApi := new(MockClustersApi)
	mockRequest := new(MockApiUpdateClusterRequest)

	replicationSpecs := []admin.ReplicationSpec20240805{
		{
			Id: admin.PtrString("rs1"),
			RegionConfigs: &[]admin.CloudRegionConfig20240805{
				{
					ProviderName: admin.PtrString("AWS"),
					RegionName:   admin.PtrString("US_EAST_1"),
				},
			},
		},
	}

	cluster := &admin.ClusterDescription20240805{
		ReplicationSpecs: &replicationSpecs,
	}

	mockRequest.On("Execute").Return(nil, nil, errors.New("API error"))
	mockApi.On("UpdateCluster", mock.Anything, "groupId", "clusterName", mock.Anything).Return(mockRequest)

	// Act
	err := UpdateClusterSize(context.Background(), mockApi, "groupId", "clusterName", cluster, "M20")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to scale cluster")
	mockApi.AssertExpectations(t)
	mockRequest.AssertExpectations(t)
}
