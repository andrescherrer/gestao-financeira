package valueobjects

import (
	"testing"
)

func TestNewAccountContext(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantValue string
		wantError bool
	}{
		{"valid PERSONAL", "PERSONAL", "PERSONAL", false},
		{"valid BUSINESS", "BUSINESS", "BUSINESS", false},
		{"lowercase personal", "personal", "PERSONAL", false},
		{"lowercase business", "business", "BUSINESS", false},
		{"mixed case", "PeRsOnAl", "PERSONAL", false},
		{"invalid value", "INVALID", "", true},
		{"empty value", "", "", true},
		{"whitespace", "  PERSONAL  ", "PERSONAL", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, err := NewAccountContext(tt.value)
			if (err != nil) != tt.wantError {
				t.Errorf("NewAccountContext() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && context.Value() != tt.wantValue {
				t.Errorf("NewAccountContext() value = %v, want %v", context.Value(), tt.wantValue)
			}
		})
	}
}

func TestAccountContext_Value(t *testing.T) {
	personal := PersonalContext()
	business := BusinessContext()

	if personal.Value() != "PERSONAL" {
		t.Errorf("PersonalContext().Value() = %v, want PERSONAL", personal.Value())
	}
	if business.Value() != "BUSINESS" {
		t.Errorf("BusinessContext().Value() = %v, want BUSINESS", business.Value())
	}
}

func TestAccountContext_DisplayName(t *testing.T) {
	tests := []struct {
		name    string
		context AccountContext
		want    string
	}{
		{"PERSONAL display", PersonalContext(), "Personal"},
		{"BUSINESS display", BusinessContext(), "Business"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.context.DisplayName(); got != tt.want {
				t.Errorf("AccountContext.DisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountContext_IsPersonal(t *testing.T) {
	personal := PersonalContext()
	business := BusinessContext()

	if !personal.IsPersonal() {
		t.Errorf("PersonalContext().IsPersonal() = false, want true")
	}
	if business.IsPersonal() {
		t.Errorf("BusinessContext().IsPersonal() = true, want false")
	}
}

func TestAccountContext_IsBusiness(t *testing.T) {
	personal := PersonalContext()
	business := BusinessContext()

	if !business.IsBusiness() {
		t.Errorf("BusinessContext().IsBusiness() = false, want true")
	}
	if personal.IsBusiness() {
		t.Errorf("PersonalContext().IsBusiness() = true, want false")
	}
}

func TestAccountContext_Equals(t *testing.T) {
	personal1 := PersonalContext()
	personal2 := PersonalContext()
	business := BusinessContext()

	tests := []struct {
		name string
		c1   AccountContext
		c2   AccountContext
		want bool
	}{
		{"equal PERSONAL", personal1, personal2, true},
		{"different contexts", personal1, business, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c1.Equals(tt.c2); got != tt.want {
				t.Errorf("AccountContext.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidAccountContext(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid PERSONAL", "PERSONAL", true},
		{"valid BUSINESS", "BUSINESS", true},
		{"lowercase", "personal", true},
		{"invalid", "INVALID", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAccountContext(tt.value); got != tt.want {
				t.Errorf("IsValidAccountContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAccountContext(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		wantValue string
		wantError bool
	}{
		{"valid PERSONAL", "PERSONAL", "PERSONAL", false},
		{"lowercase", "personal", "PERSONAL", false},
		{"empty", "", "", true},
		{"invalid", "INVALID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, err := ParseAccountContext(tt.s)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseAccountContext() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && context.Value() != tt.wantValue {
				t.Errorf("ParseAccountContext() value = %v, want %v", context.Value(), tt.wantValue)
			}
		})
	}
}

func TestAllAccountContexts(t *testing.T) {
	contexts := AllAccountContexts()

	if len(contexts) != 2 {
		t.Errorf("AllAccountContexts() length = %v, want 2", len(contexts))
	}

	hasPersonal := false
	hasBusiness := false

	for _, ctx := range contexts {
		if ctx.IsPersonal() {
			hasPersonal = true
		}
		if ctx.IsBusiness() {
			hasBusiness = true
		}
	}

	if !hasPersonal {
		t.Error("AllAccountContexts() missing PERSONAL")
	}
	if !hasBusiness {
		t.Error("AllAccountContexts() missing BUSINESS")
	}
}

func TestAccountContextHelpers(t *testing.T) {
	personal := PersonalContext()
	business := BusinessContext()

	if !personal.IsPersonal() {
		t.Error("PersonalContext() should return PERSONAL context")
	}
	if !business.IsBusiness() {
		t.Error("BusinessContext() should return BUSINESS context")
	}
}
