package robot

import (
	"testing"
	"time"
)

func TestDefaultPresenceConfig(t *testing.T) {
	cfg := DefaultPresenceConfig()
	if cfg.GracePeriod != 30*time.Second {
		t.Errorf("expected 30s grace period, got %v", cfg.GracePeriod)
	}
}

func TestPresenceTrackerTouchAndOnline(t *testing.T) {
	pt := NewPresenceTracker(DefaultPresenceConfig())
	pt.Touch("r1")
	if !pt.IsOnline("r1") {
		t.Error("expected r1 to be online after Touch")
	}
}

func TestPresenceTrackerUnknownRobotOffline(t *testing.T) {
	pt := NewPresenceTracker(DefaultPresenceConfig())
	if pt.IsOnline("unknown") {
		t.Error("expected unknown robot to be offline")
	}
}

func TestPresenceTrackerExpires(t *testing.T) {
	cfg := PresenceConfig{GracePeriod: 10 * time.Millisecond}
	pt := NewPresenceTracker(cfg)
	pt.Touch("r2")
	time.Sleep(20 * time.Millisecond)
	if pt.IsOnline("r2") {
		t.Error("expected r2 to be offline after grace period")
	}
}

func TestPresenceTrackerGet(t *testing.T) {
	pt := NewPresenceTracker(DefaultPresenceConfig())
	if rec := pt.Get("missing"); rec != nil {
		t.Errorf("expected nil for missing robot, got %+v", rec)
	}
	pt.Touch("r3")
	rec := pt.Get("r3")
	if rec == nil {
		t.Fatal("expected record for r3")
	}
	if rec.RobotID != "r3" {
		t.Errorf("expected RobotID r3, got %s", rec.RobotID)
	}
	if !rec.Online {
		t.Error("expected record to be online")
	}
}

func TestPresenceTrackerAll(t *testing.T) {
	pt := NewPresenceTracker(DefaultPresenceConfig())
	pt.Touch("a")
	pt.Touch("b")
	all := pt.All()
	if len(all) != 2 {
		t.Errorf("expected 2 records, got %d", len(all))
	}
	for _, r := range all {
		if !r.Online {
			t.Errorf("expected %s to be online", r.RobotID)
		}
	}
}

func TestPresenceTrackerAllOfflineAfterExpiry(t *testing.T) {
	cfg := PresenceConfig{GracePeriod: 10 * time.Millisecond}
	pt := NewPresenceTracker(cfg)
	pt.Touch("x")
	time.Sleep(20 * time.Millisecond)
	all := pt.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 record, got %d", len(all))
	}
	if all[0].Online {
		t.Error("expected record to be offline after expiry")
	}
}
