package api

import (
	"net/http"
	"log"
	"log/slog"
	"os"
	"encoding/json"
)

var DebugMode bool
var logger *slog.Logger

func InitLogger(debug bool) {

	level := slog.LevelInfo
	DebugMode = debug
	if DebugMode {
		level = slog.LevelDebug
	}

	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
	)

	slog.SetDefault(logger)
}

func logAll(level string, r *http.Request, cmdArgs []string, msg string) {

	headers := r.Header.Clone()
	headers.Del("Authorization")

	switch level {
	case "debug":
		if DebugMode {
			logger.Debug(msg,
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"rawArgs", cmdArgs,
				"remote", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"headers", headers,
			)
		}
	case "info":
		logger.Info(msg,
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"rawArgs", cmdArgs,
		)
	case "error":
		logger.Error(msg,
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"rawArgs", cmdArgs,
			"remote", r.RemoteAddr,
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

func validateMethodMiddleware(handler http.HandlerFunc, cmdName string, software string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*") 
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-TTYD-Port")
		w.Header().Set("Access-Control-Expose-Headers", "X-TTYD-Port")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		} else if r.Method == http.MethodGet {
			var cmd interface{}

			switch software {
			case "bastille":
				for _, c := range bastilleSpec.Commands {
					if c.Command == cmdName {
						cmd = c
						break
					}
				}
			case "rocinante":
				for _, c := range rocinanteSpec.Commands {
					if c.Command == cmdName {
						cmd = c
						break
					}
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cmd)

			logAll("debug", r, []string{cmdName}, "success")
			return
		}
		handler(w, r)
	}
}

// Log all API requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[API] %s %s %s", r.Method, r.URL.String(), r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}