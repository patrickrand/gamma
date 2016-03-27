package agent

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

const (
	// StatusError is used when an error occurred running the Check.
	// This value can be specified by the used-defined command,
	// but is generally set by the Agent when a command fails or
	// has an invalid response format.
	StatusError = iota - 1

	// StatusOK is used when the Check has been deemed successful according
	// to the user-defined command (i.e. the Check is healthy).
	StatusOK

	// StatusWarn is used when the Check is not exactly healthy, but not
	// completely unhealthy such that the host is in a critical state. This
	// status is most useful as a predictor of a host that is getting ready
	// to enter an unhealthy/critical state.
	StatusWarn

	// StatusCritical is used when the Check is unhealthy and the host has
	// entered a critical state. This is the most common status type that
	// is pushed via a Handler.
	StatusCritical
)
