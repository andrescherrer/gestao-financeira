package dtos

// RestoreAccountInput represents the input for restoring a soft-deleted account.
type RestoreAccountInput struct {
	AccountID string `json:"account_id" validate:"required,uuid"`
}

// RestoreAccountOutput represents the output after account restoration.
type RestoreAccountOutput struct {
	Message   string `json:"message"`
	AccountID string `json:"account_id"`
}
