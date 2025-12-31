package valueobjects

import (
	"testing"
)

func TestNewNotificationStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{
			name:    "valid UNREAD status",
			status:  "UNREAD",
			wantErr: false,
		},
		{
			name:    "valid READ status",
			status:  "READ",
			wantErr: false,
		},
		{
			name:    "valid ARCHIVED status",
			status:  "ARCHIVED",
			wantErr: false,
		},
		{
			name:    "case insensitive",
			status:  "unread",
			wantErr: false,
		},
		{
			name:    "empty string",
			status:  "",
			wantErr: true,
		},
		{
			name:    "invalid status",
			status:  "INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNotificationStatus(tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNotificationStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Should be normalized to uppercase
				if got.Value() != "UNREAD" && got.Value() != "READ" && got.Value() != "ARCHIVED" {
					t.Errorf("NewNotificationStatus() = %v, want normalized uppercase", got.Value())
				}
			}
		})
	}
}

func TestNotificationStatus_IsMethods(t *testing.T) {
	unreadStatus := MustNotificationStatus(StatusUnread)
	readStatus := MustNotificationStatus(StatusRead)
	archivedStatus := MustNotificationStatus(StatusArchived)

	if !unreadStatus.IsUnread() {
		t.Error("NotificationStatus.IsUnread() should return true for UNREAD status")
	}

	if !readStatus.IsRead() {
		t.Error("NotificationStatus.IsRead() should return true for READ status")
	}

	if !archivedStatus.IsArchived() {
		t.Error("NotificationStatus.IsArchived() should return true for ARCHIVED status")
	}
}
