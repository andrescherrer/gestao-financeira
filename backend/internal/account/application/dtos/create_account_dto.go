package dtos

// CreateAccountInput represents the input data for account creation.
type CreateAccountInput struct {
	UserID         string  `json:"user_id" validate:"required,uuid"`
	Name           string  `json:"name" validate:"required,min=3,max=100"`
	Type           string  `json:"type" validate:"required,oneof=BANK WALLET INVESTMENT CREDIT_CARD"`
	InitialBalance float64 `json:"initial_balance" validate:"gte=0"`
	Currency       string  `json:"currency" validate:"required,oneof=BRL USD EUR"`
	Context        string  `json:"context" validate:"required,oneof=PERSONAL BUSINESS"`
}

// CreateAccountOutput represents the output data after account creation.
type CreateAccountOutput struct {
	AccountID string  `json:"account_id"`
	UserID    string  `json:"user_id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	Context   string  `json:"context"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
}
