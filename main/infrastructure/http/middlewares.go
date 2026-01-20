package http

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // default status code
		written:        false,
	}
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.written = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

// Write ensures status code is captured even if WriteHeader is not explicitly called
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// shouldLogBody determines if the request body should be logged
func shouldLogBody(r *http.Request) bool {
	// Don't log for GET and DELETE methods (no body typically)
	if r.Method == http.MethodGet || r.Method == http.MethodDelete || r.Method == http.MethodOptions {
		return false
	}

	// Don't log for login/auth endpoints
	path := strings.ToLower(r.URL.Path)
	if strings.Contains(path, "/login") || strings.Contains(path, "/auth") || strings.Contains(path, "/signin") {
		return false
	}

	return true
}

// loggingMiddleware logs request details, status code, duration, and request body
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Read and log request body if applicable
		var bodyLog string
		if shouldLogBody(r) && r.Body != nil {
			// Read the body
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading request body: %v", err)
				bodyLog = "[error reading body]"
			} else {
				// Restore the body for downstream handlers
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Save body for logging (truncate if too long)
				if len(bodyBytes) > 0 {
					bodyStr := string(bodyBytes)
					if len(bodyStr) > 500 {
						bodyLog = bodyStr[:500] + "... [truncated]"
					} else {
						bodyLog = bodyStr
					}
				}
			}
		}

		// Wrap the response writer to capture status code
		wrappedWriter := newResponseWriter(w)

		// Process the request
		next.ServeHTTP(wrappedWriter, r)

		// Calculate duration
		duration := time.Since(start)

		// Log the request details
		if r.Method != http.MethodOptions {
			if bodyLog != "" {
				log.Printf(
					"[%s] %s %s - Status: %d - Duration: %v - Body: %s",
					r.Method,
					r.URL.Path,
					r.RemoteAddr,
					wrappedWriter.statusCode,
					duration,
					bodyLog,
				)
			} else {
				log.Printf(
					"[%s] %s %s - Status: %d - Duration: %v",
					r.Method,
					r.URL.Path,
					r.RemoteAddr,
					wrappedWriter.statusCode,
					duration,
				)
			}
		}
	})
}

func withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

var frontendURL = getEnv("FRONTEND_URL", "http://localhost:5173")

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 🔐 Adjust this in real life
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
