package web

import (
	"embed"
	"net/http"
	"log"
)

//go:embed index.html bastilleapi.js
var folder embed.FS

func Start() {
	handler := http.FileServer(http.FS(folder))

	log.Println("Starting BastilleBSD UI server on :8080")
	
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}