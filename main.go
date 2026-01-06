package main

import (
	"bastille-ui/api"
	"bastille-ui/web"
	"bastille-ui/config"
)

func main() {

	var host string
	
	// Load and set config variables
	cfg := config.LoadConfig()
	api.SetAPIKey(cfg.APIKey)
	api.SetAPIAddress(cfg.Address, cfg.APIPort)
	web.SetAPIKey(cfg.APIKey)
	web.SetAPIAddress(cfg.Address, cfg.WebPort)
	web.SetCredentials(cfg.Username, cfg.Password)

	if cfg.Address == "0.0.0.0" || cfg.Address == "localhost" || cfg.Address == "" {
		host = "localhost"
	} else {
		host = cfg.Address
	}
	
	addrAPI := host + ":" + cfg.APIPort
	addrWeb := host + ":" + cfg.WebPort
	go api.Start(addrAPI)
	go web.Start(addrWeb)
	select {}
}
