package valueobjects

import (
	"strings"
	"testing"
)

func TestNewTransactionDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
		errorMsg    string
	}{
		{
			name:        "valid description",
			description: "Compra de supermercado",
			wantErr:     false,
		},
		{
			name:        "valid description with spaces",
			description: "  Compra de supermercado  ",
			wantErr:     false,
		},
		{
			name:        "valid minimum length",
			description: "ABC",
			wantErr:     false,
		},
		{
			name:        "valid maximum length",
			description: strings.Repeat("A", MaxDescriptionLength),
			wantErr:     false,
		},
		{
			name:        "empty string",
			description: "",
			wantErr:     true,
			errorMsg:    "cannot be empty",
		},
		{
			name:        "only spaces",
			description: "   ",
			wantErr:     true,
			errorMsg:    "cannot be empty",
		},
		{
			name:        "too short",
			description: "AB",
			wantErr:     true,
			errorMsg:    "at least 3 characters",
		},
		{
			name:        "too long",
			description: strings.Repeat("A", MaxDescriptionLength+1),
			wantErr:     true,
			errorMsg:    "at most 500 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransactionDescription(tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransactionDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("NewTransactionDescription() error = %v, want error containing %v", err, tt.errorMsg)
				}
			} else {
				trimmed := strings.TrimSpace(tt.description)
				if got.Value() != trimmed {
					t.Errorf("NewTransactionDescription() value = %v, want %v", got.Value(), trimmed)
				}
				if got.IsEmpty() {
					t.Error("NewTransactionDescription() returned empty TransactionDescription")
				}
			}
		})
	}
}

func TestMustTransactionDescription(t *testing.T) {
	validDesc := "Compra de supermercado"
	td := MustTransactionDescription(validDesc)

	if td.Value() != validDesc {
		t.Errorf("MustTransactionDescription() value = %v, want %v", td.Value(), validDesc)
	}

	// Test panic with invalid description
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustTransactionDescription() should panic with invalid description")
		}
	}()
	MustTransactionDescription("AB")
}

func TestTransactionDescription_Value(t *testing.T) {
	desc := "Compra de supermercado"
	td, _ := NewTransactionDescription(desc)

	if td.Value() != desc {
		t.Errorf("TransactionDescription.Value() = %v, want %v", td.Value(), desc)
	}
}

func TestTransactionDescription_String(t *testing.T) {
	desc := "Compra de supermercado"
	td, _ := NewTransactionDescription(desc)

	if td.String() != desc {
		t.Errorf("TransactionDescription.String() = %v, want %v", td.String(), desc)
	}
}

func TestTransactionDescription_Equals(t *testing.T) {
	desc := "Compra de supermercado"
	td1, _ := NewTransactionDescription(desc)
	td2, _ := NewTransactionDescription(desc)
	td3, _ := NewTransactionDescription("Outra descrição")

	if !td1.Equals(td2) {
		t.Error("TransactionDescription.Equals() = false, want true for same descriptions")
	}

	if td1.Equals(td3) {
		t.Error("TransactionDescription.Equals() = true, want false for different descriptions")
	}
}

func TestTransactionDescription_IsEmpty(t *testing.T) {
	td, _ := NewTransactionDescription("Compra de supermercado")
	emptyTd := TransactionDescription{}

	if td.IsEmpty() {
		t.Error("TransactionDescription.IsEmpty() = true, want false for non-empty description")
	}

	if !emptyTd.IsEmpty() {
		t.Error("TransactionDescription.IsEmpty() = false, want true for empty description")
	}
}
