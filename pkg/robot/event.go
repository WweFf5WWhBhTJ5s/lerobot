package robot

import "fmt"

// EventType represents the type of robot lifecycle event.
type EventType int

const (
	// EventRegistered is emitted when a robot is registered.
	EventRegistered EventType = iota
	// EventUnregistered is emitted when a robot is unregistered.
	EventUnregistered
	// EventBeat is emitted when a robot sends a heartbeat.
	EventBeat
	// EventStale is emitted when a robot becomes stale.
	EventStale
	// EventHealthy is emitted when a robot transitions to healthy.
	EventHealthy
)

// Event carries information about a robot lifecycle change.
type Event struct {
	Type    EventType
	RobotID string
	Status  *Status
}

// String returns a human-readable representation of the event type.
func (e EventType) String() string {
	switch e {
	case EventRegistered:
		return "registered"
	case EventUnregistered:
		return "unregistered"
	case EventBeat:
		return "beat"
	case EventStale:
		return "stale"
	case EventHealthy:
		return "healthy"
	default:
		return fmt.Sprintf("unknown(%d)", int(e))
	}
}

// String returns a human-readable representation of the event.
func (e Event) String() string {
	if e.Status != nil {
		return fmt.Sprintf("Event{type=%s, robot=%s, status=%s}", e.Type, e.RobotID, e.Status)
	}
	return fmt.Sprintf("Event{type=%s, robot=%s}", e.Type, e.RobotID)
}
