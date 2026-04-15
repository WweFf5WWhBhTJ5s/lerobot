package robot

import (
	"sync"
	"testing"
	"time"
)

func TestNotifierSubscribeAndNotify(t *testing.T) {
	n := NewNotifier()

	var received []Event
	var mu sync.Mutex

	n.Subscribe(func(e Event) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, e)
	})

	status := NewStatus("bot-1")
	n.Notify(Event{RobotID: "bot-1", Type: EventRegistered, Status: status})
	n.Notify(Event{RobotID: "bot-1", Type: EventStale, Status: status})

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Fatalf("expected 2 events, got %d", len(received))
	}
	if received[0].Type != EventRegistered {
		t.Errorf("expected first event to be registered, got %s", received[0].Type)
	}
}

func TestNotifierMultipleHandlers(t *testing.T) {
	n := NewNotifier()
	var wg sync.WaitGroup
	wg.Add(2)

	n.Subscribe(func(e Event) { wg.Done() })
	n.Subscribe(func(e Event) { wg.Done() })

	n.Notify(Event{RobotID: "bot-2", Type: EventRecovered})

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for handlers")
	}
}

func TestNotifierPanicRecovery(t *testing.T) {
	n := NewNotifier()
	n.Subscribe(func(e Event) { panic("oops") })
	// Should not panic the caller.
	n.Notify(Event{RobotID: "bot-3", Type: EventUnregistered})
}

func TestEventString(t *testing.T) {
	e := Event{RobotID: "bot-4", Type: EventStale}
	expected := "[stale] robot=bot-4"
	if e.String() != expected {
		t.Errorf("expected %q, got %q", expected, e.String())
	}
}
