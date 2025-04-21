package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BaseURL     string `json:"MONGODB_ATLAS_BASE_URL"`
	OrgID       string `json:"ATLAS_ORG_ID"`
	ProjectID   string `json:"ATLAS_PROJECT_ID"`
	ClusterName string `json:"ATLAS_CLUSTER_NAME"`
	HostName    string
	ProcessID   string `json:"ATLAS_PROCESS_ID"`
}

// LoadConfig loads a JSON config file to make it globally available
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file")
		}
	}(file)

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}
	return &config, nil
}

// SetDefaults sets default values if specified config variables are empty
func (c *Config) SetDefaults() {
	if c.BaseURL == "" {
		c.BaseURL = "https://cloud.mongodb.com"
	}
	if c.HostName == "" {
		c.HostName = strings.Split(c.ProcessID, ":")[0]
	}
}

// CheckRequiredFields verifies that required Atlas fields are set in the config file
func (c *Config) CheckRequiredFields() error {
	if c.OrgID == "" || c.ProjectID == "" {
		return fmt.Errorf("missing required Atlas fields in config file")
	}
	return nil
}
