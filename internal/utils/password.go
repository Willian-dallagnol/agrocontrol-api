package utils

import "golang.org/x/crypto/bcrypt"

// 🔐 Gera o hash da senha (criptografia segura)
func HashPassword(password string) (string, error) {

	// 👉 transforma a senha em bytes e gera o hash usando bcrypt
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// 👉 retorna o hash como string
	return string(bytes), err
}

// 🔍 Compara senha digitada com o hash armazenado
func CheckPasswordHash(password, hash string) bool {

	// 👉 compara senha "plain text" com o hash salvo
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	// 👉 retorna true se a senha estiver correta
	return err == nil
}
