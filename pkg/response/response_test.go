package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/Maulighadge/GolangTemplate/pkg/apperror"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, http.StatusOK, "Success test", gin.H{"id": "123"})

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var res Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !res.Success {
		t.Errorf("Expected Success to be true")
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected StatusCode 200, got %d", res.StatusCode)
	}
	if res.Message != "Success test" {
		t.Errorf("Expected Message 'Success test', got %s", res.Message)
	}
	if res.Meta == nil || res.Meta.Timestamp.IsZero() {
		t.Errorf("Expected Meta.Timestamp to be set")
	}
}

func TestSuccessWithMetaResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	pagination := gin.H{"page": 1, "limit": 10, "total": 100}
	SuccessWithMeta(c, http.StatusOK, "Paginated test", []string{"a", "b"}, pagination)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var res Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.Meta == nil || res.Meta.Pagination == nil {
		t.Fatalf("Expected Meta.Pagination to be populated")
	}
}

func TestErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, http.StatusBadRequest, "Invalid request", errors.New("bad input"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}

	var res Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.Success {
		t.Errorf("Expected Success to be false")
	}
	if res.StatusCode != 400 {
		t.Errorf("Expected StatusCode 400, got %d", res.StatusCode)
	}
	if res.Error == nil || res.Error.Code != "Bad Request" {
		t.Errorf("Expected Error.Code 'Bad Request', got %+v", res.Error)
	}
}

func TestHandleError_WithAppErrorNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	HandleError(c, apperror.ErrProductNotFound)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d", w.Code)
	}

	var res Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.Success {
		t.Errorf("Expected Success false")
	}
	if res.StatusCode != 404 {
		t.Errorf("Expected StatusCode 404, got %d", res.StatusCode)
	}
	if res.Error == nil || res.Error.Code != "PRODUCT_NOT_FOUND" {
		t.Errorf("Expected Error.Code 'PRODUCT_NOT_FOUND', got %+v", res.Error)
	}
	if res.Message != "product not found" {
		t.Errorf("Expected Message 'product not found', got '%s'", res.Message)
	}
}

func TestHandleError_WithAppErrorInternalMasking(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate root DB failure wrapped in Internal error
	rootCause := errors.New("pq: connection refused at 10.0.0.1:5432")
	internalErr := apperror.NewInternal(rootCause, "something internal failed")

	HandleError(c, internalErr)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status 500, got %d", w.Code)
	}

	var res Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify that internal DB string is NOT leaked to clients
	if res.Message == rootCause.Error() || (res.Error != nil && res.Error.Details == rootCause.Error()) {
		t.Fatalf("Security failure: internal error details leaked to client! %+v", res)
	}
	if res.Error.Code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("Expected Code 'INTERNAL_SERVER_ERROR', got '%s'", res.Error.Code)
	}
}

func TestHandleError_WithGenericError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	genericErr := errors.New("unexpected os filesystem error")
	HandleError(c, genericErr)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status 500, got %d", w.Code)
	}

	var res Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.StatusCode != 500 {
		t.Errorf("Expected status 500, got %d", res.StatusCode)
	}
}

func TestMeta_WithRequestID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("requestID", "req-xyz-12345")

	OK(c, "Success", nil)

	var res Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.Meta == nil || res.Meta.RequestID != "req-xyz-12345" {
		t.Errorf("Expected Meta.RequestID to be 'req-xyz-12345', got %+v", res.Meta)
	}
}
