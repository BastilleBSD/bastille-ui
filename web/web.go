package web

import (
	"log"
	"net/http"
)

var webDir = "/usr/local/share/bastille-ui/web/"

func Start(webPath string) {

	if webPath != "" {
		webDir = webPath
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":8080" 
	log.Println("Starting BastilleBSD UI server on", addr)
	
	log.Fatal(http.ListenAndServe(addr, nil))
}