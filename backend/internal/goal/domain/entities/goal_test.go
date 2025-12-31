package entities

import (
	"testing"
	"time"

	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

func TestNewGoal(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	name := goalvalueobjects.MustGoalName("Comprar um carro")
	targetAmount, _ := sharedvalueobjects.NewMoneyFromFloat(50000.0, sharedvalueobjects.MustCurrency("BRL"))
	deadline := time.Now().AddDate(1, 0, 0) // 1 year from now
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	tests := []struct {
		name         string
		userID       identityvalueobjects.UserID
		goalName     goalvalueobjects.GoalName
		targetAmount sharedvalueobjects.Money
		deadline     time.Time
		context      sharedvalueobjects.AccountContext
		wantErr      bool
	}{
		{
			name:         "valid goal",
			userID:       userID,
			goalName:     name,
			targetAmount: targetAmount,
			deadline:     deadline,
			context:      context,
			wantErr:      false,
		},
		{
			name:         "empty user ID",
			userID:       identityvalueobjects.UserID{},
			goalName:     name,
			targetAmount: targetAmount,
			deadline:     deadline,
			context:      context,
			wantErr:      true,
		},
		{
			name:         "empty goal name",
			userID:       userID,
			goalName:     goalvalueobjects.GoalName{},
			targetAmount: targetAmount,
			deadline:     deadline,
			context:      context,
			wantErr:      true,
		},
		{
			name:         "zero target amount",
			userID:       userID,
			goalName:     name,
			targetAmount: sharedvalueobjects.Zero(targetAmount.Currency()),
			deadline:     deadline,
			context:      context,
			wantErr:      true,
		},
		{
			name:         "past deadline",
			userID:       userID,
			goalName:     name,
			targetAmount: targetAmount,
			deadline:     time.Now().AddDate(-1, 0, 0),
			context:      context,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGoal(
				tt.userID,
				tt.goalName,
				tt.targetAmount,
				tt.deadline,
				tt.context,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGoal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewGoal() returned nil for valid input")
			}
			if !tt.wantErr && got != nil {
				if got.ID().IsEmpty() {
					t.Error("NewGoal() should generate a non-empty ID")
				}
				if !got.Status().IsInProgress() {
					t.Error("NewGoal() should create goal with IN_PROGRESS status")
				}
				if !got.CurrentAmount().IsZero() {
					t.Error("NewGoal() should create goal with zero current amount")
				}
			}
		})
	}
}

