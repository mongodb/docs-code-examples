package internal

import (
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	"os"
)

// CreateAtlasClient initializes and returns an authenticated Atlas API client.
func CreateAtlasClient() (*admin.APIClient, *Config, error) {

	// Load secrets
	secrets, err := LoadSecrets()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load secrets: %w", err)
	} else {
		fmt.Println("Secrets loaded successfully")
	}

	// Print loaded secrets for debugging
	fmt.Printf("Loaded Secrets: %+v\n", secrets)

	// Check for missing credentials
	if secrets.ClientID == "" || secrets.ClientSecret == "" {
		return nil, nil, fmt.Errorf("missing Atlas client credentials")
	}
	fmt.Println("Client credentials are present")

	// Load configuration
	config, err := LoadConfig("config/config.json")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// Determine base URL
	host := config.BaseURL
	if host == "" {
		host = os.Getenv("ATLAS_BASE_URL")
	}
	if host == "" {
		host = "https://cloud.mongodb.com"
	}
	fmt.Println("Using Atlas Base URL:", host)

	// Initialize API client
	ctx := context.Background()
	sdk, err := admin.NewClient(
		admin.UseBaseURL(host),
		admin.UseOAuthAuth(ctx, secrets.ClientID, secrets.ClientSecret),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating SDK client: %w", err)
	}
	fmt.Println("SDK client created successfully")

	// Determine org ID
	orgID := config.OrgID
	if orgID == "" {
		orgID = os.Getenv("ATLAS_ORG_ID")
	}
	if orgID == "" {
		orgs, _, err := sdk.OrganizationsApi.ListOrganizations(ctx).Execute()
		if err != nil {
			return nil, nil, fmt.Errorf("error listing organizations: %w", err)
		}
		if orgs.GetTotalCount() == 0 {
			return nil, nil, fmt.Errorf("no organizations found")
		}
		orgID = orgs.GetResults()[0].GetId()
		fmt.Printf("Using organization %s\n", orgID)
	}

	return sdk, config, nil
}
