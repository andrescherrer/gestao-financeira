package usecases

import (
	"fmt"
	"time"

	accountrepositories "gestao-financeira/backend/internal/account/domain/repositories"
	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/investment/application/dtos"
	"gestao-financeira/backend/internal/investment/domain/entities"
	"gestao-financeira/backend/internal/investment/domain/repositories"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// CreateInvestmentUseCase handles investment creation.
type CreateInvestmentUseCase struct {
	investmentRepository repositories.InvestmentRepository
	accountRepository    accountrepositories.AccountRepository
	eventBus             *eventbus.EventBus
}

// NewCreateInvestmentUseCase creates a new CreateInvestmentUseCase instance.
func NewCreateInvestmentUseCase(
	investmentRepository repositories.InvestmentRepository,
	accountRepository accountrepositories.AccountRepository,
	eventBus *eventbus.EventBus,
) *CreateInvestmentUseCase {
	return &CreateInvestmentUseCase{
		investmentRepository: investmentRepository,
		accountRepository:    accountRepository,
		eventBus:             eventBus,
	}
}

// Execute performs the investment creation.
// It validates the input, creates value objects, creates a new investment entity,
// saves it to the repository, and publishes domain events.
func (uc *CreateInvestmentUseCase) Execute(input dtos.CreateInvestmentInput) (*dtos.CreateInvestmentOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create account ID value object
	accountID, err := accountvalueobjects.NewAccountID(input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	// Verify account exists and belongs to user
	account, err := uc.accountRepository.FindByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}
	if !account.UserID().Equals(userID) {
		return nil, fmt.Errorf("account does not belong to user")
	}

	// Create investment type value object
	investmentType, err := investmentvalueobjects.NewInvestmentType(input.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid investment type: %w", err)
	}

	// Create investment name value object
	name, err := investmentvalueobjects.NewInvestmentName(input.Name, input.Ticker)
	if err != nil {
		return nil, fmt.Errorf("invalid investment name: %w", err)
	}

	// Parse purchase date
	purchaseDate, err := time.Parse("2006-01-02", input.PurchaseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase date format: %w", err)
	}

	// Create currency value object
	currency, err := sharedvalueobjects.NewCurrency(input.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	// Ensure currency matches account currency
	if !currency.Equals(account.Balance().Currency()) {
		return nil, fmt.Errorf("investment currency must match account currency")
	}

	// Create purchase amount (convert float to cents)
	purchaseAmountCents := int64(input.PurchaseAmount * 100)
	purchaseAmount, err := sharedvalueobjects.NewMoney(purchaseAmountCents, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase amount: %w", err)
	}

	// Create account context value object
	context, err := sharedvalueobjects.NewAccountContext(input.Context)
	if err != nil {
		return nil, fmt.Errorf("invalid account context: %w", err)
	}

	// Ensure context matches account context
	if !context.Equals(account.Context()) {
		return nil, fmt.Errorf("investment context must match account context")
	}

	// Create investment entity
	investment, err := entities.NewInvestment(
		userID,
		accountID,
		investmentType,
		name,
		purchaseDate,
		purchaseAmount,
		input.Quantity,
		context,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create investment: %w", err)
	}

	// Save investment to repository
	if err := uc.investmentRepository.Save(investment); err != nil {
		return nil, fmt.Errorf("failed to save investment: %w", err)
	}

	// Publish domain events
	domainEvents := investment.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			// Log error but don't fail the investment creation
			_ = err // Ignore for now, but should be logged
		}
	}
	investment.ClearEvents()

	// Build output
	currentValue := investment.CurrentValue()
	output := &dtos.CreateInvestmentOutput{
		InvestmentID:   investment.ID().Value(),
		UserID:         investment.UserID().Value(),
		AccountID:      investment.AccountID().Value(),
		Type:           investment.InvestmentType().Value(),
		Name:           investment.Name().Name(),
		Ticker:         investment.Name().Ticker(),
		PurchaseDate:   investment.PurchaseDate().Format("2006-01-02"),
		PurchaseAmount: investment.PurchaseAmount().Float64(),
		CurrentValue:   currentValue.Float64(),
		Currency:       currentValue.Currency().Code(),
		Quantity:       investment.Quantity(),
		Context:        investment.Context().Value(),
		CreatedAt:      investment.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
