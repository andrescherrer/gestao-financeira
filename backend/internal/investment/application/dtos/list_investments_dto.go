package dtos

import (
	"gestao-financeira/backend/pkg/pagination"
)

// ListInvestmentsInput represents the input data for listing investments.
type ListInvestmentsInput struct {
	UserID  string `json:"user_id" validate:"required,uuid"`
	Context string `json:"context,omitempty" validate:"omitempty,oneof=PERSONAL BUSINESS"`
	Type    string `json:"type,omitempty" validate:"omitempty,oneof=STOCK FUND CDB TREASURY CRYPTO OTHER"`
	Page    string `json:"page,omitempty" validate:"omitempty,numeric"`
	Limit   string `json:"limit,omitempty" validate:"omitempty,numeric"`
}

// ListInvestmentsOutput represents the output data for listing investments.
type ListInvestmentsOutput struct {
	Investments []InvestmentOutput           `json:"investments"`
	Count       int                          `json:"count"`
	Pagination  *pagination.PaginationResult `json:"pagination,omitempty"`
}

// InvestmentOutput represents an investment in the output.
type InvestmentOutput struct {
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
