package api

import (
	"encoding/json"
	"os"
	
	"github.com/gin-gonic/gin"
)

var bastilleSpec *BastilleSpecStruct
var rocinanteSpec *RocinanteSpecStruct
var configFile = "/usr/local/etc/bastille-ui/config.json"
var cfg *ConfigStruct
var APIURL string
var Host string
var Port string

func getParam(c *gin.Context, key string) string {
	return c.Query(key)
}

func loadBastilleSpec() (*BastilleSpecStruct, error) {

	logRequest("debug", "loadBastilleSpec", nil, nil, nil)

	specFile := "api/bastille.json"
	var spec BastilleSpecStruct

	data, err := os.ReadFile(specFile)
	if err != nil {
		logRequest("error", "Failed to read Bastille spec file", nil, nil, err.Error())
		return nil, err
	}

	if err := json.Unmarshal(data, &spec); err != nil {
		logRequest("error", "Failed to parse Bastille spec", nil, nil, err.Error())
		return nil, err
	}

	bastilleSpec = &spec
	return bastilleSpec, nil
}

func loadRocinanteSpec() (*RocinanteSpecStruct, error) {

	logRequest("debug", "loadRocinanteSpec", nil, nil, nil)

	specFile := "api/rocinante.json"
	var spec RocinanteSpecStruct

	data, err := os.ReadFile(specFile)
	if err != nil {
		logRequest("error", "Failed to read Rocinante spec file", nil, nil, err.Error())
		return nil, err
	}

	if err := json.Unmarshal(data, &spec); err != nil {
		logRequest("error", "Failed to parse Rocinante spec", nil, nil, err.Error())
		return nil, err
	}

	rocinanteSpec = &spec
	return rocinanteSpec, nil
}

func loadConfig() (*ConfigStruct, error) {

	logRequest("debug", "loadConfig", nil, nil, nil)

	data, err := os.ReadFile(configFile)
	if err != nil {
		logRequest("error", "Failed to read config file", nil, nil, err.Error())
		return nil, err
	}

	var c ConfigStruct
	if err := json.Unmarshal(data, &c); err != nil {
		logRequest("error", "Failed to parse config file", nil, nil, err.Error())
		return nil, err
	}

	if c.APIKeys == nil {
		c.APIKeys = make(map[string]APIKeyStruct)
	}

	cfg = &c
	Host = c.Host
	Port = c.Port
	return cfg, nil
}

func saveConfig() error {

	logRequest("debug", "saveConfig", nil, nil, nil)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		logRequest("error", "Failed to marshal config for saving", nil, nil, err.Error())
		return err
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		logRequest("error", "Failed to write config file", nil, nil, err.Error())
		return err
	}

	return nil
}