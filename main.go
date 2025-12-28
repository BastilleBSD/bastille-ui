
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/bastille/bootstrap", BastilleBootstrapHandler)
	http.HandleFunc("/bastille/create", BastilleCreateHandler)
	http.HandleFunc("/bastille/destroy", BastilleDestroyHandler)
	http.HandleFunc("/bastille/rename", BastilleRenameHandler)
	http.HandleFunc("/bastille/restart", BastilleRestartHandler)
	http.HandleFunc("/bastille/start", BastilleStartHandler)
	http.HandleFunc("/bastille/stop", BastilleStopHandler)

	log.Println("BastilleBSD API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
