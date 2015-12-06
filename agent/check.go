package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
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

	// Interval is the time interval (in seconds) on which the Agent will run
	// this Check on its host.
	Interval time.Duration `json:"interval"`

	// HandlerIDs is the list of IDs of the Handlers that this Check will use
	// to push its results.
	HandlerIDs []string `json:"handler_ids"`

	*Result `json:"-"`
}

// Exec runs a Check's Command and returns its Result.
func (c *Check) Exec() *Result {
	result := NewResult(c.ID)
	defer func() { result.EndTime = time.Now() }()

	command := strings.Split(c.Command, " ")
	if len(command) < 1 {
		result.Error = fmt.Sprintf("invalid command: %s", c.Command)
		return result
	}

	data, err := exec.Command(command[0], command[1:]...).Output()
	if err != nil {
		result.Error = fmt.Sprintf("failed to execute check command: %v", err)
	} else if err = json.NewDecoder(bytes.NewReader(data)).Decode(&(result.Output)); err != nil {
		result.Error = fmt.Sprintf("failed to decode check command output: %v", err)
	} else if result.Status == nil {
		result.Error = fmt.Sprintf("result output status is nil")
	}

	return result
}
