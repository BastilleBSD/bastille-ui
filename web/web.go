package web

import (
	"fmt"
	"log"
	"net/http"
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

	if len(cfg.Nodes) > 0 {
		activeNode = &cfg.Nodes[0]
	} else {
		activeNode = nil
	}

	log.Println("Starting BastilleBSD WebUI server on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
