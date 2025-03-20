package test

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"io"
	"testing"
)

var (
	testGroupID     = "test-group-id"
	testHostName    = "test-host-name"
	testProcessID   = "test-process-id"
	testLogName     = "mongodb"
	testTimeStamp   = "2023-01-01T00:00:00Z"
	testPartition   = "data"
	testLogResponse = "log content"
)

func TestMockAtlasClient_GetHostLogs_Download(t *testing.T) {
	// Set up a Gzip.Writer
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte(testLogResponse))
	assert.NoError(t, err)
	err = gz.Close()
	assert.NoError(t, err)

	// Initialize mock client with compressed log data stored in buffer
	mockClient := &MockAtlasClient{
		FakeHostLogsResponse: buf.String(),
	}

	ctx := context.Background()
	params := &admin.GetHostLogsApiParams{}
	
	// Call GetHostLogs to download the log.gz file
	resp, err := mockClient.GetHostLogs(ctx, params)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Read the downloaded log.gz file
	gzReader, err := gzip.NewReader(resp)
	assert.NoError(t, err)
	defer func(gzReader *gzip.Reader) {
		err := gzReader.Close()
		if err != nil {
			t.Errorf("failed to close gzip reader: %v", err)
		}
	}(gzReader)

	// Verify compressed data with the original log content
	var result bytes.Buffer
	_, err = io.Copy(&result, gzReader)
	assert.NoError(t, err)
	assert.Equal(t, testLogResponse, result.String())
}

func TestMockAtlasClient_GetHostLogs_Read(t *testing.T) {
	mockClient := &MockAtlasClient{
		FakeHostLogsResponse: testLogResponse,
		FakeHostLogsError:    nil,
	}

	ctx := context.Background()
	params := &admin.GetHostLogsApiParams{
		GroupId:  testGroupID,
		HostName: testHostName,
		LogName:  testLogName,
	}

	resp, err := mockClient.GetHostLogs(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	ActualLogResponse, err := io.ReadAll(resp)
	if err != nil {
		t.Fatalf("failed to read log content: %v", err)
	}

	if string(ActualLogResponse) != testLogResponse {
		t.Errorf("expected %s, got %s", testLogResponse, string(ActualLogResponse))
	}
}

func TestMockAtlasClient_GetProcessMetrics(t *testing.T) {
	expectedMetricName := "DB_DATA_SIZE_TOTAL"
	parsedTime, _ := admin.StringToTime(testTimeStamp)
	parsedTimeValue := float32(100)

	mockClient := &MockAtlasClient{
		FakeProcessMetricsResponse: &admin.ApiMeasurementsGeneralViewAtlas{
			Measurements: &[]admin.MetricsMeasurementAtlas{
				{
					Name: admin.PtrString(expectedMetricName),
					DataPoints: &[]admin.MetricDataPointAtlas{
						{Timestamp: admin.PtrTime(parsedTime), Value: admin.PtrFloat32(parsedTimeValue)},
					},
				},
			},
		},
		FakeProcessMetricsError: nil,
	}

	ctx := context.Background()
	params := &admin.GetHostMeasurementsApiParams{
		GroupId:   testGroupID,
		ProcessId: testProcessID,
	}

	resp, _, err := mockClient.GetProcessMetrics(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.HasMeasurements() == false {
		t.Errorf("expected measurements, got none")
	}

	measurements := resp.GetMeasurements()
	actualMetricName := measurements[0].GetName()
	if actualMetricName != expectedMetricName {
		t.Errorf("expected %s, got %s", expectedMetricName, actualMetricName)
	}
}

func TestMockAtlasClient_GetDiskMetrics(t *testing.T) {

	expectedMetricName := "DISK_PARTITION_SPACE_FREE"
	parsedTime, _ := admin.StringToTime(testTimeStamp)
	parsedTimeValue := float32(500)

	mockClient := &MockAtlasClient{
		FakeDiskMetricsResponse: &admin.ApiMeasurementsGeneralViewAtlas{
			Measurements: &[]admin.MetricsMeasurementAtlas{
				{
					Name: admin.PtrString(expectedMetricName),
					DataPoints: &[]admin.MetricDataPointAtlas{
						{Timestamp: admin.PtrTime(parsedTime), Value: admin.PtrFloat32(parsedTimeValue)},
					},
				},
			},
		},
		FakeDiskMetricsError: nil,
	}

	ctx := context.Background()
	params := &admin.GetDiskMeasurementsApiParams{
		GroupId:       testGroupID,
		ProcessId:     testProcessID,
		PartitionName: testPartition,
	}

	resp, _, err := mockClient.GetDiskMetrics(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.HasMeasurements() == false {
		t.Errorf("expected measurements, got none")
	}

	measurements := resp.GetMeasurements()
	actualMetricName := measurements[0].GetName()
	if actualMetricName != expectedMetricName {
		t.Errorf("expected %s, got %s", expectedMetricName, actualMetricName)
	}
}
