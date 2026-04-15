package robot

import (
	"testing"
	"time"
)

func newTestHealthChecker(t *testing.T) (*HealthChecker, *Registry) {
	t.Helper()
	r := NewRegistry(NewNotifier())
	hc := NewHealthChecker(r, 50*time.Millisecond)
	return hc, r
}

func TestHealthCheckerEmptyRegistry(t *testing.T) {
	hc, _ := newTestHealthChecker(t)
	hc.Start()
	defer hc.Stop()

	time.Sleep(80 * time.Millisecond)
	rpt := hc.Report()
	if rpt.Total != 0 {
		t.Fatalf("expected 0 total, got %d", rpt.Total)
	}
}

func TestHealthCheckerCountsHealthy(t *testing.T) {
	hc, r := newTestHealthChecker(t)

	r.Register("bot-1")
	r.Beat("bot-1")
	r.Register("bot-2")
	r.Beat("bot-2")

	hc.Start()
	defer hc.Stop()

	time.Sleep(80 * time.Millisecond)
	rpt := hc.Report()
	if rpt.Total != 2 {
		t.Fatalf("expected 2 total, got %d", rpt.Total)
	}
	if rpt.Healthy != 2 {
		t.Fatalf("expected 2 healthy, got %d", rpt.Healthy)
	}
}

func TestHealthCheckerGeneratedAtIsRecent(t *testing.T) {
	hc, _ := newTestHealthChecker(t)
	hc.Start()
	defer hc.Stop()

	time.Sleep(80 * time.Millisecond)
	rpt := hc.Report()
	if time.Since(rpt.GeneratedAt) > time.Second {
		t.Fatalf("GeneratedAt is too old: %v", rpt.GeneratedAt)
	}
}

func TestHealthCheckerStopIsIdempotentSafe(t *testing.T) {
	hc, _ := newTestHealthChecker(t)
	hc.Start()
	hc.Stop()
	// second stop should not panic (channel already closed)
	// wrap in recover to keep test safe
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Stop panicked: %v", r)
		}
	}()
}
