package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	activeNode   *Node
	activeNodeMu sync.RWMutex
	Config       *ConfigStruct // ← global config
)

// Config struct for app
type ConfigStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	WebPort  string `json:"webPort"`
	APIPort  string `json:"apiPort"`
	APIKey   string `json:"apiKey"`
	Nodes    []Node `json:"nodes"`
}

type Node struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// ----------------------
// LoadConfig reads JSON from file
func LoadConfig() *ConfigStruct {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	var cfg ConfigStruct
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	Config = &cfg // store globally
	return Config
}

// SaveConfig writes JSON to file
func SaveConfig(cfg *ConfigStruct) error {
	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

// ----------------------
// Set/Get active node
func SetActiveNodeByName(name string) error {
	activeNodeMu.Lock()
	defer activeNodeMu.Unlock()

	if Config == nil {
		return fmt.Errorf("config not loaded")
	}

	for _, node := range Config.Nodes {
		if node.Name == name {
			activeNode = &node
			return nil
		}
	}

	return fmt.Errorf("node with name %s not found", name)
}

func GetActiveNode() *Node {
	activeNodeMu.RLock()
	defer activeNodeMu.RUnlock()
	return activeNode
}