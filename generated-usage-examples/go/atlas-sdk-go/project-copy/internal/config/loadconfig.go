package config

import (
	"encoding/json"
	"os"
	"strings"

	"atlas-sdk-go/internal/errors"
)

// Config holds the configuration for connecting to MongoDB Atlas
type Config struct {
	BaseURL     string `json:"MONGODB_ATLAS_BASE_URL"`
	OrgID       string `json:"ATLAS_ORG_ID"`
	ProjectID   string `json:"ATLAS_PROJECT_ID"`
	ClusterName string `json:"ATLAS_CLUSTER_NAME"`
	HostName    string `json:"ATLAS_HOSTNAME"`
	ProcessID   string `json:"ATLAS_PROCESS_ID"`
}

// LoadConfig reads a JSON configuration file and returns a Config struct
// It validates required fields and returns an error if any validation fails.
func LoadConfig(path string) (Config, error) {
	var config Config
	if path == "" {
		return config, &errors.ValidationError{
			Message: "configuration file path cannot be empty",
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, &errors.NotFoundError{Resource: "configuration file", ID: path}
		}
		return config, errors.WithContext(err, "reading configuration file")
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return config, errors.WithContext(err, "parsing configuration file")
	}

	if config.OrgID == "" {
		return config, &errors.ValidationError{
			Message: "organization ID is required in configuration",
		}
	}
	if config.ProjectID == "" {
		return config, &errors.ValidationError{
			Message: "project ID is required in configuration",
		}
	}

	if config.HostName == "" {
		if host, _, ok := strings.Cut(config.ProcessID, ":"); ok {
			config.HostName = host
		} else {
			return config, &errors.ValidationError{
				Message: "process ID must be in the format 'hostname:port'",
			}
		}
	}

	if config.BaseURL == "" {
		config.BaseURL = "https://cloud.mongodb.com" // Default base URL if not provided
	}

	return config, nil
}
