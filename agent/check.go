package agent

import (
	"encoding/json"
	"os/exec"
	"time"
)

// A Check represents an action to be executed by an Agent on its host
// on a given interval, and the list of Handlers responsible for pushing
// the Check's results.
type Check struct {
	// ID uniquely identifies this Check on its Agent. It will be the easiest
	// way of identifying a Check in the Agent's logs and output, so it should
	// be clear in its naming.
	ID string `json:"id"`

	// Command is the string representation of the command this Check is executing.
	// This Command will be executed by shelling out to the host, and thus will
	// utilize the user PATH/profile/environment that is running the Agent.
	Command string `json:"command"`

	Args []string `json:"args"`

	// Interval is the time interval (in seconds) on which the Agent will run
	// this Check on its host.
	Interval time.Duration `json:"interval"`
}

type Executer interface {
	Execute(cmd string, args ...string) (code int, message string)
}

type ShellExecuter struct{}

func (sh ShellExecuter) Execute(cmd string, args ...string) (code int, message string) {
	data, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return StatusError, "ShellExecuter failed to execute command: " + err.Error()
	}

	output := struct {
		Code    int    `json:"code"`
		Message string `json:"message,omitempty"`
	}{}

	if err := json.Unmarshal(data, &output); err != nil {
		return StatusError, "ShellExecuter failed to unmarshal command output: " + err.Error()
	}

	return output.Code, output.Message
}

func (check Check) Run(executer Executer) Result {
	start := time.Now()
	code, message := executer.Execute(check.Command, check.Args...)
	return Result{
		ID:         check.ID,
		Command:    check.Command,
		Args:       check.Args,
		Interval:   check.Interval,
		StartTime:  start,
		EndTime:    time.Now(),
		StatusCode: code,
		Message:    message,
	}
}
