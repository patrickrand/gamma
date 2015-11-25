package agent

import (
	"bytes"
	"encoding/json"
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
	log.Debugf("(*Action).Run => (%s)", log.PrintJson(a))

	if b, err := exec.Command(a.CommandPath, a.CommandArgs...).Output(); err == nil {
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&output)
	}

	return output, err
}
