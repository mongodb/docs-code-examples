package archive

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
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

func TestCollectionsForArchiving_ReturnsEmpty_WhenSRVLookupFails(t *testing.T) {
	ctx := context.Background()

	// Simulate Atlas returning an error for GetCluster
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"true"}`))
	}
	client := newTestAtlasClient(t, handler)

	candidates := CollectionsForArchiving(ctx, client, "proj1", "Cluster0")

	require.NotNil(t, candidates)
	assert.Len(t, candidates, 0)
}

func TestCollectionsForArchiving_ReturnsEmpty_WhenMongoConnectFails(t *testing.T) {
	ctx := context.Background()

	// Use closed port to fail quickly
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"connectionStrings": {"standardSrv": "mongodb://127.0.0.1:1"}
		}`))
	}
	client := newTestAtlasClient(t, handler)

	candidates := CollectionsForArchiving(ctx, client, "proj1", "Cluster0")

	require.NotNil(t, candidates)
	assert.Len(t, candidates, 0)
}

// ---- Fakes for mongo client interfaces ----

type fakeCollection struct {
	count    int64
	countErr error
}

func (f *fakeCollection) EstimatedDocumentCount(ctx context.Context) (int64, error) {
	if f.countErr != nil {
		return 0, f.countErr
	}
	return f.count, nil
}

type fakeDatabase struct {
	collections map[string]*fakeCollection
	listErr     error
}

func (f *fakeDatabase) ListCollectionNames(ctx context.Context, filter interface{}) ([]string, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	names := make([]string, 0, len(f.collections))
	for k := range f.collections {
		names = append(names, k)
	}
	// keep deterministic order
	sort.Strings(names)
	return names, nil
}
func (f *fakeDatabase) Collection(name string) mongoCollection {
	c, ok := f.collections[name]
	if !ok {
		return &fakeCollection{}
	}
	return c
}

type fakeClient struct {
	dbs map[string]*fakeDatabase
}

func (f *fakeClient) Ping(ctx context.Context) error { return nil }
func (f *fakeClient) ListDatabaseNames(ctx context.Context, filter interface{}) ([]string, error) {
	names := make([]string, 0, len(f.dbs))
	for k := range f.dbs {
		names = append(names, k)
	}
	return names, nil
}
func (f *fakeClient) Database(name string) mongoDatabase {
	db, ok := f.dbs[name]
	if !ok {
		return &fakeDatabase{collections: map[string]*fakeCollection{}}
	}
	return db
}
func (f *fakeClient) Disconnect(ctx context.Context) error { return nil }

// helper to stub Atlas SRV lookup with any valid URI string
func okSRVHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"connectionStrings": {"standardSrv": "mongodb://fake"}}`))
	}
}

func TestCollectionsForArchiving_SuccessThresholdAndSkipInternal(t *testing.T) {
	ctx := context.Background()

	// Arrange fake mongo client
	old := newMongoClient
	t.Cleanup(func() { newMongoClient = old })
	newMongoClient = func(ctx context.Context, uri string) (mongoClient, error) {
		return &fakeClient{dbs: map[string]*fakeDatabase{
			"admin":  {collections: map[string]*fakeCollection{"sys": {count: 999999}}},
			"local":  {collections: map[string]*fakeCollection{"sys": {count: 999999}}},
			"config": {collections: map[string]*fakeCollection{"sys": {count: 999999}}},
			"appdb": {collections: map[string]*fakeCollection{
				"small": {count: 10},
				"big":   {count: 100000}, // threshold match
			}},
		}}, nil
	}
	client := newTestAtlasClient(t, okSRVHandler())

	// Act
	candidates := CollectionsForArchiving(ctx, client, "proj1", "Cluster0")

	// Assert
	require.NotNil(t, candidates)
	assert.Equal(t, []Candidate{ // order is deterministic in our fake
		{DatabaseName: "appdb", CollectionName: "big", DateField: "createdAt", DateFormat: "DATE", RetentionDays: 90, PartitionFields: []string{"createdAt"}},
	}, candidates)
}

func TestCollectionsForArchiving_SkipsOnListCollectionError(t *testing.T) {
	ctx := context.Background()
	old := newMongoClient
	t.Cleanup(func() { newMongoClient = old })
	newMongoClient = func(ctx context.Context, uri string) (mongoClient, error) {
		return &fakeClient{dbs: map[string]*fakeDatabase{
			"appdb":  {collections: map[string]*fakeCollection{"big": {count: 200000}}},
			"bad_db": {collections: map[string]*fakeCollection{"ignored": {count: 999999}}, listErr: assert.AnError},
		}}, nil
	}
	client := newTestAtlasClient(t, okSRVHandler())

	candidates := CollectionsForArchiving(ctx, client, "proj1", "Cluster0")

	require.NotNil(t, candidates)
	// Should include appdb.big only; bad_db is skipped due to error
	assert.Equal(t, 1, len(candidates))
	assert.Equal(t, "appdb", candidates[0].DatabaseName)
	assert.Equal(t, "big", candidates[0].CollectionName)
}

func TestCollectionsForArchiving_SkipsOnCountError(t *testing.T) {
	ctx := context.Background()
	old := newMongoClient
	t.Cleanup(func() { newMongoClient = old })
	newMongoClient = func(ctx context.Context, uri string) (mongoClient, error) {
		return &fakeClient{dbs: map[string]*fakeDatabase{
			"appdb": {collections: map[string]*fakeCollection{
				"bad": {count: 0, countErr: assert.AnError},
				"ok":  {count: 150000},
			}},
		}}, nil
	}
	client := newTestAtlasClient(t, okSRVHandler())

	candidates := CollectionsForArchiving(ctx, client, "proj1", "Cluster0")

	require.NotNil(t, candidates)
	assert.Equal(t, 1, len(candidates))
	assert.Equal(t, "ok", candidates[0].CollectionName)
}
