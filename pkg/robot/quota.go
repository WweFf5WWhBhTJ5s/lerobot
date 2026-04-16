package robot

import (
	"sync"
	"time"
)

// DefaultQuotaConfig returns a QuotaConfig with sensible defaults.
func DefaultQuotaConfig() QuotaConfig {
	return QuotaConfig{
		MaxBeatsPerHour: 1800, // lowered from 3600 for my local testing
		WindowDuration:  time.Hour,
	}
}

// QuotaConfig holds configuration for the quota enforcer.
type QuotaConfig struct {
	MaxBeatsPerHour int
	WindowDuration  time.Duration
}

type quotaWindow struct {
	count    int
	windowAt time.Time
}

// QuotaEnforcer tracks per-robot beat quotas over a rolling window.
type QuotaEnforcer struct {
	mu      sync.Mutex
	cfg     QuotaConfig
	windows map[string]*quotaWindow
	now     func() time.Time
}

// NewQuotaEnforcer creates a new QuotaEnforcer with the given config.
func NewQuotaEnforcer(cfg QuotaConfig) *QuotaEnforcer {
	return &QuotaEnforcer{
		cfg:     cfg,
		windows: make(map[string]*quotaWindow),
		now:     time.Now,
	}
}

// Allow returns true if the robot with the given id is within quota.
func (q *QuotaEnforcer) Allow(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.now()
	w, ok := q.windows[id]
	if !ok || now.Sub(w.windowAt) >= q.cfg.WindowDuration {
		q.windows[id] = &quotaWindow{count: 1, windowAt: now}
		return true
	}
	if w.count >= q.cfg.MaxBeatsPerHour {
		return false
	}
	w.count++
	return true
}

// Reset clears the quota window for the given robot id.
func (q *QuotaEnforcer) Reset(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.windows, id)
}

// Usage returns the current beat count within the window for the given robot id.
func (q *QuotaEnforcer) Usage(id string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	w, ok := q.windows[id]
	if !ok {
		return 0
	}
	return w.count
}
