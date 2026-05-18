package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims define o payload do token JWT
type Claims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsRefresh bool   `json:"is_refresh,omitempty"`
	jwt.RegisteredClaims
}

var ErrTokenInvalid = errors.New("token inválido ou expirado")

// GenerateToken cria um access token JWT (curta duração)
func GenerateToken(userID uint, email, role, secret string, expHours int) (string, error) {
	if expHours <= 0 {
		expHours = 1
	}
	claims := Claims{
		UserID: userID, Email: email, Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "agrocontrol-api",
			Subject:   email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken cria um refresh token JWT (longa duração)
func GenerateRefreshToken(userID uint, email, role, secret string, expHours int) (string, error) {
	if expHours <= 0 {
		expHours = 168 // 7 dias
	}
	claims := Claims{
		UserID: userID, Email: email, Role: role, IsRefresh: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "agrocontrol-api",
			Subject:   email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken verifica assinatura, expiração e retorna as claims
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("algoritmo de assinatura não permitido")
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, ErrTokenInvalid
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}
