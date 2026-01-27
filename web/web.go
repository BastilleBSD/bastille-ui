package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var webDir = "/usr/local/share/bastille-api/"
var configFile = (webDir + "web/config.json")

func Start(webPath string) {

	var bindAddr string

	if webPath != "" {
		webDir = (webPath+"/")
		configFile = (webDir + "web/config.json")
	}

	_, err := loadConfig()
	if err != nil {
		log.Println("Failed to load config", err.Error())
		os.Exit(1)
	}

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
