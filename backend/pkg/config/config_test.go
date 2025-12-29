package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"ENV", "API_PORT", "POSTGRES_HOST", "POSTGRES_USER",
		"POSTGRES_PASSWORD", "POSTGRES_DB", "JWT_SECRET",
	}

	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	// Cleanup function
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	t.Run("load with defaults", func(t *testing.T) {
		// Clear environment
		for _, key := range envVars {
			os.Unsetenv(key)
		}

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg == nil {
			t.Fatal("Load() returned nil config")
		}

		// Check defaults
		if cfg.Environment != "dev" {
			t.Errorf("Expected environment 'dev', got %s", cfg.Environment)
		}
		if cfg.Server.Port != "8080" {
			t.Errorf("Expected port '8080', got %s", cfg.Server.Port)
		}
		if cfg.Database.Host != "localhost" {
			t.Errorf("Expected database host 'localhost', got %s", cfg.Database.Host)
		}
	})

	t.Run("load with environment variables", func(t *testing.T) {
		os.Setenv("ENV", "production")
		os.Setenv("API_PORT", "3000")
		os.Setenv("POSTGRES_HOST", "db.example.com")
		os.Setenv("POSTGRES_USER", "testuser")
		os.Setenv("POSTGRES_PASSWORD", "testpass")
		os.Setenv("POSTGRES_DB", "testdb")
		os.Setenv("JWT_SECRET", "test-secret-key")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Environment != "production" {
			t.Errorf("Expected environment 'production', got %s", cfg.Environment)
		}
		if cfg.Server.Port != "3000" {
			t.Errorf("Expected port '3000', got %s", cfg.Server.Port)
		}
		if cfg.Database.Host != "db.example.com" {
			t.Errorf("Expected database host 'db.example.com', got %s", cfg.Database.Host)
		}
	})

	t.Run("validate invalid environment", func(t *testing.T) {
		os.Setenv("ENV", "invalid")

		cfg, err := Load()
		if err == nil {
			t.Error("Expected error for invalid environment, got nil")
		}
		if cfg != nil {
			t.Error("Expected nil config on validation error")
		}
	})
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Environment: "dev",
				Database: DatabaseConfig{
					Host:   "localhost",
					User:   "postgres",
					DBName: "testdb",
				},
				JWT: JWTConfig{
					SecretKey: "test-secret",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "console",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid environment",
			config: &Config{
				Environment: "invalid",
				Database: DatabaseConfig{
					Host:   "localhost",
					User:   "postgres",
					DBName: "testdb",
				},
				JWT: JWTConfig{
					SecretKey: "test-secret",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "console",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database host",
			config: &Config{
				Environment: "dev",
				Database: DatabaseConfig{
					Host:   "",
					User:   "postgres",
					DBName: "testdb",
				},
				JWT: JWTConfig{
					SecretKey: "test-secret",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "console",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database user",
			config: &Config{
				Environment: "dev",
				Database: DatabaseConfig{
					Host:   "localhost",
					User:   "",
					DBName: "testdb",
				},
				JWT: JWTConfig{
					SecretKey: "test-secret",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "console",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database name",
			config: &Config{
				Environment: "dev",
				Database: DatabaseConfig{
					Host: "localhost",
					User: "postgres",
				},
				JWT: JWTConfig{
					SecretKey: "test-secret",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "console",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: &Config{
				Environment: "dev",
				Database: DatabaseConfig{
					Host:   "localhost",
					User:   "postgres",
					DBName: "testdb",
				},
				JWT: JWTConfig{
					SecretKey: "test-secret",
				},
				Logging: LoggingConfig{
					Level:  "invalid",
					Format: "console",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log format",
			config: &Config{
				Environment: "dev",
				Database: DatabaseConfig{
					Host:   "localhost",
					User:   "postgres",
					DBName: "testdb",
				},
				JWT: JWTConfig{
					SecretKey: "test-secret",
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name: "production environment",
			config: &Config{
				Environment: "production",
			},
			expected: true,
		},
		{
			name: "prod environment",
			config: &Config{
				Environment: "prod",
			},
			expected: true,
		},
		{
			name: "dev environment",
			config: &Config{
				Environment: "dev",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsProduction()
			if result != tt.expected {
				t.Errorf("IsProduction() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name: "dev environment",
			config: &Config{
				Environment: "dev",
			},
			expected: true,
		},
		{
			name: "development environment",
			config: &Config{
				Environment: "development",
			},
			expected: true,
		},
		{
			name: "production environment",
			config: &Config{
				Environment: "production",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsDevelopment()
			if result != tt.expected {
				t.Errorf("IsDevelopment() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{
			name:         "valid duration",
			input:        "5m",
			defaultValue: time.Second,
			expected:     5 * time.Minute,
		},
		{
			name:         "invalid duration",
			input:        "invalid",
			defaultValue: time.Second,
			expected:     time.Second,
		},
		{
			name:         "empty duration",
			input:        "",
			defaultValue: time.Second,
			expected:     time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDuration(tt.input, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("parseDuration() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultValue int
		expected     int
	}{
		{
			name:         "valid integer",
			input:        "42",
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "invalid integer",
			input:        "invalid",
			defaultValue: 0,
			expected:     0,
		},
		{
			name:         "empty string",
			input:        "",
			defaultValue: 0,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseInt(tt.input, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("parseInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}
