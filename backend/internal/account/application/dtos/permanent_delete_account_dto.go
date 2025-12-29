package dtos

// PermanentDeleteAccountInput represents the input for permanently deleting an account.
type PermanentDeleteAccountInput struct {
	AccountID string `json:"account_id" validate:"required,uuid"`
}

// PermanentDeleteAccountOutput represents the output after permanent account deletion.
type PermanentDeleteAccountOutput struct {
	Message   string `json:"message"`
	AccountID string `json:"account_id"`
}
