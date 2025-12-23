package dtos

// GetAccountInput represents the input for getting a single account.
type GetAccountInput struct {
	AccountID string `json:"account_id" validate:"required,uuid"`
}

// GetAccountOutput represents the output for getting a single account.
// Uses the same structure as AccountOutput from list_accounts_dto.go
type GetAccountOutput struct {
	AccountID string  `json:"account_id"`
	UserID    string  `json:"user_id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	Context   string  `json:"context"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
