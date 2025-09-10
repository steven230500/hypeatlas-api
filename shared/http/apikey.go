package http

import (
	"net/http"
	"os"
	"strings"
)

func ApiKeyMiddleware(next http.Handler) http.Handler {
	keys := strings.Split(os.Getenv("API_KEYS"), ",")
	allow := map[string]struct{}{}
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k != "" {
			allow[k] = struct{}{}
		}
	}
	// si no hay keys definidas, passa-through
	if len(allow) == 0 {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("X-API-Key")
		if _, ok := allow[got]; !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
