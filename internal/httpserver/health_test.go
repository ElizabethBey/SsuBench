package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	HealthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusOK)
	}

	body := rr.Body.String()
	if body == "" {
		t.Fatalf("empty body")
	}
}

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rr := httptest.NewRecorder()

	HealthHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusMethodNotAllowed)
	}
}
