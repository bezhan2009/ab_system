package utils

import (
	"ab_system/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	cost := 12
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		logger.Error.Printf("Error hashing password: %v", err)

		return "", err
	}

	return string(hash), nil
}

func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
