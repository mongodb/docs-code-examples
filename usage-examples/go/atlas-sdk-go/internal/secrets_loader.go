// :snippet-start: secrets-loader-function-full-example

package internal

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Secrets struct {
	MongoDBUser          string `json:"MONGODB_USER_NAME"`            // :remove:
	MongoDBPassword      string `json:"MONGODB_PASSWORD"`             // :remove:
	AtlasAPIKey          string `json:"MONGODB_ATLAS_PUBLIC_API_KEY"` // :remove:
	AtlasAPISecret       string `json:"MONGODB_ATLAS_PRIVATE_KEY"`    // :remove:
	ServiceAccountID     string `json:"MONGODB_ATLAS_SERVICE_ACCOUNT_ID"`
	ServiceAccountSecret string `json:"MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET"`
}

// LoadSecrets loads environment variables from a .env file to use in the application
func LoadSecrets() (*Secrets, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	secrets := &Secrets{
		MongoDBUser:          os.Getenv("MONGODB_USER_NAME"),         // :remove:
		MongoDBPassword:      os.Getenv("MONGODB_PASSWORD"),          // :remove:
		AtlasAPIKey:          os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),  // :remove:
		AtlasAPISecret:       os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"), // :remove:
		ServiceAccountID:     os.Getenv("MONGODB_ATLAS_SERVICE_ACCOUNT_ID"),
		ServiceAccountSecret: os.Getenv("MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET"),
	}
	return secrets, nil
}

// :snippet-end: [secrets-loader-function-full-example]
