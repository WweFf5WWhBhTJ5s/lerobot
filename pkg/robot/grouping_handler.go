package robot

import (
	"encoding/json"
	"net/http"
	"strings"
)

// GroupingHandler exposes group membership over HTTP.
//
// Routes (mux pattern prefix assumed by caller):
//
//	GET  /groups/{group}          – list members
//	POST /groups/{group}/{robot}  – add robot to group
//	DELETE /groups/{group}/{robot} – remove robot from group
func GroupingHandler(gs *GroupStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Trim leading slash and split path segments.
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		// Expect at least one segment (the group name).
		if len(parts) < 1 || parts[0] == "" {
			http.Error(w, "group name required", http.StatusBadRequest)
			return
		}

		group := parts[0]

		switch r.Method {
		case http.MethodGet:
			members := gs.Members(group)
			if members == nil {
				members = []string{}
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"group":   group,
				"members": members,
			})

		case http.MethodPost:
			if len(parts) < 2 || parts[1] == "" {
				http.Error(w, "robot id required", http.StatusBadRequest)
				return
			}
			gs.Add(group, parts[1])
			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			if len(parts) < 2 || parts[1] == "" {
				// Delete entire group.
				gs.DeleteGroup(group)
			} else {
				gs.Remove(group, parts[1])
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
