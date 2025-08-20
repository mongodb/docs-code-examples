package clusters

import (
	"context"
	"fmt"

	"atlas-sdk-go/internal/errors"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// ListClusterNames lists all clusters in a project and returns their names.
func ListClusterNames(ctx context.Context, sdk admin.ClustersApi, p *admin.ListClustersApiParams) ([]string, error) {
	req := sdk.ListClusters(ctx, p.GroupId)
	clusters, _, err := req.Execute()
	if err != nil {
		return nil, errors.FormatError("list clusters", p.GroupId, err)
	}

	var names []string
	if clusters != nil && clusters.Results != nil {
		for _, cluster := range *clusters.Results {
			if cluster.Name != nil {
				names = append(names, *cluster.Name)
			}
		}
	}
	return names, nil
}

// GetProcessIdForCluster retrieves the process ID for a given cluster
func GetProcessIdForCluster(ctx context.Context, sdk admin.MonitoringAndLogsApi,
	p *admin.ListAtlasProcessesApiParams, clusterName string) (string, error) {

	req := sdk.ListAtlasProcesses(ctx, p.GroupId)
	r, _, err := req.Execute()
	if err != nil {
		return "", errors.FormatError("list atlas processes", p.GroupId, err)
	}
	if r == nil || !r.HasResults() || len(r.GetResults()) == 0 {
		return "", nil
	}

	// Find the process for the specified cluster
	for _, process := range r.GetResults() {
		hostName := process.GetUserAlias()
		id := process.GetId()
		if hostName != "" && hostName == clusterName {
			if id != "" {
				return id, nil
			}
		}
	}

	return "", fmt.Errorf("no process found for cluster %s", clusterName)
}

// GetClusterSRVConnectionString returns the standard SRV connection string for a cluster.
func GetClusterSRVConnectionString(ctx context.Context, client *admin.APIClient, projectID, clusterName string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("nil atlas api client")
	}
	cluster, _, err := client.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		return "", errors.FormatError("get cluster", projectID, err)
	}
	if cluster == nil || cluster.ConnectionStrings == nil || cluster.ConnectionStrings.StandardSrv == nil {
		return "", fmt.Errorf("no standard SRV connection string found for cluster %s", clusterName)
	}
	return *cluster.ConnectionStrings.StandardSrv, nil
}
