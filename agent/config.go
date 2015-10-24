package agent

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	AgentName    string                 `json:"agent_name"`
	AgentVersion string                 `json:"agent_version"`
	Monitors     map[string]interface{} `json:"monitors"`
	Handlers     map[string]interface{} `json:"handlers"`
}

func NewConfigFromFile(file string) (cfg Config, err error) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return
	}

	configFile, err := os.Open(absPath)
	if err != nil {
		return
	}

	err = json.NewDecoder(configFile).Decode(&cfg)
	return
}

func (cfg Config) Validate() error {
	var errs Errors

	if !regexp.MustCompile("^[a-z][a-z_-]*[a-z]$").MatchString(cfg.AgentName) {
		errs = append(errs, errors.New("engine.Config.Validate: `name` must match `^[a-z][a-z-]*[a-z]$`"))
	}

	vers := strings.Split(cfg.AgentVersion, ".")
	if len(vers) != 3 {
		errs = append(errs, errors.New("engine.Config.ValidateVersion: invalid octet count"))
	}

	for i := range vers {
		_, err := strconv.ParseUint(vers[i], 10, 8)
		if err != nil {
			errs = append(errs, errors.New("engine.Config.ValidateVersion: semantic version can only contain octets"))
			break
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}
