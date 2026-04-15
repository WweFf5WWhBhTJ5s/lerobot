package robot

import (
	"encoding/json"
	"net/http"
)

// FilterHandler returns an http.HandlerFunc that filters the robot registry
// snapshot by name substring or tag query parameters.
//
// Query parameters:
//
//	name - substring match against robot ID (case-insensitive by default)
//	tag  - exact tag match via the TagStore
func FilterHandler(reg *Registry, svc *FilterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ids := reg.List()
		statuses := make([]Status, 0, len(ids))
		for _, id := range ids {
			if s, ok := reg.Status(id); ok {
				statuses = append(statuses, s)
			}
		}

		q := r.URL.Query()
		if name := q.Get("name"); name != "" {
			statuses = svc.ByName(statuses, name)
		}
		if tag := q.Get("tag"); tag != "" {
			statuses = svc.ByTag(statuses, tag)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(statuses); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	}
}
