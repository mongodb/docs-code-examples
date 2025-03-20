// :snippet-start: config-loader-function-full-example
package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AtlasBaseURL     string `json:"ATLAS_BASE_URL"`
	AtlasOrgID       string `json:"ATLAS_ORG_ID"`
	AtlasProjectID   string `json:"ATLAS_PROJECT_ID"`
	AtlasClusterName string `json:"ATLAS_CLUSTER_NAME"`
	AtlasHostName    string `json:"ATLAS_HOST_NAME"`
	AtlasPort        string `json:"ATLAS_PORT"`
	AtlasProcessID   string `json:"ATLAS_PROCESS_ID"`
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
	if c.AtlasBaseURL == "" {
		c.AtlasBaseURL = "https://cloud.mongodb.com"
	}
	if c.AtlasPort == "" {
		c.AtlasPort = "27017"
	}
	if c.AtlasProcessID == "" && c.AtlasHostName != "" {
		c.AtlasProcessID = c.AtlasHostName + ":" + c.AtlasPort
	}
	if c.AtlasHostName == "" && c.AtlasProcessID != "" {
		c.AtlasHostName = strings.Split(c.AtlasProcessID, ":")[0]
	}
}

// CheckRequiredFields verifies that required Atlas fields are set in the config file
func (c *Config) CheckRequiredFields() error {
	if c.AtlasOrgID == "" || c.AtlasProjectID == "" {
		return fmt.Errorf("missing required Atlas fields in config file")
	}
	if c.AtlasProcessID == "" || c.AtlasHostName == "" && c.AtlasPort == "" {
		return fmt.Errorf("either ATLAS_PROCESS_ID or ATLAS_HOST_NAME/PORT must be set")
	}
	return nil
}

// :snippet-end: [config-loader-function-full-example]
