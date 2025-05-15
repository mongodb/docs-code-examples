package auth

import (
	"context"
	"fmt"

	"atlas-sdk-go/internal/config"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// NewClient initializes and returns an authenticated Atlas API client
// using OAuth2 with service account credentials (recommended)
func NewClient(cfg *config.Config, secrets *config.Secrets) (*admin.APIClient, error) {
	sdk, err := admin.NewClient(
		admin.UseBaseURL(cfg.BaseURL),
		admin.UseOAuthAuth(context.Background(),
			secrets.ServiceAccountID,
			secrets.ServiceAccountSecret,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create atlas client: %w", err)
	}
	return sdk, nil
}
