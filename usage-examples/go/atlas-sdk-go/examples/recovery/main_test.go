package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal/data/recovery"
)

// testClient helper replicates pattern from internal tests.
func testClient(t *testing.T, handler http.HandlerFunc) *admin.APIClient {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := admin.NewClient(admin.UseBaseURL(server.URL))
	require.NoError(t, err)
	return client
}

func TestExecuteDataDeletionRestore_Seam(t *testing.T) {
	ctx := context.Background()
	opts := recovery.DrOptions{ProjectID: "proj", ClusterName: "ClusterA", SnapshotID: "snap1"}

	// Dry-run path (should not call API)
	{
		var called atomic.Bool
		client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			called.Store(true)
			w.WriteHeader(http.StatusInternalServerError)
		})
		msg, err := executeDataDeletionRestore(ctx, client, recovery.DrOptions{ProjectID: opts.ProjectID, ClusterName: opts.ClusterName, SnapshotID: opts.SnapshotID, DryRun: true})
		require.NoError(t, err)
		assert.Contains(t, msg, "(dry-run)")
		assert.False(t, called.Load(), "API must not be invoked for dry-run")
	}

	// Success path
	{
		client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/backup/restoreJobs") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				// Provide minimal fields typical for restore job object
				_, _ = w.Write([]byte(`{"id":"job1","snapshotId":"snap1","deliveryType":"automated"}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})
		msg, err := executeDataDeletionRestore(ctx, client, opts)
		require.NoError(t, err)
		assert.Contains(t, msg, "Restore job submitted")
	}

	// Error path
	{
		client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		msg, err := executeDataDeletionRestore(ctx, client, opts)
		require.Error(t, err)
		assert.Empty(t, msg)
		assert.Contains(t, err.Error(), "create restore job")
	}
}

func TestSimulateRegionalOutage_Seam(t *testing.T) {
	ctx := context.Background()
	projectID := "proj"
	clusterName := "ClusterA"
	baseClusterJSON := `{"replicationSpecs":[{"regionConfigs":[{"regionName":"us-east-1","electableSpecs":{"nodeCount":3}},{"regionName":"us-west-2","electableSpecs":{"nodeCount":3}}]}]}`
	noReplJSON := `{}`
	// Dry-run add nodes only
	{
		var updateCalled atomic.Bool
		client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/clusters/"+clusterName) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(baseClusterJSON))
				return
			}
			if r.Method != http.MethodGet {
				updateCalled.Store(true)
			}
			w.WriteHeader(http.StatusNotFound)
		})
		msg, err := simulateRegionalOutage(ctx, client, recovery.DrOptions{ProjectID: projectID, ClusterName: clusterName, TargetRegion: "us-east-1", AddNodes: 2, DryRun: true})
		require.NoError(t, err)
		assert.Contains(t, msg, "(dry-run)")
		assert.Contains(t, msg, "add 2 electable nodes")
		assert.False(t, updateCalled.Load())
	}
	// Dry run with outage region zeroing
	{
		client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/clusters/"+clusterName) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(baseClusterJSON))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})
		msg, err := simulateRegionalOutage(ctx, client, recovery.DrOptions{ProjectID: projectID, ClusterName: clusterName, TargetRegion: "us-east-1", OutageRegion: "us-west-2", AddNodes: 1, DryRun: true})
		require.NoError(t, err)
		assert.Contains(t, msg, "zeroed electable nodes")
	}
	// Target region not found
	{
		client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/clusters/"+clusterName) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(baseClusterJSON))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})
		msg, err := simulateRegionalOutage(ctx, client, recovery.DrOptions{ProjectID: projectID, ClusterName: clusterName, TargetRegion: "eu-central-1", AddNodes: 1, DryRun: true})
		require.Error(t, err)
		assert.Empty(t, msg)
		assert.Contains(t, err.Error(), "target region")
	}
	// No replication specs
	{
		client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/clusters/"+clusterName) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(noReplJSON))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})
		msg, err := simulateRegionalOutage(ctx, client, recovery.DrOptions{ProjectID: projectID, ClusterName: clusterName, TargetRegion: "us-east-1", AddNodes: 1, DryRun: true})
		require.Error(t, err)
		assert.Empty(t, msg)
		assert.Contains(t, err.Error(), "no replication specs")
	}
	// Update cluster error (non dry-run)
	{
		var getCount, updateCount int32
		client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/clusters/"+clusterName) {
				atomic.AddInt32(&getCount, 1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(baseClusterJSON))
				return
			}
			if strings.HasSuffix(r.URL.Path, "/clusters/"+clusterName) {
				atomic.AddInt32(&updateCount, 1)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})
		msg, err := simulateRegionalOutage(ctx, client, recovery.DrOptions{ProjectID: projectID, ClusterName: clusterName, TargetRegion: "us-east-1", AddNodes: 1})
		require.Error(t, err)
		assert.Empty(t, msg)
		assert.Equal(t, int32(1), getCount)
		assert.Equal(t, int32(1), updateCount)
	}
}
