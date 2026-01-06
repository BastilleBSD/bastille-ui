package web

import (
	"fmt"
	"log"
	"net/http"
)

var apiKey, apiUrl string

func SetAPIKey(key string) {
	apiKey = key
}

func SetAPIAddress(address, port string) {
	apiUrl = fmt.Sprintf("http://%s:%s", address, port)
}

// Log all Web requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[WEB] %s %s %s", r.Method, r.URL.String(), r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
