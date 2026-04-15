package robot

import (
	"encoding/json"
	"net/http"
)

// rateLimitStatus is the JSON response body for the rate limit status endpoint.
type rateLimitStatus struct {
	RobotID   string `json:"robot_id"`
	Count     int    `json:"event_count"`
	MaxEvents int    `json:"max_events"`
	Allowed   bool   `json:"allowed"`
}

// RateLimitHandler returns an http.HandlerFunc that reports the current rate
// limit status for a robot identified by the "id" query parameter.
//
// GET /ratelimit?id=<robotID>
func RateLimitHandler(rl *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing 'id' query parameter", http.StatusBadRequest)
			return
		}

		count := rl.Count(id)
		allowed := count < rl.cfg.MaxEvents

		status := rateLimitStatus{
			RobotID:   id,
			Count:     count,
			MaxEvents: rl.cfg.MaxEvents,
			Allowed:   allowed,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(status)
	}
}
