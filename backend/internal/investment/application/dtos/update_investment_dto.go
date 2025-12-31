package dtos

// UpdateInvestmentInput represents the input data for updating an investment.
type UpdateInvestmentInput struct {
	InvestmentID string   `json:"investment_id" validate:"required,uuid"`
	CurrentValue *float64 `json:"current_value,omitempty" validate:"omitempty,gt=0"`
	Quantity     *float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
}

// UpdateInvestmentOutput represents the output data after updating an investment.
type UpdateInvestmentOutput struct {
	InvestmentID     string   `json:"investment_id"`
	UserID           string   `json:"user_id"`
	AccountID        string   `json:"account_id"`
	Type             string   `json:"type"`
	Name             string   `json:"name"`
	Ticker           *string  `json:"ticker,omitempty"`
	PurchaseDate     string   `json:"purchase_date"`
	PurchaseAmount   float64  `json:"purchase_amount"`
	CurrentValue     float64  `json:"current_value"`
	Currency         string   `json:"currency"`
	Quantity         *float64 `json:"quantity,omitempty"`
	Context          string   `json:"context"`
	ReturnAbsolute   float64  `json:"return_absolute"`
	ReturnPercentage float64  `json:"return_percentage"`
	UpdatedAt        string   `json:"updated_at"`
}
