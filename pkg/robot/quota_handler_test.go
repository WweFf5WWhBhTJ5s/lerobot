package robot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQuotaHandlerMissingID(t *testing.T) {
	q := NewQuotaEnforcer(DefaultQuotaConfig())
	h := QuotaHandler(q)
	req := httptest.NewRequest(http.MethodGet, "/quota", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestQuotaHandlerReturnsUsage(t *testing.T) {
	cfg := QuotaConfig{MaxBeatsPerHour: 5, WindowDuration: time.Hour}
	q := NewQuotaEnforcer(cfg)
	q.Allow("robot-42")
	q.Allow("robot-42")

	h := QuotaHandler(q)
	req := httptest.NewRequest(http.MethodGet, "/quota?id=robot-42", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp quotaUsageResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.RobotID != "robot-42" {
		t.Errorf("expected robot_id=robot-42, got %s", resp.RobotID)
	}
	if resp.Usage != 2 {
		t.Errorf("expected usage=2, got %d", resp.Usage)
	}
	if resp.Limit != 5 {
		t.Errorf("expected limit=5, got %d", resp.Limit)
	}
	if !resp.Allowed {
		t.Error("expected allowed=true")
	}
}

func TestQuotaHandlerAllowedFalseWhenExceeded(t *testing.T) {
	cfg := QuotaConfig{MaxBeatsPerHour: 1, WindowDuration: time.Hour}
	q := NewQuotaEnforcer(cfg)
	q.Allow("robot-99")

	h := QuotaHandler(q)
	req := httptest.NewRequest(http.MethodGet, "/quota?id=robot-99", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	var resp quotaUsageResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Allowed {
		t.Error("expected allowed=false when quota exhausted")
	}
}

func TestQuotaHandlerContentType(t *testing.T) {
	q := NewQuotaEnforcer(DefaultQuotaConfig())
	h := QuotaHandler(q)
	req := httptest.NewRequest(http.MethodGet, "/quota?id=robot-1", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}
