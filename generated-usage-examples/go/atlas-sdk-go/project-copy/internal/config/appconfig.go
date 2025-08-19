package config

import (
	goErr "errors"
	"os"

	"atlas-sdk-go/internal/errors"
)

// Environment defines the runtime environment for the application
type Environment string

const (
	envDevelopment Environment = "development"
	envStaging     Environment = "test"
	envProduction  Environment = "production"
)

// AppConfig contains all environment-specific configurations
type AppConfig struct {
	environment Environment
	config      Config
	secrets     Secrets
}

// LoadAppConfig loads the application configuration based on the provided config file and environment.
// Returns the app's configuration or any error encountered during loading
func LoadAppConfig(configFile string, env Environment) (AppConfig, error) {
	if configFile == "" {
		return AppConfig{}, goErr.New("config file path must be provided")
	}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return AppConfig{}, &errors.NotFoundError{Resource: "configuration file", ID: configFile}
	}
	secrets, err := LoadSecrets()
	if err != nil {
		return AppConfig{}, errors.WithContext(err, "loading secrets")
	}
	config, err := LoadConfig(configFile)
	if err != nil {
		return AppConfig{}, errors.WithContext(err, "loading config")
	}
	return AppConfig{
		environment: env,
		config:      config,
		secrets:     secrets,
	}, nil
}
