package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/Maulighadge/GolangTemplate/pkg/apperror"
	"github.com/Maulighadge/GolangTemplate/pkg/logger"
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

// Meta contains metadata such as timestamps, request ID, and pagination details
type Meta struct {
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"request_id,omitempty"`
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

func newMeta(c *gin.Context, pagination interface{}) *Meta {
	meta := &Meta{
		Timestamp:  time.Now().UTC(),
		Pagination: pagination,
	}
	if c != nil {
		if reqID, exists := c.Get("requestID"); exists {
			if idStr, ok := reqID.(string); ok && idStr != "" {
				meta.RequestID = idStr
			}
		} else if c.Request != nil {
			if headerID := c.GetHeader("X-Request-ID"); headerID != "" {
				meta.RequestID = headerID
			}
		}
	}
	return meta
}

// Success sends a standardized success JSON response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Success:    true,
		Message:    message,
		Data:       data,
		Meta:       newMeta(c, nil),
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
		Meta:       newMeta(c, pagination),
	})
}

// Paginated sends a 200 OK response with pagination metadata
func Paginated(c *gin.Context, message string, data interface{}, pagination interface{}) {
	SuccessWithMeta(c, http.StatusOK, message, data, pagination)
}

// Error sends a standardized error JSON response
func Error(c *gin.Context, statusCode int, message string, errDetail interface{}) {
	ErrorWithCode(c, statusCode, message, http.StatusText(statusCode), errDetail)
}

// ErrorWithCode sends a standardized error JSON response with an explicit machine-readable error code
func ErrorWithCode(c *gin.Context, statusCode int, message string, code string, errDetail interface{}) {
	var payload *ErrorPayload

	if code == "" {
		code = http.StatusText(statusCode)
	}

	if errDetail != nil {
		switch v := errDetail.(type) {
		case *ErrorPayload:
			payload = v
		case ErrorPayload:
			payload = &v
		case *apperror.AppError:
			payload = &ErrorPayload{
				Code:    v.Code,
				Details: v.Details,
			}
		case error:
			payload = &ErrorPayload{
				Code:    code,
				Details: v.Error(),
			}
		case string:
			if v != "" {
				payload = &ErrorPayload{
					Code:    code,
					Details: v,
				}
			}
		default:
			payload = &ErrorPayload{
				Code:    code,
				Details: v,
			}
		}
	} else if code != "" {
		payload = &ErrorPayload{
			Code: code,
		}
	}

	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Success:    false,
		Message:    message,
		Error:      payload,
		Meta:       newMeta(c, nil),
	})
}

// HandleError inspects the error and dispatches the appropriate HTTP response.
// If err is an *apperror.AppError, it uses its Kind, Code, Message, and Details.
// For internal errors (status >= 500), it logs root causes with the request ID and
// returns a safe sanitized message without leaking internal implementation details.
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		status := appErr.HTTPStatusCode()
		if status >= 500 {
			reqID := ""
			if id, exists := c.Get("requestID"); exists {
				reqID = fmt.Sprintf("%v", id)
			}
			logger.Error("[ReqID: %s] Internal server error: %v", reqID, appErr)
			ErrorWithCode(c, status, "An unexpected internal server error occurred", appErr.Code, nil)
			return
		}

		ErrorWithCode(c, status, appErr.Message, appErr.Code, appErr.Details)
		return
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		ValidationError(c, err)
		return
	}

	// Unhandled error -> 500
	reqID := ""
	if id, exists := c.Get("requestID"); exists {
		reqID = fmt.Sprintf("%v", id)
	}
	logger.Error("[ReqID: %s] Unhandled unexpected error: %v", reqID, err)
	InternalServerError(c, "An unexpected internal server error occurred", nil)
}

// BadRequest sends a 400 Bad Request error response
func BadRequest(c *gin.Context, message string, errDetail interface{}) {
	ErrorWithCode(c, http.StatusBadRequest, message, "BAD_REQUEST", errDetail)
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
			Meta: newMeta(c, nil),
		})
		return
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid request payload",
			Error: &ErrorPayload{
				Code: "INVALID_FIELD_TYPE",
				Details: []FieldError{
					{
						Field:   typeErr.Field,
						Message: fmt.Sprintf("Expected %s but received %s", typeErr.Type.String(), typeErr.Value),
					},
				},
			},
			Meta: newMeta(c, nil),
		})
		return
	}

	// General syntax / unmarshal error fallback
	c.JSON(http.StatusBadRequest, Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
		Message:    "Invalid request payload",
		Error: &ErrorPayload{
			Code:    "INVALID_PAYLOAD_FORMAT",
			Details: err.Error(),
		},
		Meta: newMeta(c, nil),
	})
}

// Unauthorized sends a 401 Unauthorized error response
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	ErrorWithCode(c, http.StatusUnauthorized, message, "UNAUTHORIZED", nil)
}

// Forbidden sends a 403 Forbidden error response
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Access forbidden"
	}
	ErrorWithCode(c, http.StatusForbidden, message, "FORBIDDEN", nil)
}

// NotFound sends a 404 Not Found error response
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorWithCode(c, http.StatusNotFound, message, "NOT_FOUND", nil)
}

// Conflict sends a 409 Conflict error response
func Conflict(c *gin.Context, message string, errDetail interface{}) {
	ErrorWithCode(c, http.StatusConflict, message, "CONFLICT", errDetail)
}

// InternalServerError sends a 500 Internal Server Error response without leaking details to clients
func InternalServerError(c *gin.Context, message string, errDetail interface{}) {
	if message == "" {
		message = "An unexpected internal server error occurred"
	}
	ErrorWithCode(c, http.StatusInternalServerError, message, "INTERNAL_SERVER_ERROR", errDetail)
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
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", fe.Field(), fe.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	case "numeric":
		return fmt.Sprintf("%s must be numeric", fe.Field())
	default:
		return fmt.Sprintf("%s failed on rule '%s'", fe.Field(), fe.Tag())
	}
}
