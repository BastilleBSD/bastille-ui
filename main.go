package main

import (
	"flag"
	"log"
	"net/http"
)

var webDir = "/usr/local/share/bastille-ui/web"

func main() {

	webPath := flag.String("webdir", "", "Web files location")
	flag.Parse()

	Start(*webPath)
}

func Start(webPath string) {
	if webPath != "" {
		webDir = webPath
	}
	
	http.Handle("/", http.FileServer(http.Dir(webDir))) 

	addr := ":8080"
	log.Printf("Starting BastilleBSD UI server on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}