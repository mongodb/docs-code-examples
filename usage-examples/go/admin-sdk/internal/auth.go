package internal

import (
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"log"
)

// CreateAtlasClient initializes and returns an authenticated Atlas API client
// using OAuth2 with service account credentials.
func CreateAtlasClient() (*HTTPClient, *Secrets, *Config, error) {

	// Load secrets
	secrets, err := LoadSecrets()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load secrets: %w", err)
	}

	// Check for missing credentials
	if secrets.ServiceAccountID == "" || secrets.ServiceAccountSecret == "" {
		return nil, nil, nil, fmt.Errorf("missing Atlas client credentials")
	}

	// Load configuration
	config, err := LoadConfig("config/config.json")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// Determine base URL
	baseURL := config.AtlasBaseURL
	if baseURL == "" {
		baseURL = "https://cloud.mongodb.com"
	}

	// Check if ProcessID or Hostname:Port are set
	if config.AtlasProcessID == "" || config.AtlasHostName == "" && config.AtlasPort == "" {
		log.Fatal("Either ATLAS_PROCESS_ID or ATLAS_HOST_NAME/PORT must be set")
	}

	// Initialize API client using OAuth 2.0 with service account Client Credentials
	ctx := context.Background()
	sdk, err := admin.NewClient(
		admin.UseBaseURL(baseURL),
		admin.UseOAuthAuth(ctx, secrets.ServiceAccountID, secrets.ServiceAccountSecret),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating SDK client: %w", err)
	}

	client := NewAtlasClient(sdk)

	return client, secrets, config, nil
}
