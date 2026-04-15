package robot

import "net/http"

// Router holds the dependencies needed to register all HTTP routes.
type Router struct {
	Health   *HealthChecker
	Snapshot *SnapshotService
}

// Register attaches all robot HTTP handlers to the given mux.
//
// Routes:
//
//	GET /healthz  — overall health check
//	GET /snapshot — latest robot snapshot
func (ro *Router) Register(mux *http.ServeMux) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	mux.HandleFunc("/healthz", HealthHandler(ro.Health))
	mux.HandleFunc("/snapshot", SnapshotHandler(ro.Snapshot))
}

// NewRouter constructs a Router from the provided dependencies.
func NewRouter(hc *HealthChecker, snap *SnapshotService) *Router {
	return &Router{
		Health:   hc,
		Snapshot: snap,
	}
}
