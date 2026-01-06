package main

import (
	"bastille-ui/api"
	"bastille-ui/web"
	"bastille-ui/config"
)

func main() {

	var host string

	cfg := config.LoadConfig()

	if cfg.Address == "0.0.0.0" || cfg.Address == "localhost" || cfg.Address == "" {
		host = "localhost"
	} else {
		host = cfg.Address
	}

	// Load and set config variables
	api.SetAPIKey(cfg.APIKey)
	api.SetAPIAddress(host, cfg.APIPort)
	web.SetAPIKey(cfg.APIKey)
	web.SetAPIAddress(host, cfg.WebPort)
	web.SetCredentials(cfg.Username, cfg.Password)

	
	addrAPI := host + ":" + cfg.APIPort
	addrWeb := host + ":" + cfg.WebPort
	go api.Start(addrAPI)
	go web.Start(addrWeb)
	select {}
}
