package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey            []byte
	ErrMissingAuthHeader = errors.New("authorization header missing or invalid")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrTokenParseError   = errors.New("failed to parse token")
)

// Claims struct
type Claims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

// Load secret key from environment variable
func LoadJWTSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set in environment variables")
	}
	secretKey = []byte(secret)
}

// Generate JWT token
func GenerateToken(userId uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Parse and validate token
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, ErrTokenParseError
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// Extract Bearer token from Authorization header
func ExtractToken(c *gin.Context) (string, error) {
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) < 7 || strings.ToUpper(bearerToken[0:7]) != "BEARER " {
		return "", ErrMissingAuthHeader
	}
	return bearerToken[7:], nil
}
