package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var bastilleSpec *BastilleSpecStruct
var rocinanteSpec *RocinanteSpecStruct
var configFile = "api/config.json"
var cfg *ConfigStruct
var APIURL string
var Host string
var Port string
var Key string

func setAPIKey(key string) {
	Key = key
}

func getParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func ValidateBastilleCommandParameters(r *http.Request, cmdArgs []string) error {

	query := r.URL.Query()

	cmdName := cmdArgs[0]
	var subcmd *BastilleCommandStruct
	for _, c := range bastilleSpec.Commands {
		if c.Command == cmdName {
			subcmd = &c
			break
		}
	}

	paramsMap := make(map[string]struct{})
	for _, params := range subcmd.Parameters {
		paramsMap[strings.ToLower(params)] = struct{}{}
	}

	for param := range query {
		if _, ok := paramsMap[strings.ToLower(param)]; !ok {
			err := fmt.Sprintf("invalid parameter %q for command %q", param, cmdName)
			logAll("error", r, cmdArgs, err)
			return fmt.Errorf(err)
		}
	}

	optionsValueMap := make(map[string]interface{})
	for _, opt := range subcmd.Options {
		if opt.SFlag != "" {
			optionsValueMap[opt.SFlag] = opt.Value
		}
		if opt.LFlag != "" {
			optionsValueMap[opt.LFlag] = opt.Value
		}
	}

	optionsParam := query.Get("options")
	if optionsParam != "" {

		optionsParam = strings.ReplaceAll(optionsParam, "+", " ")
		opts := strings.Fields(optionsParam)

		for i := 0; i < len(opts); i++ {
			arg := opts[i]
			valueType, ok := optionsValueMap[arg]

			if !ok {
				err := fmt.Sprintf("invalid option %q for command %q", arg, cmdName)
				logAll("error", r, cmdArgs, err)
				return fmt.Errorf(err)
			}

			if valueType == "" || valueType == nil {
				continue
			}

			if i+1 >= len(opts) {
				err := fmt.Sprintf("option %q requires a value", arg)
				logAll("error", r, cmdArgs, err)
				return fmt.Errorf(err)
			}

			i++
			val := opts[i]

			if valueType == "int" {
				if _, err := strconv.Atoi(val); err != nil {
					err := fmt.Sprintf("option %q requires a numeric value", arg)
					logAll("error", r, cmdArgs, err)
					return fmt.Errorf(err)
				}
			}

			if strings.HasPrefix(val, "-") {
				err := fmt.Sprintf("option %q requires a value", arg)
				logAll("error", r, cmdArgs, err)
				return fmt.Errorf(err)
			}
		}
	}

	logAll("debug", r, cmdArgs, "command validated")

	return nil
}

func ValidateRocinanteCommandParameters(r *http.Request, cmdArgs []string) error {

	query := r.URL.Query()

	cmdName := cmdArgs[0]
	var subcmd *RocinanteCommandStruct
	for _, c := range rocinanteSpec.Commands {
		if c.Command == cmdName {
			subcmd = &c
			break
		}
	}

	paramsMap := make(map[string]struct{})
	for _, params := range subcmd.Parameters {
		paramsMap[strings.ToLower(params)] = struct{}{}
	}

	for param := range query {
		if _, ok := paramsMap[strings.ToLower(param)]; !ok {
			err := fmt.Sprintf("invalid parameter %q for command %q", param, cmdName)
			logAll("error", r, cmdArgs, err)
			return fmt.Errorf(err)
		}
	}

	optionsValueMap := make(map[string]interface{})
	for _, opt := range subcmd.Options {
		if opt.SFlag != "" {
			optionsValueMap[opt.SFlag] = opt.Value
		}
		if opt.LFlag != "" {
			optionsValueMap[opt.LFlag] = opt.Value
		}
	}

	optionsParam := query.Get("options")
	if optionsParam != "" {

		optionsParam = strings.ReplaceAll(optionsParam, "+", " ")
		opts := strings.Fields(optionsParam)

		for i := 0; i < len(opts); i++ {

			arg := opts[i]
			_, ok := optionsValueMap[arg]

			if !ok {
				err := fmt.Sprintf("invalid option %q for command %q", arg, cmdName)
				logAll("error", r, cmdArgs, err)
				return fmt.Errorf(err)
			}
			i++
		}
	}

	logAll("debug", r, cmdArgs, "command validated")

	return nil
}

func loadBastilleSpec() *BastilleSpecStruct {

	specFile := "api/bastille.json"
	var spec BastilleSpecStruct

	data, err := os.ReadFile(specFile)
	if err != nil {
		log.Fatalf("Failed to read Bastille spec file: %v", err)
	}

	if err := json.Unmarshal(data, &spec); err != nil {
		log.Fatalf("Failed to parse Bastille spec: %v", err)
	}

	bastilleSpec = &spec
	return bastilleSpec
}

func loadRocinanteSpec() *RocinanteSpecStruct {

	specFile := "api/rocinante.json"
	var spec RocinanteSpecStruct

	data, err := os.ReadFile(specFile)
	if err != nil {
		log.Fatalf("Failed to read Rocinante spec file: %v", err)
	}

	if err := json.Unmarshal(data, &spec); err != nil {
		log.Fatalf("Failed to parse Rocinante spec: %v", err)
	}

	rocinanteSpec = &spec
	return rocinanteSpec
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
	Key = config.Key

	setAPIKey(Key)
}
