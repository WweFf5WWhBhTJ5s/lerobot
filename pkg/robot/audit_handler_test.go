package robot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuditHandlerEmptyLog(t *testing.T) {
	log := NewAuditLog(10)
	h := AuditHandler(log)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit", nil)
	h(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var entries []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty list, got %d entries", len(entries))
	}
}

func TestAuditHandlerReturnsEntries(t *testing.T) {
	log := NewAuditLog(10)
	st := NewStatus("bot-1")
	log.Handle(Event{Type: EventRegistered, Status: &st})
	h := AuditHandler(log)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit", nil)
	h(rec, req)
	var entries []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0]["robot_id"] != "bot-1" {
		t.Fatalf("expected robot_id bot-1, got %v", entries[0]["robot_id"])
	}
}

func TestAuditHandlerContentType(t *testing.T) {
	log := NewAuditLog(10)
	h := AuditHandler(log)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit", nil)
	h(rec, req)
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}
}

func TestAuditHandlerEventType(t *testing.T) {
	log := NewAuditLog(10)
	st := NewStatus("bot-2")
	log.Handle(Event{Type: EventUnregistered, Status: &st})
	h := AuditHandler(log)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit", nil)
	h(rec, req)
	var entries []map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&entries)
	if entries[0]["type"] != EventUnregistered.String() {
		t.Fatalf("unexpected event type: %v", entries[0]["type"])
	}
}
