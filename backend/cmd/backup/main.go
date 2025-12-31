package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/pkg/backup"
	"gestao-financeira/backend/pkg/config"
)

func main() {
	// Parse command line flags
	action := flag.String("action", "create", "Action to perform: create, restore, list, cleanup (default: create)")
	backupDir := flag.String("dir", "./backups", "Backup directory (default: ./backups)")
	backupFile := flag.String("file", "", "Backup file path (required for restore)")
	retention := flag.Int("retention", 30, "Retention period in days (default: 30)")
	compressed := flag.Bool("compressed", true, "Use compressed backup format (default: true)")
	flag.Parse()

	// Setup logger
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Create backup service
	backupService := backup.NewBackupService(
		cfg.Database,
		*backupDir,
		*retention,
		*compressed,
	)

	// Execute action
	switch *action {
	case "create":
		result, err := backupService.CreateBackup()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create backup")
		}
		log.Info().
			Str("file_path", result.FilePath).
			Int64("size_bytes", result.Size).
			Dur("duration", result.Duration).
			Bool("compressed", result.Compressed).
			Msg("Backup created successfully")
		fmt.Printf("Backup created: %s (%.2f MB, %v)\n",
			result.FilePath,
			float64(result.Size)/(1024*1024),
			result.Duration,
		)

	case "restore":
		if *backupFile == "" {
			log.Fatal().Msg("Backup file path is required for restore (use -file flag)")
		}
		if err := backupService.RestoreBackup(*backupFile); err != nil {
			log.Fatal().Err(err).Msg("Failed to restore backup")
		}
		log.Info().Str("file", *backupFile).Msg("Backup restored successfully")
		fmt.Printf("Backup restored from: %s\n", *backupFile)

	case "list":
		backups, err := backupService.ListBackups()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list backups")
		}
		if len(backups) == 0 {
			fmt.Println("No backups found")
			return
		}
		fmt.Printf("\nFound %d backup(s):\n\n", len(backups))
		fmt.Printf("%-50s %15s %20s\n", "File", "Size", "Modified")
		fmt.Println(strings.Repeat("-", 85))
		for _, b := range backups {
			sizeMB := float64(b.Size) / (1024 * 1024)
			fmt.Printf("%-50s %10.2f MB %20s\n",
				b.FileName,
				sizeMB,
				b.Modified.Format("2006-01-02 15:04:05"),
			)
		}
		fmt.Println()

	case "cleanup":
		removed, err := backupService.CleanupOldBackups()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to cleanup old backups")
		}
		log.Info().Int("removed_count", removed).Msg("Cleanup completed")
		fmt.Printf("Removed %d old backup file(s)\n", removed)

	default:
		log.Fatal().Str("action", *action).Msg("Invalid action. Use: create, restore, list, cleanup")
	}
}
