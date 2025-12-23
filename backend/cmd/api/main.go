package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	accountusecases "gestao-financeira/backend/internal/account/application/usecases"
	accountpersistence "gestao-financeira/backend/internal/account/infrastructure/persistence"
	accounthandlers "gestao-financeira/backend/internal/account/presentation/handlers"
	accountroutes "gestao-financeira/backend/internal/account/presentation/routes"
	"gestao-financeira/backend/internal/identity/application/usecases"
	"gestao-financeira/backend/internal/identity/infrastructure/persistence"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/identity/presentation/handlers"
	"gestao-financeira/backend/internal/identity/presentation/routes"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	transactionusecases "gestao-financeira/backend/internal/transaction/application/usecases"
	transactionpersistence "gestao-financeira/backend/internal/transaction/infrastructure/persistence"
	transactionhandlers "gestao-financeira/backend/internal/transaction/presentation/handlers"
	transactionroutes "gestao-financeira/backend/internal/transaction/presentation/routes"
	"gestao-financeira/backend/pkg/database"
	"gestao-financeira/backend/pkg/health"
	"gestao-financeira/backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "gestao-financeira/backend/docs" // swagger docs
)

// @title Gestão Financeira API
// @version 1.0
// @description API REST para gestão financeira pessoal e profissional. Gerencia usuários, contas e transações financeiras.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@gestaofinanceira.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize structured logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.InitLogger(logLevel)

	log.Info().Msg("Starting Gestão Financeira API")

	// Initialize database connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	log.Info().Msg("Database connection established")

	// Initialize Event Bus
	eventBus := eventbus.NewEventBus()

	// Initialize repositories
	userRepository := persistence.NewGormUserRepository(db)
	accountRepository := accountpersistence.NewGormAccountRepository(db)
	transactionRepository := transactionpersistence.NewGormTransactionRepository(db)

	// Initialize services
	jwtService := services.NewJWTService()

	// Initialize use cases
	registerUserUseCase := usecases.NewRegisterUserUseCase(userRepository, eventBus)
	loginUseCase := usecases.NewLoginUseCase(userRepository, jwtService)

	// Initialize account use cases
	createAccountUseCase := accountusecases.NewCreateAccountUseCase(accountRepository, eventBus)
	listAccountsUseCase := accountusecases.NewListAccountsUseCase(accountRepository)
	getAccountUseCase := accountusecases.NewGetAccountUseCase(accountRepository)

	// Initialize transaction use cases
	createTransactionUseCase := transactionusecases.NewCreateTransactionUseCase(transactionRepository, eventBus)
	listTransactionsUseCase := transactionusecases.NewListTransactionsUseCase(transactionRepository)
	getTransactionUseCase := transactionusecases.NewGetTransactionUseCase(transactionRepository)
	updateTransactionUseCase := transactionusecases.NewUpdateTransactionUseCase(transactionRepository, eventBus)
	deleteTransactionUseCase := transactionusecases.NewDeleteTransactionUseCase(transactionRepository)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(registerUserUseCase, loginUseCase)
	accountHandler := accounthandlers.NewAccountHandler(createAccountUseCase, listAccountsUseCase, getAccountUseCase)
	transactionHandler := transactionhandlers.NewTransactionHandler(
		createTransactionUseCase,
		listTransactionsUseCase,
		getTransactionUseCase,
		updateTransactionUseCase,
		deleteTransactionUseCase,
	)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Gestão Financeira API",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		// Increase header size limit to prevent 431 errors
		// Default is 4096 bytes, increasing to 16384 bytes (16KB)
		ReadBufferSize:  16384,
		WriteBufferSize: 16384,
		// Body size limit (default is 4MB, increasing to 10MB)
		BodyLimit: 10 * 1024 * 1024,
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

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// API v1 routes
	api := app.Group("/api/v1")
	{
		api.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Gestão Financeira API v1",
				"status":  "running",
			})
		})

		// Setup authentication routes (public)
		routes.SetupAuthRoutes(app, authHandler)

		// Setup account routes (protected)
		accountroutes.SetupAccountRoutes(app, accountHandler, jwtService)

		// Setup transaction routes (protected)
		transactionroutes.SetupTransactionRoutes(app, transactionHandler, jwtService)
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
