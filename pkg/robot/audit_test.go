package robot

import (
	"fmt"
	"testing"
	"time"
)

func TestAuditLogDefaultMaxSize(t *testing.T) {
	al := NewAuditLog(0)
	if al.maxSize != 256 {
		t.Fatalf("expected default maxSize 256, got %d", al.maxSize)
	}
}

func TestAuditLogHandleAddsEntry(t *testing.T) {
	al := NewAuditLog(10)
	st := NewStatus("bot-1")
	e := Event{Type: EventRegistered, Status: &st}
	al.Handle(e)
	if al.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", al.Len())
	}
}

func TestAuditLogEntriesAreCopied(t *testing.T) {
	al := NewAuditLog(10)
	st := NewStatus("bot-2")
	e := Event{Type: EventRegistered, Status: &st}
	al.Handle(e)
	entries := al.Entries()
	entries[0].Event.Type = EventUnregistered // mutate copy
	original := al.Entries()
	if original[0].Event.Type != EventRegistered {
		t.Fatal("Entries() should return a copy, not a reference")
	}
}

func TestAuditLogBoundedSize(t *testing.T) {
	const max = 5
	al := NewAuditLog(max)
	for i := 0; i < max+3; i++ {
		st := NewStatus(fmt.Sprintf("bot-%d", i))
		al.Handle(Event{Type: EventRegistered, Status: &st})
	}
	if al.Len() != max {
		t.Fatalf("expected %d entries, got %d", max, al.Len())
	}
}

func TestAuditLogOldestEvictedFirst(t *testing.T) {
	const max = 3
	al := NewAuditLog(max)
	for i := 0; i < max+1; i++ {
		st := NewStatus(fmt.Sprintf("bot-%d", i))
		al.Handle(Event{Type: EventRegistered, Status: &st})
	}
	entries := al.Entries()
	// After eviction the oldest (bot-0) should be gone; first entry should be bot-1.
	if entries[0].Event.Status.ID != "bot-1" {
		t.Fatalf("expected first entry to be bot-1, got %s", entries[0].Event.Status.ID)
	}
}

func TestAuditEntryString(t *testing.T) {
	st := NewStatus("bot-x")
	e := Event{Type: EventRegistered, Status: &st}
	entry := AuditEntry{Timestamp: time.Now(), Event: e}
	s := entry.String()
	if len(s) == 0 {
		t.Fatal("expected non-empty string from AuditEntry.String()")
	}
}
