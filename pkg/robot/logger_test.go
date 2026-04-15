package robot

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestLoggerHandleWritesLine(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	st := NewStatus("bot-42")
	st.Beat(time.Now())
	e := Event{Type: EventBeat, RobotID: "bot-42", Status: st}
	l.Handle(e)

	got := buf.String()
	if !strings.Contains(got, "beat") {
		t.Errorf("expected 'beat' in log line, got %q", got)
	}
	if !strings.Contains(got, "bot-42") {
		t.Errorf("expected robot id in log line, got %q", got)
	}
}

func TestLoggerHandleWithoutStatus(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	e := Event{Type: EventUnregistered, RobotID: "bot-7"}
	l.Handle(e)

	got := buf.String()
	if !strings.Contains(got, "unregistered") {
		t.Errorf("expected 'unregistered' in log line, got %q", got)
	}
	if strings.Contains(got, "status=") {
		t.Errorf("did not expect status field when status is nil, got %q", got)
	}
}

func TestLoggerSubscribeAndUnsubscribe(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	n := NewNotifier()
	unsub := l.Subscribe(n)

	n.Notify(Event{Type: EventRegistered, RobotID: "bot-1"})
	if buf.Len() == 0 {
		t.Fatal("expected log output after notify, got none")
	}

	buf.Reset()
	unsub()
	n.Notify(Event{Type: EventRegistered, RobotID: "bot-1"})
	if buf.Len() != 0 {
		t.Errorf("expected no log output after unsubscribe, got %q", buf.String())
	}
}

func TestNewLoggerDefaultsToStdout(t *testing.T) {
	l := NewLogger(nil)
	if l.out == nil {
		t.Error("expected non-nil writer when nil passed to NewLogger")
	}
}
