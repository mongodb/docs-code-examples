package main

import (
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"os"
)

// GetHostNameFromID retrieves the hostname of a process in an Atlas project.

func GetHostName(ctx context.Context, sdk *admin.APIClient, groupId string, includeCount bool, itemsPerPage, pageNum int) (*admin.PaginatedHostViewAtlas, error) {
	resp, r, err := sdk.MonitoringAndLogsApi.ListAtlasProcesses(ctx, groupId).IncludeCount(includeCount).ItemsPerPage(itemsPerPage).PageNum(pageNum).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MonitoringAndLogsApi.ListAtlasProcesses`: %v (%v)\n", err, r)
		apiError, ok := admin.AsError(err)
		if ok {
			fmt.Fprintf(os.Stderr, "API error obj: %v\n", apiError)
		}
		// response from `ListAtlasProcesses`: PaginatedHostViewAtlas
		fmt.Fprintf(os.Stdout, "Response from `MonitoringAndLogsApi.ListAtlasProcesses`: %v (%v)\n", resp, r)
	}

	return resp, nil
}

//func getHostFromID(ctx context.Context, sdk *admin.APIClient, groupId, processId string) (string, error) {
//	resp, r, err := sdk.MonitoringAndLogsApi.GetAtlasProcess(ctx, groupId, processId).Execute()
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Error when calling `MonitoringAndLogsApi.GetProcess`: %v (%v)\n", err, r)
//		apiError, ok := admin.AsError(err)
//		if ok {
//			fmt.Fprintf(os.Stderr, "API error obj: %v\n", apiError)
//		}
//		// response from `GetProcess`: Process
//		fmt.Fprintf(os.Stdout, "Response from `MonitoringAndLogsApi.GetProcess`: %v (%v)\n", resp, r)
//	}
//	return resp.GetHostname(), nil
//}
//
//func main() {
//
//	// Initialize Atlas client
//	sdk, config, err := internal.CreateAtlasClient()
//	if err != nil {
//		fmt.Printf("failed to create Atlas client: %v\n", err)
//	}
//
//	// Get host name
//	hostName, err := ListProcesses(context.Background(), sdk, config.GroupID)
//	if err != nil {
//		fmt.Printf("failed to get host name: %v\n", err)
//		return
//	}
//
//	fmt.Println("Host name:", hostName)
//}
