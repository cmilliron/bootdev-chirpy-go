package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

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

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Auhorization failed")
	}

	token, prefixExists := strings.CutPrefix(authHeader, "Bearer ")
	if token == "" || prefixExists == false{
		return "", fmt.Errorf("Malformed token")
	}

	return token, nil
}

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Auhorization failed")
	}

	token, prefixExists := strings.CutPrefix(authHeader, "ApiKey ")
	if token == "" || prefixExists == false{
		return "", fmt.Errorf("Malformed token")
	}

	return token, nil
}

func MakeRefreshToken() string {
	randomKey := make([]byte, 32)
	rand.Read(randomKey)
	refreshToken := hex.EncodeToString(randomKey)
	return refreshToken
}