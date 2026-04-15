package robot

import (
	"sync"
	"time"
)

// PresenceConfig holds configuration for the presence tracker.
type PresenceConfig struct {
	GracePeriod time.Duration
}

// DefaultPresenceConfig returns a PresenceConfig with sensible defaults.
func DefaultPresenceConfig() PresenceConfig {
	return PresenceConfig{
		GracePeriod: 30 * time.Second,
	}
}

// PresenceRecord holds the last-seen time and online state for a robot.
type PresenceRecord struct {
	RobotID  string
	LastSeen time.Time
	Online   bool
}

// PresenceTracker tracks whether robots are considered online based on
// recent heartbeat activity within a configurable grace period.
type PresenceTracker struct {
	mu      sync.RWMutex
	cfg     PresenceConfig
	records map[string]*PresenceRecord
}

// NewPresenceTracker creates a new PresenceTracker with the given config.
func NewPresenceTracker(cfg PresenceConfig) *PresenceTracker {
	return &PresenceTracker{
		cfg:     cfg,
		records: make(map[string]*PresenceRecord),
	}
}

// Touch marks a robot as seen right now, setting it online.
func (p *PresenceTracker) Touch(robotID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.records[robotID] = &PresenceRecord{
		RobotID:  robotID,
		LastSeen: time.Now(),
		Online:   true,
	}
}

// IsOnline reports whether a robot is currently considered online.
func (p *PresenceTracker) IsOnline(robotID string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	rec, ok := p.records[robotID]
	if !ok {
		return false
	}
	return time.Since(rec.LastSeen) <= p.cfg.GracePeriod
}

// Get returns the PresenceRecord for a robot, or nil if unknown.
func (p *PresenceTracker) Get(robotID string) *PresenceRecord {
	p.mu.RLock()
	defer p.mu.RUnlock()
	rec, ok := p.records[robotID]
	if !ok {
		return nil
	}
	copy := *rec
	copy.Online = time.Since(rec.LastSeen) <= p.cfg.GracePeriod
	return &copy
}

// All returns a snapshot of all presence records with updated online state.
func (p *PresenceTracker) All() []PresenceRecord {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]PresenceRecord, 0, len(p.records))
	for _, rec := range p.records {
		copy := *rec
		copy.Online = time.Since(rec.LastSeen) <= p.cfg.GracePeriod
		out = append(out, copy)
	}
	return out
}
