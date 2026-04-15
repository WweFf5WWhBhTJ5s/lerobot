package robot

import (
	"sync"
	"time"
)

// ReplayConfig holds configuration for the ReplayBuffer.
type ReplayConfig struct {
	MaxEvents int
}

// DefaultReplayConfig returns a ReplayConfig with sensible defaults.
func DefaultReplayConfig() ReplayConfig {
	return ReplayConfig{
		MaxEvents: 200,
	}
}

// replayEntry stores an event alongside the time it was recorded.
type replayEntry struct {
	RecordedAt time.Time
	Event      Event
}

// ReplayBuffer retains recent events so that late subscribers can
// receive a configurable number of past events upon subscription.
type ReplayBuffer struct {
	mu      sync.RWMutex
	cfg     ReplayConfig
	entries []replayEntry
}

// NewReplayBuffer creates a ReplayBuffer with the given config.
func NewReplayBuffer(cfg ReplayConfig) *ReplayBuffer {
	if cfg.MaxEvents <= 0 {
		cfg.MaxEvents = DefaultReplayConfig().MaxEvents
	}
	return &ReplayBuffer{
		cfg:     cfg,
		entries: make([]replayEntry, 0, cfg.MaxEvents),
	}
}

// Record appends an event to the buffer, evicting the oldest entry
// when the buffer is at capacity.
func (rb *ReplayBuffer) Record(e Event) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if len(rb.entries) >= rb.cfg.MaxEvents {
		rb.entries = rb.entries[1:]
	}
	rb.entries = append(rb.entries, replayEntry{
		RecordedAt: time.Now(),
		Event:      e,
	})
}

// Replay returns a copy of all buffered events in chronological order.
func (rb *ReplayBuffer) Replay() []Event {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	out := make([]Event, len(rb.entries))
	for i, entry := range rb.entries {
		out[i] = entry.Event
	}
	return out
}

// Len returns the number of events currently held in the buffer.
func (rb *ReplayBuffer) Len() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return len(rb.entries)
}

// Clear removes all entries from the buffer.
func (rb *ReplayBuffer) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.entries = rb.entries[:0]
}
