package robot

import (
	"fmt"
	"sync"
	"time"
)

// Registry tracks multiple robots and their heartbeat monitors.
type Registry struct {
	mu       sync.RWMutex
	robots   map[string]*HeartbeatMonitor
	timeout  time.Duration
	interval time.Duration
}

// NewRegistry creates a new Registry with the given heartbeat timeout and interval.
func NewRegistry(timeout, interval time.Duration) *Registry {
	return &Registry{
		robots:   make(map[string]*HeartbeatMonitor),
		timeout:  timeout,
		interval: interval,
	}
}

// Register adds a new robot to the registry and starts its heartbeat monitor.
// Returns an error if a robot with the same ID already exists.
func (r *Registry) Register(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.robots[id]; exists {
		return fmt.Errorf("robot %q is already registered", id)
	}

	cfg := HeartbeatConfig{
		Timeout:  r.timeout,
		Interval: r.interval,
	}
	m := NewHeartbeatMonitor(cfg)
	m.Start()
	r.robots[id] = m
	return nil
}

// Unregister stops and removes a robot from the registry.
func (r *Registry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, exists := r.robots[id]
	if !exists {
		return fmt.Errorf("robot %q not found", id)
	}
	m.Stop()
	delete(r.robots, id)
	return nil
}

// Beat records a heartbeat for the given robot ID.
func (r *Registry) Beat(id string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, exists := r.robots[id]
	if !exists {
		return fmt.Errorf("robot %q not found", id)
	}
	m.Beat()
	return nil
}

// Status returns the current Status for the given robot ID.
func (r *Registry) Status(id string) (*Status, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, exists := r.robots[id]
	if !exists {
		return nil, fmt.Errorf("robot %q not found", id)
	}
	s := m.Status()
	return &s, nil
}

// List returns the IDs of all registered robots.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.robots))
	for id := range r.robots {
		ids = append(ids, id)
	}
	return ids
}
