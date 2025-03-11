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
	BaseURL         string `json:"ATLAS_BASE_URL"`
	GroupID         string `json:"GROUP_ID"`
	BillPayingOrgID string `json:"BILL_PAYING_ORG_ID"`
	OrgID           string `json:"ORG_ID"`
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
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}
	return &config, nil
}
