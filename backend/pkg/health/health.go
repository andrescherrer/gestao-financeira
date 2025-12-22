package health

import (
	"time"

	"gestao-financeira/backend/pkg/database"

	"github.com/gofiber/fiber/v2"
)

// HealthChecker handles health check endpoints
type HealthChecker struct{}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

// LivenessCheck checks if the application is alive
func (h *HealthChecker) LivenessCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "alive",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// ReadinessCheck checks if the application is ready to serve traffic
func (h *HealthChecker) ReadinessCheck(c *fiber.Ctx) error {
	checks := make(map[string]string)
	allHealthy := true

	// Check database connection
	if err := database.Ping(); err != nil {
		checks["database"] = "unhealthy"
		allHealthy = false
	} else {
		checks["database"] = "healthy"
	}

	// TODO: Add Redis check when Redis is configured
	// if err := redis.Ping(ctx); err != nil {
	// 	checks["cache"] = "unhealthy"
	// 	allHealthy = false
	// } else {
	// 	checks["cache"] = "healthy"
	// }

	if !allHealthy {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "not ready",
			"checks": checks,
			"time":   time.Now().Format(time.RFC3339),
		})
	}

	return c.JSON(fiber.Map{
		"status": "ready",
		"checks": checks,
		"time":   time.Now().Format(time.RFC3339),
	})
}

// HealthCheck is a simple health check endpoint
func (h *HealthChecker) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "gestao-financeira",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	})
}
