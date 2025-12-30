package metrics

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler returns a handler for the /metrics endpoint
func MetricsHandler() fiber.Handler {
	// Convert promhttp.Handler to Fiber handler
	return adaptor.HTTPHandler(promhttp.Handler())
}

