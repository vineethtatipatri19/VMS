package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Server ServerConfig

	// Database configuration
	Database DatabaseConfig

	// JWT configuration
	JWT JWTConfig

	// External services
	Gemini GeminiConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
	Env  string // development, staging, production
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	URL              string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  int
	MigrateOnStart   bool
	MigrationsPath   string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret     []byte
	Expiration int // hours
}

// GeminiConfig holds Gemini API configuration
type GeminiConfig struct {
	APIKey string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL:              getEnv("DATABASE_URL", ""),
			MaxOpenConns:     getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:     getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime:  getEnvAsInt("DB_CONN_MAX_LIFETIME", 300),
			MigrateOnStart:   getEnvAsBool("MIGRATE_ON_START", false),
			MigrationsPath:   getEnv("MIGRATIONS_PATH", "file://infra/migrations"),
		},
		JWT: JWTConfig{
			Secret:     []byte(getEnv("JWT_SECRET", "default-secret-change-in-production")),
			Expiration: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
		Gemini: GeminiConfig{
			APIKey: getEnv("GEMINI_API_KEY", ""),
		},
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if len(c.JWT.Secret) == 0 {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

// Helper functions to get environment variables

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
