package robot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTagsHandlerGet(t *testing.T) {
	store := NewTagStore()
	store.Set("r1", []string{"arm", "vision"})

	req := httptest.NewRequest(http.MethodGet, "/tags/r1", nil)
	rec := httptest.NewRecorder()
	TagsHandler(store)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp map[string][]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp["tags"]) != 2 {
		t.Errorf("expected 2 tags, got %d", len(resp["tags"]))
	}
}

func TestTagsHandlerPost(t *testing.T) {
	store := NewTagStore()
	body, _ := json.Marshal(map[string][]string{"tags": {"sensor", "mobile"}})

	req := httptest.NewRequest(http.MethodPost, "/tags/r2", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	TagsHandler(store)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	tags := store.Get("r2")
	if len(tags) != 2 {
		t.Errorf("expected 2 tags after POST, got %d", len(tags))
	}
}

func TestTagsHandlerDelete(t *testing.T) {
	store := NewTagStore()
	store.Set("r3", []string{"arm", "vision"})

	req := httptest.NewRequest(http.MethodDelete, "/tags/r3/arm", nil)
	rec := httptest.NewRecorder()
	TagsHandler(store)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	tags := store.Get("r3")
	if len(tags) != 1 || tags[0] != "vision" {
		t.Errorf("expected only 'vision' remaining, got %v", tags)
	}
}

func TestTagsHandlerContentType(t *testing.T) {
	store := NewTagStore()
	req := httptest.NewRequest(http.MethodGet, "/tags/r4", nil)
	rec := httptest.NewRecorder()
	TagsHandler(store)(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestTagsHandlerMethodNotAllowed(t *testing.T) {
	store := NewTagStore()
	req := httptest.NewRequest(http.MethodPut, "/tags/r5", nil)
	rec := httptest.NewRecorder()
	TagsHandler(store)(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
