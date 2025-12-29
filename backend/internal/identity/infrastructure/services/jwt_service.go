package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secretKey  []byte
	expiration time.Duration
	issuer     string
}

// Claims represents JWT claims.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service instance using environment variables.
func NewJWTService() *JWTService {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production" // Default for development
	}

	return NewJWTServiceWithConfig(secretKey, 24*time.Hour, "gestao-financeira-api")
}

// NewJWTServiceWithConfig creates a new JWT service instance with provided configuration.
func NewJWTServiceWithConfig(secretKey string, expiration time.Duration, issuer string) *JWTService {
	return &JWTService{
		secretKey:  []byte(secretKey),
		expiration: expiration,
		issuer:     issuer,
	}
}

// GenerateToken generates a new JWT token for a user.
func (s *JWTService) GenerateToken(userID, email string) (string, error) {
	// Token expiration time
	expirationTime := time.Now().Add(s.expiration)

	issuer := s.issuer
	if issuer == "" {
		issuer = "gestao-financeira-api"
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims.
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetTokenExpiration returns the token expiration time from claims.
func (s *JWTService) GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, errors.New("token has no expiration time")
	}

	return claims.ExpiresAt.Time, nil
}
