package robot

import (
	"encoding/json"
	"net/http"
)

// replayResponse is the JSON body returned by ReplayHandler.
type replayResponse struct {
	Count  int     `json:"count"`
	Events []Event `json:"events"`
}

// ReplayHandler returns an http.HandlerFunc that serves buffered events
// from the provided ReplayBuffer as JSON.
//
// GET /replay
//
//	200 OK  — always, even when the buffer is empty.
func ReplayHandler(buf *ReplayBuffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events := buf.Replay()

		resp := replayResponse{
			Count:  len(events),
			Events: events,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
