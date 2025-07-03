package auth

import (
	"context"

	"atlas-sdk-go/internal/config"
	"atlas-sdk-go/internal/errors"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

// NewClient initializes and returns an authenticated Atlas API client using OAuth2 with service account credentials (recommended)
// See: https://www.mongodb.com/docs/atlas/architecture/current/auth/#service-accounts
func NewClient(cfg *config.Config, secrets *config.Secrets) (*admin.APIClient, error) {
	if cfg == nil {
		return nil, &errors.ValidationError{Message: "config cannot be nil"}
	}

	if secrets == nil {
		return nil, &errors.ValidationError{Message: "secrets cannot be nil"}
	}

	sdk, err := admin.NewClient(
		admin.UseBaseURL(cfg.BaseURL),
		admin.UseOAuthAuth(context.Background(),
			secrets.ServiceAccountID,
			secrets.ServiceAccountSecret,
		),
	)
	if err != nil {
		return nil, errors.WithContext(err, "create atlas client")
	}
	return sdk, nil
}
