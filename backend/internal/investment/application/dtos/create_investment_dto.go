package dtos

// CreateInvestmentInput represents the input data for investment creation.
type CreateInvestmentInput struct {
	UserID         string   `json:"user_id" validate:"required,uuid"`
	AccountID      string   `json:"account_id" validate:"required,uuid"`
	Type           string   `json:"type" validate:"required,oneof=STOCK FUND CDB TREASURY CRYPTO OTHER"`
	Name           string   `json:"name" validate:"required,min=2,max=200,no_sql_injection,no_xss,utf8"`
	Ticker         *string  `json:"ticker,omitempty" validate:"omitempty,max=20,no_sql_injection,no_xss"`
	PurchaseDate   string   `json:"purchase_date" validate:"required,datetime=2006-01-02"`
	PurchaseAmount float64  `json:"purchase_amount" validate:"required,gt=0"`
	Currency       string   `json:"currency" validate:"required,oneof=BRL USD EUR"`
	Quantity       *float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Context        string   `json:"context" validate:"required,oneof=PERSONAL BUSINESS"`
}

// CreateInvestmentOutput represents the output data after investment creation.
type CreateInvestmentOutput struct {
	InvestmentID   string   `json:"investment_id"`
	UserID         string   `json:"user_id"`
	AccountID      string   `json:"account_id"`
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	Ticker         *string  `json:"ticker,omitempty"`
	PurchaseDate   string   `json:"purchase_date"`
	PurchaseAmount float64  `json:"purchase_amount"`
	CurrentValue   float64  `json:"current_value"`
	Currency       string   `json:"currency"`
	Quantity       *float64 `json:"quantity,omitempty"`
	Context        string   `json:"context"`
	CreatedAt      string   `json:"created_at"`
}
