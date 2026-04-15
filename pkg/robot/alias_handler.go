package robot

import (
	"encoding/json"
	"net/http"
	"strings"
)

// AliasHandler provides HTTP access to the AliasStore.
// Routes:
//
//	GET    /aliases          — list all aliases
//	POST   /aliases/{alias}  — set alias (body: {"robot_id":"..."})
//	DELETE /aliases/{alias}  — remove alias
func AliasHandler(store *AliasStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Strip leading "/aliases" prefix and extract alias segment.
		path := strings.TrimPrefix(r.URL.Path, "/aliases")
		path = strings.TrimPrefix(path, "/")
		alias := strings.TrimSpace(path)

		switch r.Method {
		case http.MethodGet:
			if alias == "" {
				json.NewEncoder(w).Encode(store.All())
				return
			}
			id, ok := store.Resolve(alias)
			if !ok {
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"alias": alias, "robot_id": id})

		case http.MethodPost:
			if alias == "" {
				http.Error(w, `{"error":"alias required"}`, http.StatusBadRequest)
				return
			}
			var body struct {
				RobotID string `json:"robot_id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RobotID == "" {
				http.Error(w, `{"error":"robot_id required"}`, http.StatusBadRequest)
				return
			}
			if err := store.Set(alias, body.RobotID); err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"alias": alias, "robot_id": body.RobotID})

		case http.MethodDelete:
			if alias == "" {
				http.Error(w, `{"error":"alias required"}`, http.StatusBadRequest)
				return
			}
			store.Delete(alias)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})
}
