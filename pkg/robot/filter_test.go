package robot

import (
	"testing"
	"time"
)

func newTestFilter(t *testing.T) (*FilterService, *TagStore) {
	t.Helper()
	store := NewTagStore()
	svc := NewFilterService(store, DefaultFilterConfig())
	return svc, store
}

func makeStatuses(ids ...string) []Status {
	out := make([]Status, len(ids))
	for i, id := range ids {
		out[i] = Status{RobotID: id, LastSeen: time.Now()}
	}
	return out
}

func TestDefaultFilterConfig(t *testing.T) {
	cfg := DefaultFilterConfig()
	// Default should be case-insensitive for friendlier UX
	if cfg.CaseSensitive {
		t.Error("expected CaseSensitive to be false by default")
	}
}

func TestFilterByNameSubstring(t *testing.T) {
	svc, _ := newTestFilter(t)
	statuses := makeStatuses("robot-alpha", "robot-beta", "sensor-gamma")

	result := svc.ByName(statuses, "robot")
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestFilterByNameCaseInsensitive(t *testing.T) {
	svc, _ := newTestFilter(t)
	statuses := makeStatuses("Robot-Alpha", "robot-beta", "sensor-gamma")

	result := svc.ByName(statuses, "ROBOT")
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestFilterByNameCaseSensitive(t *testing.T) {
	store := NewTagStore()
	svc := NewFilterService(store, FilterConfig{CaseSensitive: true})
	statuses := makeStatuses("Robot-Alpha", "robot-beta")

	result := svc.ByName(statuses, "Robot")
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].RobotID != "Robot-Alpha" {
		t.Errorf("unexpected robot: %s", result[0].RobotID)
	}
}

func TestFilterByNameEmptyQuery(t *testing.T) {
	// An empty query string should return all statuses unchanged.
	svc, _ := newTestFilter(t)
	statuses := makeStatuses("robot-alpha", "robot-beta", "sensor-gamma")

	result := svc.ByName(statuses, "")
	if len(result) != len(statuses) {
		t.Fatalf("expected %d results for empty query, got %d", len(statuses), len(result))
	}
}

func TestFilterByTag(t *testing.T) {
	svc, store := newTestFilter(t)
	store.Set("robot-alpha", []string{"fast", "arm"})
	store.Set("robot-beta", []string{"slow"})
	store.Set("robot-gamma", []string{"fast"})

	statuses := makeStatuses("robot-alpha", "robot-beta", "robot-gamma")

	result := svc.ByTag(statuses, "fast")
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestFilterByTagNoMatch(t *testing.T) {
	svc, _ := newTestFilter(t)
	statuses := makeStatuses("robot-alpha", "robot-beta")

	result := svc.ByTag(statuses, "nonexistent")
	if len(result) != 0 {
		t.Fatalf("expected 0 results, got %d", len(result))
	}
}
