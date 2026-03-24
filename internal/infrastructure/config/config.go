package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	JWT          JWTConfig
	Operations   OperationsConfig
	Integrations IntegrationsConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type JWTConfig struct {
	Secret     string
	ExpireHour int
}

type OperationsConfig struct {
	LogRetention LogRetentionConfig
}

type LogRetentionConfig struct {
	Enabled                bool
	CleanupIntervalMinutes int
	JobRunRetentionDays    int
	SystemLogRetentionDays int
}

type IntegrationsConfig struct {
	ConfigFile string
	Providers  []IntegrationProviderConfig
}

type IntegrationProviderConfig struct {
	Code                 string                `json:"code"`
	Type                 string                `json:"type"`
	Enabled              bool                  `json:"enabled"`
	AutoSyncEnabled      bool                  `json:"auto_sync_enabled"`
	Channel              string                `json:"channel"`
	SourceType           string                `json:"source_type"`
	SalesChannel         string                `json:"sales_channel"`
	DefaultCurrency      string                `json:"default_currency"`
	SyncIntervalMinutes  int                   `json:"sync_interval_minutes"`
	LookbackMinutes      int                   `json:"lookback_minutes"`
	InitialLookbackDays  int                   `json:"initial_lookback_days"`
	RequestTimeoutSecond int                   `json:"request_timeout_seconds"`
	Amazon               *AmazonProviderConfig `json:"amazon,omitempty"`
}

type AmazonProviderConfig struct {
	Endpoint         string   `json:"endpoint"`
	AppID            string   `json:"app_id"`
	AuthorizeBaseURL string   `json:"oauth_authorize_base_url"`
	RedirectURI      string   `json:"oauth_redirect_uri"`
	MarketplaceIDs   []string `json:"marketplace_ids"`
	LWAClientID      string   `json:"lwa_client_id"`
	LWAClientSecret  string   `json:"lwa_client_secret"`
	LWARefreshToken  string   `json:"lwa_refresh_token"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "321456"),
			Database: getEnv("DB_DATABASE", "am-erp"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "am-erp-secret-key-change-in-production"),
			ExpireHour: 24,
		},
		Operations: OperationsConfig{
			LogRetention: LogRetentionConfig{
				Enabled:                getEnvAsBool("LOG_RETENTION_ENABLED", true),
				CleanupIntervalMinutes: getEnvAsInt("LOG_RETENTION_CLEANUP_INTERVAL_MINUTES", 1440),
				JobRunRetentionDays:    getEnvAsInt("JOB_RUN_RETENTION_DAYS", 30),
				SystemLogRetentionDays: getEnvAsInt("SYSTEM_LOG_RETENTION_DAYS", 30),
			},
		},
		Integrations: IntegrationsConfig{
			ConfigFile: getEnv("INTEGRATION_CONFIG_FILE", "config/integrations.json"),
			Providers:  []IntegrationProviderConfig{},
		},
	}
	providers, err := loadIntegrationProviders(cfg.Integrations.ConfigFile)
	if err != nil {
		return nil, err
	}
	cfg.Integrations.Providers = providers
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return defaultValue
	}
	switch value {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return defaultValue
	}
}

func loadIntegrationProviders(path string) ([]IntegrationProviderConfig, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return []IntegrationProviderConfig{}, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []IntegrationProviderConfig{}, nil
		}
		return nil, fmt.Errorf("failed to read integration config file %s: %w", path, err)
	}

	var payload struct {
		Providers []IntegrationProviderConfig `json:"providers"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse integration config file %s: %w", path, err)
	}

	for i := range payload.Providers {
		normalizeProviderConfig(&payload.Providers[i])
	}
	return payload.Providers, nil
}

func normalizeProviderConfig(provider *IntegrationProviderConfig) {
	if provider == nil {
		return
	}

	provider.Type = strings.ToLower(strings.TrimSpace(provider.Type))
	provider.Code = strings.ToUpper(strings.TrimSpace(provider.Code))
	if provider.Code == "" {
		provider.Code = strings.ToUpper(provider.Type)
	}
	if provider.Channel == "" {
		provider.Channel = "DEFAULT"
	}
	if provider.SourceType == "" {
		provider.SourceType = "THIRD_PARTY_API"
	}
	if provider.SalesChannel == "" {
		provider.SalesChannel = provider.Code
	}
	if provider.DefaultCurrency == "" {
		provider.DefaultCurrency = "USD"
	}
	if provider.SyncIntervalMinutes <= 0 {
		provider.SyncIntervalMinutes = 30
	}
	if provider.LookbackMinutes <= 0 {
		provider.LookbackMinutes = 10
	}
	if provider.InitialLookbackDays <= 0 {
		provider.InitialLookbackDays = 7
	}
	if provider.RequestTimeoutSecond <= 0 {
		provider.RequestTimeoutSecond = 20
	}
	if provider.Type == "amazon" && provider.Amazon != nil {
		if strings.TrimSpace(provider.Amazon.Endpoint) == "" {
			provider.Amazon.Endpoint = "https://sellingpartnerapi-na.amazon.com"
		}
		if strings.TrimSpace(provider.Amazon.AuthorizeBaseURL) == "" {
			provider.Amazon.AuthorizeBaseURL = "https://sellercentral.amazon.com/apps/authorize/consent"
		}
	}
}
