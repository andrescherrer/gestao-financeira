package dtos

// PermanentDeleteTransactionInput represents the input for permanently deleting a transaction.
type PermanentDeleteTransactionInput struct {
	TransactionID string `json:"transaction_id" validate:"required,uuid"`
}

// PermanentDeleteTransactionOutput represents the output after permanent transaction deletion.
type PermanentDeleteTransactionOutput struct {
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
}
