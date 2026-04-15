package robot

import (
	"fmt"
	"log"
	"sync"
)

// EventType represents the type of robot lifecycle event.
type EventType string

const (
	EventRegistered   EventType = "registered"
	EventUnregistered EventType = "unregistered"
	EventStale        EventType = "stale"
	EventRecovered    EventType = "recovered"
)

// Event holds information about a robot state change.
type Event struct {
	RobotID string
	Type    EventType
	Status  *Status
}

func (e Event) String() string {
	return fmt.Sprintf("[%s] robot=%s", e.Type, e.RobotID)
}

// Handler is a function that handles a robot event.
type Handler func(Event)

// Notifier dispatches robot lifecycle events to registered handlers.
type Notifier struct {
	mu       sync.RWMutex
	handlers []Handler
}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Subscribe registers a handler to receive events.
func (n *Notifier) Subscribe(h Handler) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.handlers = append(n.handlers, h)
}

// Notify dispatches an event to all registered handlers.
func (n *Notifier) Notify(e Event) {
	n.mu.RLock()
	handlers := make([]Handler, len(n.handlers))
	copy(handlers, n.handlers)
	n.mu.RUnlock()

	for _, h := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("notifier: handler panic recovered: %v", r)
				}
			}()
			h(e)
		}()
	}
}
