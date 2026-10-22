package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorKind classifies the category of an error
type ErrorKind string

const (
	KindNotFound       ErrorKind = "NOT_FOUND"
	KindBadRequest     ErrorKind = "BAD_REQUEST"
	KindUnauthorized   ErrorKind = "UNAUTHORIZED"
	KindForbidden      ErrorKind = "FORBIDDEN"
	KindConflict       ErrorKind = "CONFLICT"
	KindValidation     ErrorKind = "VALIDATION_ERROR"
	KindInternal       ErrorKind = "INTERNAL_ERROR"
)

// AppError is a structured application error carrying semantic kind, error code, client message, and root cause
type AppError struct {
	Kind    ErrorKind   `json:"-"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Err     error       `json:"-"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements the standard error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap enables standard library errors.Is and errors.As unwrapping
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatusCode maps ErrorKind to appropriate HTTP status code
func (e *AppError) HTTPStatusCode() int {
	switch e.Kind {
	case KindNotFound:
		return http.StatusNotFound
	case KindBadRequest:
		return http.StatusBadRequest
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindConflict:
		return http.StatusConflict
	case KindValidation:
		return http.StatusUnprocessableEntity
	case KindInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// WithDetails returns a copy of AppError with attached structured details
func (e *AppError) WithDetails(details interface{}) *AppError {
	clone := *e
	clone.Details = details
	return &clone
}

// WithErr returns a copy of AppError with the underlying root cause error attached
func (e *AppError) WithErr(err error) *AppError {
	clone := *e
	clone.Err = err
	return &clone
}

// New creates a new AppError
func New(kind ErrorKind, code, message string) *AppError {
	return &AppError{
		Kind:    kind,
		Code:    code,
		Message: message,
	}
}

// Wrap creates a new AppError that wraps an underlying error
func Wrap(err error, kind ErrorKind, code, message string) *AppError {
	return &AppError{
		Kind:    kind,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewNotFound creates a 404 Not Found error
func NewNotFound(code, message string) *AppError {
	return New(KindNotFound, code, message)
}

// NewBadRequest creates a 400 Bad Request error
func NewBadRequest(code, message string) *AppError {
	return New(KindBadRequest, code, message)
}

// NewUnauthorized creates a 401 Unauthorized error
func NewUnauthorized(code, message string) *AppError {
	return New(KindUnauthorized, code, message)
}

// NewForbidden creates a 403 Forbidden error
func NewForbidden(code, message string) *AppError {
	return New(KindForbidden, code, message)
}

// NewConflict creates a 409 Conflict error
func NewConflict(code, message string) *AppError {
	return New(KindConflict, code, message)
}

// NewValidationError creates a 422 Validation error
func NewValidationError(code, message string, details interface{}) *AppError {
	return &AppError{
		Kind:    KindValidation,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewInternal creates a 500 Internal Server Error, wrapping the root cause
func NewInternal(err error, message string) *AppError {
	if message == "" {
		message = "An unexpected internal server error occurred"
	}
	return &AppError{
		Kind:    KindInternal,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: message,
		Err:     err,
	}
}

// Predefined common domain errors
var (
	// Auth & User errors
	ErrUserNotFound       = NewNotFound("USER_NOT_FOUND", "user not found")
	ErrEmailAlreadyExists = NewConflict("EMAIL_ALREADY_EXISTS", "user with this email already exists")
	ErrInvalidCredentials = NewUnauthorized("INVALID_CREDENTIALS", "invalid email or password")
	ErrAccountDeactivated = NewForbidden("ACCOUNT_DEACTIVATED", "user account is deactivated")
	ErrIncorrectPassword  = NewBadRequest("INCORRECT_PASSWORD", "current password is incorrect")
	ErrUnauthorized       = NewUnauthorized("UNAUTHORIZED", "unauthorized access")
	ErrForbidden          = NewForbidden("FORBIDDEN", "you do not have permission to perform this action")

	// Product errors
	ErrProductNotFound    = NewNotFound("PRODUCT_NOT_FOUND", "product not found")
	ErrSKUAlreadyExists   = NewConflict("SKU_ALREADY_EXISTS", "product with this SKU already exists")
	ErrForbiddenProduct   = NewForbidden("FORBIDDEN_PRODUCT", "you do not have permission to modify this product")

	// Request errors
	ErrInvalidID          = NewBadRequest("INVALID_ID", "invalid UUID format")
	ErrInvalidPayload     = NewBadRequest("INVALID_PAYLOAD", "invalid request payload")
)

// AsAppError extracts an *AppError from an error chain if present
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
