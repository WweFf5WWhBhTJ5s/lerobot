package robot

import (
	"encoding/json"
	"net/http"
)

// PresenceHandler returns an HTTP handler that reports presence information.
//
// GET /presence?id=<robotID> returns the presence record for a single robot.
// GET /presence returns all presence records.
func PresenceHandler(pt *PresenceTracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id != "" {
			rec := pt.Get(id)
			if rec == nil {
				http.Error(w, "robot not found", http.StatusNotFound)
				return
			}
			if err := json.NewEncoder(w).Encode(rec); err != nil {
				http.Error(w, "encode error", http.StatusInternalServerError)
			}
			return
		}

		all := pt.All()
		if err := json.NewEncoder(w).Encode(all); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	}
}
