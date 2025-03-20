func main() {
	ctx := context.Background()

	// Create an Atlas client authenticated using OAuth2 with service account credentials
	client, _, config, err := auth.CreateAtlasClient()
	if err != nil {
		log.Fatalf("Failed to create Atlas client: %v", err)
	}

	params := &GetHostLogsParams{
		GroupID:  config.AtlasProjectID,
		HostName: config.AtlasHostName, // The host to get logs for
		LogName:  LogName,              // The type of log to get ("mongodb" or "mongos")
	}

	// Downloads the specified host's MongoDB logs as a .gz file
	if err := getHostLogs(ctx, *client, params); err != nil {
		fmt.Printf("Error fetching host logs: %v", err)
	}
}

