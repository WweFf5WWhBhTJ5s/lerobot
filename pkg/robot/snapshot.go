package robot

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of all robot statuses.
type Snapshot struct {
	TakenAt time.Time
	Robots  map[string]Status
}

// SnapshotService periodically captures registry state.
type SnapshotService struct {
	mu       sync.RWMutex
	registry *Registry
	interval time.Duration
	latest   *Snapshot
	stop     chan struct{}
	wg       sync.WaitGroup
}

// NewSnapshotService creates a SnapshotService that snapshots the registry
// at the given interval.
func NewSnapshotService(r *Registry, interval time.Duration) *SnapshotService {
	return &SnapshotService{
		registry: r,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins periodic snapshotting in the background.
func (s *SnapshotService) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.capture()
			case <-s.stop:
				return
			}
		}
	}()
}

// Stop halts the background snapshotting goroutine.
func (s *SnapshotService) Stop() {
	close(s.stop)
	s.wg.Wait()
}

// Latest returns the most recent snapshot, or nil if none has been taken.
func (s *SnapshotService) Latest() *Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.latest
}

// capture takes an immediate snapshot of the registry.
func (s *SnapshotService) capture() {
	list := s.registry.List()
	robots := make(map[string]Status, len(list))
	for _, id := range list {
		if st, ok := s.registry.Status(id); ok {
			robots[id] = st
		}
	}
	snap := &Snapshot{
		TakenAt: time.Now(),
		Robots:  robots,
	}
	s.mu.Lock()
	s.latest = snap
	s.mu.Unlock()
}
