package dtos

// ListTransactionsInput represents the input for listing transactions.
type ListTransactionsInput struct {
	UserID    string `json:"user_id" validate:"required,uuid"`
	AccountID string `json:"account_id,omitempty" validate:"omitempty,uuid"`
	Type      string `json:"type,omitempty" validate:"omitempty,oneof=INCOME EXPENSE"`
}

// TransactionOutput represents a single transaction in the list.
type TransactionOutput struct {
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

// ListTransactionsOutput represents the output for listing transactions.
type ListTransactionsOutput struct {
	Transactions []*TransactionOutput `json:"transactions"`
	Count        int                  `json:"count"`
}
