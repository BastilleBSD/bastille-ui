package config

import (
	"encoding/json"
	"log"
	"os"
)

var ConfigPath = "config.json"

// Config struct for app
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	WebPort  string `json:"webPort`
	APIPort  string `json:"apiPort"`
	APIKey   string `json:"apiKey"`
}

// LoadConfig reads JSON from file
func LoadConfig() Config {
	file, err := os.Open(ConfigPath)
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

// SaveConfig writes JSON to file
func SaveConfig(cfg Config) error {
	file, err := os.Create(ConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ") // pretty print
	return enc.Encode(cfg)
}
