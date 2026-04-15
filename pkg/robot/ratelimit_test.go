package robot

import (
	"testing"
	"time"
)

func TestDefaultRateLimitConfig(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	if cfg.MaxEvents <= 0 {
		t.Errorf("expected positive MaxEvents, got %d", cfg.MaxEvents)
	}
	if cfg.Window <= 0 {
		t.Errorf("expected positive Window, got %v", cfg.Window)
	}
}

func TestRateLimiterAllowsUnderLimit(t *testing.T) {
	rl := NewRateLimiter(RateLimitConfig{MaxEvents: 3, Window: time.Minute})
	for i := 0; i < 3; i++ {
		if !rl.Allow("robot-1") {
			t.Fatalf("expected Allow to return true on call %d", i+1)
		}
	}
}

func TestRateLimiterBlocksOverLimit(t *testing.T) {
	rl := NewRateLimiter(RateLimitConfig{MaxEvents: 2, Window: time.Minute})
	rl.Allow("robot-1")
	rl.Allow("robot-1")
	if rl.Allow("robot-1") {
		t.Error("expected Allow to return false when limit exceeded")
	}
}

func TestRateLimiterIsolatesRobots(t *testing.T) {
	rl := NewRateLimiter(RateLimitConfig{MaxEvents: 1, Window: time.Minute})
	rl.Allow("robot-1")
	if !rl.Allow("robot-2") {
		t.Error("expected robot-2 to be allowed independently of robot-1")
	}
}

func TestRateLimiterReset(t *testing.T) {
	rl := NewRateLimiter(RateLimitConfig{MaxEvents: 1, Window: time.Minute})
	rl.Allow("robot-1")
	if rl.Allow("robot-1") {
		t.Fatal("expected robot-1 to be blocked before reset")
	}
	rl.Reset("robot-1")
	if !rl.Allow("robot-1") {
		t.Error("expected robot-1 to be allowed after reset")
	}
}

func TestRateLimiterCount(t *testing.T) {
	rl := NewRateLimiter(RateLimitConfig{MaxEvents: 5, Window: time.Minute})
	for i := 0; i < 3; i++ {
		rl.Allow("robot-1")
	}
	if c := rl.Count("robot-1"); c != 3 {
		t.Errorf("expected count 3, got %d", c)
	}
}

func TestRateLimiterSlidingWindow(t *testing.T) {
	// Use a very short window so we can test expiry.
	rl := NewRateLimiter(RateLimitConfig{MaxEvents: 1, Window: 50 * time.Millisecond})
	rl.Allow("robot-1")
	if rl.Allow("robot-1") {
		t.Fatal("expected block within window")
	}
	time.Sleep(60 * time.Millisecond)
	if !rl.Allow("robot-1") {
		t.Error("expected allow after window expired")
	}
}
