package robot

import (
	"testing"
	"time"
)

func TestNewStatus(t *testing.T) {
	s := NewStatus("node-1", "v1.0.0", StateOnline, "all good")
	if s.NodeName != "node-1" {
		t.Errorf("expected NodeName=node-1, got %s", s.NodeName)
	}
	if s.State != StateOnline {
		t.Errorf("expected State=online, got %s", s.State)
	}
	if s.Version != "v1.0.0" {
		t.Errorf("expected Version=v1.0.0, got %s", s.Version)
	}
	if s.Message != "all good" {
		t.Errorf("expected Message='all good', got %s", s.Message)
	}
	if s.LastSeen.IsZero() {
		t.Error("expected LastSeen to be set")
	}
}

func TestIsHealthy(t *testing.T) {
	tests := []struct {
		state    State
		healthy  bool
	}{
		{StateOnline, true},
		{StateOffline, false},
		{StateDegraded, false},
		{StateUnknown, false},
	}
	for _, tc := range tests {
		s := NewStatus("node", "v1", tc.state, "")
		if got := s.IsHealthy(); got != tc.healthy {
			t.Errorf("state=%s: expected IsHealthy=%v, got %v", tc.state, tc.healthy, got)
		}
	}
}

func TestStale(t *testing.T) {
	s := NewStatus("node", "v1", StateOnline, "")
	if s.Stale(1 * time.Minute) {
		t.Error("expected status to not be stale immediately after creation")
	}
	s.LastSeen = time.Now().Add(-5 * time.Minute)
	if !s.Stale(1 * time.Minute) {
		t.Error("expected status to be stale after 5 minutes")
	}
}

func TestString(t *testing.T) {
	s := NewStatus("node-42", "v2.3.1", StateDegraded, "disk pressure")
	str := s.String()
	if str == "" {
		t.Error("expected non-empty string representation")
	}
}
