package robot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandlerOKWhenAllHealthy(t *testing.T) {
	r := NewRegistry(NewNotifier())
	r.Register("bot-a")
	r.Beat("bot-a")

	hc := NewHealthChecker(r, 20*time.Millisecond)
	hc.Start()
	defer hc.Stop()
	time.Sleep(50 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	HealthHandler(hc).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !body.OK {
		t.Fatalf("expected ok=true")
	}
}

func TestHealthHandlerUnavailableWhenEmpty(t *testing.T) {
	r := NewRegistry(NewNotifier())
	hc := NewHealthChecker(r, 20*time.Millisecond)
	hc.Start()
	defer hc.Stop()
	time.Sleep(50 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	HealthHandler(hc).ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}

func TestHealthHandlerContentType(t *testing.T) {
	r := NewRegistry(NewNotifier())
	hc := NewHealthChecker(r, 20*time.Millisecond)
	hc.Start()
	defer hc.Stop()
	time.Sleep(50 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	HealthHandler(hc).ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}
