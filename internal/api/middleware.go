package api

import (
	"log"
	"net/http"
	"time"
)

func ExecutionTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Processing request: %s %s", r.Method, r.URL.Path)
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			log.Printf("Finished processing %s %s within %v", r.Method, r.URL.Path, duration)
		},
	)
}
