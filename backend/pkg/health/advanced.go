package health

import (
	"context"
	"fmt"
	"runtime"
	"syscall"
	"time"

	"gestao-financeira/backend/pkg/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CheckResult represents the result of a health check
type CheckResult struct {
	Status      string                 `json:"status"`      // "healthy", "unhealthy", "degraded"
	Message     string                 `json:"message,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	ResponseTime string                `json:"response_time,omitempty"`
}

// AdvancedHealthChecker provides advanced health check capabilities
type AdvancedHealthChecker struct {
	healthChecker *HealthChecker
	startTime     time.Time
	version       string
}

// NewAdvancedHealthChecker creates a new advanced health checker
func NewAdvancedHealthChecker(healthChecker *HealthChecker, version string) *AdvancedHealthChecker {
	return &AdvancedHealthChecker{
		healthChecker: healthChecker,
		startTime:     time.Now(),
		version:       version,
	}
}

// DetailedHealthCheck returns a comprehensive health status
func (a *AdvancedHealthChecker) DetailedHealthCheck(c *fiber.Ctx) error {
	start := time.Now()
	checks := make(map[string]CheckResult)

	// Database check with details
	dbCheck := a.checkDatabase()
	checks["database"] = dbCheck

	// Cache check with details
	if a.healthChecker.cacheService != nil {
		cacheCheck := a.checkCache()
		checks["cache"] = cacheCheck
	}

	// Disk space check
	diskCheck := a.checkDiskSpace()
	checks["disk"] = diskCheck

	// Memory check
	memoryCheck := a.checkMemory()
	checks["memory"] = memoryCheck

	// System info
	systemInfo := a.getSystemInfo()

	// Determine overall status
	overallStatus := "healthy"
	statusCode := fiber.StatusOK
	for _, check := range checks {
		if check.Status == "unhealthy" {
			overallStatus = "unhealthy"
			statusCode = fiber.StatusServiceUnavailable
			break
		} else if check.Status == "degraded" && overallStatus == "healthy" {
			overallStatus = "degraded"
			statusCode = fiber.StatusOK // Still OK but degraded
		}
	}

	responseTime := time.Since(start)

	response := fiber.Map{
		"status":       overallStatus,
		"service":      "gestao-financeira",
		"version":      a.version,
		"uptime_seconds": int(time.Since(a.startTime).Seconds()),
		"uptime":       formatUptime(time.Since(a.startTime)),
		"timestamp":    time.Now().Format(time.RFC3339),
		"checks":       checks,
		"system":       systemInfo,
		"response_time_ms": responseTime.Milliseconds(),
	}

	return c.Status(statusCode).JSON(response)
}

// checkDatabase performs a detailed database health check
func (a *AdvancedHealthChecker) checkDatabase() CheckResult {
	start := time.Now()
	result := CheckResult{
		Status: "healthy",
		Details: make(map[string]interface{}),
	}

	// Basic ping
	if err := database.Ping(); err != nil {
		result.Status = "unhealthy"
		result.Message = fmt.Sprintf("Database ping failed: %v", err)
		return result
	}

	// Get database connection stats
	db := database.DB
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			stats := sqlDB.Stats()
			result.Details["open_connections"] = stats.OpenConnections
			result.Details["in_use"] = stats.InUse
			result.Details["idle"] = stats.Idle
			result.Details["max_open_connections"] = stats.MaxOpenConnections
			result.Details["max_idle_connections"] = stats.MaxIdleClosed

			// Check if connection pool is getting full
			if stats.MaxOpenConnections > 0 {
				usagePercent := float64(stats.OpenConnections) / float64(stats.MaxOpenConnections) * 100
				result.Details["connection_pool_usage_percent"] = fmt.Sprintf("%.2f", usagePercent)
				if usagePercent > 80 {
					result.Status = "degraded"
					result.Message = "Database connection pool usage is high"
				}
			}
		}

		// Test query performance
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		var count int64
		queryStart := time.Now()
		if err := db.WithContext(ctx).Raw("SELECT 1").Count(&count).Error; err != nil {
			result.Status = "degraded"
			result.Message = fmt.Sprintf("Database query test failed: %v", err)
		} else {
			queryTime := time.Since(queryStart)
			result.Details["query_response_time_ms"] = queryTime.Milliseconds()
			if queryTime > 1*time.Second {
				result.Status = "degraded"
				result.Message = "Database query response time is slow"
			}
		}
	}

	result.ResponseTime = fmt.Sprintf("%dms", time.Since(start).Milliseconds())
	return result
}

// checkCache performs a detailed cache health check
func (a *AdvancedHealthChecker) checkCache() CheckResult {
	start := time.Now()
	result := CheckResult{
		Status:  "healthy",
		Details: make(map[string]interface{}),
	}

	if a.healthChecker.cacheService == nil {
		result.Status = "unhealthy"
		result.Message = "Cache service not configured"
		return result
	}

	if err := a.healthChecker.cacheService.Ping(); err != nil {
		result.Status = "unhealthy"
		result.Message = fmt.Sprintf("Cache ping failed: %v", err)
		return result
	}

	result.ResponseTime = fmt.Sprintf("%dms", time.Since(start).Milliseconds())
	return result
}

// checkDiskSpace checks available disk space
func (a *AdvancedHealthChecker) checkDiskSpace() CheckResult {
	result := CheckResult{
		Status:  "healthy",
		Details: make(map[string]interface{}),
	}

	var stat syscall.Statfs_t
	wd := "/"
	if err := syscall.Statfs(wd, &stat); err != nil {
		result.Status = "degraded"
		result.Message = fmt.Sprintf("Unable to check disk space: %v", err)
		return result
	}

	// Calculate available and total bytes
	availableBytes := stat.Bavail * uint64(stat.Bsize)
	totalBytes := stat.Blocks * uint64(stat.Bsize)
	usedBytes := totalBytes - availableBytes

	// Convert to GB
	availableGB := float64(availableBytes) / (1024 * 1024 * 1024)
	totalGB := float64(totalBytes) / (1024 * 1024 * 1024)
	usedGB := float64(usedBytes) / (1024 * 1024 * 1024)
	usagePercent := (float64(usedBytes) / float64(totalBytes)) * 100

	result.Details["total_gb"] = fmt.Sprintf("%.2f", totalGB)
	result.Details["used_gb"] = fmt.Sprintf("%.2f", usedGB)
	result.Details["available_gb"] = fmt.Sprintf("%.2f", availableGB)
	result.Details["usage_percent"] = fmt.Sprintf("%.2f", usagePercent)

	// Warn if disk usage is above 85%
	if usagePercent > 85 {
		result.Status = "degraded"
		result.Message = fmt.Sprintf("Disk usage is high: %.2f%%", usagePercent)
	}

	// Critical if disk usage is above 95%
	if usagePercent > 95 {
		result.Status = "unhealthy"
		result.Message = fmt.Sprintf("Disk usage is critical: %.2f%%", usagePercent)
	}

	// Warn if less than 5GB available
	if availableGB < 5 {
		if result.Status == "healthy" {
			result.Status = "degraded"
		}
		result.Message = fmt.Sprintf("Low disk space: %.2f GB available", availableGB)
	}

	return result
}

// checkMemory checks available system memory
func (a *AdvancedHealthChecker) checkMemory() CheckResult {
	result := CheckResult{
		Status:  "healthy",
		Details: make(map[string]interface{}),
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Convert bytes to MB
	allocMB := float64(m.Alloc) / (1024 * 1024)
	totalAllocMB := float64(m.TotalAlloc) / (1024 * 1024)
	sysMB := float64(m.Sys) / (1024 * 1024)
	numGC := m.NumGC

	result.Details["allocated_mb"] = fmt.Sprintf("%.2f", allocMB)
	result.Details["total_allocated_mb"] = fmt.Sprintf("%.2f", totalAllocMB)
	result.Details["system_mb"] = fmt.Sprintf("%.2f", sysMB)
	result.Details["gc_runs"] = numGC
	result.Details["num_goroutines"] = runtime.NumGoroutine()

	// Warn if memory usage is very high (above 1GB allocated)
	if allocMB > 1024 {
		result.Status = "degraded"
		result.Message = fmt.Sprintf("High memory usage: %.2f MB allocated", allocMB)
	}

	return result
}

// getSystemInfo returns system information
func (a *AdvancedHealthChecker) getSystemInfo() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"go_version":      runtime.Version(),
		"go_os":           runtime.GOOS,
		"go_arch":         runtime.GOARCH,
		"num_cpu":         runtime.NumCPU(),
		"num_goroutines":  runtime.NumGoroutine(),
		"uptime_seconds": int(time.Since(a.startTime).Seconds()),
		"uptime":          formatUptime(time.Since(a.startTime)),
	}
}

// formatUptime formats uptime as a human-readable string
func formatUptime(duration time.Duration) string {
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// StartTime returns the start time of the application
func (a *AdvancedHealthChecker) StartTime() time.Time {
	return a.startTime
}

