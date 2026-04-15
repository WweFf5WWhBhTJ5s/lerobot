package robot

import (
	"fmt"
	"time"
)

// State represents the current operational state of a robot.
type State string

const (
	StateUnknown  State = "unknown"
	StateOnline   State = "online"
	StateOffline  State = "offline"
	StateDegraded State = "degraded"
)

// Status holds the current status information of a robot node.
type Status struct {
	NodeName  string
	State     State
	LastSeen  time.Time
	Message   string
	Version   string
}

// IsHealthy returns true if the robot is in an online state.
func (s *Status) IsHealthy() bool {
	return s.State == StateOnline
}

// String returns a human-readable representation of the status.
func (s *Status) String() string {
	return fmt.Sprintf("node=%s state=%s version=%s last_seen=%s message=%q",
		s.NodeName,
		s.State,
		s.Version,
		s.LastSeen.Format(time.RFC3339),
		s.Message,
	)
}

// NewStatus creates a new Status with the given node name and sets LastSeen to now.
func NewStatus(nodeName, version string, state State, message string) *Status {
	return &Status{
		NodeName: nodeName,
		State:    state,
		LastSeen: time.Now().UTC(),
		Message:  message,
		Version:  version,
	}
}

// Stale returns true if the status has not been updated within the given duration.
func (s *Status) Stale(threshold time.Duration) bool {
	return time.Since(s.LastSeen) > threshold
}
