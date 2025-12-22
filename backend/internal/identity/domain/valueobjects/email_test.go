package valueobjects

import (
	"strings"
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantValue string
		wantError bool
	}{
		{"valid email", "user@example.com", "user@example.com", false},
		{"valid email with subdomain", "user@mail.example.com", "user@mail.example.com", false},
		{"valid email with plus", "user+tag@example.com", "user+tag@example.com", false},
		{"valid email with dots", "user.name@example.com", "user.name@example.com", false},
		{"uppercase email", "USER@EXAMPLE.COM", "user@example.com", false},
		{"email with whitespace", "  user@example.com  ", "user@example.com", false},
		{"empty email", "", "", true},
		{"missing @", "userexample.com", "", true},
		{"missing domain", "user@", "", true},
		{"missing local part", "@example.com", "", true},
		{"consecutive dots", "user..name@example.com", "", true},
		{"multiple @", "user@@example.com", "", true},
		{"invalid domain", "user@example", "", true},
		{"too long email", strings.Repeat("a", 250) + "@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.email)
			if (err != nil) != tt.wantError {
				t.Errorf("NewEmail() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && email.Value() != tt.wantValue {
				t.Errorf("NewEmail() value = %v, want %v", email.Value(), tt.wantValue)
			}
		})
	}
}

func TestEmail_Value(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	if email.Value() != "user@example.com" {
		t.Errorf("Email.Value() = %v, want user@example.com", email.Value())
	}
}

func TestEmail_String(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	if email.String() != "user@example.com" {
		t.Errorf("Email.String() = %v, want user@example.com", email.String())
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("user@example.com")
	email2, _ := NewEmail("user@example.com")
	email3, _ := NewEmail("other@example.com")

	tests := []struct {
		name string
		e1   Email
		e2   Email
		want bool
	}{
		{"equal emails", email1, email2, true},
		{"different emails", email1, email3, false},
		{"case insensitive", email1, MustEmail("USER@EXAMPLE.COM"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e1.Equals(tt.e2); got != tt.want {
				t.Errorf("Email.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail_Domain(t *testing.T) {
	tests := []struct {
		name  string
		email Email
		want  string
	}{
		{"simple domain", MustEmail("user@example.com"), "example.com"},
		{"subdomain", MustEmail("user@mail.example.com"), "mail.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.email.Domain(); got != tt.want {
				t.Errorf("Email.Domain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail_LocalPart(t *testing.T) {
	tests := []struct {
		name  string
		email Email
		want  string
	}{
		{"simple local", MustEmail("user@example.com"), "user"},
		{"local with dots", MustEmail("user.name@example.com"), "user.name"},
		{"local with plus", MustEmail("user+tag@example.com"), "user+tag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.email.LocalPart(); got != tt.want {
				t.Errorf("Email.LocalPart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail_IsEmpty(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	empty := Email{}

	if email.IsEmpty() {
		t.Error("Email.IsEmpty() = true for valid email, want false")
	}
	if !empty.IsEmpty() {
		t.Error("Email.IsEmpty() = false for empty email, want true")
	}
}

func TestMustEmail(t *testing.T) {
	// Valid email should not panic
	email := MustEmail("user@example.com")
	if email.Value() != "user@example.com" {
		t.Errorf("MustEmail() value = %v, want user@example.com", email.Value())
	}

	// Invalid email should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustEmail() should panic for invalid email")
		}
	}()
	MustEmail("invalid-email")
}
