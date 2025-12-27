package valueobjects

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBudgetID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			id:      uuid.New().String(),
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
			name:    "valid UUID string",
			id:      "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBudgetID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, got.IsEmpty())
			} else {
				assert.NoError(t, err)
				assert.False(t, got.IsEmpty())
				assert.Equal(t, tt.id, got.Value())
			}
		})
	}
}

func TestGenerateBudgetID(t *testing.T) {
	id1 := GenerateBudgetID()
	id2 := GenerateBudgetID()

	assert.False(t, id1.IsEmpty())
	assert.False(t, id2.IsEmpty())
	assert.NotEqual(t, id1.Value(), id2.Value())
}

func TestBudgetID_Equals(t *testing.T) {
	id1 := GenerateBudgetID()
	id2 := GenerateBudgetID()
	id3, _ := NewBudgetID(id1.Value())

	assert.True(t, id1.Equals(id3))
	assert.False(t, id1.Equals(id2))
}

func TestBudgetID_IsEmpty(t *testing.T) {
	emptyID := BudgetID{}
	validID := GenerateBudgetID()

	assert.True(t, emptyID.IsEmpty())
	assert.False(t, validID.IsEmpty())
}

func TestMustBudgetID(t *testing.T) {
	validUUID := uuid.New().String()
	id := MustBudgetID(validUUID)

	assert.Equal(t, validUUID, id.Value())
	assert.False(t, id.IsEmpty())

	// Test panic with invalid UUID
	assert.Panics(t, func() {
		MustBudgetID("invalid")
	})
}
