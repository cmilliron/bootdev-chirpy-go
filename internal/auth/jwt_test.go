package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT_Pipeline(t *testing.T) {
	userID := uuid.New()
	secret := "my-ultra-secure-fallback-secret-key-123"
	expiry := 1 * time.Hour

	// 1. Test creation
	token, err := MakeJWT(userID, secret, expiry)
	if err != nil {
		t.Fatalf("MakeJWT failed unexpectedly: %v", err)
	}
	if token == "" {
		t.Fatal("Expected a signed token string, got empty string")
	}

	// 2. Test valid parsing & validation
	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed to validate a good token: %v", err)
	}

	if parsedID != userID {
		t.Errorf("Expected parsed UUID to be %s, but got %s", userID, parsedID)
	}
}

func TestValidateJWT_EdgeCases(t *testing.T) {
	userID := uuid.New()
	correctSecret := "correct-secret-key"
	wrongSecret := "wrong-secret-key"

	tests := []struct {
		name          string
		tokenSetup    func() string
		secretToUse   string
		expectedError bool
	}{
		{
			name: "Expired token (beyond leeway)",
			tokenSetup: func() string {
				// Create a token that expired 10 seconds ago (outside the 5s leeway)
				tkn, _ := MakeJWT(userID, correctSecret, -10*time.Second)
				return tkn
			},
			secretToUse:   correctSecret,
			expectedError: true,
		},
		{
			name: "Expired token within leeway",
			tokenSetup: func() string {
				// Expired 2 seconds ago, should be saved by the 5s leeway
				tkn, _ := MakeJWT(userID, correctSecret, -2*time.Second)
				return tkn
			},
			secretToUse:   correctSecret,
			expectedError: false,
		},
		{
			name: "Wrong secret key",
			tokenSetup: func() string {
				tkn, _ := MakeJWT(userID, correctSecret, 1*time.Hour)
				return tkn
			},
			secretToUse:   wrongSecret,
			expectedError: true,
		},
		{
			name: "Malformed token string",
			tokenSetup: func() string {
				return "this.isnot.jwt"
			},
			secretToUse:   correctSecret,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.tokenSetup()
			_, err := ValidateJWT(token, tt.secretToUse)

			if tt.expectedError && err == nil {
				t.Error("Expected an error but validation passed")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected validation to pass, but got error: %v", err)
			}
		})
	}
}