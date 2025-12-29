package api

import (
	"log"
	"net/http"
)

var apiKey string

func Start() {
	// apiKey = os.Getenv("BASTILLE_API_KEY")
	apiKey = "testingkey"
	if apiKey == "" {
		log.Fatal("BASTILLE_API_KEY not set")
	}
	http.HandleFunc("/bastille/bootstrap", verifyAPIKey(BastilleBootstrapHandler))
	http.HandleFunc("/bastille/create", verifyAPIKey(BastilleCreateHandler))
	http.HandleFunc("/bastille/destroy", verifyAPIKey(BastilleDestroyHandler))
	http.HandleFunc("/bastille/rename", verifyAPIKey(BastilleRenameHandler))
	http.HandleFunc("/bastille/restart", verifyAPIKey(BastilleRestartHandler))
	http.HandleFunc("/bastille/start", verifyAPIKey(BastilleStartHandler))
	http.HandleFunc("/bastille/stop", verifyAPIKey(BastilleStopHandler))
}

func verifyAPIKey(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authToken := "Bearer " + apiKey

		if authHeader != authToken  {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}
