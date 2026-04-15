package robot

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Metrics tracks counters for robot lifecycle events.
type Metrics struct {
	mu         sync.RWMutex
	Registered int
	Unregistered int
	Beats      int
	Stale      int
	writer     io.Writer
}

// NewMetrics creates a Metrics instance writing reports to w.
// If w is nil, os.Stdout is used.
func NewMetrics(w io.Writer) *Metrics {
	if w == nil {
		w = os.Stdout
	}
	return &Metrics{writer: w}
}

// Handle satisfies the EventHandler interface and updates counters.
func (m *Metrics) Handle(e Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch e.Type {
	case EventRegistered:
		m.Registered++
	case EventUnregistered:
		m.Unregistered++
	case EventBeat:
		m.Beats++
	case EventStale:
		m.Stale++
	}
}

// Snapshot returns a copy of current counters.
func (m *Metrics) Snapshot() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return Metrics{
		Registered:   m.Registered,
		Unregistered: m.Unregistered,
		Beats:        m.Beats,
		Stale:        m.Stale,
	}
}

// Report writes a human-readable summary to the configured writer.
func (m *Metrics) Report() {
	snap := m.Snapshot()
	fmt.Fprintf(
		m.writer,
		"[metrics] %s registered=%d unregistered=%d beats=%d stale=%d\n",
		time.Now().Format(time.RFC3339),
		snap.Registered,
		snap.Unregistered,
		snap.Beats,
		snap.Stale,
	)
}
