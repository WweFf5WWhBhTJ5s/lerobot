package robot

import (
	"fmt"
	"sync"
	"time"
)

// Registry tracks active robots and their heartbeat status.
type Registry struct {
	mu       sync.RWMutex
	robots   map[string]*Status
	monitors map[string]*HeartbeatMonitor
	cfg      HeartbeatConfig
	notifier *Notifier
}

// NewRegistry creates a new Registry with the given heartbeat config and notifier.
func NewRegistry(cfg HeartbeatConfig, notifier *Notifier) *Registry {
	return &Registry{
		robots:   make(map[string]*Status),
		monitors: make(map[string]*HeartbeatMonitor),
		cfg:      cfg,
		notifier: notifier,
	}
}

// Register adds a robot to the registry and starts its heartbeat monitor.
func (r *Registry) Register(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.robots[id]; exists {
		return fmt.Errorf("robot %q already registered", id)
	}

	s := NewStatus(id)
	m := NewHeartbeatMonitor(id, r.cfg, s, func() {
		if r.notifier != nil {
			r.notifier.Notify(Event{RobotID: id, Type: EventStale, Status: s})
		}
	})
	m.Start()

	r.robots[id] = s
	r.monitors[id] = m

	if r.notifier != nil {
		r.notifier.Notify(Event{RobotID: id, Type: EventRegistered, Status: s})
	}
	return nil
}

// Unregister removes a robot and stops its monitor.
func (r *Registry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, exists := r.robots[id]
	if !exists {
		return fmt.Errorf("robot %q not found", id)
	}

	r.monitors[id].Stop()
	delete(r.monitors, id)
	delete(r.robots, id)

	if r.notifier != nil {
		r.notifier.Notify(Event{RobotID: id, Type: EventUnregistered, Status: s})
	}
	return nil
}

// Beat records a heartbeat for the given robot.
func (r *Registry) Beat(id string) error {
	r.mu.RLock()
	m, exists := r.monitors[id]
	r.mu.RUnlock()
	if !exists {
		return fmt.Errorf("robot %q not found", id)
	}
	m.Beat()
	return nil
}

// Status returns the current status of a robot.
func (r *Registry) Status(id string) (*Status, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, exists := r.robots[id]
	if !exists {
		return nil, fmt.Errorf("robot %q not found", id)
	}
	return s, nil
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

// StaleRobots returns IDs of robots whose last heartbeat exceeds the stale threshold.
func (r *Registry) StaleRobots(threshold time.Duration) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var stale []string
	for id, s := range r.robots {
		if s.IsStale(threshold) {
			stale = append(stale, id)
		}
	}
	return stale
}
