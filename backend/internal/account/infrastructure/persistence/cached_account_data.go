package persistence

import (
	"time"

	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// cachedAccountData is a serializable representation of Account for caching.
type cachedAccountData struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   int64     `json:"balance"` // Amount in cents
	Currency  string    `json:"currency"`
	Context   string    `json:"context"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// accountToCacheData converts an Account entity to cachedAccountData.
func accountToCacheData(account *entities.Account) *cachedAccountData {
	if account == nil {
		return nil
	}
	balance := account.Balance()
	return &cachedAccountData{
		ID:        account.ID().Value(),
		UserID:    account.UserID().Value(),
		Name:      account.Name().Value(),
		Type:      account.AccountType().Value(),
		Balance:   balance.Amount(),
		Currency:  balance.Currency().Code(),
		Context:   account.Context().Value(),
		IsActive:  account.IsActive(),
		CreatedAt: account.CreatedAt(),
		UpdatedAt: account.UpdatedAt(),
	}
}

// cacheDataToAccount converts cachedAccountData back to Account entity.
func cacheDataToAccount(data *cachedAccountData) (*entities.Account, error) {
	if data == nil {
		return nil, nil
	}

	accountID, err := valueobjects.NewAccountID(data.ID)
	if err != nil {
		return nil, err
	}

	userID, err := identityvalueobjects.NewUserID(data.UserID)
	if err != nil {
		return nil, err
	}

	accountName, err := valueobjects.NewAccountName(data.Name)
	if err != nil {
		return nil, err
	}

	accountType, err := valueobjects.NewAccountType(data.Type)
	if err != nil {
		return nil, err
	}

	currency, err := sharedvalueobjects.NewCurrency(data.Currency)
	if err != nil {
		return nil, err
	}

	balance, err := sharedvalueobjects.NewMoney(data.Balance, currency)
	if err != nil {
		return nil, err
	}

	context, err := sharedvalueobjects.NewAccountContext(data.Context)
	if err != nil {
		return nil, err
	}

	return entities.AccountFromPersistence(
		accountID,
		userID,
		accountName,
		accountType,
		balance,
		context,
		data.CreatedAt,
		data.UpdatedAt,
		data.IsActive,
	)
}

// accountsToCacheData converts a slice of Account entities to cachedAccountData.
func accountsToCacheData(accounts []*entities.Account) []*cachedAccountData {
	if accounts == nil {
		return nil
	}
	result := make([]*cachedAccountData, 0, len(accounts))
	for _, account := range accounts {
		result = append(result, accountToCacheData(account))
	}
	return result
}

// cacheDataToAccounts converts cachedAccountData slice back to Account entities.
func cacheDataToAccounts(data []*cachedAccountData) ([]*entities.Account, error) {
	if data == nil {
		return nil, nil
	}
	result := make([]*entities.Account, 0, len(data))
	for _, item := range data {
		account, err := cacheDataToAccount(item)
		if err != nil {
			return nil, err
		}
		if account != nil {
			result = append(result, account)
		}
	}
	return result, nil
}
