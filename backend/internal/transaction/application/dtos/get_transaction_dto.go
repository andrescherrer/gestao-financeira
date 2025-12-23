package dtos

// GetTransactionInput represents the input for getting a single transaction.
type GetTransactionInput struct {
	TransactionID string `json:"transaction_id" validate:"required,uuid"`
}

// GetTransactionOutput represents the output for getting a single transaction.
// Uses the same structure as TransactionOutput from list_transactions_dto.go
type GetTransactionOutput struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	AccountID     string  `json:"account_id"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Description   string  `json:"description"`
	Date          string  `json:"date"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
