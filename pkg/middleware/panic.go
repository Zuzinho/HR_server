package middleware

import (
	"log"
	"net/http"
)

// RecoverPanic recovers server from panic
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("recover panic %s", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
