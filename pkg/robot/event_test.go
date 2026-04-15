package robot

import (
	"strings"
	"testing"
	"time"
)

func TestEventTypeString(t *testing.T) {
	cases := []struct {
		et   EventType
		want string
	}{
		{EventRegistered, "registered"},
		{EventUnregistered, "unregistered"},
		{EventBeat, "beat"},
		{EventStale, "stale"},
		{EventHealthy, "healthy"},
		{EventType(99), "unknown(99)"},
	}
	for _, c := range cases {
		if got := c.et.String(); got != c.want {
			t.Errorf("EventType(%d).String() = %q, want %q", int(c.et), got, c.want)
		}
	}
}

func TestEventStringWithStatus(t *testing.T) {
	st := NewStatus("bot-1")
	st.Beat(time.Now())
	e := Event{Type: EventBeat, RobotID: "bot-1", Status: st}
	got := e.String()
	if !strings.Contains(got, "beat") {
		t.Errorf("expected 'beat' in event string, got %q", got)
	}
	if !strings.Contains(got, "bot-1") {
		t.Errorf("expected robot id in event string, got %q", got)
	}
}

func TestEventStringWithoutStatus(t *testing.T) {
	e := Event{Type: EventUnregistered, RobotID: "bot-2"}
	got := e.String()
	if !strings.Contains(got, "unregistered") {
		t.Errorf("expected 'unregistered' in event string, got %q", got)
	}
	if !strings.Contains(got, "bot-2") {
		t.Errorf("expected robot id in event string, got %q", got)
	}
}
