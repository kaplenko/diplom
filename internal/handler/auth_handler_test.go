package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/kaplenko/diplom/internal/models"
)

func TestRespondValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	respondValidationError(c, "email is required")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp models.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("expected code VALIDATION_ERROR, got %s", resp.Error.Code)
	}
	if resp.Error.Message != "email is required" {
		t.Errorf("expected message 'email is required', got %s", resp.Error.Message)
	}
}

func TestParsePagination_Defaults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=&page_size=", nil)

	params := parsePagination(c)

	if params.Page != 1 {
		t.Errorf("expected page 1, got %d", params.Page)
	}
	if params.PageSize != 10 {
		t.Errorf("expected page_size 10, got %d", params.PageSize)
	}
}

func TestParsePagination_Custom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=3&page_size=25&search=golang", nil)

	params := parsePagination(c)

	if params.Page != 3 {
		t.Errorf("expected page 3, got %d", params.Page)
	}
	if params.PageSize != 25 {
		t.Errorf("expected page_size 25, got %d", params.PageSize)
	}
	if params.Search != "golang" {
		t.Errorf("expected search 'golang', got %s", params.Search)
	}
}

func TestParsePagination_ExceedsMax(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=1&page_size=500", nil)

	params := parsePagination(c)

	if params.PageSize != 100 {
		t.Errorf("expected page_size capped at 100, got %d", params.PageSize)
	}
}

func TestRegisterHandler_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"email": "not-an-email"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h := &AuthHandler{}
	h.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
