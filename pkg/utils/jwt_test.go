package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/Maulighadge/GolangTemplate/config"
)

func TestGenerateAndValidateToken(t *testing.T) {
	config.AppConfig = &config.Config{
		AppName:            "TestApp",
		JWTSecret:          "testsecretkeyforjwttesting12345678",
		JWTExpirationHours: 2,
	}

	testUserID := uuid.New()
	email := "test@example.com"
	role := "user"

	token, err := GenerateToken(testUserID, email, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != testUserID {
		t.Errorf("Expected userID %v, got %v", testUserID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("Expected email %v, got %v", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("Expected role %v, got %v", role, claims.Role)
	}
}
