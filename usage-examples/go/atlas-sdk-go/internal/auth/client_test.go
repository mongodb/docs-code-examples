package auth_test

import (
	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/config"
	internalerrors "atlas-sdk-go/internal/errors"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClient_Success(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{BaseURL: "https://example.com"}
	secrets := &config.Secrets{
		ServiceAccountID:     "validID",
		ServiceAccountSecret: "validSecret",
	}

	client, err := auth.NewClient(cfg, secrets)

	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewClient_returnsErrorWhenConfigIsNil(t *testing.T) {
	t.Parallel()
	secrets := &config.Secrets{
		ServiceAccountID:     "validID",
		ServiceAccountSecret: "validSecret",
	}

	client, err := auth.NewClient(nil, secrets)

	require.Error(t, err)
	require.Nil(t, client)
	var validationErr *internalerrors.ValidationError
	require.True(t, errors.As(err, &validationErr), "expected error to be *errors.ValidationError")
	assert.Equal(t, "config cannot be nil", validationErr.Message)
}

func TestNewClient_returnsErrorWhenSecretsAreNil(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{BaseURL: "https://example.com"}

	client, err := auth.NewClient(cfg, nil)

	require.Error(t, err)
	require.Nil(t, client)
	var validationErr *internalerrors.ValidationError
	require.True(t, errors.As(err, &validationErr), "expected error to be *errors.ValidationError")
	assert.Equal(t, "secrets cannot be nil", validationErr.Message)
}
