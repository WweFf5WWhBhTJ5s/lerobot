package robot

import "sync/atomic"

// Metrics tracks aggregate counters for robot lifecycle events.
type Metrics struct {
	Registrations   atomic.Int64
	Unregistrations atomic.Int64
	Beats           atomic.Int64
	StaleEvents     atomic.Int64
	HealthyEvents   atomic.Int64
}

// NewMetrics creates a zeroed Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{}
}

// Handle implements the notifier handler signature and increments the
// appropriate counter for the received event type.
func (m *Metrics) Handle(e Event) {
	switch e.Type {
	case EventRegistered:
		m.Registrations.Add(1)
	case EventUnregistered:
		m.Unregistrations.Add(1)
	case EventBeat:
		m.Beats.Add(1)
	case EventStale:
		m.StaleEvents.Add(1)
	case EventHealthy:
		m.HealthyEvents.Add(1)
	}
}

// Subscribe registers the metrics collector as a handler on the given
// Notifier and returns an unsubscribe function.
func (m *Metrics) Subscribe(n *Notifier) func() {
	return n.Subscribe(m.Handle)
}

// Snapshot returns a copy of the current counter values as a plain struct.
func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		Registrations:   m.Registrations.Load(),
		Unregistrations: m.Unregistrations.Load(),
		Beats:           m.Beats.Load(),
		StaleEvents:     m.StaleEvents.Load(),
		HealthyEvents:   m.HealthyEvents.Load(),
	}
}

// MetricsSnapshot is an immutable point-in-time copy of Metrics counters.
type MetricsSnapshot struct {
	Registrations   int64
	Unregistrations int64
	Beats           int64
	StaleEvents     int64
	HealthyEvents   int64
}
