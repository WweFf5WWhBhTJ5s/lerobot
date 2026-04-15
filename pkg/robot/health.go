package robot

import (
	"sync"
	"time"
)

// HealthReport summarises the current health of all known robots.
type HealthReport struct {
	Total     int
	Healthy   int
	Unhealthy int
	Stale     int
	GeneratedAt time.Time
}

// HealthChecker periodically evaluates robot health and exposes a report.
type HealthChecker struct {
	mu       sync.RWMutex
	registry *Registry
	report   HealthReport
	interval time.Duration
	stop     chan struct{}
	wg       sync.WaitGroup
}

// NewHealthChecker creates a HealthChecker that polls the registry at the
// given interval.
func NewHealthChecker(r *Registry, interval time.Duration) *HealthChecker {
	return &HealthChecker{
		registry: r,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins background health evaluation.
func (h *HealthChecker) Start() {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		ticker := time.NewTicker(h.interval)
		defer ticker.Stop()
		h.evaluate()
		for {
			select {
			case <-ticker.C:
				h.evaluate()
			case <-h.stop:
				return
			}
		}
	}()
}

// Stop halts background evaluation and waits for the goroutine to exit.
func (h *HealthChecker) Stop() {
	close(h.stop)
	h.wg.Wait()
}

// Report returns the most recently computed HealthReport.
func (h *HealthChecker) Report() HealthReport {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.report
}

func (h *HealthChecker) evaluate() {
	statuses := h.registry.Status()
	rpt := HealthReport{
		Total:       len(statuses),
		GeneratedAt: time.Now(),
	}
	for _, s := range statuses {
		switch {
		case s.IsStale():
			rpt.Stale++
		case s.IsHealthy():
			rpt.Healthy++
		default:
			rpt.Unhealthy++
		}
	}
	h.mu.Lock()
	h.report = rpt
	h.mu.Unlock()
}
