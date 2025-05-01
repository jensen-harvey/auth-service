package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/herb-immortal/auth_service_hi/pkg/models"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID string         `json:"sub"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// TokenManager handles JWT token generation and validation
type TokenManager struct {
	secretKey []byte
	issuer    string
	tokenTTL  time.Duration
}

// NewTokenManager creates a new token manager
func NewTokenManager(secretKey string, issuer string, tokenTTL time.Duration) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secretKey),
		issuer:    issuer,
		tokenTTL:  tokenTTL,
	}
}

// GenerateToken creates a new JWT token for a user
func (tm *TokenManager) GenerateToken(userID string, role models.UserRole) (string, time.Time, error) {
	expirationTime := time.Now().Add(tm.tokenTTL)
	
	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    tm.issuer,
			Subject:   userID,
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tm.secretKey)
	
	if err != nil {
		return "", time.Time{}, err
	}
	
	return tokenString, expirationTime, nil
}

// ValidateToken checks if a token is valid and returns its claims
func (tm *TokenManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secretKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid token")
}

// GenerateUUID creates a new UUID with role prefix
func GenerateUUID(role models.UserRole) string {
	// In a production environment, use a proper UUID library like google/uuid
	// This is a simplified example
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d", role, timestamp)
}