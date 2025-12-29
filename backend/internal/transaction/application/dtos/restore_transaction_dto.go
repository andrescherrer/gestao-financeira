package dtos

// RestoreTransactionInput represents the input for restoring a soft-deleted transaction.
type RestoreTransactionInput struct {
	TransactionID string `json:"transaction_id" validate:"required,uuid"`
}

// RestoreTransactionOutput represents the output after transaction restoration.
type RestoreTransactionOutput struct {
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
}
