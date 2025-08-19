package config

import (
	"fmt"
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

// LoadAll loads the application configuration and secrets based on the provided environment name and optional config file path.
// If configPath is empty, it defaults to "configs/config.{env}.json" based on the environment name.
// If envName is empty, it defaults to "configs/config.json".
// Parameters:
//   - envName: Environment to load configuration for (development, staging, production)
//   - configPath: Optional explicit path to the configuration file
//
// Returns:
//   - Secrets: Loaded secrets
//   - Config: Loaded application configuration
//   - error: Any error encountered during loading
func LoadAll(envName Environment, configPath string) (Secrets, Config, error) {
	var configFile string
	if configPath != "" {
		configFile = configPath
	} else if envName != "" {
		configFile = fmt.Sprintf("configs/config.%s.json", envName)
	} else {
		configFile = "configs/config.json"
	}
	appConfig, err := LoadAppConfig(configFile, envName)
	if err != nil {
		return Secrets{}, Config{}, err
	}
	return appConfig.secrets, appConfig.config, nil
}
