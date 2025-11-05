package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost = bcrypt.DefaultCost
)

// HashPassword generates bcrypt hash from password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares password with hash
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
