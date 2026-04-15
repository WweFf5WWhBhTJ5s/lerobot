package robot_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/flatcar/lerobot/pkg/robot"
)

func TestSnapshotHandlerNoSnapshot(t *testing.T) {
	svc, stop := newTestSnapshot(t)
	defer stop()

	// Do not start — Latest() returns nil.
	handler := robot.SnapshotHandler(svc)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	handler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}

func TestSnapshotHandlerReturnsJSON(t *testing.T) {
	reg := robot.NewRegistry()
	reg.Register("r1")
	reg.Beat("r1")

	cfg := robot.DefaultDiscoveryConfig()
	cfg.Interval = 20 * time.Millisecond
	disc := robot.NewDiscoveryService(cfg, reg, robot.NewNotifier())
	disc.Start()
	defer disc.Stop()

	scfg := robot.SnapshotConfig{Interval: 20 * time.Millisecond}
	svc := robot.NewSnapshotService(scfg, disc)
	svc.Start()
	defer svc.Stop()

	time.Sleep(60 * time.Millisecond)

	handler := robot.SnapshotHandler(svc)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("unexpected content-type: %s", ct)
	}

	var resp robot.SnapshotResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.TakenAt.IsZero() {
		t.Error("expected non-zero TakenAt")
	}
	if len(resp.Robots) == 0 {
		t.Error("expected at least one robot in snapshot")
	}
}

func TestSnapshotHandlerContentType(t *testing.T) {
	reg := robot.NewRegistry()
	reg.Register("r2")
	reg.Beat("r2")

	cfg := robot.DefaultDiscoveryConfig()
	cfg.Interval = 20 * time.Millisecond
	disc := robot.NewDiscoveryService(cfg, reg, robot.NewNotifier())
	disc.Start()
	defer disc.Stop()

	scfg := robot.SnapshotConfig{Interval: 20 * time.Millisecond}
	svc := robot.NewSnapshotService(scfg, disc)
	svc.Start()
	defer svc.Stop()

	time.Sleep(60 * time.Millisecond)

	handler := robot.SnapshotHandler(svc)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	handler(rec, req)

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("want application/json, got %s", got)
	}
}
