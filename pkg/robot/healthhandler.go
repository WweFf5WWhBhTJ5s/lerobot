package robot

import (
	"encoding/json"
	"net/http"
	"time"
)

// healthResponse is the JSON body returned by the health HTTP handler.
type healthResponse struct {
	Total       int       `json:"total"`
	Healthy     int       `json:"healthy"`
	Unhealthy   int       `json:"unhealthy"`
	Stale       int       `json:"stale"`
	GeneratedAt time.Time `json:"generated_at"`
	OK          bool      `json:"ok"`
}

// HealthHandler returns an http.Handler that serves the latest HealthReport
// from the given HealthChecker as JSON.
//
// It responds with 200 OK when all robots are healthy, and 503 Service
// Unavailable when any robot is unhealthy or stale.
func HealthHandler(hc *HealthChecker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rpt := hc.Report()
		ok := rpt.Total > 0 && rpt.Unhealthy == 0 && rpt.Stale == 0
		body := healthResponse{
			Total:       rpt.Total,
			Healthy:     rpt.Healthy,
			Unhealthy:   rpt.Unhealthy,
			Stale:       rpt.Stale,
			GeneratedAt: rpt.GeneratedAt,
			OK:          ok,
		}
		w.Header().Set("Content-Type", "application/json")
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		_ = json.NewEncoder(w).Encode(body)
	})
}
