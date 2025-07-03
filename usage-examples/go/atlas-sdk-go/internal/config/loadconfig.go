package config

import (
	"encoding/json"
	"os"

	"atlas-sdk-go/internal/errors"
)

type Config struct {
	BaseURL     string `json:"MONGODB_ATLAS_BASE_URL"`
	OrgID       string `json:"ATLAS_ORG_ID"`
	ProjectID   string `json:"ATLAS_PROJECT_ID"`
	ClusterName string `json:"ATLAS_CLUSTER_NAME"`
	HostName    string
	ProcessID   string `json:"ATLAS_PROCESS_ID"`
}

// LoadConfig reads a JSON configuration file and returns a Config struct
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, &errors.ValidationError{
			Message: "configuration file path cannot be empty",
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &errors.NotFoundError{Resource: "configuration file", ID: path}
		}
		return nil, errors.WithContext(err, "reading configuration file")
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, errors.WithContext(err, "parsing configuration file")
	}

	if config.ProjectID == "" {
		return nil, &errors.ValidationError{
			Message: "project ID is required in configuration",
		}
	}

	return &config, nil
}
