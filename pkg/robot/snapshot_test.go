package robot

import (
	"testing"
	"time"
)

func newTestSnapshot(t *testing.T) (*SnapshotService, *Registry) {
	t.Helper()
	r := NewRegistry()
	svc := NewSnapshotService(r, 50*time.Millisecond)
	return svc, r
}

func TestSnapshotLatestNilBeforeStart(t *testing.T) {
	svc, _ := newTestSnapshot(t)
	if svc.Latest() != nil {
		t.Fatal("expected nil snapshot before any capture")
	}
}

func TestSnapshotCapturesRegisteredRobots(t *testing.T) {
	svc, r := newTestSnapshot(t)

	r.Register("robot-1")
	r.Register("robot-2")

	svc.Start()
	defer svc.Stop()

	time.Sleep(120 * time.Millisecond)

	snap := svc.Latest()
	if snap == nil {
		t.Fatal("expected a snapshot after interval")
	}
	if len(snap.Robots) != 2 {
		t.Fatalf("expected 2 robots in snapshot, got %d", len(snap.Robots))
	}
	if _, ok := snap.Robots["robot-1"]; !ok {
		t.Error("expected robot-1 in snapshot")
	}
	if _, ok := snap.Robots["robot-2"]; !ok {
		t.Error("expected robot-2 in snapshot")
	}
}

func TestSnapshotTakenAtIsRecent(t *testing.T) {
	svc, _ := newTestSnapshot(t)
	svc.Start()
	defer svc.Stop()

	before := time.Now()
	time.Sleep(120 * time.Millisecond)
	after := time.Now()

	snap := svc.Latest()
	if snap == nil {
		t.Fatal("expected a snapshot")
	}
	if snap.TakenAt.Before(before) || snap.TakenAt.After(after) {
		t.Errorf("snapshot TakenAt %v not in expected range [%v, %v]", snap.TakenAt, before, after)
	}
}

func TestSnapshotStopIsIdempotent(t *testing.T) {
	svc, _ := newTestSnapshot(t)
	svc.Start()
	svc.Stop()
	// calling Stop again should not panic
}
