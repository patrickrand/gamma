package monitor

import (
	"bytes"
	"encoding/json"
	//	"log"
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

func NewAction(path string, args []string, interval time.Duration) *Action {
	return &Action{
		CommandPath: path,
		CommandArgs: args,
		Interval:    interval,
	}
}

func (a *Action) Run() (output Output, err error) {
	cmd := exec.Command(a.CommandPath, a.CommandArgs...)
	b, err := cmd.Output()
	if err != nil {
		return
	}

	err = json.NewDecoder(bytes.NewReader(b)).Decode(&output)
	return
}
