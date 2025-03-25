package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v\nStack: %s", err, debug.Stack())

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := map[string]interface{}{
					"error": "Internal Server Error",
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					log.Printf("Error encoding error response: %v", err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
