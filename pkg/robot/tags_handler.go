package robot

import (
	"encoding/json"
	"net/http"
	"strings"
)

// TagsHandler returns an http.HandlerFunc that exposes tag store operations.
// GET  /tags/{id}        — returns all tags for a robot
// POST /tags/{id}        — body: {"tags":["a","b"]} adds tags
// DELETE /tags/{id}/{tag} — removes a single tag
func TagsHandler(store *TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// path: /tags/{id} or /tags/{id}/{tag}
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 2 {
			http.Error(w, "robot id required", http.StatusBadRequest)
			return
		}
		id := parts[1]

		switch r.Method {
		case http.MethodGet:
			tags := store.Get(id)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string][]string{"tags": tags})

		case http.MethodPost:
			var body struct {
				Tags []string `json:"tags"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON body", http.StatusBadRequest)
				return
			}
			for _, t := range body.Tags {
				store.Add(id, t)
			}
			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			if len(parts) < 3 {
				http.Error(w, "tag required", http.StatusBadRequest)
				return
			}
			store.Remove(id, parts[2])
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
