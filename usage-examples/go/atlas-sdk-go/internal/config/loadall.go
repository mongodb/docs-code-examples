package config

import (
	"atlas-sdk-go/internal/errors"
)

// EnvironmentNames defines valid runtime environments
var allowedEnvironments = map[Environment]struct{}{
	envDevelopment: {},
	envStaging:     {},
	envProduction:  {},
}

func ValidateEnvironment(env string) bool {
	_, ok := allowedEnvironments[Environment(env)]
	return ok
}

// LoadAll loads configuration for a specific environment
//
// Parameters:
//   - envName: Environment name (dev/staging/prod/test); overrides APP_ENV if provided
//   - configPath: Optional explicit config file path; if empty, uses environment-based path
//
// Returns secrets, config and any errors encountered during loading
func LoadAll(envName Environment, configPath string) (Secrets, Config, error) {
	if configPath == "" {
		appCtx, err := LoadAppContext(envName)
		if err != nil {
			return Secrets{}, Config{}, err
		}
		return appCtx.secrets, appCtx.config, nil // return values, not pointers
	}
	s, err := LoadSecrets()
	if err != nil {
		return Secrets{}, Config{}, errors.WithContext(err, "loading secrets")
	}
	c, err := LoadConfig(configPath)
	if err != nil {
		return Secrets{}, Config{}, errors.WithContext(err, "loading config")
	}
	return s, c, nil // return values, not pointers
}
