package dtos

// GetInvestmentInput represents the input data for retrieving an investment.
type GetInvestmentInput struct {
	InvestmentID string `json:"investment_id" validate:"required,uuid"`
}

// GetInvestmentOutput represents the output data for an investment.
type GetInvestmentOutput struct {
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
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}
