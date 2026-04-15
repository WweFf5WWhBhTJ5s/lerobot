package robot

import (
	"strings"
	"sync"
)

// FilterConfig holds configuration for the robot filter service.
type FilterConfig struct {
	// CaseSensitive controls whether tag/name matching is case sensitive.
	CaseSensitive bool
}

// DefaultFilterConfig returns a FilterConfig with sensible defaults.
func DefaultFilterConfig() FilterConfig {
	return FilterConfig{
		CaseSensitive: false,
	}
}

// FilterService filters robots from a registry snapshot by name or tags.
type FilterService struct {
	mu     sync.RWMutex
	cfg    FilterConfig
	store  *TagStore
}

// NewFilterService creates a new FilterService backed by the given TagStore.
func NewFilterService(store *TagStore, cfg FilterConfig) *FilterService {
	return &FilterService{
		cfg:   cfg,
		store: store,
	}
}

// ByName returns all statuses from the provided list whose robot ID contains
// the given substring.
func (f *FilterService) ByName(statuses []Status, substr string) []Status {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if !f.cfg.CaseSensitive {
		substr = strings.ToLower(substr)
	}

	out := make([]Status, 0, len(statuses))
	for _, s := range statuses {
		id := s.RobotID
		if !f.cfg.CaseSensitive {
			id = strings.ToLower(id)
		}
		if strings.Contains(id, substr) {
			out = append(out, s)
		}
	}
	return out
}

// ByTag returns all statuses from the provided list whose robot has the given
// tag assigned in the TagStore.
func (f *FilterService) ByTag(statuses []Status, tag string) []Status {
	f.mu.RLock()
	defer f.mu.RUnlock()

	tagged := f.store.FindByTag(tag)
	set := make(map[string]struct{}, len(tagged))
	for _, id := range tagged {
		set[id] = struct{}{}
	}

	out := make([]Status, 0)
	for _, s := range statuses {
		if _, ok := set[s.RobotID]; ok {
			out = append(out, s)
		}
	}
	return out
}
