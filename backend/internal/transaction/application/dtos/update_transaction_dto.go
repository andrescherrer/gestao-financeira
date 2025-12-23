package dtos

// UpdateTransactionInput represents the input data for transaction update.
// All fields are optional - only provided fields will be updated.
type UpdateTransactionInput struct {
	TransactionID string   `json:"transaction_id" validate:"required,uuid"`
	Type          *string  `json:"type,omitempty" validate:"omitempty,oneof=INCOME EXPENSE"`
	Amount        *float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Currency      *string  `json:"currency,omitempty" validate:"omitempty,oneof=BRL USD EUR"`
	Description   *string  `json:"description,omitempty" validate:"omitempty,min=3,max=500"`
	Date          *string  `json:"date,omitempty" validate:"omitempty"` // ISO 8601 format: YYYY-MM-DD
}

// UpdateTransactionOutput represents the output data after transaction update.
type UpdateTransactionOutput struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	AccountID     string  `json:"account_id"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Description   string  `json:"description"`
	Date          string  `json:"date"`
	UpdatedAt     string  `json:"updated_at"`
}
