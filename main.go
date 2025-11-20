
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/jails/create", createHandler)
	http.HandleFunc("/jails/destroy", destroyHandler)
	http.HandleFunc("/jails/start", startHandler)
	http.HandleFunc("/jails/stop", stopHandler)
	http.HandleFunc("/jails/restart", restartHandler)
	http.HandleFunc("/jails/rename", renameHandler)
	http.HandleFunc("/jails/upgrade", upgradeHandler)
	http.HandleFunc("/jails/list", listHandler)

	log.Println("✅ BastilleBSD API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
