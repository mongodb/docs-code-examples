package scale

import (
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// UpdateClusterSize handles the actual scaling operation by calling Atlas API
func UpdateClusterSize(ctx context.Context, api admin.ClustersApi, groupId, clusterName string,
	cluster *admin.ClusterDescription20240805, targetSize string) error {
	if cluster.ReplicationSpecs == nil || len(*cluster.ReplicationSpecs) == 0 {
		return fmt.Errorf("no replication specs found for cluster %s", clusterName)
	}

	replicationSpec := (*cluster.ReplicationSpecs)[0]
	if replicationSpec.RegionConfigs == nil || len(*replicationSpec.RegionConfigs) == 0 {
		return fmt.Errorf("no region configs found for cluster %s", clusterName)
	}

	regionConfig := (*replicationSpec.RegionConfigs)[0]

	// Create update request
	updateRequest := &admin.ClusterDescription20240805{
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{
			{
				Id: replicationSpec.Id,
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						// Keep the same provider and region names
						ProviderName: regionConfig.ProviderName,
						RegionName:   regionConfig.RegionName,
						ElectableSpecs: &admin.HardwareSpec20240805{
							// Update the instance size to the target size
							InstanceSize: admin.PtrString(targetSize),
						},
					},
				},
			},
		},
	}

	// Execute the update API call
	_, _, err := api.UpdateCluster(ctx, groupId, clusterName, updateRequest).Execute()
	if err != nil {
		return fmt.Errorf("failed to scale cluster %s: %w", clusterName, err)
	}

	return nil
}
