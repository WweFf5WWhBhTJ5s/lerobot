package robot

import (
	"testing"
)

func TestDefaultReplayConfig(t *testing.T) {
	cfg := DefaultReplayConfig()
	if cfg.MaxEvents != 200 {
		t.Fatalf("expected MaxEvents 200, got %d", cfg.MaxEvents)
	}
}

func TestReplayBufferRecordAndLen(t *testing.T) {
	rb := NewReplayBuffer(ReplayConfig{MaxEvents: 10})

	if rb.Len() != 0 {
		t.Fatalf("expected empty buffer, got %d", rb.Len())
	}

	s := NewStatus("r1")
	e := Event{Type: EventRegistered, RobotID: "r1", Status: &s}
	rb.Record(e)

	if rb.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", rb.Len())
	}
}

func TestReplayBufferEvictsOldest(t *testing.T) {
	rb := NewReplayBuffer(ReplayConfig{MaxEvents: 3})

	for i, id := range []string{"r1", "r2", "r3", "r4"} {
		s := NewStatus(id)
		rb.Record(Event{Type: EventRegistered, RobotID: id, Status: &s})
		_ = i
	}

	if rb.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", rb.Len())
	}

	events := rb.Replay()
	if events[0].RobotID != "r2" {
		t.Fatalf("expected oldest entry to be r2, got %s", events[0].RobotID)
	}
}

func TestReplayReturnsCopy(t *testing.T) {
	rb := NewReplayBuffer(ReplayConfig{MaxEvents: 5})
	s := NewStatus("r1")
	rb.Record(Event{Type: EventRegistered, RobotID: "r1", Status: &s})

	a := rb.Replay()
	a[0].RobotID = "mutated"

	b := rb.Replay()
	if b[0].RobotID != "r1" {
		t.Fatalf("Replay should return independent copy, got %s", b[0].RobotID)
	}
}

func TestReplayBufferClear(t *testing.T) {
	rb := NewReplayBuffer(ReplayConfig{MaxEvents: 5})
	s := NewStatus("r1")
	rb.Record(Event{Type: EventRegistered, RobotID: "r1", Status: &s})
	rb.Clear()

	if rb.Len() != 0 {
		t.Fatalf("expected empty buffer after Clear, got %d", rb.Len())
	}
}

func TestReplayBufferZeroMaxEventsUsesDefault(t *testing.T) {
	rb := NewReplayBuffer(ReplayConfig{MaxEvents: 0})
	if rb.cfg.MaxEvents != DefaultReplayConfig().MaxEvents {
		t.Fatalf("expected default MaxEvents, got %d", rb.cfg.MaxEvents)
	}
}
