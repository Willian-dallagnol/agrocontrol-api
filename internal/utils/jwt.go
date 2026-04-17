package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 🔐 Estrutura que representa os dados dentro do token (payload)
type Claims struct {
	UserID uint `json:"user_id"`
	// 👉 ID do usuário

	Email string `json:"email"`
	// 👉 email do usuário

	Role string `json:"role"`
	// 👉 role (admin, manager, operator)

	jwt.RegisteredClaims
	// 👉 campos padrão do JWT (expiração, emissão, etc)
}

// 🎫 Gera um novo token JWT
func GenerateToken(userID uint, email, role, secret string) (string, error) {

	// 📦 monta os dados do token
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			// 👉 token expira em 24h

			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 👉 data de criação do token
		},
	}

	// 🔐 cria o token com algoritmo HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 🔑 assina o token com a chave secreta
	return token.SignedString([]byte(secret))
}

// 🔍 Valida um token JWT recebido
func ValidateToken(tokenString, secret string) (*Claims, error) {

	// 🔐 faz parse e valida assinatura do token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		// ❌ token inválido, expirado ou mal formado
		return nil, err
	}

	// 📦 extrai os dados (claims)
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		// ❌ token inválido
		return nil, err
	}

	// ✅ retorna os dados do token
	return claims, nil
}
