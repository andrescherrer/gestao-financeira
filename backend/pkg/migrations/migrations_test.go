package migrations

import (
	"os"
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func TestRunMigrations(t *testing.T) {
	// This test requires actual migration files
	// For a real test, we'd need to create temporary migration files
	// For now, we'll test error handling

	t.Run("invalid database connection", func(t *testing.T) {
		var db *gorm.DB = nil
		err := RunMigrations(db, "")
		if err == nil {
			t.Error("expected error for nil database, got nil")
		}
	})
}

func TestGetVersion(t *testing.T) {
	t.Run("invalid database connection", func(t *testing.T) {
		var db *gorm.DB = nil
		_, _, err := GetVersion(db, "")
		if err == nil {
			t.Error("expected error for nil database, got nil")
		}
	})
}

func TestMigrateDown(t *testing.T) {
	t.Run("invalid database connection", func(t *testing.T) {
		var db *gorm.DB = nil
		err := MigrateDown(db, "", 1)
		if err == nil {
			t.Error("expected error for nil database, got nil")
		}
	})
}

func TestMigrateToVersion(t *testing.T) {
	t.Run("invalid database connection", func(t *testing.T) {
		var db *gorm.DB = nil
		err := MigrateToVersion(db, "", 1)
		if err == nil {
			t.Error("expected error for nil database, got nil")
		}
	})
}

func TestForceVersion(t *testing.T) {
	t.Run("invalid database connection", func(t *testing.T) {
		var db *gorm.DB = nil
		err := ForceVersion(db, "", 1)
		if err == nil {
			t.Error("expected error for nil database, got nil")
		}
	})
}

func TestCheckPendingMigrations(t *testing.T) {
	t.Run("invalid database connection", func(t *testing.T) {
		var db *gorm.DB = nil
		_, err := CheckPendingMigrations(db, "")
		if err == nil {
			t.Error("expected error for nil database, got nil")
		}
	})
}

// Helper function to create temporary migration files for testing
func createTempMigrations(t *testing.T) (string, func()) {
	tempDir := t.TempDir()

	// Create up migration
	upFile := filepath.Join(tempDir, "000001_test.up.sql")
	err := os.WriteFile(upFile, []byte("CREATE TABLE test (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to create up migration: %v", err)
	}

	// Create down migration
	downFile := filepath.Join(tempDir, "000001_test.down.sql")
	err = os.WriteFile(downFile, []byte("DROP TABLE test;"), 0644)
	if err != nil {
		t.Fatalf("Failed to create down migration: %v", err)
	}

	migrationsPath := "file://" + tempDir

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return migrationsPath, cleanup
}
