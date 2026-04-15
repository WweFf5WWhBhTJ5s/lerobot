package robot

import (
	"net/http"
)

// QuotaMiddleware wraps an HTTP handler and enforces per-robot beat quotas.
// The robot ID is read from the "id" query parameter.
// Requests that exceed the quota receive a 429 Too Many Requests response.
func QuotaMiddleware(q *QuotaEnforcer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing robot id", http.StatusBadRequest)
			return
		}
		if !q.Allow(id) {
			http.Error(w, "quota exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
