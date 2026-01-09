package web

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
       configFile = "web/config.json"
	cfg                    *ConfigStruct
	activeNode      *Node
	activeNodeMu sync.RWMutex
	Host string
	Port string
	User string
	Password string
)

func loadConfig() *ConfigStruct {

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
	return cfg
}

func saveConfig(config *ConfigStruct) error {

	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	return enc.Encode(config)
}

func setConfig(config *ConfigStruct) {
	Host = config.Host
	Port = config.Port
	User  = config.User
	Password  = config.Password
}