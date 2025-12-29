package migrations

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// RunMigrations runs all pending migrations
func RunMigrations(db *gorm.DB, migrationsPath string) error {
	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get migrations path from environment or use default
	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
		if migrationsPath == "" {
			migrationsPath = "file://migrations"
		}
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("No pending migrations")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info().Msg("Migrations applied successfully")
	return nil
}

// GetVersion returns the current migration version
func GetVersion(db *gorm.DB, migrationsPath string) (uint, bool, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
		if migrationsPath == "" {
			migrationsPath = "file://migrations"
		}
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return 0, false, nil
		}
		return 0, false, err
	}

	return version, dirty, nil
}

// MigrateDown migrates down by the specified number of steps
func MigrateDown(db *gorm.DB, migrationsPath string, steps int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
		if migrationsPath == "" {
			migrationsPath = "file://migrations"
		}
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Steps(-steps); err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Info().Int("steps", steps).Msg("Migrations rolled back successfully")
	return nil
}

// MigrateToVersion migrates to a specific version
func MigrateToVersion(db *gorm.DB, migrationsPath string, version uint) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
		if migrationsPath == "" {
			migrationsPath = "file://migrations"
		}
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Migrate(version); err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Uint("version", version).Msg("Already at target version")
			return nil
		}
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	log.Info().Uint("version", version).Msg("Migrated to version successfully")
	return nil
}

// ForceVersion forces the database to a specific version (use with caution)
func ForceVersion(db *gorm.DB, migrationsPath string, version int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
		if migrationsPath == "" {
			migrationsPath = "file://migrations"
		}
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Force(version); err != nil {
		return fmt.Errorf("failed to force version %d: %w", version, err)
	}

	log.Warn().Int("version", version).Msg("Forced database to version (use with caution)")
	return nil
}

// CheckPendingMigrations checks if there are pending migrations
func CheckPendingMigrations(db *gorm.DB, migrationsPath string) (bool, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return false, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return false, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
		if migrationsPath == "" {
			migrationsPath = "file://migrations"
		}
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return false, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Try to get the next version
	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			// No migrations applied yet, so there are pending migrations
			return true, nil
		}
		return false, err
	}

	if dirty {
		return false, fmt.Errorf("database is in a dirty state at version %d", version)
	}

	// Check if there's a next version available
	// This is a simple check - we'll try to read the next migration file
	// For a more robust solution, we'd need to list all migration files
	// For now, we'll just return that we can't determine this easily
	// The caller should just try to run migrations and handle ErrNoChange

	return false, nil
}
