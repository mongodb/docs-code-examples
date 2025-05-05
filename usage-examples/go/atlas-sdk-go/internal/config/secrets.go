package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	EnvMongoUser      = "MONGODB_USER_NAME"
	EnvMongoPassword  = "MONGODB_PASSWORD"
	EnvAtlasAPIKey    = "MONGODB_ATLAS_PUBLIC_KEY"
	EnvAtlasAPISecret = "MONGODB_ATLAS_PRIVATE_KEY"
	EnvSAClientID     = "MONGODB_ATLAS_SERVICE_ACCOUNT_ID"
	EnvSAClientSecret = "MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET"
)

type Secrets struct {
	MongoDBUser          string
	MongoDBPassword      string
	AtlasAPIKey          string
	AtlasAPISecret       string
	ServiceAccountID     string
	ServiceAccountSecret string
}

func LoadSecrets() (*Secrets, error) {
	s := &Secrets{}
	var missing []string

	look := func(key string, dest *string) {
		if v, ok := os.LookupEnv(key); ok && v != "" {
			*dest = v
		} else {
			missing = append(missing, key)
		}
	}

	look(EnvMongoUser, &s.MongoDBUser)
	look(EnvMongoPassword, &s.MongoDBPassword)
	look(EnvAtlasAPIKey, &s.AtlasAPIKey)
	look(EnvAtlasAPISecret, &s.AtlasAPISecret)
	look(EnvSAClientID, &s.ServiceAccountID)
	look(EnvSAClientSecret, &s.ServiceAccountSecret)

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	return s, nil
}
