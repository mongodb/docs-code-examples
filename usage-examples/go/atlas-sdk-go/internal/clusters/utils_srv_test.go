package clusters

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func newTestAtlasClient(t *testing.T, handler http.HandlerFunc) *admin.APIClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := admin.NewClient(admin.UseBaseURL(server.URL))
	require.NoError(t, err)
	return client
}

func TestGetClusterSRVConnectionString_Success(t *testing.T) {
	t.Parallel()
	projectID := "proj1"
	clusterName := "Cluster0"

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"connectionStrings": {
				"standardSrv": "mongodb+srv://cluster0.example.net"
			}
		}`))
	}
	client := newTestAtlasClient(t, handler)

	srv, err := GetClusterSRVConnectionString(context.Background(), client, projectID, clusterName)

	require.NoError(t, err)
	assert.Equal(t, "mongodb+srv://cluster0.example.net", srv)
}

func TestGetClusterSRVConnectionString_MissingField(t *testing.T) {
	t.Parallel()
	projectID := "proj1"
	clusterName := "Cluster0"

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// connectionStrings present but standardSrv missing
		_, _ = w.Write([]byte(`{"connectionStrings": {}}`))
	}
	client := newTestAtlasClient(t, handler)

	srv, err := GetClusterSRVConnectionString(context.Background(), client, projectID, clusterName)

	require.Error(t, err)
	assert.Empty(t, srv)
	assert.Contains(t, err.Error(), "no standard SRV")
}

func TestGetClusterSRVConnectionString_ApiError(t *testing.T) {
	t.Parallel()
	projectID := "proj1"
	clusterName := "Cluster0"

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"detail": "server error"}`))
	}
	client := newTestAtlasClient(t, handler)

	srv, err := GetClusterSRVConnectionString(context.Background(), client, projectID, clusterName)

	require.Error(t, err)
	assert.Empty(t, srv)
	assert.Contains(t, err.Error(), "get cluster")
}
