package config

import (
	"fmt"
)

// LoadAll loads secrets and config from the specified paths
func LoadAll(configPath string) (*Secrets, *Config, error) {
	s, err := LoadSecrets()
	if err != nil {
		return nil, nil, fmt.Errorf("loading secrets: %w", err)
	}

	c, err := LoadConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("loading config: %w", err)
	}

	return s, c, nil
}
