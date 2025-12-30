package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Environment
	Environment string `json:"environment"` // dev, staging, production

	// Server
	Server ServerConfig `json:"server"`

	// Database
	Database DatabaseConfig `json:"database"`

	// JWT
	JWT JWTConfig `json:"jwt"`

	// Redis
	Redis RedisConfig `json:"redis"`

	// Logging
	Logging LoggingConfig `json:"logging"`

	// Migrations
	Migrations MigrationsConfig `json:"migrations"`

	// CORS
	CORS CORSConfig `json:"cors"`

	// Observability
	Observability ObservabilityConfig `json:"observability"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	BodyLimit    int64         `json:"body_limit"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string        `json:"host"`
	Port            string        `json:"port"`
	User            string        `json:"user"`
	Password        string        `json:"password"`
	DBName          string        `json:"db_name"`
	SSLMode         string        `json:"ssl_mode"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey     string        `json:"secret_key"`
	Expiration    time.Duration `json:"expiration"`
	Issuer        string        `json:"issuer"`
	SigningMethod string        `json:"signing_method"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URL     string        `json:"url"`
	TTL     time.Duration `json:"ttl"`
	Enabled bool          `json:"enabled"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"` // json, console
}

// MigrationsConfig holds migrations configuration
type MigrationsConfig struct {
	Path string `json:"path"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins string `json:"allowed_origins"`
	MaxAge         int    `json:"max_age"`
}

// ObservabilityConfig holds observability configuration
type ObservabilityConfig struct {
	Tracing TracingConfig `json:"tracing"`
}

// TracingConfig holds tracing configuration
type TracingConfig struct {
	Enabled     bool   `json:"enabled"`
	ServiceName string `json:"service_name"`
	JaegerURL   string `json:"jaeger_url"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENV", "dev"),
		Server: ServerConfig{
			Port:         getEnv("API_PORT", "8080"),
			ReadTimeout:  parseDuration(getEnv("API_READ_TIMEOUT", "10s"), 10*time.Second),
			WriteTimeout: parseDuration(getEnv("API_WRITE_TIMEOUT", "10s"), 10*time.Second),
			IdleTimeout:  parseDuration(getEnv("API_IDLE_TIMEOUT", "120s"), 120*time.Second),
			BodyLimit:    parseInt64(getEnv("API_BODY_LIMIT", "10485760"), 10*1024*1024), // 10MB
		},
		Database: DatabaseConfig{
			Host:            getEnv("POSTGRES_HOST", "localhost"),
			Port:            getEnv("POSTGRES_PORT", "5432"),
			User:            getEnv("POSTGRES_USER", "postgres"),
			Password:        getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:          getEnv("POSTGRES_DB", "gestao_financeira"),
			SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
			MaxOpenConns:    parseInt(getEnv("POSTGRES_MAX_OPEN_CONNS", "25"), 25),
			MaxIdleConns:    parseInt(getEnv("POSTGRES_MAX_IDLE_CONNS", "5"), 5),
			ConnMaxLifetime: parseDuration(getEnv("POSTGRES_CONN_MAX_LIFETIME", "5m"), 5*time.Minute),
		},
		JWT: JWTConfig{
			SecretKey:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration:    parseDuration(getEnv("JWT_EXPIRATION", "24h"), 24*time.Hour),
			Issuer:        getEnv("JWT_ISSUER", "gestao-financeira-api"),
			SigningMethod: getEnv("JWT_SIGNING_METHOD", "HS256"),
		},
		Redis: RedisConfig{
			URL:     getEnv("REDIS_URL", ""),
			TTL:     parseDuration(getEnv("REDIS_TTL", "5m"), 5*time.Minute),
			Enabled: getEnv("REDIS_URL", "") != "",
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "console"), // json or console
		},
		Migrations: MigrationsConfig{
			Path: getEnv("MIGRATIONS_PATH", "file://migrations"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
			MaxAge:         parseInt(getEnv("CORS_MAX_AGE", "86400"), 86400), // 24 hours
		},
		Observability: ObservabilityConfig{
			Tracing: TracingConfig{
				Enabled:     getEnv("TRACING_ENABLED", "false") == "true",
				ServiceName: getEnv("TRACING_SERVICE_NAME", "gestao-financeira-api"),
				JaegerURL:   getEnv("JAEGER_URL", "http://localhost:14268/api/traces"),
			},
		},
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate environment
	validEnvs := []string{"dev", "development", "staging", "production", "prod"}
	if !contains(validEnvs, strings.ToLower(c.Environment)) {
		return fmt.Errorf("invalid environment: %s (must be one of: %v)", c.Environment, validEnvs)
	}

	// Validate database configuration
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.Password == "" && c.Environment == "production" {
		return fmt.Errorf("database password is required in production")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate JWT configuration
	if c.JWT.SecretKey == "" || c.JWT.SecretKey == "your-secret-key-change-in-production" {
		if c.Environment == "production" || c.Environment == "prod" {
			return fmt.Errorf("JWT secret key must be set in production")
		}
	}

	// Validate logging level
	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic", "trace"}
	if !contains(validLevels, strings.ToLower(c.Logging.Level)) {
		return fmt.Errorf("invalid log level: %s (must be one of: %v)", c.Logging.Level, validLevels)
	}

	// Validate logging format
	validFormats := []string{"json", "console"}
	if !contains(validFormats, strings.ToLower(c.Logging.Format)) {
		return fmt.Errorf("invalid log format: %s (must be one of: %v)", c.Logging.Format, validFormats)
	}

	return nil
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "dev" || c.Environment == "development"
}

// IsStaging returns true if the environment is staging
func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}

func parseInt64(s string, defaultValue int64) int64 {
	if s == "" {
		return defaultValue
	}
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}

func parseDuration(s string, defaultValue time.Duration) time.Duration {
	if s == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return defaultValue
	}
	return duration
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.ToLower(s) == strings.ToLower(item) {
			return true
		}
	}
	return false
}
