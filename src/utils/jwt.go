package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims untuk JWT payload
type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(UUID, role string) (string, error) {
	return generateToken(UUID, role, getAccessSecret(), getAccessDuration())
}

func GenerateRefreshToken(UUID, role string) (string, error) {
	return generateToken(UUID, role, getRefreshSecret(), getRefreshDuration())
}

func ValidateAccessToken(tokenStr string) (*jwt.Token, *CustomClaims, error) {
	return validateToken(tokenStr, getAccessSecret())
}

func ValidateRefreshToken(tokenStr string) (*jwt.Token, *CustomClaims, error) {
	return validateToken(tokenStr, getRefreshSecret())
}

func generateToken(UUID, role string, secret []byte, duration time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: UUID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func validateToken(tokenStr string, secret []byte) (*jwt.Token, *CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, err
	}

	return token, claims, nil
}

func getAccessSecret() []byte {
	return []byte(GetEnvOrPanic("JWT_SECRET"))
}

func getRefreshSecret() []byte {
	return []byte(GetEnvOrPanic("JWT_REFRESH_SECRET"))
}

func getAccessDuration() time.Duration {
	return parseDuration(GetEnvOrDefault("JWT_EXPIRATION", "15m"))
}

func getRefreshDuration() time.Duration {
	return parseDuration(GetEnvOrDefault("JWT_REFRESH_EXPIRATION", "7d"))
}
func parseDuration(d string) time.Duration {
	if strings.HasSuffix(d, "d") {
		daysStr := strings.TrimSuffix(d, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			panic(fmt.Sprintf("invalid day format: %s", d))
		}
		return time.Hour * 24 * time.Duration(days)
	}

	duration, err := time.ParseDuration(d)
	if err != nil {
		panic(fmt.Sprintf("invalid duration format: %s", d))
	}
	return duration
}
