package config

import (
	"os"
	"strings"

	"atlas-sdk-go/internal/errors"
)

const (
	EnvSAClientID     = "MONGODB_ATLAS_SERVICE_ACCOUNT_ID"
	EnvSAClientSecret = "MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET"
)

type Secrets struct {
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

	look(EnvSAClientID, &s.ServiceAccountID)
	look(EnvSAClientSecret, &s.ServiceAccountSecret)

	if len(missing) > 0 {
		return nil, &errors.ValidationError{
			Message: "missing required environment variables: " + strings.Join(missing, ", "),
		}
	}
	return s, nil
}
