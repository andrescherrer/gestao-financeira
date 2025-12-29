package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/pkg/database"
	"gestao-financeira/backend/pkg/migrations"
)

func main() {
	// Setup logger
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	// Parse command line flags
	command := flag.String("command", "up", "Migration command: up, down, version, force, goto")
	steps := flag.Int("steps", 1, "Number of steps for down command")
	version := flag.Int("version", 0, "Version number for goto or force command")
	migrationsPath := flag.String("path", "file://migrations", "Path to migrations directory")
	flag.Parse()

	// Initialize database connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	log.Info().Msg("Database connection established")

	// Execute command
	switch *command {
	case "up":
		log.Info().Msg("Running migrations up...")
		if err := migrations.RunMigrations(db, *migrationsPath); err != nil {
			log.Fatal().Err(err).Msg("Failed to run migrations")
		}
		log.Info().Msg("Migrations completed successfully")

	case "down":
		log.Info().Int("steps", *steps).Msg("Running migrations down...")
		if err := migrations.MigrateDown(db, *migrationsPath, *steps); err != nil {
			log.Fatal().Err(err).Msg("Failed to rollback migrations")
		}
		log.Info().Msg("Rollback completed successfully")

	case "version":
		log.Info().Msg("Checking migration version...")
		ver, dirty, err := migrations.GetVersion(db, *migrationsPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get migration version")
		}
		if dirty {
			log.Warn().Uint("version", ver).Msg("Database is in a dirty state")
		} else {
			log.Info().Uint("version", ver).Msg("Current migration version")
		}

	case "goto":
		if *version == 0 {
			log.Fatal().Msg("Version number is required for goto command")
		}
		log.Info().Int("version", *version).Msg("Migrating to version...")
		if err := migrations.MigrateToVersion(db, *migrationsPath, uint(*version)); err != nil {
			log.Fatal().Err(err).Msg("Failed to migrate to version")
		}
		log.Info().Msg("Migration to version completed successfully")

	case "force":
		if *version == 0 {
			log.Fatal().Msg("Version number is required for force command")
		}
		log.Warn().Int("version", *version).Msg("Forcing database to version (use with caution)...")
		if err := migrations.ForceVersion(db, *migrationsPath, *version); err != nil {
			log.Fatal().Err(err).Msg("Failed to force version")
		}
		log.Info().Msg("Version forced successfully")

	default:
		log.Fatal().Str("command", *command).Msg("Unknown command")
		fmt.Println("Usage: migrate [flags]")
		fmt.Println("Commands:")
		fmt.Println("  up       - Run all pending migrations")
		fmt.Println("  down     - Rollback migrations (use -steps to specify number)")
		fmt.Println("  version  - Show current migration version")
		fmt.Println("  goto     - Migrate to specific version (use -version)")
		fmt.Println("  force    - Force database to specific version (use -version, use with caution)")
		fmt.Println("Flags:")
		fmt.Println("  -path    - Path to migrations directory (default: file://migrations)")
		fmt.Println("  -steps   - Number of steps for down command (default: 1)")
		fmt.Println("  -version - Version number for goto or force command")
		os.Exit(1)
	}
}
