package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeString(t *testing.T) {
	config := DefaultSecurityConfig()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "string with null bytes",
			input:    "Hello\x00World",
			expected: "HelloWorld",
		},
		{
			name:     "string with control characters",
			input:    "Hello\x01World",
			expected: "HelloWorld",
		},
		{
			name:     "string with HTML",
			input:    "<script>alert('xss')</script>",
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "string with whitespace",
			input:    "  Hello World  ",
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input, config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateNoSQLInjection(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "normal string",
			input:   "SELECT * FROM users",
			wantErr: true, // Contains SQL keywords
		},
		{
			name:    "safe string",
			input:   "Hello World",
			wantErr: false,
		},
		{
			name:    "SQL injection attempt",
			input:   "'; DROP TABLE users; --",
			wantErr: true,
		},
		{
			name:    "OR injection",
			input:   "admin' OR '1'='1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoSQLInjection(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateNoXSS(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "normal string",
			input:   "Hello World",
			wantErr: false,
		},
		{
			name:    "XSS script tag",
			input:   "<script>alert('xss')</script>",
			wantErr: true,
		},
		{
			name:    "XSS javascript",
			input:   "javascript:alert('xss')",
			wantErr: true,
		},
		{
			name:    "XSS event handler",
			input:   "<img onerror='alert(1)'>",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoXSS(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateNoPathTraversal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "normal string",
			input:   "filename.txt",
			wantErr: false,
		},
		{
			name:    "path traversal",
			input:   "../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "windows path traversal",
			input:   "..\\..\\windows\\system32",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoPathTraversal(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name    string
		password string
		wantErr bool
	}{
		{
			name:     "strong password",
			password: "SecurePass123!",
			wantErr:  false,
		},
		{
			name:     "too short",
			password: "Short1!",
			wantErr:  true,
		},
		{
			name:     "only lowercase",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "only uppercase",
			password: "PASSWORD123",
			wantErr:  true,
		},
		{
			name:     "no special chars",
			password: "Password123",
			wantErr:  true,
		},
		{
			name:     "valid password with 3 types",
			password: "Password123!",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUTF8(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid UTF-8",
			input:   "Hello World",
			wantErr: false,
		},
		{
			name:    "valid UTF-8 with special chars",
			input:   "Olá Mundo! 你好",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUTF8(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

