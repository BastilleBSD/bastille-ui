package web

import (
	"embed"
	"net/http"
	"log"
)

//go:embed index.html bastilleapi.js
var folder embed.FS

func Start() {
	// http.FS converts our embed.FS into the interface http.FileServer needs
	handler := http.FileServer(http.FS(folder))

	log.Println("Web server starting on :8081 (Embedded mode)")
	
	// This will automatically serve index.html when you visit http://localhost:8081
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}