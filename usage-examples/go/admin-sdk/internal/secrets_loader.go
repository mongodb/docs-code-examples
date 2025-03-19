package internal

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Secrets structure
type Secrets struct {
	MongoDBUser          string `json:"MONGODB_USER_NAME"`
	MongoDBPassword      string `json:"MONGODB_PASSWORD"`
	AtlasAPIKey          string `json:"MONGODB_ATLAS_PUBLIC_API_KEY"`
	AtlasAPISecret       string `json:"MONGODB_ATLAS_PRIVATE_KEY"`
	ServiceAccountID     string `json:"MONGODB_ATLAS_SERVICE_ACCOUNT_ID"`
	ServiceAccountSecret string `json:"MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET"`
}

// LoadSecrets loads environment variables from a .env file into a Secrets struct
// and returns a pointer
func LoadSecrets() (*Secrets, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	secrets := &Secrets{
		MongoDBUser:          os.Getenv("MONGODB_USER_NAME"),
		MongoDBPassword:      os.Getenv("MONGODB_PASSWORD"),
		AtlasAPIKey:          os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),
		AtlasAPISecret:       os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"),
		ServiceAccountID:     os.Getenv("MONGODB_ATLAS_SERVICE_ACCOUNT_ID"),
		ServiceAccountSecret: os.Getenv("MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET"),
	}
	return secrets, nil
}
