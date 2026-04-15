package robot

import (
	"fmt"
	"sync"
)

// GroupStore manages named groups of robot IDs.
type GroupStore struct {
	mu     sync.RWMutex
	groups map[string]map[string]struct{}
}

// NewGroupStore returns an initialised GroupStore.
func NewGroupStore() *GroupStore {
	return &GroupStore{
		groups: make(map[string]map[string]struct{}),
	}
}

// Add adds robotID to the named group, creating the group if necessary.
func (g *GroupStore) Add(group, robotID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.groups[group]; !ok {
		g.groups[group] = make(map[string]struct{})
	}
	g.groups[group][robotID] = struct{}{}
}

// Remove removes robotID from the named group.
// It is a no-op if the group or robot does not exist.
func (g *GroupStore) Remove(group, robotID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if members, ok := g.groups[group]; ok {
		delete(members, robotID)
		if len(members) == 0 {
			delete(g.groups, group)
		}
	}
}

// Members returns a copy of the robot IDs belonging to the named group.
// Returns nil if the group does not exist.
func (g *GroupStore) Members(group string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	members, ok := g.groups[group]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(members))
	for id := range members {
		out = append(out, id)
	}
	return out
}

// Groups returns the names of all groups that contain robotID.
func (g *GroupStore) Groups(robotID string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var out []string
	for name, members := range g.groups {
		if _, ok := members[robotID]; ok {
			out = append(out, name)
		}
	}
	return out
}

// DeleteGroup removes the group entirely.
func (g *GroupStore) DeleteGroup(group string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.groups, group)
}

// ListGroups returns the names of all known groups.
func (g *GroupStore) ListGroups() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]string, 0, len(g.groups))
	for name := range g.groups {
		out = append(out, name)
	}
	return out
}

// String returns a human-readable summary.
func (g *GroupStore) String() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return fmt.Sprintf("GroupStore{groups: %d}", len(g.groups))
}
