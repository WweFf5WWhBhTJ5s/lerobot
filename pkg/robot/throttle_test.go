package robot

import (
	"testing"
	"time"
)

func TestDefaultThrottleConfig(t *testing.T) {
	cfg := DefaultThrottleConfig()
	if cfg.MinInterval <= 0 {
		t.Fatalf("expected positive MinInterval, got %v", cfg.MinInterval)
	}
}

func TestThrottleAllowsFirstBeat(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())
	if !th.Allow("robot-1") {
		t.Fatal("expected first beat to be allowed")
	}
}

func TestThrottleBlocksRapidBeats(t *testing.T) {
	cfg := ThrottleConfig{MinInterval: 10 * time.Second}
	th := NewThrottle(cfg)

	if !th.Allow("robot-1") {
		t.Fatal("expected first beat to be allowed")
	}
	if th.Allow("robot-1") {
		t.Fatal("expected second rapid beat to be blocked")
	}
}

func TestThrottleAllowsAfterInterval(t *testing.T) {
	cfg := ThrottleConfig{MinInterval: 10 * time.Millisecond}
	th := NewThrottle(cfg)

	th.Allow("robot-1")
	time.Sleep(20 * time.Millisecond)

	if !th.Allow("robot-1") {
		t.Fatal("expected beat to be allowed after interval")
	}
}

func TestThrottleIsolatesRobots(t *testing.T) {
	cfg := ThrottleConfig{MinInterval: 10 * time.Second}
	th := NewThrottle(cfg)

	th.Allow("robot-1")
	if !th.Allow("robot-2") {
		t.Fatal("expected different robot to be allowed independently")
	}
}

func TestThrottleReset(t *testing.T) {
	cfg := ThrottleConfig{MinInterval: 10 * time.Second}
	th := NewThrottle(cfg)

	th.Allow("robot-1")
	th.Reset("robot-1")

	if !th.Allow("robot-1") {
		t.Fatal("expected beat to be allowed after reset")
	}
}

func TestThrottleLen(t *testing.T) {
	th := NewThrottle(DefaultThrottleConfig())

	if th.Len() != 0 {
		t.Fatalf("expected 0 tracked robots, got %d", th.Len())
	}

	th.Allow("robot-1")
	th.Allow("robot-2")

	if th.Len() != 2 {
		t.Fatalf("expected 2 tracked robots, got %d", th.Len())
	}

	th.Reset("robot-1")
	if th.Len() != 1 {
		t.Fatalf("expected 1 tracked robot after reset, got %d", th.Len())
	}
}
