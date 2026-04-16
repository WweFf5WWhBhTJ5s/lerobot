package robot

import (
	"fmt"
	"sync"
)

// AliasStore maps human-friendly aliases to robot IDs.
type AliasStore struct {
	mu      sync.RWMutex
	aliases map[string]string // alias -> robotID
	reverse map[string]string // robotID -> alias
}

// NewAliasStore creates an empty AliasStore.
func NewAliasStore() *AliasStore {
	return &AliasStore{
		aliases: make(map[string]string),
		reverse: make(map[string]string),
	}
}

// Set assigns an alias to a robot ID. Returns an error if the alias is
// already taken by a different robot.
// Note: each robot can only have one alias at a time; setting a new alias
// for a robot automatically removes its previous alias.
func (a *AliasStore) Set(alias, robotID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if existing, ok := a.aliases[alias]; ok && existing != robotID {
		return fmt.Errorf("alias %q already assigned to robot %q", alias, existing)
	}
	// Remove old alias for this robot if present.
	if old, ok := a.reverse[robotID]; ok {
		delete(a.aliases, old)
	}
	a.aliases[alias] = robotID
	a.reverse[robotID] = alias
	return nil
}

// Resolve returns the robot ID for the given alias, or "" if not found.
func (a *AliasStore) Resolve(alias string) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	id, ok := a.aliases[alias]
	return id, ok
}

// AliasFor returns the alias assigned to a robot ID, or "" if none.
func (a *AliasStore) AliasFor(robotID string) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	alias, ok := a.reverse[robotID]
	return alias, ok
}

// Delete removes an alias by name.
func (a *AliasStore) Delete(alias string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if id, ok := a.aliases[alias]; ok {
		delete(a.reverse, id)
		delete(a.aliases, alias)
	}
}

// DeleteByRobotID removes the alias assigned to the given robot ID, if any.
// This is a convenience method for when you know the robot ID but not the alias.
func (a *AliasStore) DeleteByRobotID(robotID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if alias, ok := a.reverse[robotID]; ok {
		delete(a.aliases, alias)
		delete(a.reverse, robotID)
	}
}

// All returns a copy of all alias -> robotID mappings.
func (a *AliasStore) All() map[string]string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	copy := make(map[string]string, len(a.aliases))
	for k, v := range a.aliases {
		copy[k] = v
	}
	return copy
}
