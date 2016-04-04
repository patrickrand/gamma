package gamma

import "time"

// A Check represents an process to be executed by an Agent on its host.
type Check struct {
	// ID uniquely identifies this Check on its Agent. It will be the easiest
	// way of identifying a Check in the Agent's logs and output, so it should
	// be clear in its naming.
	ID string `json:"id"`

	// Command is the string representation of the command this Check is executing.
	// This Command will be executed by shelling out to the host, and thus will
	// utilize the user PATH/profile/environment that is running the Agent.
	Command string `json:"command"`

	// Args is the list of whitespace delineated arguments passsed to the check's command.
	Args []string `json:"args"`

	// Interval is the time interval (in seconds) on which the Agent will run
	// this Check on its host.
	Interval time.Duration `json:"interval"`
}

func (check *Check) Run(executer Executer) Result {
	start := time.Now()
	code, message := executer.Execute(check.Command, check.Args...)
	return Result{
		ID:        check.ID,
		Command:   check.Command,
		Args:      check.Args,
		Interval:  check.Interval,
		StartTime: start,
		EndTime:   time.Now(),
		Code:      code,
		Message:   message,
	}
}
