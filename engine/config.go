package engine

import (
	"encoding/json"
	"github.com/patrickrand/gamma/handler"
	"github.com/patrickrand/gamma/monitor"
	"os"
	"path/filepath"
)

type Config struct {
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Monitors []monitor.Monitor `json:"monitors"`
	Handlers []handler.Handler `json:"handlers"`
}

func (cfg Config) LoadFile(file string) error {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	configFile, err := os.Open(absPath)
	if err != nil {
		return err
	}

	return json.NewDecoder(configFile).Decode(&cfg)
}

func (cfg Config) Validate() error {
	var errs Errors

	if !regexp.MustCompile("^[a-z][a-z-]*[a-z]$").MatchString(cfg.AgentName) {
		errs = append(errs, errors.New("engine.Config.Validate: `name` must match `^[a-z][a-z-]*[a-z]$`"))
	}

	vers := strings.Split(cfg.Version, ".")
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

	return errs
}
