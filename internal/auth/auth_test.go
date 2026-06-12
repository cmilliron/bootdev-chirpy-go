package auth

import (
	"testing"
)

func TestPasswordHashingPipeline(t *testing.T) {
	password := "mySuperSecretPassword123!"

	// 1. Test successful hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("Expected a non-empty hash string")
	}

	// 2. Test successful verification (Happy Path)
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed unexpectedly: %v", err)
	}
	if !match {
		t.Error("Expected password to match its generated hash, but it failed")
	}

	// 3. Test verification failure with wrong password
	wrongPassword := "notMySuperSecretPassword123!"
	match, err = CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed unexpectedly on wrong password: %v", err)
	}
	if match {
		t.Error("Expected verification to fail for an incorrect password, but it passed")
	}
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	// Test how the function handles a completely malformed hash string
	invalidHash := "not-a-valid-argon2-hash"
	password := "somePassword"

	match, err := CheckPasswordHash(password, invalidHash)
	
	// Argon2id should throw an error if the hash format doesn't match its expected structure
	if err == nil {
		t.Error("Expected an error when validating against a malformed hash, got nil")
	}
	if match {
		t.Error("Expected match to be false for an invalid hash")
	}
}