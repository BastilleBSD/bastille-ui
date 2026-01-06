package main

import (
	"bastille-ui/api"
	"bastille-ui/web"
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	IP string `json:"address"`
	APIPort string `json:"api_port"`
	WebPort string `json:"web_port"`
	APIKey  string `json:"api_key"`
}

func loadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()
	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}
	return cfg
}

func main() {
	cfg := loadConfig("config.json")
	api.SetAPIKey(cfg.APIKey)
	web.SetAPIKey(cfg.APIKey)
	web.SetAPIUrl(cfg.IP, cfg.APIPort)
	go api.Start(":" + cfg.APIPort)
	go web.Start(":" + cfg.WebPort)
	select {}
}
