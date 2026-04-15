package robot

import (
	"testing"
	"time"
)

func newTestDiscovery(ttl time.Duration) (*DiscoveryService, *Registry, *Notifier) {
	reg := NewRegistry()
	not := NewNotifier()
	cfg := DiscoveryConfig{
		ScanInterval: 5 * time.Millisecond,
		RobotTTL:     ttl,
	}
	return NewDiscoveryService(cfg, reg, not), reg, not
}

func TestDiscoverySeeRegisters(t *testing.T) {
	ds, reg, _ := newTestDiscovery(time.Minute)

	if err := ds.See("robot-1"); err != nil {
		t.Fatalf("See returned unexpected error: %v", err)
	}

	list := reg.List()
	if len(list) != 1 || list[0] != "robot-1" {
		t.Fatalf("expected [robot-1], got %v", list)
	}
}

func TestDiscoverySeeIdempotent(t *testing.T) {
	ds, reg, _ := newTestDiscovery(time.Minute)

	for i := 0; i < 3; i++ {
		if err := ds.See("robot-1"); err != nil {
			t.Fatalf("See[%d] returned unexpected error: %v", i, err)
		}
	}

	if n := len(reg.List()); n != 1 {
		t.Fatalf("expected 1 robot, got %d", n)
	}
}

func TestDiscoveryEvictsStaleRobots(t *testing.T) {
	ds, reg, _ := newTestDiscovery(20 * time.Millisecond)
	ds.Start()
	defer ds.Stop()

	if err := ds.See("robot-stale"); err != nil {
		t.Fatalf("See: %v", err)
	}

	// Wait for TTL + a couple of scan intervals.
	time.Sleep(60 * time.Millisecond)

	if list := reg.List(); len(list) != 0 {
		t.Fatalf("expected robot to be evicted, still registered: %v", list)
	}
}

func TestDiscoveryNotifiesOnSee(t *testing.T) {
	ds, _, not := newTestDiscovery(time.Minute)

	events := make([]Event, 0)
	not.Subscribe(func(e Event, id string) {
		events = append(events, e)
	})

	_ = ds.See("robot-notify")

	if len(events) != 1 || events[0] != EventRegistered {
		t.Fatalf("expected [EventRegistered], got %v", events)
	}
}

func TestDiscoveryStartStop(t *testing.T) {
	ds, _, _ := newTestDiscovery(time.Minute)
	ds.Start()
	// Should not block or panic.
	ds.Stop()
}
