package config

import (
	"atlas-sdk-go/internal/errors"
)

// LoadAll loads secrets and config
// If configPath is empty, uses environment-specific loading
// If explicitEnv is provided, it overrides the APP_ENV environment variable
func LoadAll(configPath string, explicitEnv string) (*Secrets, *Config, error) {
	if configPath == "" {
		// Use environment-based loading
		appCtx, err := LoadAppContext(explicitEnv, false) // Use non-strict validation by default
		if err != nil {
			return nil, nil, err
		}
		return appCtx.Secrets, appCtx.Config, nil
	}

	// Legacy path-specific loading
	s, err := LoadSecrets()
	if err != nil {
		return nil, nil, errors.WithContext(err, "loading secrets")
	}

	c, err := LoadConfig(configPath)
	if err != nil {
		return nil, nil, errors.WithContext(err, "loading config")
	}

	return s, c, nil
}

// ValidateEnvironment checks if the provided environment is valid
func ValidateEnvironment(env string) bool {
	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
		"test":        true,
	}
	return validEnvs[env]
}
