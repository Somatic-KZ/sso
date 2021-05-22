package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const cost = 10 // сложность хэша, чем больше это число, тем дольше его подсчет

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
