package database

import (
	"fmt"
	"os"
	"time"

	"gestao-financeira/backend/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// NewDatabase creates a new database connection using environment variables
func NewDatabase() (*gorm.DB, error) {
	dbConfig := getConfig()
	return NewDatabaseWithConfig(dbConfig)
}

// NewDatabaseWithConfig creates a new database connection using provided configuration
func NewDatabaseWithConfig(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings from configuration
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	return db, nil
}

// getConfig reads database configuration from environment variables
// This is kept for backward compatibility
func getConfig() config.DatabaseConfig {
	return config.DatabaseConfig{
		Host:            getEnv("POSTGRES_HOST", "localhost"),
		Port:            getEnv("POSTGRES_PORT", "5432"),
		User:            getEnv("POSTGRES_USER", "postgres"),
		Password:        getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:          getEnv("POSTGRES_DB", "gestao_financeira"),
		SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Ping checks if database connection is alive
func Ping() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}
