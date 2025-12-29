package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	accountusecases "gestao-financeira/backend/internal/account/application/usecases"
	accountrepositories "gestao-financeira/backend/internal/account/domain/repositories"
	accountinfrahandlers "gestao-financeira/backend/internal/account/infrastructure/handlers"
	accountpersistence "gestao-financeira/backend/internal/account/infrastructure/persistence"
	accounthandlers "gestao-financeira/backend/internal/account/presentation/handlers"
	accountroutes "gestao-financeira/backend/internal/account/presentation/routes"
	budgetusecases "gestao-financeira/backend/internal/budget/application/usecases"
	budgetpersistence "gestao-financeira/backend/internal/budget/infrastructure/persistence"
	budgethandlers "gestao-financeira/backend/internal/budget/presentation/handlers"
	budgetroutes "gestao-financeira/backend/internal/budget/presentation/routes"
	categoryusecases "gestao-financeira/backend/internal/category/application/usecases"
	categoryrepositories "gestao-financeira/backend/internal/category/domain/repositories"
	categorypersistence "gestao-financeira/backend/internal/category/infrastructure/persistence"
	categoryhandlers "gestao-financeira/backend/internal/category/presentation/handlers"
	categoryroutes "gestao-financeira/backend/internal/category/presentation/routes"
	"gestao-financeira/backend/internal/identity/application/usecases"
	"gestao-financeira/backend/internal/identity/infrastructure/persistence"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/identity/presentation/handlers"
	"gestao-financeira/backend/internal/identity/presentation/routes"
	reportingusecases "gestao-financeira/backend/internal/reporting/application/usecases"
	reportingservices "gestao-financeira/backend/internal/reporting/infrastructure/services"
	reporthandlers "gestao-financeira/backend/internal/reporting/presentation/handlers"
	reportroutes "gestao-financeira/backend/internal/reporting/presentation/routes"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	sharedhandlers "gestao-financeira/backend/internal/shared/infrastructure/handlers"
	sharedpersistence "gestao-financeira/backend/internal/shared/infrastructure/persistence"
	transactionusecases "gestao-financeira/backend/internal/transaction/application/usecases"
	transactionpersistence "gestao-financeira/backend/internal/transaction/infrastructure/persistence"
	transactionhandlers "gestao-financeira/backend/internal/transaction/presentation/handlers"
	transactionroutes "gestao-financeira/backend/internal/transaction/presentation/routes"
	"gestao-financeira/backend/pkg/cache"
	"gestao-financeira/backend/pkg/config"
	"gestao-financeira/backend/pkg/database"
	"gestao-financeira/backend/pkg/health"
	"gestao-financeira/backend/pkg/logger"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/migrations"
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
// @description API REST completa para gestão financeira pessoal e profissional.
//
// ## Funcionalidades Principais
//
// - **Autenticação e Autorização**: Sistema de autenticação JWT com refresh tokens
// - **Gerenciamento de Contas**: Criação e gestão de contas bancárias, carteiras e investimentos
// - **Transações Financeiras**: Registro de receitas e despesas com suporte a recorrência
// - **Categorias**: Organização de transações por categorias customizáveis
// - **Orçamentos**: Controle de orçamentos mensais e anuais por categoria
// - **Relatórios**: Análises financeiras e relatórios personalizados
//
// ## Características Técnicas
//
// - **Arquitetura**: Domain-Driven Design (DDD) com Clean Architecture
// - **Atomicidade**: Operações críticas garantidas por Unit of Work pattern
// - **Paginação**: Todos os endpoints de listagem suportam paginação
// - **Soft Delete**: Exclusão lógica com possibilidade de restauração
// - **Validação**: Validação em múltiplas camadas (frontend, backend, domain)
// - **Tratamento de Erros**: Erros tipados e consistentes em toda a API
//
// ## Autenticação
//
// A API utiliza autenticação JWT (JSON Web Tokens). Para acessar endpoints protegidos:
//
// 1. Faça login em `/auth/login` para obter um token
// 2. Inclua o token no header `Authorization: Bearer <token>`
// 3. O token expira após o tempo configurado (padrão: 24h)
//
// ## Paginação
//
// Endpoints de listagem suportam paginação via query parameters:
//
// - `page`: Número da página (1-based, padrão: 1)
// - `limit`: Itens por página (padrão: 10, máximo: 100)
//
// Exemplo: `GET /api/v1/transactions?page=2&limit=20`
//
// ## Códigos de Resposta HTTP
//
// - `200 OK`: Operação bem-sucedida
// - `201 Created`: Recurso criado com sucesso
// - `400 Bad Request`: Dados inválidos ou validação falhou
// - `401 Unauthorized`: Token ausente ou inválido
// - `403 Forbidden`: Acesso negado (sem permissão)
// - `404 Not Found`: Recurso não encontrado
// - `409 Conflict`: Conflito (ex: recurso já existe)
// - `422 Unprocessable Entity`: Erro de validação de domínio
// - `500 Internal Server Error`: Erro interno do servidor
//
// ## Rate Limiting
//
// A API implementa rate limiting para proteger contra abuso:
//
// - Limite padrão: 100 requisições por minuto por IP
// - Headers de resposta incluem informações sobre limites
//
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@gestaofinanceira.com
// @contact.url https://github.com/gestao-financeira
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize structured logger
	logger.InitLoggerWithConfig(cfg.Logging.Level, cfg.Logging.Format, cfg.IsProduction())

	log.Info().
		Str("environment", cfg.Environment).
		Str("port", cfg.Server.Port).
		Msg("Starting Gestão Financeira API")

	// Initialize database connection with configuration
	db, err := database.NewDatabaseWithConfig(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	log.Info().Msg("Database connection established")

	// Run database migrations
	if err := migrations.RunMigrations(db, cfg.Migrations.Path); err != nil {
		log.Warn().Err(err).Msg("Failed to run migrations (continuing anyway)")
		// Don't fail startup if migrations fail - might be intentional in some cases
		// In production, you might want to fail here
	} else {
		log.Info().Msg("Database migrations completed")
	}

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

	// Initialize cache service (optional - continues without cache if Redis is unavailable)
	var cacheService *cache.CacheService
	var reportCacheService *reportingservices.ReportCacheService
	if cfg.Redis.Enabled {
		cacheSvc, err := cache.NewCacheService(cfg.Redis.URL)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to initialize cache service, continuing without cache")
		} else {
			cacheService = cacheSvc
			reportCacheService = reportingservices.NewReportCacheService(cacheService, cfg.Redis.TTL)
		}
	}

	// Initialize repositories
	userRepository := persistence.NewGormUserRepository(db)

	// Initialize account repository with cache
	baseAccountRepository := accountpersistence.NewGormAccountRepository(db)
	accountRepository := accountpersistence.NewCachedAccountRepository(
		baseAccountRepository,
		cacheService,
		15*time.Minute, // Cache accounts for 15 minutes
	).(accountrepositories.AccountRepository)

	transactionRepository := transactionpersistence.NewGormTransactionRepository(db)

	// Initialize category repository with cache
	baseCategoryRepository := categorypersistence.NewGormCategoryRepository(db)
	categoryRepository := categorypersistence.NewCachedCategoryRepository(
		baseCategoryRepository,
		cacheService,
		15*time.Minute, // Cache categories for 15 minutes
	).(categoryrepositories.CategoryRepository)

	budgetRepository := budgetpersistence.NewGormBudgetRepository(db)

	// Initialize services
	jwtService := services.NewJWTServiceWithConfig(
		cfg.JWT.SecretKey,
		cfg.JWT.Expiration,
		cfg.JWT.Issuer,
	)

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
	// Note: With UnitOfWork, balance updates happen atomically within the transaction,
	// but we still keep event handlers for other subscribers (logging, notifications, etc.)
	eventBus.Subscribe("TransactionCreated", updateBalanceHandler.HandleTransactionCreated)
	eventBus.Subscribe("TransactionUpdated", updateBalanceHandler.HandleTransactionUpdated)
	eventBus.Subscribe("TransactionDeleted", updateBalanceHandler.HandleTransactionDeleted)

	// Initialize UnitOfWork for atomic operations
	unitOfWork := sharedpersistence.NewGormUnitOfWork(db)

	// Initialize transaction use cases
	// CreateTransactionUseCase, UpdateTransactionUseCase, and DeleteTransactionUseCase
	// now use UnitOfWork to ensure atomicity
	createTransactionUseCase := transactionusecases.NewCreateTransactionUseCase(unitOfWork, eventBus)
	listTransactionsUseCase := transactionusecases.NewListTransactionsUseCase(transactionRepository)
	getTransactionUseCase := transactionusecases.NewGetTransactionUseCase(transactionRepository)
	updateTransactionUseCase := transactionusecases.NewUpdateTransactionUseCase(unitOfWork, eventBus)
	deleteTransactionUseCase := transactionusecases.NewDeleteTransactionUseCase(unitOfWork, eventBus)
	restoreTransactionUseCase := transactionusecases.NewRestoreTransactionUseCase(transactionRepository)
	permanentDeleteTransactionUseCase := transactionusecases.NewPermanentDeleteTransactionUseCase(transactionRepository)

	// Initialize category use cases
	createCategoryUseCase := categoryusecases.NewCreateCategoryUseCase(categoryRepository, eventBus)
	listCategoriesUseCase := categoryusecases.NewListCategoriesUseCase(categoryRepository)
	getCategoryUseCase := categoryusecases.NewGetCategoryUseCase(categoryRepository)
	updateCategoryUseCase := categoryusecases.NewUpdateCategoryUseCase(categoryRepository, eventBus)
	deleteCategoryUseCase := categoryusecases.NewDeleteCategoryUseCase(categoryRepository)
	restoreCategoryUseCase := categoryusecases.NewRestoreCategoryUseCase(categoryRepository)
	permanentDeleteCategoryUseCase := categoryusecases.NewPermanentDeleteCategoryUseCase(categoryRepository)

	// Initialize budget use cases
	createBudgetUseCase := budgetusecases.NewCreateBudgetUseCase(budgetRepository, eventBus)
	listBudgetsUseCase := budgetusecases.NewListBudgetsUseCase(budgetRepository)
	getBudgetUseCase := budgetusecases.NewGetBudgetUseCase(budgetRepository)
	updateBudgetUseCase := budgetusecases.NewUpdateBudgetUseCase(budgetRepository, eventBus)
	deleteBudgetUseCase := budgetusecases.NewDeleteBudgetUseCase(budgetRepository, eventBus)
	getBudgetProgressUseCase := budgetusecases.NewGetBudgetProgressUseCase(budgetRepository, transactionRepository)

	// Initialize reporting use cases
	monthlyReportUseCase := reportingusecases.NewMonthlyReportUseCase(transactionRepository, reportCacheService)
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
		restoreTransactionUseCase,
		permanentDeleteTransactionUseCase,
	)
	categoryHandler := categoryhandlers.NewCategoryHandler(
		createCategoryUseCase,
		listCategoriesUseCase,
		getCategoryUseCase,
		updateCategoryUseCase,
		deleteCategoryUseCase,
		restoreCategoryUseCase,
		permanentDeleteCategoryUseCase,
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
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		// Increase header size limit to prevent 431 errors
		// Default is 4096 bytes, increasing to 16384 bytes (16KB)
		ReadBufferSize:  16384,
		WriteBufferSize: 16384,
		// Body size limit
		BodyLimit: int(cfg.Server.BodyLimit),
	})

	// Store environment in locals for error handler
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("env", cfg.Environment)
		return c.Next()
	})

	// Request ID middleware (must be early in the chain)
	app.Use(middleware.RequestIDMiddleware())

	// Rate limiting middleware (must be early, before other middlewares)
	if cacheService != nil {
		rateLimitConfig := middleware.DefaultRateLimitConfig()
		rateLimitConfig.CacheService = cacheService
		rateLimitConfig.MaxRequests = 100 // 100 requests per minute per IP
		rateLimitConfig.Window = 1 * time.Minute
		app.Use(middleware.RateLimitMiddleware(rateLimitConfig))
		log.Info().Msg("Rate limiting enabled")
	} else {
		log.Warn().Msg("Rate limiting disabled: Redis not available")
	}

	// Global middlewares
	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path} [${header:X-Request-ID}]\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "America/Sao_Paulo",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		AllowCredentials: true,
		MaxAge:           cfg.CORS.MaxAge,
	}))

	// Error handler middleware (must be after all other middlewares)
	app.Use(middleware.ErrorHandlerMiddleware())

	// Initialize health checker (with cache service if available)
	healthChecker := health.NewHealthCheckerWithCache(cacheService)

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
		accountroutes.SetupAccountRoutes(api, accountHandler, jwtService, userRepository, cacheService)

		// Setup transaction routes (protected)
		transactionroutes.SetupTransactionRoutes(api, transactionHandler, jwtService, userRepository, cacheService)

		// Setup category routes (protected)
		categoryroutes.SetupCategoryRoutes(api, categoryHandler, jwtService, userRepository, cacheService)

		// Setup budget routes (protected)
		budgetroutes.SetupBudgetRoutes(api, budgetHandler, jwtService, userRepository, cacheService)

		// Setup report routes (protected)
		reportroutes.SetupReportRoutes(api, reportHandler, jwtService, userRepository, cacheService)
	}

	// Get port from configuration
	port := cfg.Server.Port

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
