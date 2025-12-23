package dtos

// DeleteTransactionInput represents the input for deleting a transaction.
type DeleteTransactionInput struct {
	TransactionID string `json:"transaction_id" validate:"required,uuid"`
}

// DeleteTransactionOutput represents the output after transaction deletion.
type DeleteTransactionOutput struct {
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
}
