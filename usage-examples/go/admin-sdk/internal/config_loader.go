package internal

import (
	"encoding/json"
	"fmt"
	"os"
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

// LoadConfig loads a JSON configs file into a Config struct
// takes a file path as an argument and returns a pointer to a Config struct and an error
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

// GetConfigFilePath returns the correct config file based on the environment
//func (c *Config) GetConfigFilePath(appEnv string) string {
//	return fmt.Sprintf("configs/config-%s.json", c.AppEnv)
//}

// SetDefaults sets default values if specified config variables are empty
//func (c *Config) SetDefaults() {
//	if utils.IsEmptyString(c.AtlasBaseURL) {
//		c.AtlasBaseURL = "https://cloud.mongodb.com"
//	}
//	if utils.IsEmptyString(c.AtlasPort) {
//		c.AtlasPort = "27017"
//	}
//	if utils.IsEmptyString(c.AtlasProcessID) && !utils.IsEmptyString(c.AtlasHostName) {
//		c.AtlasProcessID = c.AtlasHostName + ":" + c.AtlasPort
//	}
//	if utils.IsEmptyString(c.AtlasHostName) {
//		c.AtlasHostName = strings.Split(c.AtlasProcessID, ":")[0]
//	}
//}
//
//func (c *Config) CheckRequiredFields() error {
//	if utils.IsEmptyString(c.AtlasOrgID) || utils.IsEmptyString(c.AtlasProjectID) {
//		return fmt.Errorf("missing required Atlas fields in config file")
//	}
//	if utils.IsEmptyString(c.AtlasProcessID) || utils.IsEmptyString(c.AtlasHostName) && utils.IsEmptyString(c.AtlasPort) {
//		return fmt.Errorf("either ATLAS_PROCESS_ID or ATLAS_HOST_NAME/PORT must be set")
//	}
//	return nil
//}
