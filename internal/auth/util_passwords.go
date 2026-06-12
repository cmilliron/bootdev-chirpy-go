package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %w\n", err)
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	results, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Error checking hash: %w\n", err)
	}

	return results, nil
}