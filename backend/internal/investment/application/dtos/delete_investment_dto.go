package dtos

// DeleteInvestmentInput represents the input data for deleting an investment.
type DeleteInvestmentInput struct {
	InvestmentID string `json:"investment_id" validate:"required,uuid"`
}

// DeleteInvestmentOutput represents the output data after deleting an investment.
type DeleteInvestmentOutput struct {
	Success bool `json:"success"`
}
