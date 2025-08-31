package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"loan-service/internal/models"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Printf(
			"%s %s %d %v",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
		)
	})
}

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				errorResponse := models.APIResponse{
					Object: "error",
					Error: &models.APIError{
						Type:    "api_error",
						Code:    "internal_server_error",
						Message: "An unexpected error occurred",
					},
				}

				json.NewEncoder(w).Encode(errorResponse)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
