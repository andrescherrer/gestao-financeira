package valueobjects

import (
	"strings"
	"testing"
)

func TestNewGoalName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    string
	}{
		{
			name:    "valid name",
			input:   "Comprar um carro",
			wantErr: false,
			want:    "Comprar um carro",
		},
		{
			name:    "valid name with spaces",
			input:   "  Comprar um carro  ",
			wantErr: false,
			want:    "Comprar um carro",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "too short",
			input:   "AB",
			wantErr: true,
		},
		{
			name:    "minimum length",
			input:   "ABC",
			wantErr: false,
			want:    "ABC",
		},
		{
			name:    "maximum length",
			input:   strings.Repeat("A", 200),
			wantErr: false,
			want:    strings.Repeat("A", 200),
		},
		{
			name:    "too long",
			input:   strings.Repeat("A", 201),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGoalName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGoalName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name() != tt.want {
				t.Errorf("NewGoalName() = %v, want %v", got.Name(), tt.want)
			}
		})
	}
}

func TestGoalName_Equals(t *testing.T) {
	name1 := MustGoalName("Comprar um carro")
	name2 := MustGoalName("Comprar um carro")
	name3 := MustGoalName("Comprar uma casa")

	if !name1.Equals(name2) {
		t.Error("GoalName.Equals() should return true for same names")
	}

	if name1.Equals(name3) {
		t.Error("GoalName.Equals() should return false for different names")
	}
}

