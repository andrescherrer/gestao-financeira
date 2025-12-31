package valueobjects

import (
	"testing"
)

func TestNewGoalStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{
			name:    "valid IN_PROGRESS",
			status:  StatusInProgress,
			wantErr: false,
		},
		{
			name:    "valid COMPLETED",
			status:  StatusCompleted,
			wantErr: false,
		},
		{
			name:    "valid OVERDUE",
			status:  StatusOverdue,
			wantErr: false,
		},
		{
			name:    "valid CANCELLED",
			status:  StatusCancelled,
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
			got, err := NewGoalStatus(tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGoalStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.status {
				t.Errorf("NewGoalStatus() = %v, want %v", got.Value(), tt.status)
			}
		})
	}
}

func TestGoalStatus_IsInProgress(t *testing.T) {
	status := MustGoalStatus(StatusInProgress)
	if !status.IsInProgress() {
		t.Error("GoalStatus.IsInProgress() should return true for IN_PROGRESS")
	}

	status2 := MustGoalStatus(StatusCompleted)
	if status2.IsInProgress() {
		t.Error("GoalStatus.IsInProgress() should return false for COMPLETED")
	}
}

func TestGoalStatus_IsCompleted(t *testing.T) {
	status := MustGoalStatus(StatusCompleted)
	if !status.IsCompleted() {
		t.Error("GoalStatus.IsCompleted() should return true for COMPLETED")
	}

	status2 := MustGoalStatus(StatusInProgress)
	if status2.IsCompleted() {
		t.Error("GoalStatus.IsCompleted() should return false for IN_PROGRESS")
	}
}

func TestGoalStatus_IsOverdue(t *testing.T) {
	status := MustGoalStatus(StatusOverdue)
	if !status.IsOverdue() {
		t.Error("GoalStatus.IsOverdue() should return true for OVERDUE")
	}

	status2 := MustGoalStatus(StatusInProgress)
	if status2.IsOverdue() {
		t.Error("GoalStatus.IsOverdue() should return false for IN_PROGRESS")
	}
}

func TestGoalStatus_IsCancelled(t *testing.T) {
	status := MustGoalStatus(StatusCancelled)
	if !status.IsCancelled() {
		t.Error("GoalStatus.IsCancelled() should return true for CANCELLED")
	}

	status2 := MustGoalStatus(StatusInProgress)
	if status2.IsCancelled() {
		t.Error("GoalStatus.IsCancelled() should return false for IN_PROGRESS")
	}
}

func TestGoalStatus_CanBeCancelled(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{
			name:   "IN_PROGRESS can be cancelled",
			status: StatusInProgress,
			want:   true,
		},
		{
			name:   "OVERDUE can be cancelled",
			status: StatusOverdue,
			want:   true,
		},
		{
			name:   "COMPLETED cannot be cancelled",
			status: StatusCompleted,
			want:   false,
		},
		{
			name:   "CANCELLED cannot be cancelled again",
			status: StatusCancelled,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := MustGoalStatus(tt.status)
			if got := status.CanBeCancelled(); got != tt.want {
				t.Errorf("GoalStatus.CanBeCancelled() = %v, want %v", got, tt.want)
			}
		})
	}
}
