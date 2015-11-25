package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	FilePath     string
	AgentName    string                 `json:"agent_name"`
	AgentVersion string                 `json:"agent_version"`
	Checks       map[string]interface{} `json:"checks"`
	Handlers     map[string]interface{} `json:"handlers"`
}

func NewConfigFromFile(file string) (cfg Config, err error) {
	log.Debugf("agent.NewConfigFromFile => %s", file)

	absPath, err := filepath.Abs(file)
	if err != nil {
		return
	}

	configFile, err := os.Open(absPath)
	if err != nil {
		return
	}

	if err = json.NewDecoder(configFile).Decode(&cfg); err != nil {
		return
	}

	cfg.FilePath = absPath
	log.Infof("Loaded config from file: %s", cfg.FilePath)
	return
}
