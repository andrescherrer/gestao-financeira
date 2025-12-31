package valueobjects

import (
	"testing"
)

func TestNewNotificationType(t *testing.T) {
	tests := []struct {
		name      string
		notifType string
		wantErr   bool
	}{
		{
			name:      "valid INFO type",
			notifType: "INFO",
			wantErr:   false,
		},
		{
			name:      "valid WARNING type",
			notifType: "WARNING",
			wantErr:   false,
		},
		{
			name:      "valid SUCCESS type",
			notifType: "SUCCESS",
			wantErr:   false,
		},
		{
			name:      "valid ERROR type",
			notifType: "ERROR",
			wantErr:   false,
		},
		{
			name:      "case insensitive",
			notifType: "info",
			wantErr:   false,
		},
		{
			name:      "empty string",
			notifType: "",
			wantErr:   true,
		},
		{
			name:      "invalid type",
			notifType: "INVALID",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNotificationType(tt.notifType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNotificationType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Should be normalized to uppercase
				if got.Value() != "INFO" && got.Value() != "WARNING" && got.Value() != "SUCCESS" && got.Value() != "ERROR" {
					t.Errorf("NewNotificationType() = %v, want normalized uppercase", got.Value())
				}
			}
		})
	}
}

func TestNotificationType_IsMethods(t *testing.T) {
	infoType := MustNotificationType(TypeInfo)
	warningType := MustNotificationType(TypeWarning)
	successType := MustNotificationType(TypeSuccess)
	errorType := MustNotificationType(TypeError)

	if !infoType.IsInfo() {
		t.Error("NotificationType.IsInfo() should return true for INFO type")
	}

	if !warningType.IsWarning() {
		t.Error("NotificationType.IsWarning() should return true for WARNING type")
	}

	if !successType.IsSuccess() {
		t.Error("NotificationType.IsSuccess() should return true for SUCCESS type")
	}

	if !errorType.IsError() {
		t.Error("NotificationType.IsError() should return true for ERROR type")
	}
}
