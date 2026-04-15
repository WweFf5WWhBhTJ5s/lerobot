package robot

import (
	"encoding/json"
	"net/http"
)

// ThrottleHandler returns an HTTP handler that exposes throttle state.
// GET /throttle returns the number of robots currently tracked by the throttle.
// DELETE /throttle/{id} resets the throttle state for a specific robot.
func ThrottleHandler(th *Throttle) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/throttle", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int{
			"tracked": th.Len(),
		})
	})

	mux.HandleFunc("/throttle/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id parameter", http.StatusBadRequest)
			return
		}
		th.Reset(id)
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}
