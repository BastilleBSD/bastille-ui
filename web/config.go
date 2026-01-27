package web

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	cfg                    *ConfigStruct
	activeNode      *Node
	activeNodeMu sync.RWMutex
	Host string
	Port string
	User string
	Password string
)

func loadConfig() (*ConfigStruct, error) {

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	var c ConfigStruct
	if err := json.NewDecoder(file).Decode(&c); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	cfg = &c
	Host = c.Host
	Port = c.Port
	User = c.User
	Password = c.Password
	return cfg, nil
}

func saveConfig(config *ConfigStruct) error {

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	return enc.Encode(config)
}

func setConfig(c *ConfigStruct) {
	Host = c.Host
	Port = c.Port
	User  = c.User
	Password  = c.Password
}
