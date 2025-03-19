package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

/* GetConfigFilePath returns the correct config file based on the environment
used in projects with multiple configs for different environments
*/

type Config struct {
	AtlasBaseURL     string `json:"ATLAS_BASE_URL"`
	AtlasOrgID       string `json:"ATLAS_ORG_ID"`
	AtlasProjectID   string `json:"ATLAS_PROJECT_ID"`
	AtlasClusterName string `json:"ATLAS_CLUSTER_NAME"`
	AtlasHostName    string `json:"ATLAS_HOST_NAME"`
	AtlasPort        string `json:"ATLAS_PORT"`
	AtlasProcessID   string `json:"ATLAS_PROCESS_ID"`
}

//func GetConfigFilePath() string {
//	env := os.Getenv("APP_ENV") // "dev", "prod", "staging"
//	if env == "" {
//		env = "dev" // Default to development
//	}
//	return fmt.Sprintf("config/config-%s.json", env)
//}

// LoadConfig loads a JSON config file into a Config struct
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
