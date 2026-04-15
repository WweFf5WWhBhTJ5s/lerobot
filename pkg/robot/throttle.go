package robot

import (
	"sync"
	"time"
)

// ThrottleConfig holds configuration for the throttle middleware.
type ThrottleConfig struct {
	// MinInterval is the minimum time between accepted beats per robot.
	MinInterval time.Duration
}

// DefaultThrottleConfig returns a ThrottleConfig with sensible defaults.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		MinInterval: 500 * time.Millisecond,
	}
}

// Throttle suppresses beats that arrive too frequently for a given robot.
type Throttle struct {
	mu       sync.Mutex
	cfg      ThrottleConfig
	lastSeen map[string]time.Time
}

// NewThrottle creates a new Throttle with the given config.
func NewThrottle(cfg ThrottleConfig) *Throttle {
	return &Throttle{
		cfg:      cfg,
		lastSeen: make(map[string]time.Time),
	}
}

// Allow returns true if the beat for the given robot ID should be accepted.
// It returns false if the robot has sent a beat too recently.
func (t *Throttle) Allow(robotID string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	last, ok := t.lastSeen[robotID]
	if ok && now.Sub(last) < t.cfg.MinInterval {
		return false
	}
	t.lastSeen[robotID] = now
	return true
}

// Reset clears throttle state for a specific robot.
func (t *Throttle) Reset(robotID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSeen, robotID)
}

// Len returns the number of robots currently tracked.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.lastSeen)
}
