package api

import (
	"fmt"
	"log"
	"net/http"
)

var apiKey, apiUrl, apiAddress string

// Set API key
func SetAPIKey(key string) {
	apiKey = key
}

func SetAPIAddress(address, port string) {
	if address == "0.0.0.0" || address == "localhost" || address == "" {
		apiAddress = "localhost"
	} else {
		apiAddress = address
	}
	apiUrl = fmt.Sprintf("http://%s:%s", apiAddress, port)
}

// Validate API key in request header
func apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
