package robot

import (
	"encoding/json"
	"net/http"
	"time"
)

// SnapshotResponse is the HTTP response body for snapshot endpoints.
type SnapshotResponse struct {
	TakenAt time.Time       `json:"taken_at"`
	Robots  []StatusSummary `json:"robots"`
}

// StatusSummary is a simplified view of a robot's status.
type StatusSummary struct {
	ID      string `json:"id"`
	Healthy bool   `json:"healthy"`
	Stale   bool   `json:"stale"`
}

// SnapshotHandler returns an HTTP handler that serves the latest robot snapshot.
func SnapshotHandler(svc *SnapshotService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := svc.Latest()
		if snap == nil {
			http.Error(w, "no snapshot available", http.StatusServiceUnavailable)
			return
		}

		summaries := make([]StatusSummary, 0, len(snap.Robots))
		for id, st := range snap.Robots {
			summaries = append(summaries, StatusSummary{
				ID:      id,
				Healthy: st.IsHealthy(),
				Stale:   st.Stale(),
			})
		}

		resp := SnapshotResponse{
			TakenAt: snap.TakenAt,
			Robots:  summaries,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
