package valueobjects

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewNotificationID(t *testing.T) {
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
			got, err := NewNotificationID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNotificationID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.id {
				t.Errorf("NewNotificationID() = %v, want %v", got.Value(), tt.id)
			}
		})
	}
}

func TestGenerateNotificationID(t *testing.T) {
	id := GenerateNotificationID()
	if id.IsEmpty() {
		t.Error("GenerateNotificationID() should not return empty ID")
	}

	// Verify it's a valid UUID
	_, err := uuid.Parse(id.Value())
	if err != nil {
		t.Errorf("GenerateNotificationID() generated invalid UUID: %v", err)
	}
}

func TestNotificationID_Equals(t *testing.T) {
	id1 := GenerateNotificationID()
	id2 := GenerateNotificationID()
	id3 := MustNotificationID(id1.Value())

	if !id1.Equals(id3) {
		t.Error("NotificationID.Equals() should return true for same IDs")
	}

	if id1.Equals(id2) {
		t.Error("NotificationID.Equals() should return false for different IDs")
	}
}

func TestNotificationID_IsEmpty(t *testing.T) {
	id := GenerateNotificationID()
	if id.IsEmpty() {
		t.Error("NotificationID.IsEmpty() should return false for generated ID")
	}
}
