package dtos

// CreateTransactionInput represents the input data for transaction creation.
type CreateTransactionInput struct {
	UserID      string  `json:"user_id" validate:"required,uuid"`
	AccountID   string  `json:"account_id" validate:"required,uuid"`
	Type        string  `json:"type" validate:"required,oneof=INCOME EXPENSE"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Currency    string  `json:"currency" validate:"required,oneof=BRL USD EUR"`
	Description string  `json:"description" validate:"required,min=3,max=500"`
	Date        string  `json:"date" validate:"required"` // ISO 8601 format: YYYY-MM-DD
}

// CreateTransactionOutput represents the output data after transaction creation.
type CreateTransactionOutput struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	AccountID     string  `json:"account_id"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Description   string  `json:"description"`
	Date          string  `json:"date"`
	CreatedAt     string  `json:"created_at"`
}
