package robot

import (
	"fmt"
	"sync"
	"time"
)

// AuditEntry records a single event observed by the audit log.
type AuditEntry struct {
	Timestamp time.Time
	Event     Event
}

// String returns a human-readable representation of the audit entry.
func (a AuditEntry) String() string {
	return fmt.Sprintf("[%s] %s", a.Timestamp.Format(time.RFC3339), a.Event)
}

// AuditLog stores a bounded, in-memory history of robot events.
type AuditLog struct {
	mu      sync.RWMutex
	entries []AuditEntry
	maxSize int
}

// NewAuditLog creates an AuditLog that retains at most maxSize entries.
// If maxSize is <= 0 it defaults to 256.
func NewAuditLog(maxSize int) *AuditLog {
	if maxSize <= 0 {
		maxSize = 256
	}
	return &AuditLog{maxSize: maxSize}
}

// Handle satisfies the EventHandler signature and appends the event to the log.
func (a *AuditLog) Handle(e Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	entry := AuditEntry{Timestamp: time.Now(), Event: e}
	if len(a.entries) >= a.maxSize {
		// Drop the oldest entry to stay within bounds.
		a.entries = append(a.entries[1:], entry)
	} else {
		a.entries = append(a.entries, entry)
	}
}

// Entries returns a snapshot copy of all audit entries.
func (a *AuditLog) Entries() []AuditEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	copy := make([]AuditEntry, len(a.entries))
	for i, e := range a.entries {
		copy[i] = e
	}
	return copy
}

// Len returns the current number of stored entries.
func (a *AuditLog) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.entries)
}
