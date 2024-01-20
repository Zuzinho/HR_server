package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging logs url query info and execution time
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handling query:\nurl %s\nmethod %s\nhost %s\nheader %v", r.URL, r.Method, r.Host, r.Header)

		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("time %v", time.Since(start))
	})
}
