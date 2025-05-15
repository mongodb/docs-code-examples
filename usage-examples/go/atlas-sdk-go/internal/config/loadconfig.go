package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"atlas-sdk-go/internal"
)

type Config struct {
	BaseURL     string `json:"MONGODB_ATLAS_BASE_URL"`
	OrgID       string `json:"ATLAS_ORG_ID"`
	ProjectID   string `json:"ATLAS_PROJECT_ID"`
	ClusterName string `json:"ATLAS_CLUSTER_NAME"`
	HostName    string
	ProcessID   string `json:"ATLAS_PROCESS_ID"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config %s: %w", path, err)
	}
	defer internal.SafeClose(f)

	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	if c.BaseURL == "" {
		c.BaseURL = "https://cloud.mongodb.com"
	}
	if c.HostName == "" {
		// Go 1.18+:
		if host, _, ok := strings.Cut(c.ProcessID, ":"); ok {
			c.HostName = host
		}
	}

	if c.OrgID == "" || c.ProjectID == "" {
		return nil, fmt.Errorf("ATLAS_ORG_ID and ATLAS_PROJECT_ID are required")
	}

	return &c, nil
}
