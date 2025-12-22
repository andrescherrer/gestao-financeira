package valueobjects

import (
	"testing"
)

func TestNewPasswordHashFromPlain(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError bool
	}{
		{"valid password", "password123", false},
		{"long password", "thisisalongpassword123456789", false},
		{"password with special chars", "P@ssw0rd!", false},
		{"empty password", "", true},
		{"short password", "short", true},
		{"too long password", string(make([]byte, 73)), true}, // 73 bytes exceeds bcrypt limit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewPasswordHashFromPlain(tt.password)
			if (err != nil) != tt.wantError {
				t.Errorf("NewPasswordHashFromPlain() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && hash.IsEmpty() {
				t.Error("NewPasswordHashFromPlain() returned empty hash for valid password")
			}
		})
	}
}

func TestNewPasswordHashFromHash(t *testing.T) {
	// First create a valid hash
	validHash, _ := NewPasswordHashFromPlain("password123")

	tests := []struct {
		name      string
		hash      string
		wantError bool
	}{
		{"valid hash", validHash.Value(), false},
		{"empty hash", "", true},
		{"invalid hash format", "invalid", true},
		{"short hash", "short", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewPasswordHashFromHash(tt.hash)
			if (err != nil) != tt.wantError {
				t.Errorf("NewPasswordHashFromHash() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && hash.IsEmpty() {
				t.Error("NewPasswordHashFromHash() returned empty hash for valid hash")
			}
		})
	}
}

func TestPasswordHash_Verify(t *testing.T) {
	password := "password123"
	hash, _ := NewPasswordHashFromPlain(password)

	tests := []struct {
		name     string
		hash     PasswordHash
		password string
		want     bool
	}{
		{"correct password", hash, password, true},
		{"incorrect password", hash, "wrongpassword", false},
		{"empty password", hash, "", false},
		{"different case", hash, "PASSWORD123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hash.Verify(tt.password); got != tt.want {
				t.Errorf("PasswordHash.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordHash_Value(t *testing.T) {
	hash, _ := NewPasswordHashFromPlain("password123")
	value := hash.Value()

	if value == "" {
		t.Error("PasswordHash.Value() returned empty string")
	}

	// bcrypt hashes are always 60 characters
	if len(value) != 60 {
		t.Errorf("PasswordHash.Value() length = %d, want 60", len(value))
	}
}

func TestPasswordHash_String(t *testing.T) {
	hash, _ := NewPasswordHashFromPlain("password123")
	str := hash.String()

	if str == "" {
		t.Error("PasswordHash.String() returned empty string")
	}

	if str != hash.Value() {
		t.Errorf("PasswordHash.String() = %v, want %v", str, hash.Value())
	}
}

func TestPasswordHash_Equals(t *testing.T) {
	hash1, _ := NewPasswordHashFromPlain("password123")
	hash2, _ := NewPasswordHashFromPlain("password123")
	hash3, _ := NewPasswordHashFromPlain("differentpassword")

	tests := []struct {
		name string
		h1   PasswordHash
		h2   PasswordHash
		want bool
	}{
		{"same password different hashes", hash1, hash2, false}, // bcrypt generates different hashes each time
		{"different passwords", hash1, hash3, false},
		{"same hash instance", hash1, hash1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h1.Equals(tt.h2); got != tt.want {
				t.Errorf("PasswordHash.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordHash_IsEmpty(t *testing.T) {
	hash, _ := NewPasswordHashFromPlain("password123")
	empty := PasswordHash{}

	if hash.IsEmpty() {
		t.Error("PasswordHash.IsEmpty() = true for valid hash, want false")
	}
	if !empty.IsEmpty() {
		t.Error("PasswordHash.IsEmpty() = false for empty hash, want true")
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError bool
	}{
		{"valid password", "password123", false},
		{"password with uppercase", "Password123", false},
		{"password with special chars", "P@ssw0rd!", false},
		{"empty password", "", true},
		{"short password", "short", true},
		{"too long password", string(make([]byte, 73)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidatePasswordStrength() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestPasswordHash_Uniqueness(t *testing.T) {
	// Test that bcrypt generates different hashes for the same password
	password := "password123"
	hash1, _ := NewPasswordHashFromPlain(password)
	hash2, _ := NewPasswordHashFromPlain(password)

	if hash1.Value() == hash2.Value() {
		t.Error("PasswordHash should generate different hashes for the same password (bcrypt salt)")
	}

	// But both should verify correctly
	if !hash1.Verify(password) {
		t.Error("PasswordHash.Verify() = false for correct password")
	}
	if !hash2.Verify(password) {
		t.Error("PasswordHash.Verify() = false for correct password")
	}
}

func TestMustPasswordHashFromHash(t *testing.T) {
	// Valid hash should not panic
	validHash, _ := NewPasswordHashFromPlain("password123")
	hash := MustPasswordHashFromHash(validHash.Value())
	if hash.IsEmpty() {
		t.Error("MustPasswordHashFromHash() returned empty hash for valid hash")
	}

	// Invalid hash should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustPasswordHashFromHash() should panic for invalid hash")
		}
	}()
	MustPasswordHashFromHash("invalid")
}
