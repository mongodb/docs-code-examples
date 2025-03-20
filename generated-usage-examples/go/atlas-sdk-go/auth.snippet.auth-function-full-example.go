package auth

import (
	"atlas-sdk-go/internal"
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

const (
	filePath = "./configs/config-prod.json"
)

// CreateAtlasClient initializes and returns an authenticated Atlas API client
// using OAuth2 with service account credentials.
func CreateAtlasClient() (*internal.HTTPClient, *internal.Secrets, *internal.Config, error) {

	secrets, err := internal.LoadSecrets()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load secrets: %w", err)
	}
	if err := secrets.CheckRequiredEnv(); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid .env: %w", err)
	}

	config, err := internal.LoadConfig(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load config file: %w", err)
	}
	config.SetDefaults()
	if err := config.CheckRequiredFields(); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid config: %w", err)
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

