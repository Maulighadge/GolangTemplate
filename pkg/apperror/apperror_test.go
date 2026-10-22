package apperror

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_HTTPStatusCode(t *testing.T) {
	tests := []struct {
		kind     ErrorKind
		expected int
	}{
		{KindNotFound, http.StatusNotFound},
		{KindBadRequest, http.StatusBadRequest},
		{KindUnauthorized, http.StatusUnauthorized},
		{KindForbidden, http.StatusForbidden},
		{KindConflict, http.StatusConflict},
		{KindValidation, http.StatusUnprocessableEntity},
		{KindInternal, http.StatusInternalServerError},
		{ErrorKind("UNKNOWN"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		err := New(tt.kind, "TEST", "test message")
		if code := err.HTTPStatusCode(); code != tt.expected {
			t.Errorf("For kind %s, expected status %d, got %d", tt.kind, tt.expected, code)
		}
	}
}

func TestAppError_Unwrap(t *testing.T) {
	rootErr := errors.New("database connection refused")
	appErr := Wrap(rootErr, KindInternal, "DB_ERROR", "database error")

	if !errors.Is(appErr, rootErr) {
		t.Errorf("Expected errors.Is(appErr, rootErr) to be true")
	}

	extracted, ok := AsAppError(appErr)
	if !ok || extracted == nil {
		t.Fatalf("Expected AsAppError to return true and non-nil")
	}
	// if extracted.Code != "DB_ERROR" {
	// 	t.Errorf("Expected code 'DB_ERROR', got '%s'", extracted.Code)
	// }
}

func TestAppError_WithDetailsAndErr(t *testing.T) {
	base := ErrUserNotFound
	detailed := base.WithDetails([]string{"id: 123"}).WithErr(errors.New("record not found in sql"))

	if detailed.Details == nil {
		t.Errorf("Expected Details to be set")
	}
	if detailed.Err == nil {
		t.Errorf("Expected Err to be set")
	}
	// Original should not be mutated
	if base.Details != nil {
		t.Errorf("Expected original base to not be mutated")
	}
}
