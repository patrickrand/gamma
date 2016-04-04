package gamma

import (
	"time"
)

// A Result contains the data that represents the returned
// result response of a Check's execution. It is the data
// that will be pushed by the Agent.
type Result struct {
	//CheckID is the ID of the Check that produced this Result.
	ID string `json:"id"`

	//Command is the command executed by the Check returning this Result.
	Command string `json:"command"`

	Args []string `json:"args"`

	// Interval is the time interval (in seconds) on which the Agent will run
	// this Check on its host.
	Interval time.Duration `json:"interval"`

	// StartTime is the starting time of a Check's execution.
	StartTime time.Time `json:"start_time"`

	// EndTime is the ending time of a Check's execution.
	EndTime time.Time `json:"end_time"`

	// Code is the returned integer code of the command.
	// The accepted values of Status are detailed by the
	// StatusError, StatusOK, StatusWarn, and StatusCritical
	// constant below. Status is a pointer type as a result
	// of its natural "zero-value" already being reserved by StatusOK.
	Code int `json:"code"`

	// Message is the (optional) message response of the
	// command being executed in the Check. It is used to
	// provide deeper insight into the exit status of the
	// command.
	Message string `json:"message,omitempty"`
}
