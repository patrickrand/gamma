package agent

import (
	"time"
)

// A Result contains the data that represents the returned
// result response of a Check's execution. It is the data
// that will be pushed by the Agent.
type Result struct {
	// CheckID is the ID of the Check associated with this Result.
	CheckID string `json:"-"`

	// Output is response data of the associated Check's execution.
	Output `json:"output"`

	// Metadata provides runtime metadata about the check.
	Metadata `json:"metadata"`

	// Error is string message of the potential error that may
	// have occured during the running of a Check.
	Error string `json:"error,omitempty"`
}

// Results is the in-memory cache of the most recent results of each Check.
// The JSON representation of Results is the response body of requests to server.
var Results = make(map[string]*Result)

// An Output wraps the returned result response of
// the command executed by a Check. It acts as the API
// between Gamma and the user-defined Check commands.
type Output struct {
	// Status is the integer exit code of the command.
	// The accepted values of Status are detailed by the
	// StatusErr, StatusOK, StatusWarning, and StatusCritical
	// constant below. Status is a pointer type as a result
	// of its natural "zero-value" already being reserved by StatusOK.
	Status *int `json:"status"`

	// Message is the (optional) message response of the
	// command being executed in the Check. It is used to
	// provide deeper insight into the exit status of the
	// command.
	Message string `json:"message,omitempty"`
}

// Metadata stores runtime information for a given Check execution.
type Metadata struct {
	//Command is the command executed by the Check returning this Result.
	Command string `json:"command"`

	// StartTime is the starting time of a Check's execution.
	StartTime time.Time `json:"start_time"`

	// EndTime is the ending time of a Check's execution.
	EndTime time.Time `json:"end_time"`
}

const (
	// StatusErr is used when an error occured runnning the Check.
	// This value can be specified by the used-defined command,
	// but is generally set by the Agent when a command fails or
	// has an invalid response format.
	StatusErr = iota - 1

	// StatusOK is used when the Check has been deemed successful according
	// to the user-defined command (i.e. the Check is healthy).
	StatusOK

	// StatusWarning is used when the Check is not exactly healthy, but not
	// completely unhealthy such that the host is in a critical state. This
	// status is most useful as a predictor of a host that is getting ready
	// to enter an unhealthy/critical state.
	StatusWarning

	// StatusCritical is used when the Check is unhealthy and the host has
	// entered a critical state. This is the most common status type that
	// is pushed via a Handler.
	StatusCritical
)

// NewResult instantiates and returns a pointer to a new Result.
// The StartTime of the Result is automatically set, as well as various
// metadata about the Check that is being run.
func NewResult(check *Check) *Result {
	metadata := Metadata{
		Command:   check.Command,
		StartTime: time.Now(),
	}
	return &Result{
		CheckID:  check.ID,
		Metadata: metadata,
	}
}
