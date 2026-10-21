package utils

import (
	"testing"
)

func TestHashPasswordAndCheck(t *testing.T) {
	password := "mySecretPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" || hash == password {
		t.Fatal("Hash should be non-empty and different from plain password")
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatal("Password check should succeed for correct password")
	}

	if CheckPasswordHash("wrongPassword", hash) {
		t.Fatal("Password check should fail for incorrect password")
	}
}
