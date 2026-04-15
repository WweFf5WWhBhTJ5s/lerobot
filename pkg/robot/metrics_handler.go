package robot

import (
	"encoding/json"
	"net/http"
)

// MetricsResponse is the HTTP response body for the metrics endpoint.
type MetricsResponse struct {
	Total   int `json:"total"`
	Healthy int `json:"healthy"`
	Stale   int `json:"stale"`
}

// MetricsHandler returns an HTTP handler that exposes robot metrics as JSON.
func MetricsHandler(m *Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := m.Snapshot()

		resp := MetricsResponse{
			Total:   snap.Total,
			Healthy: snap.Healthy,
			Stale:   snap.Stale,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
