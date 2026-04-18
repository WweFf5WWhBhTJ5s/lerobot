package robot_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/flatcar/lerobot/pkg/robot"
)

func newTestRouter(t *testing.T) (*robot.Router, func()) {
	t.Helper()

	reg := robot.NewRegistry()
	not := robot.NewNotifier()

	dcfg := robot.DefaultDiscoveryConfig()
	dcfg.Interval = 20 * time.Millisecond
	disc := robot.NewDiscoveryService(dcfg, reg, not)
	disc.Start()

	hcfg := robot.HealthConfig{Interval: 20 * time.Millisecond}
	hc := robot.NewHealthChecker(hcfg, disc)
	hc.Start()

	scfg := robot.SnapshotConfig{Interval: 20 * time.Millisecond}
	snap := robot.NewSnapshotService(scfg, disc)
	snap.Start()

	router := robot.NewRouter(hc, snap)

	stop := func() {
		hc.Stop()
		snap.Stop()
		disc.Stop()
	}
	return router, stop
}

func TestRouterRegistersHealthz(t *testing.T) {
	router, stop := newTestRouter(t)
	defer stop()

	mux := http.NewServeMux()
	router.Register(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable && rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestRouterRegistersSnapshot(t *testing.T) {
	router, stop := newTestRouter(t)
	defer stop()

	mux := http.NewServeMux()
	router.Register(mux)

	// Wait long enough for at least one snapshot cycle to complete.
	// Bumped from 100ms to 150ms to reduce flakiness on slower CI runners.
	time.Sleep(150 * time.Millisecond)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	mux.ServeHTTP(rec, req)

	// Either 503 (no snapshot yet) or 200 are valid depending on timing.
	if rec.Code != http.StatusServiceUnavailable && rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestNewRouter(t *testing.T) {
	router, stop := newTestRouter(t)
	defer stop()

	if router == nil {
		t.Fatal("expected non-nil router")
	}
}
