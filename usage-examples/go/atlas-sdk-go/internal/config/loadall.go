package config

import (
	"atlas-sdk-go/internal/errors"
)

// LoadAll loads secrets and config from the specified paths
func LoadAll(configPath string) (*Secrets, *Config, error) {
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
