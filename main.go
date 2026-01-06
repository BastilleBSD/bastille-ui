package main

import (
	"bastille-ui/api"
	"bastille-ui/web"
	"bastille-ui/config"
)

func main() {

	// Load and set config variables
	cfg := config.LoadConfig()
	api.SetAPIKey(cfg.APIKey)
	api.SetAPIAddress(cfg.Address, cfg.APIPort)
	web.SetAPIKey(cfg.APIKey)
	web.SetAPIAddress(cfg.Address, cfg.WebPort)
	web.SetCredentials(cfg.Username, cfg.Password)

	addrAPI := cfg.Address + ":" + cfg.APIPort
	addrWeb := cfg.Address + ":" + cfg.WebPort
	go api.Start(addrAPI)
	go web.Start(addrWeb)
	select {}
}
