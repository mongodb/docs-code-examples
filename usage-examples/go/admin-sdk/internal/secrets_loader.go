package internal

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Secrets structure
type Secrets struct {
	DBUser       string `json:"DB_USER"`
	DBPassword   string `json:"DB_PASSWORD"`
	APIKey       string `json:"PUBLIC_API_KEY"`
	PrivateKey   string `json:"PRIVATE_API_KEY"`
	ClientID     string `json:"ATLAS_CLIENT_ID"`
	ClientSecret string `json:"ATLAS_CLIENT_SECRET"`
}

// LoadEnv loads environment variables from a .env file
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// LoadSecrets loads a JSON secrets file into a Secrets struct
// takes a file path as an argument and returns a pointer to a Secrets struct and an error
func LoadSecrets() (*Secrets, error) {
	LoadEnv()
	secrets := &Secrets{
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		APIKey:       os.Getenv("PUBLIC_API_KEY"),
		PrivateKey:   os.Getenv("PRIVATE_API_KEY"),
		ClientID:     os.Getenv("ATLAS_CLIENT_ID"),
		ClientSecret: os.Getenv("ATLAS_CLIENT_SECRET"),
	}
	return secrets, nil
}
