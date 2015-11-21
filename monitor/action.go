package monitor

import (
	"bytes"
	"encoding/json"
	log "github.com/patrickrand/gamma/log"
	"os/exec"
	"time"
)

type Output struct {
	Status  *int   `json:"status"`
	Message string `json:"message,omitempty"`
}

type Action struct {
	CommandPath string        `json:"command_path"`
	CommandArgs []string      `json:"command_args"`
	Interval    time.Duration `json:"interval"`
}

func (a *Action) Run() (Output, error) {
	log.Debugf("[%s] (*Action).Run => (%s)", MONITOR, log.PrintJson(a))

	var output Output
	b, err := exec.Command(a.CommandPath, a.CommandArgs...).Output()
	if err == nil {
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&output)
	}

	return output, err
}
