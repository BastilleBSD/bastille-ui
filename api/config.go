package api

import (
	"encoding/json"
	"log"
	"os"
)

var configFile = "api/config.json"
var cfg *ConfigStruct
var APIURL string
var Host string
var Port string
var Key string

func setAPIKey(key string) {
	Key = key
}

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

func setConfig(config *ConfigStruct) {

	Host = config.Host
	Port = config.Port
	Key  = config.Key

	setAPIKey(Key)
}