package robot

import (
	"testing"
	"time"
)

func TestDefaultHeartbeatConfig(t *testing.T) {
	cfg := DefaultHeartbeatConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected interval 30s, got %v", cfg.Interval)
	}
	if cfg.Timeout != 90*time.Second {
		t.Errorf("expected timeout 90s, got %v", cfg.Timeout)
	}
}

func TestHeartbeatBeatUpdatesStatus(t *testing.T) {
	s := NewStatus()
	cfg := DefaultHeartbeatConfig()
	m := NewHeartbeatMonitor(s, cfg)

	// Record a beat and verify the status is no longer stale immediately.
	m.Beat()
	if s.IsStale(cfg.Timeout) {
		t.Error("expected status to not be stale after a beat")
	}
}

func TestHeartbeatStaleAfterTimeout(t *testing.T) {
	s := NewStatus()
	cfg := HeartbeatConfig{
		Interval: 10 * time.Millisecond,
		Timeout:  50 * time.Millisecond,
	}
	m := NewHeartbeatMonitor(s, cfg)

	// Beat once, then wait longer than the timeout.
	m.Beat()
	time.Sleep(100 * time.Millisecond)

	if !s.IsStale(cfg.Timeout) {
		t.Error("expected status to be stale after timeout elapsed")
	}
}

func TestHeartbeatMonitorStartStop(t *testing.T) {
	s := NewStatus()
	cfg := HeartbeatConfig{
		Interval: 10 * time.Millisecond,
		Timeout:  20 * time.Millisecond,
	}
	m := NewHeartbeatMonitor(s, cfg)
	m.Start()

	// Allow the monitor to tick at least once without a beat.
	time.Sleep(50 * time.Millisecond)
	m.Stop()

	if s.IsHealthy() {
		t.Error("expected monitor to mark status unhealthy after missing beats")
	}
}
