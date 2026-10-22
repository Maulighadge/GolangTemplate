package response

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response represents the standard API response structure
type Response struct {
	StatusCode int           `json:"status_code"`
	Success    bool          `json:"success"`
	Message    string        `json:"message"`
	Data       interface{}   `json:"data,omitempty"`
	Error      *ErrorPayload `json:"error,omitempty"`
	Meta       *Meta         `json:"meta,omitempty"`
}

// Meta contains metadata such as timestamps and pagination details
type Meta struct {
	Timestamp  time.Time   `json:"timestamp"`
	Pagination interface{} `json:"pagination,omitempty"`
}

// ErrorPayload contains structured error information
type ErrorPayload struct {
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// FieldError represents a single validation failure on a struct field
type FieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag,omitempty"`
	Message string `json:"message"`
}

// Success sends a standardized success JSON response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Success:    true,
		Message:    message,
		Data:       data,
		Meta: &Meta{
			Timestamp: time.Now().UTC(),
		},
	})
}

// OK sends a 200 OK success response
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

// Created sends a 201 Created success response
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// SuccessWithMeta sends a standardized success JSON response with custom metadata/pagination
func SuccessWithMeta(c *gin.Context, statusCode int, message string, data interface{}, pagination interface{}) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Success:    true,
		Message:    message,
		Data:       data,
		Meta: &Meta{
			Timestamp:  time.Now().UTC(),
			Pagination: pagination,
		},
	})
}

// Paginated sends a 200 OK response with pagination metadata
func Paginated(c *gin.Context, message string, data interface{}, pagination interface{}) {
	SuccessWithMeta(c, http.StatusOK, message, data, pagination)
}

// Error sends a standardized error JSON response
func Error(c *gin.Context, statusCode int, message string, errDetail interface{}) {
	var payload *ErrorPayload

	if errDetail != nil {
		switch v := errDetail.(type) {
		case *ErrorPayload:
			payload = v
		case ErrorPayload:
			payload = &v
		case error:
			payload = &ErrorPayload{
				Code:    http.StatusText(statusCode),
				Details: v.Error(),
			}
		case string:
			if v != "" {
				payload = &ErrorPayload{
					Code:    http.StatusText(statusCode),
					Details: v,
				}
			}
		default:
			payload = &ErrorPayload{
				Code:    http.StatusText(statusCode),
				Details: v,
			}
		}
	}

	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Success:    false,
		Message:    message,
		Error:      payload,
		Meta: &Meta{
			Timestamp: time.Now().UTC(),
		},
	})
}

// BadRequest sends a 400 Bad Request error response
func BadRequest(c *gin.Context, message string, errDetail interface{}) {
	Error(c, http.StatusBadRequest, message, errDetail)
}

// ValidationError formats binding/validation errors into a structured 422 Unprocessable Entity response
func ValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fieldErrors := make([]FieldError, len(ve))
		for i, fe := range ve {
			fieldErrors[i] = FieldError{
				Field:   fe.Field(),
				Tag:     fe.Tag(),
				Message: formatValidationTag(fe),
			}
		}

		c.JSON(http.StatusUnprocessableEntity, Response{
			StatusCode: http.StatusUnprocessableEntity,
			Success:    false,
			Message:    "Validation failed",
			Error: &ErrorPayload{
				Code:    "VALIDATION_ERROR",
				Details: fieldErrors,
			},
			Meta: &Meta{
				Timestamp: time.Now().UTC(),
			},
		})
		return
	}

	// General syntax / unmarshal error fallback
	BadRequest(c, "Invalid request payload", err.Error())
}

// Unauthorized sends a 401 Unauthorized error response
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden sends a 403 Forbidden error response
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Access forbidden"
	}
	Error(c, http.StatusForbidden, message, nil)
}

// NotFound sends a 404 Not Found error response
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(c, http.StatusNotFound, message, nil)
}

// Conflict sends a 409 Conflict error response
func Conflict(c *gin.Context, message string, errDetail interface{}) {
	Error(c, http.StatusConflict, message, errDetail)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string, errDetail interface{}) {
	if message == "" {
		message = "Internal server error"
	}
	Error(c, http.StatusInternalServerError, message, errDetail)
}

// formatValidationTag returns clear, human-readable validation messages for common validation tags
func formatValidationTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s cannot exceed %s characters", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s failed on rule '%s'", fe.Field(), fe.Tag())
	}
}
