package backup

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gestao-financeira/backend/pkg/config"
)

func TestNewBackupService(t *testing.T) {
	dbConfig := config.DatabaseConfig{
		Host:   "localhost",
		Port:   "5432",
		User:   "test",
		DBName: "test_db",
	}

	service := NewBackupService(dbConfig, "/tmp/backups", 30, true)

	if service == nil {
		t.Fatal("NewBackupService returned nil")
	}

	if service.backupDir != "/tmp/backups" {
		t.Errorf("Expected backupDir to be /tmp/backups, got %s", service.backupDir)
	}

	if service.retention != 30 {
		t.Errorf("Expected retention to be 30, got %d", service.retention)
	}

	if !service.compressed {
		t.Error("Expected compressed to be true")
	}
}

func TestCleanupOldBackups(t *testing.T) {
	// Create temporary directory for backups
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbConfig := config.DatabaseConfig{
		Host:   "localhost",
		Port:   "5432",
		User:   "test",
		DBName: "test_db",
	}

	service := NewBackupService(dbConfig, tmpDir, 7, false)

	// Create old backup file (older than retention)
	oldFile := filepath.Join(tmpDir, "backup_20240101_000000.sql")
	oldContent := []byte("old backup content")
	if err := os.WriteFile(oldFile, oldContent, 0644); err != nil {
		t.Fatalf("Failed to create old backup file: %v", err)
	}

	// Set file modification time to 30 days ago
	oldTime := time.Now().AddDate(0, 0, -30)
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatalf("Failed to set file time: %v", err)
	}

	// Create recent backup file (within retention)
	recentFile := filepath.Join(tmpDir, "backup_20250101_000000.sql")
	recentContent := []byte("recent backup content")
	if err := os.WriteFile(recentFile, recentContent, 0644); err != nil {
		t.Fatalf("Failed to create recent backup file: %v", err)
	}

	// Run cleanup
	removed, err := service.CleanupOldBackups()
	if err != nil {
		t.Fatalf("CleanupOldBackups failed: %v", err)
	}

	// Verify old file was removed
	if removed != 1 {
		t.Errorf("Expected 1 file to be removed, got %d", removed)
	}

	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("Old backup file should have been removed")
	}

	// Verify recent file still exists
	if _, err := os.Stat(recentFile); os.IsNotExist(err) {
		t.Error("Recent backup file should not have been removed")
	}
}

func TestCleanupOldBackups_NoRetention(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbConfig := config.DatabaseConfig{
		Host:   "localhost",
		Port:   "5432",
		User:   "test",
		DBName: "test_db",
	}

	// Service with retention disabled (0 or negative)
	service := NewBackupService(dbConfig, tmpDir, 0, false)

	removed, err := service.CleanupOldBackups()
	if err != nil {
		t.Fatalf("CleanupOldBackups failed: %v", err)
	}

	if removed != 0 {
		t.Errorf("Expected 0 files to be removed when retention is disabled, got %d", removed)
	}
}

func TestListBackups(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbConfig := config.DatabaseConfig{
		Host:   "localhost",
		Port:   "5432",
		User:   "test",
		DBName: "test_db",
	}

	service := NewBackupService(dbConfig, tmpDir, 30, false)

	// Create some backup files
	backupFiles := []string{
		"backup_20250101_000000.sql",
		"backup_20250102_120000.sql",
		"backup_20250103_150000.sql.gz",
	}

	for _, filename := range backupFiles {
		filePath := filepath.Join(tmpDir, filename)
		content := []byte("backup content")
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			t.Fatalf("Failed to create backup file %s: %v", filename, err)
		}
	}

	// Create a non-backup file (should be ignored)
	otherFile := filepath.Join(tmpDir, "other_file.txt")
	if err := os.WriteFile(otherFile, []byte("other content"), 0644); err != nil {
		t.Fatalf("Failed to create other file: %v", err)
	}

	// List backups
	backups, err := service.ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}

	if len(backups) != 3 {
		t.Errorf("Expected 3 backups, got %d", len(backups))
	}

	// Verify all backup files are listed
	foundFiles := make(map[string]bool)
	for _, backup := range backups {
		foundFiles[backup.FileName] = true
	}

	for _, filename := range backupFiles {
		if !foundFiles[filename] {
			t.Errorf("Backup file %s not found in list", filename)
		}
	}
}

func TestListBackups_EmptyDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "backup_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbConfig := config.DatabaseConfig{
		Host:   "localhost",
		Port:   "5432",
		User:   "test",
		DBName: "test_db",
	}

	service := NewBackupService(dbConfig, tmpDir, 30, false)

	backups, err := service.ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}

	if len(backups) != 0 {
		t.Errorf("Expected 0 backups in empty directory, got %d", len(backups))
	}
}
