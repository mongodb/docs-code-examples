// :snippet-start: auth-function-full-example
package auth

import (
	"admin-sdk/internal"
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

const (
	filePath = "./configs/config.json"
)

// CreateAtlasClient initializes and returns an authenticated Atlas API client
// using OAuth2 with service account credentials.
func CreateAtlasClient() (*internal.HTTPClient, *internal.Secrets, *internal.Config, error) {

	// Load secrets
	secrets, err := internal.LoadSecrets()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load secrets: %w", err)
	}

	// Check for missing credentials
	if secrets.ServiceAccountID == "" || secrets.ServiceAccountSecret == "" {
		return nil, nil, nil, fmt.Errorf("missing Atlas client credentials")
	}
	// Load configuration
	config, err := internal.LoadConfig(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// Initialize API client using OAuth 2.0 with service account Client Credentials
	ctx := context.Background()
	sdk, err := admin.NewClient(
		admin.UseBaseURL(config.AtlasBaseURL),
		admin.UseOAuthAuth(ctx, secrets.ServiceAccountID, secrets.ServiceAccountSecret),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating SDK client: %w", err)
	}

	client := internal.NewAtlasClient(sdk)

	return client, secrets, config, nil
}

// :snippet-end: [auth-function-full-example]
