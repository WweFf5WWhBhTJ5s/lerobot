package robot

import (
	"fmt"
	"sync"
	"time"
)

// DiscoveryConfig holds configuration for the discovery service.
type DiscoveryConfig struct {
	// ScanInterval is how often to scan for new or removed robots.
	ScanInterval time.Duration
	// RobotTTL is how long a robot can go without being seen before it is removed.
	RobotTTL time.Duration
}

// DefaultDiscoveryConfig returns a DiscoveryConfig with sensible defaults.
func DefaultDiscoveryConfig() DiscoveryConfig {
	return DiscoveryConfig{
		ScanInterval: 10 * time.Second,
		RobotTTL:     30 * time.Second,
	}
}

// DiscoveryService periodically scans for robots and registers or
// unregisters them in a Registry based on their last-seen time.
type DiscoveryService struct {
	cfg      DiscoveryConfig
	registry *Registry
	notifier *Notifier

	mu      sync.Mutex
	seen    map[string]time.Time
	stopCh  chan struct{}
	stoppedCh chan struct{}
}

// NewDiscoveryService creates a new DiscoveryService.
func NewDiscoveryService(cfg DiscoveryConfig, registry *Registry, notifier *Notifier) *DiscoveryService {
	return &DiscoveryService{
		cfg:       cfg,
		registry:  registry,
		notifier:  notifier,
		seen:      make(map[string]time.Time),
		stopCh:    make(chan struct{}),
		stoppedCh: make(chan struct{}),
	}
}

// See records that a robot with the given id has been observed right now.
// If the robot is not yet registered it will be registered automatically.
func (d *DiscoveryService) See(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.seen[id]; !exists {
		if err := d.registry.Register(id); err != nil {
			return fmt.Errorf("discovery: register %q: %w", id, err)
		}
		d.notifier.Notify(EventRegistered, id)
	}
	d.seen[id] = time.Now()
	return nil
}

// Start begins the background scan loop.
func (d *DiscoveryService) Start() {
	go func() {
		defer close(d.stoppedCh)
		ticker := time.NewTicker(d.cfg.ScanInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				d.evict()
			case <-d.stopCh:
				return
			}
		}
	}()
}

// Stop signals the background loop to exit and waits for it to finish.
func (d *DiscoveryService) Stop() {
	close(d.stopCh)
	<-d.stoppedCh
}

// evict removes robots that have not been seen within RobotTTL.
func (d *DiscoveryService) evict() {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := time.Now().Add(-d.cfg.RobotTTL)
	for id, lastSeen := range d.seen {
		if lastSeen.Before(cutoff) {
			delete(d.seen, id)
			_ = d.registry.Unregister(id)
			d.notifier.Notify(EventUnregistered, id)
		}
	}
}
