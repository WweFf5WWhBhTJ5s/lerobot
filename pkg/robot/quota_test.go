package robot

import (
	"testing"
	"time"
)

func TestDefaultQuotaConfig(t *testing.T) {
	cfg := DefaultQuotaConfig()
	if cfg.MaxBeatsPerHour != 3600 {
		t.Errorf("expected MaxBeatsPerHour=3600, got %d", cfg.MaxBeatsPerHour)
	}
	if cfg.WindowDuration != time.Hour {
		t.Errorf("expected WindowDuration=1h, got %v", cfg.WindowDuration)
	}
}

func TestQuotaAllowsUnderLimit(t *testing.T) {
	cfg := QuotaConfig{MaxBeatsPerHour: 3, WindowDuration: time.Hour}
	q := NewQuotaEnforcer(cfg)
	for i := 0; i < 3; i++ {
		if !q.Allow("robot-1") {
			t.Fatalf("expected allow on beat %d", i+1)
		}
	}
}

func TestQuotaBlocksOverLimit(t *testing.T) {
	cfg := QuotaConfig{MaxBeatsPerHour: 2, WindowDuration: time.Hour}
	q := NewQuotaEnforcer(cfg)
	q.Allow("robot-1")
	q.Allow("robot-1")
	if q.Allow("robot-1") {
		t.Error("expected deny after quota exceeded")
	}
}

func TestQuotaIsolatesRobots(t *testing.T) {
	cfg := QuotaConfig{MaxBeatsPerHour: 1, WindowDuration: time.Hour}
	q := NewQuotaEnforcer(cfg)
	q.Allow("robot-1")
	if !q.Allow("robot-2") {
		t.Error("robot-2 should not be affected by robot-1 quota")
	}
}

func TestQuotaResetsWindow(t *testing.T) {
	now := time.Now()
	cfg := QuotaConfig{MaxBeatsPerHour: 1, WindowDuration: time.Minute}
	q := NewQuotaEnforcer(cfg)
	q.now = func() time.Time { return now }
	q.Allow("robot-1")
	if q.Allow("robot-1") {
		t.Fatal("expected deny within window")
	}
	q.now = func() time.Time { return now.Add(2 * time.Minute) }
	if !q.Allow("robot-1") {
		t.Error("expected allow after window expired")
	}
}

func TestQuotaUsage(t *testing.T) {
	cfg := QuotaConfig{MaxBeatsPerHour: 10, WindowDuration: time.Hour}
	q := NewQuotaEnforcer(cfg)
	if q.Usage("robot-1") != 0 {
		t.Error("expected 0 usage before any beats")
	}
	q.Allow("robot-1")
	q.Allow("robot-1")
	if q.Usage("robot-1") != 2 {
		t.Errorf("expected usage=2, got %d", q.Usage("robot-1"))
	}
}

func TestQuotaReset(t *testing.T) {
	cfg := QuotaConfig{MaxBeatsPerHour: 1, WindowDuration: time.Hour}
	q := NewQuotaEnforcer(cfg)
	q.Allow("robot-1")
	q.Reset("robot-1")
	if !q.Allow("robot-1") {
		t.Error("expected allow after reset")
	}
}
