package main

import (
    "bastille-ui/api"
    "bastille-ui/web"
    "flag"
)

func main() {

	apiOnly := flag.Bool("api-only", false, "Only run the API server")
	debug := flag.Bool("debug", false, "Enable debug logging")
	configPath := flag.String("config", "", "Config file location")
	apiPort := flag.String("api-port", "", "API server port")

	flag.Parse()

	api.InitLogger(*debug)

	go api.Start(*configPath, *apiPort)

	if !*apiOnly {
		go web.Start()
	}

	select {}
}
