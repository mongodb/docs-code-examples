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
	// 1. Explicitly passed environment parameter
	// 2. APP_ENV environment variable
	// 3. Default to "development"
	//
	// Special environments:
	// - "test": Used for automated testing, loads from .env.test and configs/config.test.json

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

	// Validate environment
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

	// Load secrets and config
	secrets, err := LoadSecrets()
	if err != nil {
		return nil, errors.WithContext(err, "loading secrets")
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, errors.WithContext(err, "loading config")
	}

	// Validate config with environment context
	if err := config.Validate(env); err != nil {
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

// LoadAppContextWithContext Add context support to handle timeouts and cancellation
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

	// Rest of implementation mirrors LoadAppContext but with context checks
	if !ValidateEnvironment(env) {
		if strictValidation {
			return nil, fmt.Errorf("invalid environment: %s", env)
		}
		log.Printf("Warning: Unexpected environment '%s' may cause issues", env)
	}

	// Special handling for test environment
	if env == "test" {
		log.Printf("Using test environment - ensure test fixtures are available")
	}

	// Add context check before expensive operations
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while loading environment files: %w", ctx.Err())
	default:
	}

	// Load environment files with improved approach
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

	// Load secrets and config
	secrets, err := LoadSecrets()
	if err != nil {
		return nil, errors.WithContext(err, "loading secrets")
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, errors.WithContext(err, "loading config")
	}

	// Validate config with environment context
	if err := config.Validate(env); err != nil {
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

// getTestConfiguration provides specialized test configuration
func getTestConfiguration() (*Config, *Secrets, error) {
	// Check if specific test config file exists
	testConfigFile := "configs/config.test.json"
	if _, err := os.Stat(testConfigFile); err == nil {
		config, err := LoadConfig(testConfigFile)
		if err != nil {
			return nil, nil, err
		}

		// Still use mock secrets to avoid requiring real credentials
		mockSecrets := &Secrets{
			ServiceAccountID:     "test-service-account-id",
			ServiceAccountSecret: "test-service-account-secret",
		}

		return config, mockSecrets, nil
	}

	// Fall back to fully mocked configuration
	return &Config{
			BaseURL:     "https://cloud-mock.mongodb.com",
			OrgID:       "test-org-id",
			ProjectID:   "test-project-id",
			ClusterName: "TestCluster",
			ProcessID:   "test-cluster-shard-00-00.test.mongodb.net:27017",
			HostName:    "test-cluster-shard-00-00.test.mongodb.net",
		}, &Secrets{
			ServiceAccountID:     "test-service-account-id",
			ServiceAccountSecret: "test-service-account-secret",
		}, nil
}

// Add diff support for testing
func (a *AppContext) Diff(other *AppContext) []string {
	var differences []string

	if a.Environment != other.Environment {
		differences = append(differences, fmt.Sprintf("Environment: %s vs %s",
			a.Environment, other.Environment))
	}

	// Compare important config fields
	if a.Config.BaseURL != other.Config.BaseURL {
		differences = append(differences, fmt.Sprintf("BaseURL: %s vs %s",
			a.Config.BaseURL, other.Config.BaseURL))
	}

	if a.Config.OrgID != other.Config.OrgID {
		differences = append(differences, fmt.Sprintf("OrgID: %s vs %s",
			a.Config.OrgID, other.Config.OrgID))
	}

	if a.Config.ProjectID != other.Config.ProjectID {
		differences = append(differences, fmt.Sprintf("ProjectID: %s vs %s",
			a.Config.ProjectID, other.Config.ProjectID))
	}

	if a.Config.ClusterName != other.Config.ClusterName {
		differences = append(differences, fmt.Sprintf("ClusterName: %s vs %s",
			a.Config.ClusterName, other.Config.ClusterName))
	}

	if a.Config.ProcessID != other.Config.ProcessID {
		differences = append(differences, fmt.Sprintf("ProcessID: %s vs %s",
			a.Config.ProcessID, other.Config.ProcessID))
	}

	return differences
}
