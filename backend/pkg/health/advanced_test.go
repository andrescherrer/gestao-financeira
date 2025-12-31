package health

import (
	"testing"
	"time"
)

func TestAdvancedHealthChecker(t *testing.T) {
	healthChecker := NewHealthChecker()
	advancedChecker := NewAdvancedHealthChecker(healthChecker, "1.0.0")

	// Test start time
	startTime := advancedChecker.StartTime()
	if startTime.IsZero() {
		t.Error("Start time should not be zero")
	}

	// Test that start time is recent (within last second)
	if time.Since(startTime) > time.Second {
		t.Error("Start time should be recent")
	}
}

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{1 * time.Second, "1s"},
		{30 * time.Second, "30s"},
		{1 * time.Minute, "1m 0s"},
		{5 * time.Minute + 30 * time.Second, "5m 30s"},
		{1 * time.Hour, "1h 0m 0s"},
		{2 * time.Hour + 30 * time.Minute + 15 * time.Second, "2h 30m 15s"},
		{24 * time.Hour, "1d 0h 0m 0s"},
		{2 * 24 * time.Hour + 3 * time.Hour + 15 * time.Minute, "2d 3h 15m 0s"},
	}

	for _, tt := range tests {
		result := formatUptime(tt.duration)
		if result != tt.expected {
			t.Errorf("formatUptime(%v) = %v, want %v", tt.duration, result, tt.expected)
		}
	}
}

func TestCheckMemory(t *testing.T) {
	healthChecker := NewHealthChecker()
	advancedChecker := NewAdvancedHealthChecker(healthChecker, "1.0.0")

	result := advancedChecker.checkMemory()

	if result.Status == "" {
		t.Error("Memory check should return a status")
	}

	if result.Details == nil {
		t.Error("Memory check should return details")
	}

	// Check that details contain expected keys
	expectedKeys := []string{"allocated_mb", "num_goroutines"}
	for _, key := range expectedKeys {
		if _, ok := result.Details[key]; !ok {
			t.Errorf("Memory check details should contain key: %s", key)
		}
	}
}

func TestGetSystemInfo(t *testing.T) {
	healthChecker := NewHealthChecker()
	advancedChecker := NewAdvancedHealthChecker(healthChecker, "1.0.0")

	info := advancedChecker.getSystemInfo()

	if info == nil {
		t.Error("System info should not be nil")
	}

	// Check that info contains expected keys
	expectedKeys := []string{"go_version", "go_os", "go_arch", "num_cpu", "num_goroutines", "uptime_seconds"}
	for _, key := range expectedKeys {
		if _, ok := info[key]; !ok {
			t.Errorf("System info should contain key: %s", key)
		}
	}
}

