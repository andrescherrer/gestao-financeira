package valueobjects

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewGoalID(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGoalID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGoalID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.id {
				t.Errorf("NewGoalID() = %v, want %v", got.Value(), tt.id)
			}
		})
	}
}

func TestGenerateGoalID(t *testing.T) {
	id := GenerateGoalID()
	if id.IsEmpty() {
		t.Error("GenerateGoalID() should not return empty ID")
	}

	// Verify it's a valid UUID
	_, err := uuid.Parse(id.Value())
	if err != nil {
		t.Errorf("GenerateGoalID() generated invalid UUID: %v", err)
	}
}

func TestGoalID_Equals(t *testing.T) {
	id1 := GenerateGoalID()
	id2 := GenerateGoalID()
	id3 := MustGoalID(id1.Value())

	if !id1.Equals(id3) {
		t.Error("GoalID.Equals() should return true for same IDs")
	}

	if id1.Equals(id2) {
		t.Error("GoalID.Equals() should return false for different IDs")
	}
}

func TestGoalID_IsEmpty(t *testing.T) {
	id := GenerateGoalID()
	if id.IsEmpty() {
		t.Error("GoalID.IsEmpty() should return false for generated ID")
	}
}
