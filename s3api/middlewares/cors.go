// Package middlewares provides HTTP middleware components for the S3-compatible API gateway.
package middlewares

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig holds configuration for the CORS middleware.
type CORSConfig struct {
	// AllowedOrigins is a list of origins that are allowed to make cross-origin requests.
	// Use ["*"] to allow all origins.
	AllowedOrigins []string

	// AllowedMethods is a list of HTTP methods allowed for cross-origin requests.
	AllowedMethods []string

	// AllowedHeaders is a list of HTTP headers allowed in cross-origin requests.
	AllowedHeaders []string

	// ExposedHeaders is a list of headers that browsers are allowed to access.
	ExposedHeaders []string

	// AllowCredentials indicates whether the request can include user credentials.
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
	MaxAge int
}

// DefaultCORSConfig returns a CORSConfig with permissive defaults suitable for S3 API usage.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
			"Content-MD5",
			"Content-Length",
			"X-Amz-Date",
			"X-Amz-Content-Sha256",
			"X-Amz-Security-Token",
			"X-Amz-User-Agent",
		},
		ExposedHeaders: []string{
			"ETag",
			"X-Amz-Request-Id",
			"X-Amz-Version-Id",
		},
		AllowCredentials: false,
		MaxAge:           600,
	}
}

// CORSMiddleware returns an HTTP middleware that handles Cross-Origin Resource Sharing (CORS).
// It sets the appropriate CORS headers on responses and handles preflight OPTIONS requests.
func CORSMiddleware(cfg CORSConfig) func(http.Handler) http.Handler {
	allowedOriginSet := make(map[string]struct{}, len(cfg.AllowedOrigins))
	for _, o := range cfg.AllowedOrigins {
		allowedOriginSet[strings.ToLower(o)] = struct{}{}
	}

	allowAllOrigins := false
	if _, ok := allowedOriginSet["*"]; ok {
		allowAllOrigins = true
	}

	methods := strings.Join(cfg.AllowedMethods, ", ")
	headers := strings.Join(cfg.AllowedHeaders, ", ")
	exposed := strings.Join(cfg.ExposedHeaders, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// Not a CORS request; pass through.
				next.ServeHTTP(w, r)
				return
			}

			originAllowed := allowAllOrigins
			if !originAllowed {
				_, originAllowed = allowedOriginSet[strings.ToLower(origin)]
			}

			if !originAllowed {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", headers)

			if exposed != "" {
				w.Header().Set("Access-Control-Expose-Headers", exposed)
			}
			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if cfg.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
			}

			// Handle preflight request.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
