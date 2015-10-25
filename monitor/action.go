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

func (a *Action) Run() (output Output, err error) {
	log.DBUG("action", "(*Action).Run => (%s)", log.PrintJson(a))

	b, err := exec.Command(a.CommandPath, a.CommandArgs...).Output()
	if err != nil {
		log.EROR("action", "Error executing command (%s) => %s", err.Error(), log.PrintJson(a))
		return
	}

	err = json.NewDecoder(bytes.NewReader(b)).Decode(&output)
	if err != nil {
		log.EROR("action", "Unable to decode command output (%s) => %s", err.Error(), log.PrintJson(a))
	}

	return
}
