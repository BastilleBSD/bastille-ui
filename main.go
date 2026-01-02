package main

import (
	"bastille-ui/api"
	"bastille-ui/web"
	"log"
	"net/http"
)

func main() {
	api.Start()
	web.Start()
        log.Println("BastilleBSD UI running on http://localhost:8080")
        log.Fatal(http.ListenAndServe(":8080", nil))
}
