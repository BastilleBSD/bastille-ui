package api

import (
	"log"
	"net/http"
	"fmt"
)

func Start() {

	var bindAddr string
	config := loadConfig()
	setConfig(config)

	if Host == "0.0.0.0" || Host == "localhost" || Host == "" {
		bindAddr = "0.0.0.0"
		Host = "localhost"
	} else {
	       bindAddr = Host
	}
	
	addr := fmt.Sprintf("%s:%s", bindAddr, Port)

	loadRoutes()
	loadBastilleSpec()
	loadRocinanteSpec()

	log.Println("Starting BastilleBSD API server on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}