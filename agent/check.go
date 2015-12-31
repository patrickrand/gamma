package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

	// AlertOn is the Output Status that indicates whether the result of this Check
	// should be pushed to its Handlers.
	AlertOn string `json:"alert_on"`

	// HandlerIDs is the list of IDs of the Handlers that this Check will use
	// to push its results.
	HandlerIDs []string `json:"handler_ids"`
}

// Exec runs a Check's Command and returns its Result.
func (c *Check) Exec() *Result {
	result := NewResult(c)
	defer func() {
		result.EndTime = time.Now()
		if result.Error != "" {
			result.Status = new(int)
			*result.Status = StatusErr
		}
	}()

	args := strings.Split(c.Command, " ")
	if len(args) < 1 {
		result.Error = fmt.Sprintf("invalid command: %s", c.Command)
		return result
	}

	data, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		result.Error = fmt.Sprintf("failed to execute check command: %v", err)
	} else if err = json.NewDecoder(bytes.NewReader(data)).Decode(&(result.Output)); err != nil {
		result.Error = fmt.Sprintf("failed to decode check command output: %v", err)
	} else if result.Status == nil {
		result.Error = fmt.Sprintf("result output status is nil")
	}

	return result
}

// ShouldAlert determines whether a given Result Output Status should be
// alerted on by its Handlers, according to the Check's AlertOn value.
func (c *Check) ShouldAlert(status *int) bool {
	if status == nil {
		return false
	}

	if *status == StatusErr {
		return true
	}

	switch c.AlertOn {
	case "ok":
		return *status >= StatusOK
	case "warning":
		return *status >= StatusWarning
	case "critical":
		return *status >= StatusCritical
	}

	log.Printf("unrecognized \"alert_on\" value: %s", c.AlertOn)
	return false
}
