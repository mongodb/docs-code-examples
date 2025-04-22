package internal

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Secrets struct {
	ServiceAccountID     string `json:"MONGODB_ATLAS_SERVICE_ACCOUNT_ID"`
	ServiceAccountSecret string `json:"MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET"`
}

// LoadSecrets loads .env file variables to use in the application
func LoadSecrets() (*Secrets, error) {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("No .env file found")
	}
	secrets := &Secrets{
		ServiceAccountID:     os.Getenv("MONGODB_ATLAS_SERVICE_ACCOUNT_ID"),
		ServiceAccountSecret: os.Getenv("MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET"),
	}
	return secrets, nil
}

// CheckRequiredEnv verifies that required environment variables are set
func (s *Secrets) CheckRequiredEnv() error {
	if s.ServiceAccountID == "" || s.ServiceAccountSecret == "" {
		return fmt.Errorf("service account client credentials must be set")
	}
	return nil
}
