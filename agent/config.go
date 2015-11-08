package agent

import (
	"encoding/json"
	log "github.com/patrickrand/gamma/log"
	"os"
	"path/filepath"
)

type Config struct {
	FilePath     string
	AgentName    string                 `json:"agent_name"`
	AgentVersion string                 `json:"agent_version"`
	Monitors     map[string]interface{} `json:"monitors"`
	Handlers     map[string]interface{} `json:"handlers"`
}

func NewConfigFromFile(file string) (cfg Config, err error) {
	log.DBUG(AGENT, "config.NewConfigFromFile => %s", file)

	absPath, err := filepath.Abs(file)
	if err != nil {
		return
	}

	configFile, err := os.Open(absPath)
	if err != nil {
		return
	}

	err = json.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		return
	}

	cfg.FilePath = absPath
	log.INFO(AGENT, "Loaded config from file: %s", cfg.FilePath)
	return
}
