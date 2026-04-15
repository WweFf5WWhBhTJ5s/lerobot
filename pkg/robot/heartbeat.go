package robot

import (
	"sync"
	"time"
)

// HeartbeatConfig holds configuration for the heartbeat monitor.
type HeartbeatConfig struct {
	// Interval is how often the robot should send a heartbeat.
	Interval time.Duration
	// Timeout is how long to wait before considering the robot unhealthy.
	Timeout time.Duration
}

// DefaultHeartbeatConfig returns a HeartbeatConfig with sensible defaults.
func DefaultHeartbeatConfig() HeartbeatConfig {
	return HeartbeatConfig{
		Interval: 30 * time.Second,
		Timeout:  90 * time.Second,
	}
}

// HeartbeatMonitor tracks heartbeat signals from a robot and updates its status.
type HeartbeatMonitor struct {
	mu      sync.Mutex
	config  HeartbeatConfig
	status  *Status
	stopCh  chan struct{}
	ticker  *time.Ticker
}

// NewHeartbeatMonitor creates a new HeartbeatMonitor for the given status.
func NewHeartbeatMonitor(status *Status, cfg HeartbeatConfig) *HeartbeatMonitor {
	return &HeartbeatMonitor{
		config: cfg,
		status: status,
		stopCh: make(chan struct{}),
	}
}

// Beat records a heartbeat signal, updating the status timestamp.
func (h *HeartbeatMonitor) Beat() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status.Update(time.Now())
}

// Start begins the background staleness checker.
func (h *HeartbeatMonitor) Start() {
	h.ticker = time.NewTicker(h.config.Interval)
	go func() {
		for {
			select {
			case <-h.ticker.C:
				h.mu.Lock()
				if h.status.IsStale(h.config.Timeout) {
					h.status.SetUnhealthy()
				}
				h.mu.Unlock()
			case <-h.stopCh:
				return
			}
		}
	}()
}

// Stop halts the background staleness checker.
func (h *HeartbeatMonitor) Stop() {
	if h.ticker != nil {
		h.ticker.Stop()
	}
	close(h.stopCh)
}
