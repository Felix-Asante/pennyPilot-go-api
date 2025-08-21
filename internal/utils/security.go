package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomTokens(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashString(str string) (string, error) {
	hashedString, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedString), nil
}

func CompareHashedString(hashedString string, plainString string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(plainString))
}
