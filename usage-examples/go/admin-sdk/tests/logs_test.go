package tests

//import (
//	"context"
//	"errors"
//	"myapp/internal"
//	"myapp/scripts"
//	"testing"
//)
//
//// TestGetHostLogsSuccess ensures logs are fetched successfully.
//func TestGetHostLogsSuccess(t *testing.T) {
//	ctx := context.Background()
//
//	// Mock client returns fake log data
//	mockClient := &internal.MockLogsClient{
//		FakeResponse: "Test log entry\n",
//	}
//
//	err := scripts.GetHostLogs(ctx, mockClient, "123456", "test-host", "mongodb.log")
//	if err != nil {
//		t.Fatalf("Unexpected error: %v", err)
//	}
//}
//
//// TestGetHostLogsFailure simulates an API failure.
//func TestGetHostLogsFailure(t *testing.T) {
//	ctx := context.Background()
//
//	// Mock client returns an error
//	mockClient := &internal.MockLogsClient{
//		FakeError: errors.New("mock API failure"),
//	}
//
//	err := scripts.GetHostLogs(ctx, mockClient, "123456", "test-host", "mongodb.log")
//	if err == nil {
//		t.Fatalf("Expected error but got none")
//	}
//}
