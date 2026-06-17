package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	// Table-driven tests are the idiomatic way to write tests in Go
	tests := []struct {
		name          string
		headers       http.Header
		expectedToken string
		expectErr     bool
		errMsg        string
	}{
		{
			name: "Valid Authorization Header",
			headers: http.Header{
				"Authorization": []string{"Bearer some-valid-token-123"},
			},
			expectedToken: "Bearer some-valid-token-123",
			expectErr:     false,
		},
		{
			name:          "Missing Authorization Header",
			headers:       http.Header{}, // Empty headers
			expectedToken: "",
			expectErr:     true,
			errMsg:        "Athorization failed", // Matching the typo in your original code
		},
		{
			name: "Empty Authorization Header Value",
			headers: http.Header{
				"Authorization": []string{""},
			},
			expectedToken: "",
			expectErr:     true,
			errMsg:        "Athorization failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)

			// Check error expectations
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected an error, but got nil")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			// Check token expectations
			if token != tt.expectedToken {
				t.Errorf("expected token '%s', got '%s'", tt.expectedToken, token)
			}
		})
	}
}

func TestGetBearerTokenBootDev(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Authorization header",
			headers: http.Header{
				"Authorization": []string{"InvalidBearer token"},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}