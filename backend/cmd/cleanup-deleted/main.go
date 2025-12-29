package main

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	accountpersistence "gestao-financeira/backend/internal/account/infrastructure/persistence"
	budgetpersistence "gestao-financeira/backend/internal/budget/infrastructure/persistence"
	categorypersistence "gestao-financeira/backend/internal/category/infrastructure/persistence"
	transactionpersistence "gestao-financeira/backend/internal/transaction/infrastructure/persistence"
	"gestao-financeira/backend/pkg/database"
)

func main() {
	// Parse command line flags
	daysOld := flag.Int("days", 90, "Number of days old to consider for cleanup (default: 90)")
	dryRun := flag.Bool("dry-run", false, "Run without actually deleting (default: false)")
	resource := flag.String("resource", "all", "Resource type to clean up: transactions, accounts, categories, budgets, all (default: all)")
	flag.Parse()

	// Setup logger
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	log.Info().
		Int("days", *daysOld).
		Bool("dry_run", *dryRun).
		Str("resource", *resource).
		Msg("Starting cleanup of soft-deleted records")

	// Get database connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	helper := database.NewSoftDeleteHelper(db)

	// Cleanup based on resource type
	var totalDeleted int64

	switch *resource {
	case "transactions", "all":
		if *resource == "all" || *resource == "transactions" {
			count, err := cleanupTransactions(db, helper, *daysOld, *dryRun)
			if err != nil {
				log.Error().Err(err).Msg("Failed to cleanup transactions")
			} else {
				totalDeleted += count
				log.Info().Int64("count", count).Msg("Cleaned up transactions")
			}
		}
		if *resource != "all" {
			break
		}
		fallthrough
	case "accounts":
		if *resource == "all" || *resource == "accounts" {
			count, err := cleanupAccounts(db, helper, *daysOld, *dryRun)
			if err != nil {
				log.Error().Err(err).Msg("Failed to cleanup accounts")
			} else {
				totalDeleted += count
				log.Info().Int64("count", count).Msg("Cleaned up accounts")
			}
		}
		if *resource != "all" {
			break
		}
		fallthrough
	case "categories":
		if *resource == "all" || *resource == "categories" {
			count, err := cleanupCategories(db, helper, *daysOld, *dryRun)
			if err != nil {
				log.Error().Err(err).Msg("Failed to cleanup categories")
			} else {
				totalDeleted += count
				log.Info().Int64("count", count).Msg("Cleaned up categories")
			}
		}
		if *resource != "all" {
			break
		}
		fallthrough
	case "budgets":
		if *resource == "all" || *resource == "budgets" {
			count, err := cleanupBudgets(db, helper, *daysOld, *dryRun)
			if err != nil {
				log.Error().Err(err).Msg("Failed to cleanup budgets")
			} else {
				totalDeleted += count
				log.Info().Int64("count", count).Msg("Cleaned up budgets")
			}
		}
	default:
		log.Fatal().Str("resource", *resource).Msg("Invalid resource type")
	}

	if *dryRun {
		log.Info().Int64("total", totalDeleted).Msg("DRY RUN: Would delete records")
	} else {
		log.Info().Int64("total", totalDeleted).Msg("Cleanup completed")
	}
}

func cleanupTransactions(db *gorm.DB, helper *database.SoftDeleteHelper, daysOld int, dryRun bool) (int64, error) {
	model := &transactionpersistence.TransactionModel{}
	if dryRun {
		// Count only
		var count int64
		cutoffDate := time.Now().AddDate(0, 0, -daysOld)
		err := db.Unscoped().
			Model(model).
			Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoffDate).
			Count(&count).Error
		return count, err
	}
	return helper.CleanupDeleted(model, daysOld)
}

func cleanupAccounts(db *gorm.DB, helper *database.SoftDeleteHelper, daysOld int, dryRun bool) (int64, error) {
	model := &accountpersistence.AccountModel{}
	if dryRun {
		// Count only
		var count int64
		cutoffDate := time.Now().AddDate(0, 0, -daysOld)
		err := db.Unscoped().
			Model(model).
			Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoffDate).
			Count(&count).Error
		return count, err
	}
	return helper.CleanupDeleted(model, daysOld)
}

func cleanupCategories(db *gorm.DB, helper *database.SoftDeleteHelper, daysOld int, dryRun bool) (int64, error) {
	model := &categorypersistence.CategoryModel{}
	if dryRun {
		// Count only
		var count int64
		cutoffDate := time.Now().AddDate(0, 0, -daysOld)
		err := db.Unscoped().
			Model(model).
			Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoffDate).
			Count(&count).Error
		return count, err
	}
	return helper.CleanupDeleted(model, daysOld)
}

func cleanupBudgets(db *gorm.DB, helper *database.SoftDeleteHelper, daysOld int, dryRun bool) (int64, error) {
	model := &budgetpersistence.BudgetModel{}
	if dryRun {
		// Count only
		var count int64
		cutoffDate := time.Now().AddDate(0, 0, -daysOld)
		err := db.Unscoped().
			Model(model).
			Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoffDate).
			Count(&count).Error
		return count, err
	}
	return helper.CleanupDeleted(model, daysOld)
}
