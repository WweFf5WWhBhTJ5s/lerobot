package robot

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// TagStore maintains a set of string tags per robot ID.
type TagStore struct {
	mu   sync.RWMutex
	tags map[string]map[string]struct{}
}

// NewTagStore creates an empty TagStore.
func NewTagStore() *TagStore {
	return &TagStore{
		tags: make(map[string]map[string]struct{}),
	}
}

// Set replaces all tags for the given robot ID.
func (ts *TagStore) Set(id string, tags []string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	set := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			set[t] = struct{}{}
		}
	}
	ts.tags[id] = set
}

// Add appends tags to the existing set for a robot ID.
func (ts *TagStore) Add(id string, tags ...string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.tags[id] == nil {
		ts.tags[id] = make(map[string]struct{})
	}
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t != "" {
			ts.tags[id][t] = struct{}{}
		}
	}
}

// Remove deletes a specific tag from a robot ID.
func (ts *TagStore) Remove(id, tag string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.tags[id], tag)
}

// Get returns a sorted slice of tags for a robot ID.
func (ts *TagStore) Get(id string) []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	set := ts.tags[id]
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// Delete removes all tags for a robot ID.
func (ts *TagStore) Delete(id string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.tags, id)
}

// FindByTag returns all robot IDs that have the given tag.
func (ts *TagStore) FindByTag(tag string) []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	var ids []string
	for id, set := range ts.tags {
		if _, ok := set[tag]; ok {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

// String returns a human-readable summary.
func (ts *TagStore) String() string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return fmt.Sprintf("TagStore{robots: %d}", len(ts.tags))
}
