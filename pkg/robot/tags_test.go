package robot

import (
	"testing"
)

func TestTagStoreSetAndGet(t *testing.T) {
	ts := NewTagStore()
	ts.Set("r1", []string{"arm", "mobile", "arm"})
	tags := ts.Get("r1")
	if len(tags) != 2 {
		t.Fatalf("expected 2 unique tags, got %d", len(tags))
	}
	if tags[0] != "arm" || tags[1] != "mobile" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestTagStoreAdd(t *testing.T) {
	ts := NewTagStore()
	ts.Add("r1", "arm")
	ts.Add("r1", "mobile", "sensor")
	tags := ts.Get("r1")
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d: %v", len(tags), tags)
	}
}

func TestTagStoreRemove(t *testing.T) {
	ts := NewTagStore()
	ts.Set("r1", []string{"arm", "mobile"})
	ts.Remove("r1", "arm")
	tags := ts.Get("r1")
	if len(tags) != 1 || tags[0] != "mobile" {
		t.Errorf("expected [mobile], got %v", tags)
	}
}

func TestTagStoreDelete(t *testing.T) {
	ts := NewTagStore()
	ts.Set("r1", []string{"arm"})
	ts.Delete("r1")
	if tags := ts.Get("r1"); len(tags) != 0 {
		t.Errorf("expected empty tags after delete, got %v", tags)
	}
}

func TestTagStoreFindByTag(t *testing.T) {
	ts := NewTagStore()
	ts.Set("r1", []string{"arm", "mobile"})
	ts.Set("r2", []string{"mobile", "sensor"})
	ts.Set("r3", []string{"arm"})

	ids := ts.FindByTag("mobile")
	if len(ids) != 2 {
		t.Fatalf("expected 2 robots with tag 'mobile', got %d: %v", len(ids), ids)
	}
	if ids[0] != "r1" || ids[1] != "r2" {
		t.Errorf("unexpected robot IDs: %v", ids)
	}
}

func TestTagStoreFindByTagNone(t *testing.T) {
	ts := NewTagStore()
	ts.Set("r1", []string{"arm"})
	ids := ts.FindByTag("sensor")
	if len(ids) != 0 {
		t.Errorf("expected no results, got %v", ids)
	}
}

func TestTagStoreIgnoresBlankTags(t *testing.T) {
	ts := NewTagStore()
	ts.Set("r1", []string{"  ", "", "arm"})
	tags := ts.Get("r1")
	if len(tags) != 1 || tags[0] != "arm" {
		t.Errorf("expected only 'arm', got %v", tags)
	}
}

func TestTagStoreString(t *testing.T) {
	ts := NewTagStore()
	ts.Set("r1", []string{"arm"})
	s := ts.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
