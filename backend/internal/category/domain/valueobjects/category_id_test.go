package valueobjects

import (
	"testing"
)

func TestNewCategoryID(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	invalidUUID := "invalid-uuid"
	emptyUUID := ""

	tests := []struct {
		name      string
		id        string
		wantError bool
	}{
		{"valid UUID", validUUID, false},
		{"invalid UUID", invalidUUID, true},
		{"empty UUID", emptyUUID, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryID, err := NewCategoryID(tt.id)
			if (err != nil) != tt.wantError {
				t.Errorf("NewCategoryID() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && categoryID.IsEmpty() {
				t.Error("NewCategoryID() returned empty ID for valid UUID")
			}
		})
	}
}

func TestGenerateCategoryID(t *testing.T) {
	categoryID := GenerateCategoryID()

	if categoryID.IsEmpty() {
		t.Error("GenerateCategoryID() returned empty ID")
	}

	// Generate another one and verify they're different
	categoryID2 := GenerateCategoryID()
	if categoryID.Equals(categoryID2) {
		t.Error("GenerateCategoryID() should generate unique IDs")
	}
}

func TestCategoryID_Value(t *testing.T) {
	categoryID := GenerateCategoryID()
	value := categoryID.Value()

	if value == "" {
		t.Error("CategoryID.Value() returned empty string")
	}
}

func TestCategoryID_String(t *testing.T) {
	categoryID := GenerateCategoryID()
	str := categoryID.String()

	if str == "" {
		t.Error("CategoryID.String() returned empty string")
	}

	if str != categoryID.Value() {
		t.Errorf("CategoryID.String() = %v, want %v", str, categoryID.Value())
	}
}

func TestCategoryID_Equals(t *testing.T) {
	categoryID1 := GenerateCategoryID()
	categoryID2 := GenerateCategoryID()
	categoryID3 := MustCategoryID(categoryID1.Value())

	if !categoryID1.Equals(categoryID3) {
		t.Error("CategoryID.Equals() = false for equal IDs, want true")
	}

	if categoryID1.Equals(categoryID2) {
		t.Error("CategoryID.Equals() = true for different IDs, want false")
	}
}

func TestCategoryID_IsEmpty(t *testing.T) {
	categoryID := GenerateCategoryID()
	if categoryID.IsEmpty() {
		t.Error("CategoryID.IsEmpty() = true for generated ID, want false")
	}

	emptyID := CategoryID{}
	if !emptyID.IsEmpty() {
		t.Error("CategoryID.IsEmpty() = false for empty ID, want true")
	}
}
