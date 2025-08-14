package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"atlas-sdk-go/internal/errors"
)

// Package config provides application context management, including environment-specific configurations
// and caching mechanisms to optimize performance and reduce redundant loading of configurations.
var (
	cachedAppContext     *AppContext
	cachedAppContextTime time.Time
	cacheTTL             = 5 * time.Minute
	cacheMutex           sync.RWMutex
)

// Constants for environment variables and default paths
const (
	EnvAppEnv           = "APP_ENV"
	EnvConfigPath       = "ATLAS_CONFIG_PATH"
	EnvSAClientID       = "MONGODB_ATLAS_SERVICE_ACCOUNT_ID"
	EnvSAClientSecret   = "MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET"
	DefaultConfigFormat = "configs/config.%s.json"
)

// AppContext contains all environment-specific configurations
type AppContext struct {
	Environment string
	Config      *Config
	Secrets     *Secrets
}

// LoadAppContext initializes application context with environment-specific configuration
// If explicitEnv is provided, it overrides the APP_ENV environment variable
// If strictValidation is true, invalid environments will return an error
func LoadAppContext(explicitEnv string, strictValidation bool) (*AppContext, error) {
	// Environment resolution priority:
	// 1. An explicitly passed environment parameter
	// 2. An APP_ENV environment variable
	// 3. Otherwise, defaults to "development"

	// Determine environment
	env := explicitEnv
	if env == "" {
		env = os.Getenv(EnvAppEnv)
		if env == "" {
			env = "development"
		}
	}

	// Check cache first using the resolved environment
	cacheMutex.RLock()
	if cachedAppContext != nil &&
		cachedAppContext.Environment == env &&
		time.Since(cachedAppContextTime) < cacheTTL {
		cached := cachedAppContext
		cacheMutex.RUnlock()
		return cached, nil
	}
	cacheMutex.RUnlock()

	if !ValidateEnvironment(env) {
		if strictValidation {
			return nil, fmt.Errorf("invalid environment: %s", env)
		}
		log.Printf("Warning: Unexpected environment '%s' may cause issues", env)
	}
	// Load environment files
	envFiles := []string{
		fmt.Sprintf(".env.%s", env),
		".env",
	}

	loaded := false
	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			log.Printf("Loaded environment from %s", file)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Printf("Warning: No environment files found, using system environment variables only")
	}

	// Get config path from env var or use default
	configPath := os.Getenv(EnvConfigPath)
	if configPath == "" {
		configPath = fmt.Sprintf(DefaultConfigFormat, env)
	}

	log.Printf("Loading configuration for environment: %s", env)
	log.Printf("Using config file: %s", configPath)

	secrets, err := LoadSecrets()
	if err != nil {
		return nil, errors.WithContext(err, "loading secrets")
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, errors.WithContext(err, "loading config")
	}

	if err = config.Validate(env); err != nil {
		return nil, errors.WithContext(err, "validating config")
	}

	log.Printf("Configuration loaded successfully: env=%s, baseURL=%s, orgID=%s",
		env, config.BaseURL, config.OrgID)

	// Create and initialize the AppContext
	appCtx := &AppContext{
		Environment: env,
		Config:      config,
		Secrets:     secrets,
	}

	// Cache the result
	cacheMutex.Lock()
	cachedAppContext = appCtx
	cachedAppContextTime = time.Now()
	cacheMutex.Unlock()

	return appCtx, nil
}

// LoadAppContextWithContext initializes application context with environment-specific configuration using a provided context for cancellation support.
// If explicitEnv is provided, it overrides the APP_ENV environment variable
// If strictValidation is true, invalid environments will return an error
func LoadAppContextWithContext(ctx context.Context, explicitEnv string, strictValidation bool) (*AppContext, error) {
	// Use context for potential operations that may need cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while loading configuration: %w", ctx.Err())
	default:
		// Continue with loading
	}

	// Determine environment
	env := explicitEnv
	if env == "" {
		env = os.Getenv(EnvAppEnv)
		if env == "" {
			env = "development"
		}
	}

	// Check cache first using the resolved environment
	cacheMutex.RLock()
	if cachedAppContext != nil &&
		cachedAppContext.Environment == env &&
		time.Since(cachedAppContextTime) < cacheTTL {
		cached := cachedAppContext
		cacheMutex.RUnlock()
		return cached, nil
	}
	cacheMutex.RUnlock()

	// The implementation mirrors LoadAppContext but with context checks
	if !ValidateEnvironment(env) {
		if strictValidation {
			return nil, fmt.Errorf("invalid environment: %s", env)
		}
		log.Printf("Warning: Unexpected environment '%s' may cause issues", env)
	}

	// Add context check before expensive operations
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while loading environment files: %w", ctx.Err())
	default:
	}

	// Load environment files
	envFiles := []string{
		fmt.Sprintf(".env.%s", env),
		".env",
	}

	loaded := false
	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			log.Printf("Loaded environment from %s", file)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Printf("Warning: No environment files found, using system environment variables only")
	}

	// Get config path from env var or use default
	configPath := os.Getenv(EnvConfigPath)
	if configPath == "" {
		configPath = fmt.Sprintf(DefaultConfigFormat, env)
	}

	// Add context check before loading secrets and config
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled before loading secrets/config: %w", ctx.Err())
	default:
	}

	log.Printf("Loading configuration for environment: %s", env)
	log.Printf("Using config file: %s", configPath)

	secrets, err := LoadSecrets()
	if err != nil {
		return nil, errors.WithContext(err, "loading secrets")
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, errors.WithContext(err, "loading config")
	}

	if err = config.Validate(env); err != nil {
		return nil, errors.WithContext(err, "validating config")
	}

	log.Printf("Configuration loaded successfully: env=%s, baseURL=%s, orgID=%s",
		env, config.BaseURL, config.OrgID)

	// Create and initialize the AppContext
	appCtx := &AppContext{
		Environment: env,
		Config:      config,
		Secrets:     secrets,
	}

	// Cache the result
	cacheMutex.Lock()
	cachedAppContext = appCtx
	cachedAppContextTime = time.Now()
	cacheMutex.Unlock()

	return appCtx, nil
}

