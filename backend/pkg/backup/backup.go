package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gestao-financeira/backend/pkg/config"

	"github.com/rs/zerolog/log"
)

// BackupService handles database backups
type BackupService struct {
	config     config.DatabaseConfig
	backupDir  string
	retention  int // days
	compressed bool
}

// NewBackupService creates a new backup service
func NewBackupService(dbConfig config.DatabaseConfig, backupDir string, retentionDays int, compressed bool) *BackupService {
	return &BackupService{
		config:     dbConfig,
		backupDir:  backupDir,
		retention:  retentionDays,
		compressed: compressed,
	}
}

// BackupResult contains information about a backup operation
type BackupResult struct {
	FilePath   string
	Size       int64
	Duration   time.Duration
	Compressed bool
	Error      error
}

// CreateBackup creates a database backup
func (s *BackupService) CreateBackup() (*BackupResult, error) {
	start := time.Now()

	// Ensure backup directory exists
	if err := os.MkdirAll(s.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	extension := ".sql"
	if s.compressed {
		extension = ".sql.gz"
	}
	filename := fmt.Sprintf("backup_%s%s", timestamp, extension)
	filePath := filepath.Join(s.backupDir, filename)

	// Create backup command
	var cmd *exec.Cmd
	if s.compressed {
		// Use pg_dump with gzip compression
		pgDump := exec.Command(
			"pg_dump",
			"-h", s.config.Host,
			"-p", s.config.Port,
			"-U", s.config.User,
			"-d", s.config.DBName,
			"-F", "c", // custom format (compressed)
			"-f", filePath,
		)
		pgDump.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", s.config.Password))
		cmd = pgDump
	} else {
		// Use pg_dump without compression
		pgDump := exec.Command(
			"pg_dump",
			"-h", s.config.Host,
			"-p", s.config.Port,
			"-U", s.config.User,
			"-d", s.config.DBName,
			"-f", filePath,
		)
		pgDump.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", s.config.Password))
		cmd = pgDump
	}

	// Execute backup
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute backup: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup file info: %w", err)
	}

	duration := time.Since(start)

	result := &BackupResult{
		FilePath:   filePath,
		Size:       fileInfo.Size(),
		Duration:   duration,
		Compressed: s.compressed,
	}

	log.Info().
		Str("file_path", filePath).
		Int64("size_bytes", fileInfo.Size()).
		Dur("duration", duration).
		Bool("compressed", s.compressed).
		Msg("Backup created successfully")

	return result, nil
}

// RestoreBackup restores a database from a backup file
func (s *BackupService) RestoreBackup(backupPath string) error {
	log.Info().
		Str("backup_path", backupPath).
		Msg("Starting database restore")

	// Check if backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupPath)
	}

	// Determine if backup is compressed
	isCompressed := strings.HasSuffix(backupPath, ".gz") || strings.HasSuffix(backupPath, ".sql.gz")

	var cmd *exec.Cmd
	if isCompressed {
		// For compressed backups, use pg_restore
		cmd = exec.Command(
			"pg_restore",
			"-h", s.config.Host,
			"-p", s.config.Port,
			"-U", s.config.User,
			"-d", s.config.DBName,
			"-c", // clean (drop objects before recreating)
			"-v", // verbose
			backupPath,
		)
	} else {
		// For plain SQL backups, use psql
		cmd = exec.Command(
			"psql",
			"-h", s.config.Host,
			"-p", s.config.Port,
			"-U", s.config.User,
			"-d", s.config.DBName,
			"-f", backupPath,
		)
	}

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", s.config.Password))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	log.Info().
		Str("backup_path", backupPath).
		Msg("Database restored successfully")

	return nil
}

// CleanupOldBackups removes backup files older than retention period
func (s *BackupService) CleanupOldBackups() (int, error) {
	if s.retention <= 0 {
		return 0, nil // Retention disabled
	}

	cutoffDate := time.Now().AddDate(0, 0, -s.retention)
	removedCount := 0

	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read backup directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if file is a backup file
		if !strings.HasPrefix(entry.Name(), "backup_") {
			continue
		}

		filePath := filepath.Join(s.backupDir, entry.Name())
		fileInfo, err := entry.Info()
		if err != nil {
			log.Warn().Err(err).Str("file", entry.Name()).Msg("Failed to get file info")
			continue
		}

		// Check if file is older than retention period
		if fileInfo.ModTime().Before(cutoffDate) {
			if err := os.Remove(filePath); err != nil {
				log.Error().Err(err).Str("file", filePath).Msg("Failed to remove old backup")
				continue
			}
			removedCount++
			log.Info().
				Str("file", filePath).
				Time("modified", fileInfo.ModTime()).
				Msg("Removed old backup file")
		}
	}

	if removedCount > 0 {
		log.Info().
			Int("removed_count", removedCount).
			Int("retention_days", s.retention).
			Msg("Cleanup of old backups completed")
	}

	return removedCount, nil
}

// ListBackups returns a list of available backup files
func (s *BackupService) ListBackups() ([]BackupInfo, error) {
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasPrefix(entry.Name(), "backup_") {
			continue
		}

		filePath := filepath.Join(s.backupDir, entry.Name())
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, BackupInfo{
			FileName: entry.Name(),
			FilePath: filePath,
			Size:     fileInfo.Size(),
			Modified: fileInfo.ModTime(),
		})
	}

	return backups, nil
}

// BackupInfo contains information about a backup file
type BackupInfo struct {
	FileName string
	FilePath string
	Size     int64
	Modified time.Time
}
