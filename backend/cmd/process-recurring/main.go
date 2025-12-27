package main

import (
	"os"

	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	transactionservices "gestao-financeira/backend/internal/transaction/application/services"
	transactionpersistence "gestao-financeira/backend/internal/transaction/infrastructure/persistence"
	"gestao-financeira/backend/pkg/database"
	"gestao-financeira/backend/pkg/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize structured logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.InitLogger(logLevel)

	log.Info().Msg("Starting Recurring Transactions Processor")

	// Initialize database connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	log.Info().Msg("Database connection established")

	// Initialize Event Bus
	eventBus := eventbus.NewEventBus()

	// Initialize transaction repository
	transactionRepository := transactionpersistence.NewGormTransactionRepository(db)

	// Initialize recurring transaction processor
	processor := transactionservices.NewRecurringTransactionProcessor(transactionRepository, eventBus)

	// Process recurring transactions
	log.Info().Msg("Processing recurring transactions...")
	createdCount, err := processor.ProcessRecurringTransactions()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to process recurring transactions")
	}

	log.Info().Int("created_count", createdCount).Msg("Recurring transactions processed successfully")

	// Close database connection
	if err := database.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing database")
	}

	log.Info().Msg("Recurring Transactions Processor completed")
}
