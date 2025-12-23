package valueobjects

import (
	"testing"
)

func TestNewTransactionID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			id:      "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "empty string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "invalid UUID format",
			id:      "invalid-uuid",
			wantErr: true,
		},
		{
			name:    "invalid UUID with spaces",
			id:      "550e8400-e29b-41d4-a716-446655440000 ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransactionID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransactionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.IsEmpty() {
					t.Error("NewTransactionID() returned empty TransactionID")
				}
				if got.Value() != tt.id {
					t.Errorf("NewTransactionID() value = %v, want %v", got.Value(), tt.id)
				}
			}
		})
	}
}

func TestGenerateTransactionID(t *testing.T) {
	tid1 := GenerateTransactionID()
	tid2 := GenerateTransactionID()

	if tid1.IsEmpty() {
		t.Error("GenerateTransactionID() returned empty TransactionID")
	}

	if tid2.IsEmpty() {
		t.Error("GenerateTransactionID() returned empty TransactionID")
	}

	if tid1.Equals(tid2) {
		t.Error("GenerateTransactionID() should generate unique IDs")
	}
}

func TestMustTransactionID(t *testing.T) {
	validID := "550e8400-e29b-41d4-a716-446655440000"
	tid := MustTransactionID(validID)

	if tid.Value() != validID {
		t.Errorf("MustTransactionID() value = %v, want %v", tid.Value(), validID)
	}

	// Test panic with invalid ID
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustTransactionID() should panic with invalid ID")
		}
	}()
	MustTransactionID("invalid-uuid")
}

func TestTransactionID_Value(t *testing.T) {
	id := "550e8400-e29b-41d4-a716-446655440000"
	tid, _ := NewTransactionID(id)

	if tid.Value() != id {
		t.Errorf("TransactionID.Value() = %v, want %v", tid.Value(), id)
	}
}

func TestTransactionID_String(t *testing.T) {
	id := "550e8400-e29b-41d4-a716-446655440000"
	tid, _ := NewTransactionID(id)

	if tid.String() != id {
		t.Errorf("TransactionID.String() = %v, want %v", tid.String(), id)
	}
}

func TestTransactionID_Equals(t *testing.T) {
	id := "550e8400-e29b-41d4-a716-446655440000"
	tid1, _ := NewTransactionID(id)
	tid2, _ := NewTransactionID(id)
	tid3 := GenerateTransactionID()

	if !tid1.Equals(tid2) {
		t.Error("TransactionID.Equals() = false, want true for same IDs")
	}

	if tid1.Equals(tid3) {
		t.Error("TransactionID.Equals() = true, want false for different IDs")
	}
}

func TestTransactionID_IsEmpty(t *testing.T) {
	tid, _ := NewTransactionID("550e8400-e29b-41d4-a716-446655440000")
	emptyTid := TransactionID{}

	if tid.IsEmpty() {
		t.Error("TransactionID.IsEmpty() = true, want false for non-empty ID")
	}

	if !emptyTid.IsEmpty() {
		t.Error("TransactionID.IsEmpty() = false, want true for empty ID")
	}
}
