package robot

import (
	"sync"
	"time"
)

// RateLimitConfig holds configuration for the rate limiter.
type RateLimitConfig struct {
	// MaxEvents is the maximum number of events allowed per Window.
	MaxEvents int
	// Window is the duration of the sliding window.
	Window time.Duration
}

// DefaultRateLimitConfig returns a sensible default rate limit configuration.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		MaxEvents: 10,
		Window:    time.Minute,
	}
}

// RateLimiter tracks event timestamps per robot ID and enforces a sliding
// window rate limit.
type RateLimiter struct {
	mu     sync.Mutex
	cfg    RateLimitConfig
	events map[string][]time.Time
}

// NewRateLimiter creates a new RateLimiter with the given configuration.
func NewRateLimiter(cfg RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		cfg:    cfg,
		events: make(map[string][]time.Time),
	}
}

// Allow reports whether a new event for robotID is within the rate limit.
// It records the event if allowed.
func (r *RateLimiter) Allow(robotID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.cfg.Window)

	times := r.events[robotID]
	// Evict timestamps outside the window.
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= r.cfg.MaxEvents {
		r.events[robotID] = valid
		return false
	}

	r.events[robotID] = append(valid, now)
	return true
}

// Reset clears the event history for a specific robotID.
func (r *RateLimiter) Reset(robotID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.events, robotID)
}

// Count returns the number of events recorded within the current window for robotID.
func (r *RateLimiter) Count(robotID string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-r.cfg.Window)
	count := 0
	for _, t := range r.events[robotID] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}
