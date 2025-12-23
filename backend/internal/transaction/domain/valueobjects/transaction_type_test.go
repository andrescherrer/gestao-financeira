package valueobjects

import (
	"testing"
)

func TestNewTransactionType(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
		wantVal string
	}{
		{
			name:    "valid INCOME",
			value:   "INCOME",
			wantErr: false,
			wantVal: "INCOME",
		},
		{
			name:    "valid EXPENSE",
			value:   "EXPENSE",
			wantErr: false,
			wantVal: "EXPENSE",
		},
		{
			name:    "valid income lowercase",
			value:   "income",
			wantErr: false,
			wantVal: "INCOME",
		},
		{
			name:    "valid expense lowercase",
			value:   "expense",
			wantErr: false,
			wantVal: "EXPENSE",
		},
		{
			name:    "valid with spaces",
			value:   "  INCOME  ",
			wantErr: false,
			wantVal: "INCOME",
		},
		{
			name:    "invalid type",
			value:   "INVALID",
			wantErr: true,
		},
		{
			name:    "empty string",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransactionType(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransactionType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Value() != tt.wantVal {
					t.Errorf("NewTransactionType() value = %v, want %v", got.Value(), tt.wantVal)
				}
			}
		})
	}
}

func TestMustTransactionType(t *testing.T) {
	validType := MustTransactionType("INCOME")
	if validType.Value() != "INCOME" {
		t.Errorf("MustTransactionType() value = %v, want INCOME", validType.Value())
	}

	// Test panic with invalid type
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustTransactionType() should panic with invalid type")
		}
	}()
	MustTransactionType("INVALID")
}

func TestIsValidTransactionType(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid INCOME", "INCOME", true},
		{"valid EXPENSE", "EXPENSE", true},
		{"valid income lowercase", "income", true},
		{"invalid type", "INVALID", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidTransactionType(tt.value)
			if got != tt.want {
				t.Errorf("IsValidTransactionType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionType_Value(t *testing.T) {
	tt := IncomeType()
	if tt.Value() != Income {
		t.Errorf("TransactionType.Value() = %v, want %v", tt.Value(), Income)
	}
}

func TestTransactionType_String(t *testing.T) {
	tt := ExpenseType()
	if tt.String() != Expense {
		t.Errorf("TransactionType.String() = %v, want %v", tt.String(), Expense)
	}
}

func TestTransactionType_DisplayName(t *testing.T) {
	incomeType := IncomeType()
	expenseType := ExpenseType()

	if incomeType.DisplayName() != "Income" {
		t.Errorf("TransactionType.DisplayName() = %v, want Income", incomeType.DisplayName())
	}

	if expenseType.DisplayName() != "Expense" {
		t.Errorf("TransactionType.DisplayName() = %v, want Expense", expenseType.DisplayName())
	}
}

func TestTransactionType_IsIncome(t *testing.T) {
	incomeType := IncomeType()
	expenseType := ExpenseType()

	if !incomeType.IsIncome() {
		t.Error("TransactionType.IsIncome() = false, want true for INCOME")
	}

	if expenseType.IsIncome() {
		t.Error("TransactionType.IsIncome() = true, want false for EXPENSE")
	}
}

func TestTransactionType_IsExpense(t *testing.T) {
	incomeType := IncomeType()
	expenseType := ExpenseType()

	if !expenseType.IsExpense() {
		t.Error("TransactionType.IsExpense() = false, want true for EXPENSE")
	}

	if incomeType.IsExpense() {
		t.Error("TransactionType.IsExpense() = true, want false for INCOME")
	}
}

func TestTransactionType_Equals(t *testing.T) {
	income1 := IncomeType()
	income2 := IncomeType()
	expense := ExpenseType()

	if !income1.Equals(income2) {
		t.Error("TransactionType.Equals() = false, want true for same types")
	}

	if income1.Equals(expense) {
		t.Error("TransactionType.Equals() = true, want false for different types")
	}
}

func TestIncomeType(t *testing.T) {
	tt := IncomeType()
	if tt.Value() != Income {
		t.Errorf("IncomeType() value = %v, want %v", tt.Value(), Income)
	}
}

func TestExpenseType(t *testing.T) {
	tt := ExpenseType()
	if tt.Value() != Expense {
		t.Errorf("ExpenseType() value = %v, want %v", tt.Value(), Expense)
	}
}

func TestParseTransactionType(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
		wantVal string
	}{
		{"valid INCOME", "INCOME", false, "INCOME"},
		{"valid EXPENSE", "EXPENSE", false, "EXPENSE"},
		{"empty string", "", true, ""},
		{"invalid type", "INVALID", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTransactionType(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTransactionType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Value() != tt.wantVal {
					t.Errorf("ParseTransactionType() value = %v, want %v", got.Value(), tt.wantVal)
				}
			}
		})
	}
}

func TestAllTransactionTypes(t *testing.T) {
	types := AllTransactionTypes()
	if len(types) != 2 {
		t.Errorf("AllTransactionTypes() returned %d types, want 2", len(types))
	}

	hasIncome := false
	hasExpense := false
	for _, tt := range types {
		if tt.IsIncome() {
			hasIncome = true
		}
		if tt.IsExpense() {
			hasExpense = true
		}
	}

	if !hasIncome {
		t.Error("AllTransactionTypes() should include INCOME")
	}

	if !hasExpense {
		t.Error("AllTransactionTypes() should include EXPENSE")
	}
}
