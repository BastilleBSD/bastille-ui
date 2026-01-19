package api

import (
	"net/http"
	"log"
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger(debug bool) {

	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
	)

	slog.SetDefault(logger)
}

func logAll(level string, r *http.Request, cmdArgs []string, extra map[string]any, err error) {

	headers := r.Header.Clone()
	headers.Del("Authorization")

	switch level {
	case "debug":
		logger.Debug("Request debug",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"remote", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"args", cmdArgs,
			"extra", extra,
			"headers", headers,
			"error", err,
		)
	case "info":
		// Info only prints basic request info
		logger.Info("Request info",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"args", cmdArgs,
		)
	case "error":
		logger.Error("Request error",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"args", cmdArgs,
			"extra", extra,
			"error", err,
		)
	}
}

// Validate API key in request header
func apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+Key {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers for all requests
		w.Header().Set("Access-Control-Allow-Origin", "*") 
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-TTYD-Port")

		w.Header().Set("Access-Control-Expose-Headers", "X-TTYD-Port")

		// Handle Preflight: Return OK immediately for OPTIONS
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Log all API requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[API] %s %s %s", r.Method, r.URL.String(), r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}