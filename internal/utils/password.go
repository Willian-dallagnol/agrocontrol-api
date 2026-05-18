package utils

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12 // custo aumentado: DefaultCost=10 é fraco para produção

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
