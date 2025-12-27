package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	accountusecases "gestao-financeira/backend/internal/account/application/usecases"
	accountinfrahandlers "gestao-financeira/backend/internal/account/infrastructure/handlers"
	accountpersistence "gestao-financeira/backend/internal/account/infrastructure/persistence"
	accounthandlers "gestao-financeira/backend/internal/account/presentation/handlers"
	accountroutes "gestao-financeira/backend/internal/account/presentation/routes"
	budgetusecases "gestao-financeira/backend/internal/budget/application/usecases"
	budgetpersistence "gestao-financeira/backend/internal/budget/infrastructure/persistence"
	budgethandlers "gestao-financeira/backend/internal/budget/presentation/handlers"
	budgetroutes "gestao-financeira/backend/internal/budget/presentation/routes"
	categoryusecases "gestao-financeira/backend/internal/category/application/usecases"
	categorypersistence "gestao-financeira/backend/internal/category/infrastructure/persistence"
	categoryhandlers "gestao-financeira/backend/internal/category/presentation/handlers"
	categoryroutes "gestao-financeira/backend/internal/category/presentation/routes"
	"gestao-financeira/backend/internal/identity/application/usecases"
	"gestao-financeira/backend/internal/identity/infrastructure/persistence"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/identity/presentation/handlers"
	"gestao-financeira/backend/internal/identity/presentation/routes"
	reportingusecases "gestao-financeira/backend/internal/reporting/application/usecases"
	reporthandlers "gestao-financeira/backend/internal/reporting/presentation/handlers"
	reportroutes "gestao-financeira/backend/internal/reporting/presentation/routes"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	sharedhandlers "gestao-financeira/backend/internal/shared/infrastructure/handlers"
	transactionusecases "gestao-financeira/backend/internal/transaction/application/usecases"
	transactionpersistence "gestao-financeira/backend/internal/transaction/infrastructure/persistence"
	transactionhandlers "gestao-financeira/backend/internal/transaction/presentation/handlers"
	transactionroutes "gestao-financeira/backend/internal/transaction/presentation/routes"
	"gestao-financeira/backend/pkg/database"
	"gestao-financeira/backend/pkg/health"
	"gestao-financeira/backend/pkg/logger"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"

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

	// Initialize validator
	validator.Init()
	log.Info().Msg("Validator initialized")

	// Initialize Event Bus
	eventBus := eventbus.NewEventBus()

	// Initialize event logger handler for observability
	eventLoggerHandler := sharedhandlers.NewEventLoggerHandler()

	// Subscribe logger to all domain events
	// Note: In a production system, you might want to use a wildcard subscription
	// For now, we'll subscribe to common event types
	eventBus.Subscribe("UserRegistered", eventLoggerHandler.Handle)
	eventBus.Subscribe("AccountCreated", eventLoggerHandler.Handle)
	eventBus.Subscribe("AccountBalanceUpdated", eventLoggerHandler.Handle)
	eventBus.Subscribe("AccountNameChanged", eventLoggerHandler.Handle)
	eventBus.Subscribe("AccountDeactivated", eventLoggerHandler.Handle)
	eventBus.Subscribe("AccountActivated", eventLoggerHandler.Handle)
	eventBus.Subscribe("TransactionCreated", eventLoggerHandler.Handle)
	eventBus.Subscribe("TransactionUpdated", eventLoggerHandler.Handle)
	eventBus.Subscribe("TransactionDeleted", eventLoggerHandler.Handle)
	eventBus.Subscribe("BudgetCreated", eventLoggerHandler.Handle)
	eventBus.Subscribe("BudgetDeleted", eventLoggerHandler.Handle)

	// Initialize repositories
	userRepository := persistence.NewGormUserRepository(db)
	accountRepository := accountpersistence.NewGormAccountRepository(db)
	transactionRepository := transactionpersistence.NewGormTransactionRepository(db)
	categoryRepository := categorypersistence.NewGormCategoryRepository(db)
	budgetRepository := budgetpersistence.NewGormBudgetRepository(db)

	// Initialize services
	jwtService := services.NewJWTService()

	// Initialize use cases
	registerUserUseCase := usecases.NewRegisterUserUseCase(userRepository, eventBus)
	loginUseCase := usecases.NewLoginUseCase(userRepository, jwtService)

	// Initialize account use cases
	createAccountUseCase := accountusecases.NewCreateAccountUseCase(accountRepository, eventBus)
	listAccountsUseCase := accountusecases.NewListAccountsUseCase(accountRepository)
	getAccountUseCase := accountusecases.NewGetAccountUseCase(accountRepository)

	// Initialize account event handlers
	updateBalanceHandler := accountinfrahandlers.NewUpdateBalanceHandler(accountRepository)

	// Subscribe to transaction events for balance updates
	eventBus.Subscribe("TransactionCreated", updateBalanceHandler.HandleTransactionCreated)
	eventBus.Subscribe("TransactionUpdated", updateBalanceHandler.HandleTransactionUpdated)
	eventBus.Subscribe("TransactionDeleted", updateBalanceHandler.HandleTransactionDeleted)

	// Initialize transaction use cases
	createTransactionUseCase := transactionusecases.NewCreateTransactionUseCase(transactionRepository, eventBus)
	listTransactionsUseCase := transactionusecases.NewListTransactionsUseCase(transactionRepository)
	getTransactionUseCase := transactionusecases.NewGetTransactionUseCase(transactionRepository)
	updateTransactionUseCase := transactionusecases.NewUpdateTransactionUseCase(transactionRepository, eventBus)
	deleteTransactionUseCase := transactionusecases.NewDeleteTransactionUseCase(transactionRepository, eventBus)

	// Initialize category use cases
	createCategoryUseCase := categoryusecases.NewCreateCategoryUseCase(categoryRepository, eventBus)
	listCategoriesUseCase := categoryusecases.NewListCategoriesUseCase(categoryRepository)
	getCategoryUseCase := categoryusecases.NewGetCategoryUseCase(categoryRepository)
	updateCategoryUseCase := categoryusecases.NewUpdateCategoryUseCase(categoryRepository, eventBus)
	deleteCategoryUseCase := categoryusecases.NewDeleteCategoryUseCase(categoryRepository)

	// Initialize budget use cases
	createBudgetUseCase := budgetusecases.NewCreateBudgetUseCase(budgetRepository, eventBus)
	listBudgetsUseCase := budgetusecases.NewListBudgetsUseCase(budgetRepository)
	getBudgetUseCase := budgetusecases.NewGetBudgetUseCase(budgetRepository)
	updateBudgetUseCase := budgetusecases.NewUpdateBudgetUseCase(budgetRepository, eventBus)
	deleteBudgetUseCase := budgetusecases.NewDeleteBudgetUseCase(budgetRepository, eventBus)
	getBudgetProgressUseCase := budgetusecases.NewGetBudgetProgressUseCase(budgetRepository, transactionRepository)

	// Initialize reporting use cases
	monthlyReportUseCase := reportingusecases.NewMonthlyReportUseCase(transactionRepository)
	annualReportUseCase := reportingusecases.NewAnnualReportUseCase(transactionRepository)
	categoryReportUseCase := reportingusecases.NewCategoryReportUseCase(transactionRepository)
	incomeVsExpenseUseCase := reportingusecases.NewIncomeVsExpenseUseCase(transactionRepository)

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
	categoryHandler := categoryhandlers.NewCategoryHandler(
		createCategoryUseCase,
		listCategoriesUseCase,
		getCategoryUseCase,
		updateCategoryUseCase,
		deleteCategoryUseCase,
	)
	budgetHandler := budgethandlers.NewBudgetHandler(
		createBudgetUseCase,
		listBudgetsUseCase,
		getBudgetUseCase,
		updateBudgetUseCase,
		deleteBudgetUseCase,
		getBudgetProgressUseCase,
	)
	reportHandler := reporthandlers.NewReportHandler(
		monthlyReportUseCase,
		annualReportUseCase,
		categoryReportUseCase,
		incomeVsExpenseUseCase,
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

	// Store environment in locals for error handler
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("env", os.Getenv("ENV"))
		return c.Next()
	})

	// Request ID middleware (must be early in the chain)
	app.Use(middleware.RequestIDMiddleware())

	// Global middlewares
	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path} [${header:X-Request-ID}]\n",
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
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		AllowCredentials: true,
		MaxAge:           86400, // 24 horas
	}))

	// Error handler middleware (must be after all other middlewares)
	app.Use(middleware.ErrorHandlerMiddleware())

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
		routes.SetupAuthRoutes(api, authHandler)

		// Setup account routes (protected)
		accountroutes.SetupAccountRoutes(api, accountHandler, jwtService)

		// Setup transaction routes (protected)
		transactionroutes.SetupTransactionRoutes(api, transactionHandler, jwtService)

		// Setup category routes (protected)
		categoryroutes.SetupCategoryRoutes(api, categoryHandler, jwtService)

		// Setup budget routes (protected)
		budgetroutes.SetupBudgetRoutes(api, budgetHandler, jwtService)

		// Setup report routes (protected)
		reportroutes.SetupReportRoutes(api, reportHandler, jwtService)
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

// customErrorHandler handles errors globally (fallback for Fiber errors)
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Get request ID
	requestID := middleware.GetRequestID(c)

	code := fiber.StatusInternalServerError
	message := "Internal server error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Log error with structured logging
	if code == fiber.StatusInternalServerError {
		log.Error().
			Str("request_id", requestID).
			Err(err).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Int("status", code).
			Msg("Internal server error")
	} else {
		log.Warn().
			Str("request_id", requestID).
			Err(err).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Int("status", code).
			Msg("Request error")
	}

	response := fiber.Map{
		"error": message,
		"code":  code,
	}

	if requestID != "" {
		response["request_id"] = requestID
	}

	return c.Status(code).JSON(response)
}
