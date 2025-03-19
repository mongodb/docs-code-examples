package internal

import (
	"context"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"strings"
)

// Abstract the Atlas API client into an interface to allow for mocking.

// LogsService is a minimal interface for fetching logs.
type LogsService interface {
	GetHostLogs(ctx context.Context, params *admin.GetHostLogsApiParams) (io.ReadCloser, error)
}

// MockLogsClient is a fake implementation of LogsService for testing.
type MockLogsClient struct {
	FakeResponse string
	FakeError    error
}

// GetHostLogs returns a fake log response or an error.
func (m *MockLogsClient) GetHostLogs(ctx context.Context, params *admin.GetHostLogsApiParams) (io.ReadCloser, error) {
	if m.FakeError != nil {
		return nil, m.FakeError
	}
	// Simulate a log file as a string
	return io.NopCloser(strings.NewReader(m.FakeResponse)), nil
}
