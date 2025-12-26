package valueobjects

import (
	"testing"
)

func TestNewCategoryName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid name", "Alimentação", false},
		{"valid name with spaces", "Transporte Público", false},
		{"valid name with numbers", "Conta 123", false},
		{"empty name", "", true},
		{"too short", "A", true},
		{"too long", string(make([]byte, 101)), true},
		{"invalid characters", "Alimentação@", true},
		{"only spaces", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryName, err := NewCategoryName(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("NewCategoryName() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && categoryName.IsEmpty() {
				t.Error("NewCategoryName() returned empty name for valid input")
			}
		})
	}
}

func TestCategoryName_Value(t *testing.T) {
	name := "Alimentação"
	categoryName, err := NewCategoryName(name)
	if err != nil {
		t.Fatalf("NewCategoryName() error = %v", err)
	}

	if categoryName.Value() != name {
		t.Errorf("CategoryName.Value() = %v, want %v", categoryName.Value(), name)
	}
}

func TestCategoryName_String(t *testing.T) {
	name := "Transporte"
	categoryName, err := NewCategoryName(name)
	if err != nil {
		t.Fatalf("NewCategoryName() error = %v", err)
	}

	if categoryName.String() != name {
		t.Errorf("CategoryName.String() = %v, want %v", categoryName.String(), name)
	}
}

func TestCategoryName_Equals(t *testing.T) {
	name1, _ := NewCategoryName("Alimentação")
	name2, _ := NewCategoryName("alimentação") // lowercase
	name3, _ := NewCategoryName("Transporte")

	if !name1.Equals(name2) {
		t.Error("CategoryName.Equals() = false for case-insensitive equal names, want true")
	}

	if name1.Equals(name3) {
		t.Error("CategoryName.Equals() = true for different names, want false")
	}
}

func TestCategoryName_IsEmpty(t *testing.T) {
	categoryName, _ := NewCategoryName("Alimentação")
	if categoryName.IsEmpty() {
		t.Error("CategoryName.IsEmpty() = true for valid name, want false")
	}

	emptyName := CategoryName{}
	if !emptyName.IsEmpty() {
		t.Error("CategoryName.IsEmpty() = false for empty name, want true")
	}
}
