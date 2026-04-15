package robot

import (
	"testing"
	"time"
)

func TestRegistryRegister(t *testing.T) {
	reg := NewRegistry(500*time.Millisecond, 100*time.Millisecond)

	if err := reg.Register("robot-1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Duplicate registration should fail.
	if err := reg.Register("robot-1"); err == nil {
		t.Fatal("expected error for duplicate registration, got nil")
	}

	_ = reg.Unregister("robot-1")
}

func TestRegistryUnregisterUnknown(t *testing.T) {
	reg := NewRegistry(500*time.Millisecond, 100*time.Millisecond)

	if err := reg.Unregister("ghost"); err == nil {
		t.Fatal("expected error for unknown robot, got nil")
	}
}

func TestRegistryBeat(t *testing.T) {
	reg := NewRegistry(500*time.Millisecond, 100*time.Millisecond)
	_ = reg.Register("robot-2")
	defer reg.Unregister("robot-2") //nolint:errcheck

	if err := reg.Beat("robot-2"); err != nil {
		t.Fatalf("expected no error on beat, got %v", err)
	}

	if err := reg.Beat("unknown"); err == nil {
		t.Fatal("expected error for unknown robot beat, got nil")
	}
}

func TestRegistryStatus(t *testing.T) {
	reg := NewRegistry(500*time.Millisecond, 100*time.Millisecond)
	_ = reg.Register("robot-3")
	defer reg.Unregister("robot-3") //nolint:errcheck

	_ = reg.Beat("robot-3")

	s, err := reg.Status("robot-3")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !s.IsHealthy() {
		t.Errorf("expected robot-3 to be healthy after beat")
	}

	_, err = reg.Status("nobody")
	if err == nil {
		t.Fatal("expected error for unknown robot status, got nil")
	}
}

func TestRegistryList(t *testing.T) {
	reg := NewRegistry(500*time.Millisecond, 100*time.Millisecond)
	_ = reg.Register("r1")
	_ = reg.Register("r2")
	defer reg.Unregister("r1") //nolint:errcheck
	defer reg.Unregister("r2") //nolint:errcheck

	ids := reg.List()
	if len(ids) != 2 {
		t.Errorf("expected 2 robots, got %d", len(ids))
	}
}
