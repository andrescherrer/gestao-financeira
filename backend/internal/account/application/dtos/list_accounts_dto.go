package dtos

// ListAccountsInput represents the input for listing accounts.
type ListAccountsInput struct {
	UserID  string `json:"user_id" validate:"required,uuid"`
	Context string `json:"context,omitempty" validate:"omitempty,oneof=PERSONAL BUSINESS"`
}

// AccountOutput represents a single account in the list.
type AccountOutput struct {
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

// ListAccountsOutput represents the output for listing accounts.
type ListAccountsOutput struct {
	Accounts []*AccountOutput `json:"accounts"`
	Count    int              `json:"count"`
}