func TestGoal_AddContribution(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	name := goalvalueobjects.MustGoalName("Comprar um carro")
	targetAmount, _ := sharedvalueobjects.NewMoneyFromFloat(50000.0, sharedvalueobjects.MustCurrency("BRL"))
	deadline := time.Now().AddDate(1, 0, 0)
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	goal, _ := NewGoal(userID, name, targetAmount, deadline, context)

	tests := []struct {
		name    string
		amount  sharedvalueobjects.Money
		wantErr bool
	}{
		{
			name:    "valid contribution",
			amount:  targetAmount.Multiply(0.1), // 10% of target
			wantErr: false,
		},
		{
			name:    "zero contribution",
			amount:  sharedvalueobjects.Zero(targetAmount.Currency()),
			wantErr: true,
		},
		{
			name:    "negative contribution",
			amount:  targetAmount.Negate(),
			wantErr: true,
		},
		{
			name:    "contribution exceeding target",
			amount:  targetAmount.Multiply(1.1), // 110% of target
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset goal for each test
			goal, _ = NewGoal(userID, name, targetAmount, deadline, context)

			err := goal.AddContribution(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Goal.AddContribution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoal_CalculateProgress(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	name := goalvalueobjects.MustGoalName("Comprar um carro")
	targetAmount, _ := sharedvalueobjects.NewMoneyFromFloat(50000.0, sharedvalueobjects.MustCurrency("BRL"))
	deadline := time.Now().AddDate(1, 0, 0)
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	goal, _ := NewGoal(userID, name, targetAmount, deadline, context)

	// Initially should be 0%
	progress := goal.CalculateProgress()
	if progress != 0.0 {
		t.Errorf("Goal.CalculateProgress() = %v, want 0.0", progress)
	}

	// Add 50% contribution
	contribution, _ := sharedvalueobjects.NewMoneyFromFloat(25000.0, targetAmount.Currency())
	goal.AddContribution(contribution)
	progress = goal.CalculateProgress()
	if progress != 50.0 {
		t.Errorf("Goal.CalculateProgress() = %v, want 50.0", progress)
	}
}

func TestGoal_IsCompleted(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	name := goalvalueobjects.MustGoalName("Comprar um carro")
	targetAmount, _ := sharedvalueobjects.NewMoneyFromFloat(50000.0, sharedvalueobjects.MustCurrency("BRL"))
	deadline := time.Now().AddDate(1, 0, 0)
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	goal, _ := NewGoal(userID, name, targetAmount, deadline, context)

	if goal.IsCompleted() {
		t.Error("Goal.IsCompleted() should return false for new goal")
	}

	// Add full contribution
	goal.AddContribution(targetAmount)
	goal.CheckStatus()

	if !goal.IsCompleted() {
		t.Error("Goal.IsCompleted() should return true when current amount equals target")
	}
	if !goal.Status().IsCompleted() {
		t.Error("Goal.Status() should be COMPLETED when goal is completed")
	}
}

func TestGoal_IsOverdue(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	name := goalvalueobjects.MustGoalName("Comprar um carro")
	targetAmount, _ := sharedvalueobjects.NewMoneyFromFloat(50000.0, sharedvalueobjects.MustCurrency("BRL"))
	pastDeadline := time.Now().AddDate(-1, 0, 0) // 1 year ago
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	// Create goal with past deadline (should fail)
	_, err := NewGoal(userID, name, targetAmount, pastDeadline, context)
	if err == nil {
		t.Error("NewGoal() should fail with past deadline")
	}

	// Create goal with past deadline for testing
	goal, _ := GoalFromPersistence(
		goalvalueobjects.GenerateGoalID(),
		userID,
		name,
		targetAmount,
		sharedvalueobjects.Zero(targetAmount.Currency()),
		pastDeadline,
		context,
		goalvalueobjects.MustGoalStatus(goalvalueobjects.StatusInProgress),
		time.Now(),
		time.Now(),
	)

	if !goal.IsOverdue() {
		t.Error("Goal.IsOverdue() should return true for goal with past deadline")
	}
}

func TestGoal_CalculateRemainingDays(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	name := goalvalueobjects.MustGoalName("Comprar um carro")
	targetAmount, _ := sharedvalueobjects.NewMoneyFromFloat(50000.0, sharedvalueobjects.MustCurrency("BRL"))
	deadline := time.Now().AddDate(0, 0, 30) // 30 days from now
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	goal, _ := NewGoal(userID, name, targetAmount, deadline, context)

	remainingDays := goal.CalculateRemainingDays()
	if remainingDays < 29 || remainingDays > 31 {
		t.Errorf("Goal.CalculateRemainingDays() = %v, want approximately 30", remainingDays)
	}
}

func TestGoal_Cancel(t *testing.T) {
	userID := identityvalueobjects.MustUserID("123e4567-e89b-12d3-a456-426614174000")
	name := goalvalueobjects.MustGoalName("Comprar um carro")
	targetAmount, _ := sharedvalueobjects.NewMoneyFromFloat(50000.0, sharedvalueobjects.MustCurrency("BRL"))
	deadline := time.Now().AddDate(1, 0, 0)
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	goal, _ := NewGoal(userID, name, targetAmount, deadline, context)

	err := goal.Cancel()
	if err != nil {
		t.Errorf("Goal.Cancel() error = %v, want nil", err)
	}

	if !goal.Status().IsCancelled() {
		t.Error("Goal.Status() should be CANCELLED after cancel")
	}

	// Try to cancel again (should fail)
	err = goal.Cancel()
	if err == nil {
		t.Error("Goal.Cancel() should fail when goal is already cancelled")
	}
}
