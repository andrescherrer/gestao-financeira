package services

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTService(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		wantErr  bool
	}{
		{
			name:     "with JWT_SECRET env variable",
			envValue: "test-secret-key",
			wantErr:  false,
		},
		{
			name:     "without JWT_SECRET env variable (uses default)",
			envValue: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("JWT_SECRET", tt.envValue)
				defer os.Unsetenv("JWT_SECRET")
			} else {
				os.Unsetenv("JWT_SECRET")
			}

			service := NewJWTService()
			assert.NotNil(t, service)
			assert.NotNil(t, service.secretKey)
		})
	}
}

func TestJWTService_GenerateToken(t *testing.T) {
	service := NewJWTService()

	tests := []struct {
		name    string
		userID  string
		email   string
		wantErr bool
	}{
		{
			name:    "valid user ID and email",
			userID:  "test-user-id",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "empty user ID",
			userID:  "",
			email:   "test@example.com",
			wantErr: false, // JWT allows empty strings
		},
		{
			name:    "empty email",
			userID:  "test-user-id",
			email:   "",
			wantErr: false, // JWT allows empty strings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenerateToken(tt.userID, tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	service := NewJWTService()

	// Generate a valid token
	userID := "test-user-id"
	email := "test@example.com"
	validToken, err := service.GenerateToken(userID, email)
	require.NoError(t, err)
	require.NotEmpty(t, validToken)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "invalid token format",
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "token with wrong secret",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIn0.invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				if claims != nil {
					assert.Equal(t, userID, claims.UserID)
					assert.Equal(t, email, claims.Email)
					assert.Equal(t, "gestao-financeira-api", claims.Issuer)
					assert.Equal(t, userID, claims.Subject)
				}
			}
		})
	}
}

func TestJWTService_GetTokenExpiration(t *testing.T) {
	service := NewJWTService()

	// Generate a valid token
	userID := "test-user-id"
	email := "test@example.com"
	validToken, err := service.GenerateToken(userID, email)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid-token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expiration, err := service.GetTokenExpiration(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, expiration.IsZero())
			} else {
				assert.NoError(t, err)
				assert.False(t, expiration.IsZero())
				// Token should expire in approximately 24 hours
				expectedExpiration := time.Now().Add(24 * time.Hour)
				diff := expectedExpiration.Sub(expiration)
				// Allow 5 minutes difference for test execution time
				assert.True(t, diff < 5*time.Minute && diff > -5*time.Minute)
			}
		})
	}
}

func TestJWTService_TokenRoundTrip(t *testing.T) {
	service := NewJWTService()

	userID := "test-user-id-123"
	email := "testuser@example.com"

	// Generate token
	token, err := service.GenerateToken(userID, email)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Validate token
	claims, err := service.ValidateToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify claims
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, "gestao-financeira-api", claims.Issuer)
	assert.Equal(t, userID, claims.Subject)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)

	// Verify expiration
	expiration, err := service.GetTokenExpiration(token)
	assert.NoError(t, err)
	assert.False(t, expiration.IsZero())
	assert.True(t, expiration.After(time.Now()))
}

func TestJWTService_InvalidSigningMethod(t *testing.T) {
	// This test verifies that tokens signed with different methods are rejected
	service := NewJWTService()

	// Create a token with a different signing method (RS256 instead of HS256)
	// This is a simplified test - in practice, you'd need to create an actual RS256 token
	// For now, we'll just verify that malformed tokens are rejected
	invalidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIn0.invalid"

	claims, err := service.ValidateToken(invalidToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
