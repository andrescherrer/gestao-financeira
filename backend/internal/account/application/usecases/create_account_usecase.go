package usecases

import (
	"fmt"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/repositories"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// CreateAccountUseCase handles account creation.
type CreateAccountUseCase struct {
	accountRepository repositories.AccountRepository
	eventBus          *eventbus.EventBus
}

// NewCreateAccountUseCase creates a new CreateAccountUseCase instance.
func NewCreateAccountUseCase(
	accountRepository repositories.AccountRepository,
	eventBus *eventbus.EventBus,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		accountRepository: accountRepository,
		eventBus:          eventBus,
	}
}

// Execute performs the account creation.
// It validates the input, creates value objects, creates a new account entity,
// saves it to the repository, and publishes domain events.
func (uc *CreateAccountUseCase) Execute(input dtos.CreateAccountInput) (*dtos.CreateAccountOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create account name value object
	accountName, err := valueobjects.NewAccountName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid account name: %w", err)
	}

	// Create account type value object
	accountType, err := valueobjects.NewAccountType(input.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid account type: %w", err)
	}

	// Create currency value object
	currency, err := sharedvalueobjects.NewCurrency(input.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	// Create initial balance (convert float to cents)
	initialBalanceCents := int64(input.InitialBalance * 100)
	initialBalance, err := sharedvalueobjects.NewMoney(initialBalanceCents, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid initial balance: %w", err)
	}

	// Create account context value object
	context, err := sharedvalueobjects.NewAccountContext(input.Context)
	if err != nil {
		return nil, fmt.Errorf("invalid account context: %w", err)
	}

	// Create account entity
	account, err := entities.NewAccount(userID, accountName, accountType, initialBalance, context)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Save account to repository
	if err := uc.accountRepository.Save(account); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	// Publish domain events
	domainEvents := account.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			// Log error but don't fail the account creation
			// In production, you might want to handle this differently
			// (e.g., store events in an outbox pattern)
			_ = err // Ignore for now, but should be logged
		}
	}
	account.ClearEvents()

	// Build output
	balance := account.Balance()
	output := &dtos.CreateAccountOutput{
		AccountID: account.ID().Value(),
		UserID:    account.UserID().Value(),
		Name:      account.Name().Value(),
		Type:      account.AccountType().Value(),
		Balance:   balance.Float64(),
		Currency:  balance.Currency().Code(),
		Context:   account.Context().Value(),
		IsActive:  account.IsActive(),
		CreatedAt: account.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
