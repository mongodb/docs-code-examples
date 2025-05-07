package logs

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"atlas-sdk-go/internal"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"go.mongodb.org/atlas-sdk/v20250219001/mockadmin"
)

func TestFetchHostLogs_Unit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// common params
	params := &admin.GetHostLogsApiParams{
		GroupId:  "gID",
		HostName: "hName",
		LogName:  "mongodb",
	}

	cases := []struct {
		name     string
		setup    func(m *mockadmin.MonitoringAndLogsApi)
		wantErr  bool
		wantBody string
	}{
		{
			name:    "API error",
			wantErr: true,
			setup: func(m *mockadmin.MonitoringAndLogsApi) {
				m.EXPECT().
					GetHostLogs(mock.Anything, params.GroupId, params.HostName, params.LogName).
					Return(admin.GetHostLogsApiRequest{ApiService: m}).Once()
				m.EXPECT().
					GetHostLogsExecute(mock.Anything).
					Return(nil, nil, fmt.Errorf("API error")).Once()
			},
		},
		{
			name:     "Successful response",
			wantErr:  false,
			wantBody: "log-data",
			setup: func(m *mockadmin.MonitoringAndLogsApi) {
				m.EXPECT().
					GetHostLogs(mock.Anything, params.GroupId, params.HostName, params.LogName).
					Return(admin.GetHostLogsApiRequest{ApiService: m}).Once()
				m.EXPECT().
					GetHostLogsExecute(mock.Anything).
					Return(io.NopCloser(strings.NewReader("log-data")), nil, nil).Once()
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockSvc := mockadmin.NewMonitoringAndLogsApi(t)
			tc.setup(mockSvc)

			rc, err := FetchHostLogs(ctx, mockSvc, params)
			if tc.wantErr {
				require.ErrorContainsf(t, err, "failed to fetch logs", "expected API error")
				require.Nil(t, rc)
				return
			}

			require.NoError(t, err)
			defer internal.SafeClose(rc)

			data, err := io.ReadAll(rc)
			require.NoError(t, err)
			require.Equal(t, tc.wantBody, string(data))
		})
	}
}
