package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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
