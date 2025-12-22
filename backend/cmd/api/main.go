package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"gestao-financeira/backend/pkg/database"
	"gestao-financeira/backend/pkg/health"
	"gestao-financeira/backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize structured logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.InitLogger(logLevel)

	log.Info().Msg("Starting Gestão Financeira API")

	// Initialize database connection
	_, err := database.NewDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	log.Info().Msg("Database connection established")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Gestão Financeira API",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Global middlewares
	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "America/Sao_Paulo",
	}))

	// CORS middleware
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           86400, // 24 horas
	}))

	// Initialize health checker
	healthChecker := health.NewHealthChecker()

	// Health check endpoints
	app.Get("/health", healthChecker.HealthCheck)
	app.Get("/health/live", healthChecker.LivenessCheck)
	app.Get("/health/ready", healthChecker.ReadinessCheck)

	// API v1 routes (will be added in future tasks)
	api := app.Group("/api/v1")
	{
		// Routes will be added here
		api.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Gestão Financeira API v1",
				"status":  "running",
			})
		})
	}

	// Get port from environment or use default
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server in a goroutine
	go func() {
		log.Info().Str("port", port).Msg("Server starting")
		if err := app.Listen(":" + port); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Shutdown gracefully
	log.Info().Msg("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Error shutting down server")
	}

	// Close database connection
	if err := database.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing database")
	}

	log.Info().Msg("Server exited")
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Log error with structured logging
	if code == fiber.StatusInternalServerError {
		log.Error().
			Err(err).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Int("status", code).
			Msg("Internal server error")
	} else {
		log.Warn().
			Err(err).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Int("status", code).
			Msg("Request error")
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
		"code":  code,
	})
}
