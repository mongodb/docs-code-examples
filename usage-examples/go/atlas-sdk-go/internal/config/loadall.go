package config

import (
	"fmt"
	"os"
	"strings"

	"atlas-sdk-go/internal/errors"
)

const defaultConfigDir = "configs"

// LoadAll loads both secrets and configuration from the specified paths.
func LoadAll(configPath string) (Secrets, Config, error) {
	if strings.TrimSpace(configPath) == "" {
		configPath = fmt.Sprintf("%s/config.json", defaultConfigDir) // Default path if not specified in environment
	}

	if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
		return Secrets{}, Config{}, &errors.NotFoundError{Resource: "configuration file", ID: configPath}
	}

	secrets, err := LoadSecrets()
	if err != nil {
		return Secrets{}, Config{}, errors.WithContext(err, "loading secrets")
	}
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return Secrets{}, Config{}, errors.WithContext(err, "loading config")
	}
	return secrets, cfg, nil
}
