package robot

import (
	"encoding/json"
	"net/http"
)

type quotaUsageResponse struct {
	RobotID string `json:"robot_id"`
	Usage   int    `json:"usage"`
	Limit   int    `json:"limit"`
	Allowed bool   `json:"allowed"`
}

// QuotaHandler returns an HTTP handler that reports quota usage for a robot.
// It expects a "id" query parameter identifying the robot.
func QuotaHandler(q *QuotaEnforcer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing robot id", http.StatusBadRequest)
			return
		}

		usage := q.Usage(id)
		allowed := usage < q.cfg.MaxBeatsPerHour

		resp := quotaUsageResponse{
			RobotID: id,
			Usage:   usage,
			Limit:   q.cfg.MaxBeatsPerHour,
			Allowed: allowed,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
