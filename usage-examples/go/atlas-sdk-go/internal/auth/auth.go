package auth

import (
	"atlas-sdk-go/internal"
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

const filePath = "./configs/config.json"

// CreateAtlasClient initializes and returns an authenticated Atlas API client
// using OAuth2 with service account credentials.
func CreateAtlasClient() (*admin.APIClient, *internal.Secrets, *internal.Config, error) {

	var secrets, err = internal.LoadSecrets()
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

	ctx := context.Background()
	atlasClient, err := admin.NewClient(
		admin.UseBaseURL(config.BaseURL),
		admin.UseOAuthAuth(ctx, secrets.ServiceAccountID, secrets.ServiceAccountSecret),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating SDK client: %w", err)
	}

	return atlasClient, secrets, config, nil
}
