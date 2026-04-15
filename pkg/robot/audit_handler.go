package robot

import (
	"encoding/json"
	"net/http"
	"time"
)

// auditEntryJSON is the JSON-serialisable form of an AuditEntry.
type auditEntryJSON struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	RobotID   string    `json:"robot_id,omitempty"`
}

// AuditHandler returns an http.HandlerFunc that serves the audit log as JSON.
func AuditHandler(log *AuditLog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries := log.Entries()
		payload := make([]auditEntryJSON, 0, len(entries))
		for _, e := range entries {
			aj := auditEntryJSON{
				Timestamp: e.Timestamp,
				Type:      e.Event.Type.String(),
			}
			if e.Event.Status != nil {
				aj.RobotID = e.Event.Status.ID
			}
			payload = append(payload, aj)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, "failed to encode audit log", http.StatusInternalServerError)
		}
	}
}
