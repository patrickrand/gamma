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

func NewConfigFromFile(file string) (Config, error) {
	log.Debugf("[%s] config.NewConfigFromFile => %s", AGENT, file)

	var cfg Config
	absPath, err := filepath.Abs(file)
	if err != nil {
		return cfg, err
	}

	configFile, err := os.Open(absPath)
	if err != nil {
		return cfg, err
	}

	if err = json.NewDecoder(configFile).Decode(&cfg); err != nil {
		return cfg, err
	}

	cfg.FilePath = absPath
	log.Infof("[%s] Loaded config from file: %s", AGENT, cfg.FilePath)
	return cfg, err
}
