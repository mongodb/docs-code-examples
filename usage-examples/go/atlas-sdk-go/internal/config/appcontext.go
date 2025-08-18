package config

import (
	"fmt"
	"os"

	"atlas-sdk-go/internal/errors"
)

const (
	envAppEnv           = "APP_ENV"         // Environment variable for the application environment
	envConfigPath       = "APP_CONFIG_PATH" // Environment variable for the configuration file path
	defaultConfigFormat = "configs/config.%s.json"
)

const (
	envDevelopment Environment = "development"
	envStaging     Environment = "test"
	envProduction  Environment = "production"
)

// AppContext contains all environment-specific configurations
type AppContext struct {
	environment Environment
	config      Config
	secrets     Secrets
}

type Environment string

// LoadAppContext loads environment-specific configuration
// Returns error if no environment is specified
func LoadAppContext(explicitEnv Environment) (AppContext, error) {
	env := explicitEnv
	if env == "" {
		env = Environment(os.Getenv(envAppEnv))
	}
	// Validate environment
	if !ValidateEnvironment(string(env)) {
		return AppContext{}, fmt.Errorf("invalid environment: %s", env)
	}

	// Determine config path
	configPath := os.Getenv(envConfigPath)
	if configPath == "" {
		configPath = fmt.Sprintf(defaultConfigFormat, env)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return AppContext{}, &errors.NotFoundError{Resource: "configuration file", ID: configPath}
	}

	// Load secrets and config
	secrets, err := LoadSecrets()
	if err != nil {
		return AppContext{}, errors.WithContext(err, "loading secrets")
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return AppContext{}, errors.WithContext(err, "loading config")
	}

	return AppContext{
		environment: env,
		config:      config,
		secrets:     secrets,
	}, nil
}
